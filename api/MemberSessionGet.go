package api

import (
	"encoding/json"
	"github.com/user/apiserver/app"
	"net/http"
	"strings"
)

func MemberSessionGet(r *http.Request) (int, []byte) {

	var o MemberSession

	//----------------------------------------------------------
	// get the sessionID from the URL
	//----------------------------------------------------------
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 2 {
		return 400, BuildErrorMessage(400, "Invalid Input Format", 2006)
	}
	sessionID := tmp[1]

	//----------------------------------------------------------
	// run the query to get the records from the database
	//----------------------------------------------------------
	sql := "SELECT m.memberID, m.firstName, m.lastName, m.email, m.zipCode, m.gender, " +
		"m.facebookUserID, m.twitterUsername, ms.loginType " +
		"FROM MEMBER m " +
		"JOIN MEMBER_SESSION ms ON m.memberID = ms.memberID " +
		"WHERE ms.sessionID = ? AND ms.logoffTime IS NULL"
	app.Debug(sql)
	row := app.Db.QueryRow(sql, sessionID)

	//----------------------------------------------------------
	// parse the database row into the output variable
	//----------------------------------------------------------
	err := row.Scan(&o.MemberID, &o.FirstName, &o.LastName, &o.Email, &o.ZipCode, &o.Gender, &o.FacebookID, &o.Twitter, &o.LoginType)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return 404, BuildErrorMessage(404, "Unknown SessionID", 2006)
		} else {
			return 403, BuildErrorMessage(403, "Invalid Input Format", 2006)
		}

	}
	o.SessionID = sessionID

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
