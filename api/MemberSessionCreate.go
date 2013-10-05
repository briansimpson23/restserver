package api

import (
	"../app"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func MemberSessionCreate(r *http.Request) (int, []byte) {

	var o MemberSession

	//----------------------------------------------------------
	// get the sessionID from the URL
	//----------------------------------------------------------
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 2 {
		app.Error(fmt.Sprintf("MemberSessionCreate() - INVALID_INPUT_FORMAT"))

		e := make(map[string]string)
		e["ERROR"] = "INVALID_INPUT_FORMAT"
		b, _ := json.Marshal(e)
		return 400, b
	}
	sessionID := tmp[1]

	//----------------------------------------------------------
	// get the payload from the HTTP request and parse it
	// into the needed variables
	//----------------------------------------------------------
	loginType := r.FormValue("loginType")
	input1 := r.FormValue("input1")
	input2 := r.FormValue("input2")
	remoteIP := r.FormValue("remoteIP")

	app.Debug(fmt.Sprintf("loginType[%s]", loginType))
	app.Debug(fmt.Sprintf("sessionID[%s]", sessionID))
	app.Debug(fmt.Sprintf("   input1[%s]", input1))
	app.Debug(fmt.Sprintf("   input1[%s]", input2))
	app.Debug(fmt.Sprintf(" remoteIP[%s]", remoteIP))

	//-------------------------------------------------------------
	// run the query to get the records from the database
	// we are doing this to validate the username and password
	//-------------------------------------------------------------
	if loginType == "local" {
		sql := "SELECT memberID, firstName, lastName, email, zipCode, gender, facebookUserID, twitterUsername FROM MEMBER "
		sql += "WHERE email = ? AND localPassword = ?"
		app.Debug(sql)
		row := app.Db.QueryRow(sql, input1, input2)
		err := row.Scan(&o.MemberID, &o.FirstName, &o.LastName, &o.Email, &o.ZipCode, &o.Gender, &o.FacebookID, &o.Twitter)
		if err != nil {
			e := make(map[string]string)
			if err.Error() == "sql: no rows in result set" {
				app.Error("INVALID_CREDENTIALS")
				e["ERROR"] = "INVALID_CREDENTIALS"
			} else {
				app.Error(err.Error())
				e["ERROR"] = fmt.Sprintf("%s", err)
			}

			b, _ := json.Marshal(e)
			return 403, b
		}
	} else {
		sql := "SELECT memberID, firstName, lastName, email, zipCode, gender, facebookUserID, twitterUsername FROM MEMBER "
		sql += "WHERE facebookUserID = ?"
		app.Debug(sql)
		row := app.Db.QueryRow(sql, input1)
		err := row.Scan(&o.MemberID, &o.FirstName, &o.LastName, &o.Email, &o.ZipCode, &o.Gender, &o.FacebookID, &o.Twitter)
		if err != nil {
			e := make(map[string]string)
			if err.Error() == "sql: no rows in result set" {
				app.Error("INVALID_CREDENTIALS")
				e["ERROR"] = "INVALID_CREDENTIALS"
			} else {
				app.Error(err.Error())
				e["ERROR"] = fmt.Sprintf("%s", err)
			}

			b, _ := json.Marshal(e)
			return 403, b
		}
	}

	//---------------------------------------------------------------
	// parse the database row into the output variable.
	//
	// if the query is empty then return the "INVALID_CREDENTIALS"
	//---------------------------------------------------------------
	o.SessionID = sessionID
	o.LoginType = loginType

	//----------------------------------------------------------
	// create the session entry
	//----------------------------------------------------------
	sql := "INSERT INTO MEMBER_SESSION " +
		"SET " +
		"sessionID = ?, memberID = ?, loginTime = now(), logOffTime = null, lastActiveDate = now(), " +
		"remoteIPAddress = ?, loginType=?"
	app.Debug(sql)
	_, err := app.Db.Exec(sql, o.SessionID, o.MemberID, remoteIP, loginType)
	if err != nil {
		//TODO - need to add better error handling here
		panic(err)
	}

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = 200
	resp.Status = "ok"
	resp.Data = o

	output, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return 0, output
}
