package api

import (
	"encoding/json"
	"fmt"
	"github.com/user/apiserver/app"
	"net/http"
	"strings"
)

func MemberCreate(r *http.Request) (int, []byte) {

	//----------------------------------------------------------
	// get the payload from the HTTP request and parse it
	// into the needed variables
	//----------------------------------------------------------
	m := Member{}
	m.FirstName = strings.TrimSpace(r.FormValue("firstName"))
	m.LastName = strings.TrimSpace(r.FormValue("lastName"))
	m.Email = strings.TrimSpace(r.FormValue("email"))
	m.ZipCode = strings.TrimSpace(r.FormValue("zipCode"))
	m.Gender = strings.TrimSpace(r.FormValue("gender"))
	m.FacebookID = strings.TrimSpace(r.FormValue("facebookID"))
	m.FacebookCity = strings.TrimSpace(r.FormValue("facebookCity"))
	m.FacebookLocationID = strings.TrimSpace(r.FormValue("facebookLocationID"))
	password := strings.TrimSpace(r.FormValue("password"))

	// determine the type of registration:  local, facebook, or twitter
	regType := "local"
	if m.FacebookID != "" {
		regType = "facebook"
	}

	//----------------------------------------------------------------------------
	// THIS IS ALL DEBUG OUTPUT
	//
	app.Debug("    +------------------------------------------------------")
	app.Debug("    | MEMBER CREATE DEBUG INFORMATION")
	app.Debug("    |")
	app.Debug("    |    INPUT")
	app.Debug(fmt.Sprintf("    |           FirstName : %s", m.FirstName))
	app.Debug(fmt.Sprintf("    |            LastName : %s", m.LastName))
	app.Debug(fmt.Sprintf("    |               Email : %s", m.Email))
	app.Debug(fmt.Sprintf("    |             ZipCode : %s", m.ZipCode))
	app.Debug(fmt.Sprintf("    |              Gender : %s", m.Gender))
	app.Debug(fmt.Sprintf("    |            Password : %s", password))
	app.Debug(fmt.Sprintf("    |          FacebookID : %s", m.FacebookID))
	app.Debug(fmt.Sprintf("    |        FacebookCity : %s", m.FacebookCity))
	app.Debug(fmt.Sprintf("    |  FacebookLocationID : %s", m.FacebookLocationID))
	app.Debug("    |")
	app.Debug("    |    OTHER")
	app.Debug(fmt.Sprintf("    |             regType : %s", regType))
	app.Debug("    |")
	app.Debug("    +------------------------------------------------------")
	//
	// END ALL THE DEBUG OUTPUT
	//----------------------------------------------------------------------------

	//-------------------------------------------------------------------------------------------------------------
	//-- validate that we've received everything all needed input
	//-------------------------------------------------------------------------------------------------------------
	if m.FirstName == "" {
		app.Debug("missing firstName")
		return 400, BuildErrorMessage(400, "Missing First Name", 2008)
	}
	if m.LastName == "" {
		app.Debug("missing LastName")
		return 400, BuildErrorMessage(400, "Missing Last Name", 2008)
	}
	if m.Email == "" {
		app.Debug("missing Email Address")
		return 400, BuildErrorMessage(400, "Missing Email Address", 2008)
	}
	if m.Gender == "" {
		app.Debug("missing Gender")
		return 400, BuildErrorMessage(400, "Missing Gender", 2008)
	}

	//-------------------------------------------------------------------------------------------------------------
	//-- check to see if this email already exists.
	//-------------------------------------------------------------------------------------------------------------
	sql := "SELECT memberID FROM MEMBER WHERE email = ?"
	app.Debug(sql)
	row := app.Db.QueryRow(sql, m.Email)
	err := row.Scan(&m.MemberID)
	if err != nil {
		// if the error is anything except "no rows found then something unplanned is wrong"
		if err.Error() != "sql: no rows in result set" {
			app.Error(fmt.Sprint("%s", err))
			return 500, BuildErrorMessage(500, "Unexpected error", 2008)
		}
	}

	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	//
	// Handle everything if this is a local registration
	//TODO putting this all in one file is a bad idea.  It will make this registration function too large and hard to maintain.
	//
	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	if regType == "local" {

		if m.ZipCode == "" {
			app.Debug("missing ZipCode")
			return 400, BuildErrorMessage(400, "Missing ZipCode", 2008)
		}
		if password == "" {
			app.Debug("missing Password")
			return 400, BuildErrorMessage(400, "Missing Password", 2008)
		}

		// if this is a local registration then we've got a problem if the memberID is > 0
		if m.MemberID > 0 {
			app.Error(fmt.Sprintf("email address already exists - [%s]", m.Email))
			return 409, BuildErrorMessage(409, "Duplicate Email Address", 2008)
		}

		//everything is OK so insert the new member record
		sql = "INSERT INTO MEMBER " +
			"SET " +
			"firstName = ?, lastName = ?, email = ?, zipCode = ?, gender = ?, localPassword = ?, " +
			"memberStatusID = 1, marketingSourceID = 1, lastLoginDate = now(), dateAdded = now() "
		//s, _ := app.Db.Prepare(sql)
		app.Debug(sql)
		result, err := app.Db.Exec(sql, m.FirstName, m.LastName, m.Email, m.ZipCode, m.Gender, password)
		if err != nil {
			panic(err)
		}
		m.MemberID, _ = result.LastInsertId()

		// convert the output type into json and return it
		o := make(map[string]Member)
		o["Member"] = m

		output, err := json.Marshal(o)
		if err != nil {
			panic(err)
		}

		return 0, output
	}

	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	//
	// Handle everything if this is a facebook registration
	//TODO putting this all in one file is a bad idea.  It will make this registration function too large and hard to maintain.
	//
	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	//-------------------------------------------------------------------------
	if regType == "facebook" {

		//TODO validate that we've received all the needed input
		// check for facebookUserID, facebookCity, facebookLocationID
		// we should have already checked for everything else above
		if m.FacebookID == "" {
			app.Debug("missing FacebookID")
			return 400, BuildErrorMessage(400, "Missing FacebookID", 2008)
		}
		if m.FacebookCity == "" {
			app.Debug("missing FacebookCity")
			return 400, BuildErrorMessage(400, "Missing FacebookCity", 2008)
		}
		if m.FacebookLocationID == "" {
			app.Debug("missing FacebookLocationID")
			return 400, BuildErrorMessage(400, "Missing FacebookLocationID", 2008)
		}

		// check to see if we have the facebookUserID on file.  We will compare it against the memberID we loaded
		// earlier when checking to see if we already had the email address on file.
		var memberID_1 int64

		sql := "SELECT memberID FROM MEMBER WHERE facebookUserID = ?"
		row := app.Db.QueryRow(sql, m.FacebookID)
		err := row.Scan(&memberID_1)
		if err != nil {
			// if the error is anything except "no rows found then something unplanned is wrong"
			if err.Error() != "sql: no rows in result set" {
				app.Error(fmt.Sprint("%s", err))
				return 500, BuildErrorMessage(500, "Unexpected error", 2008)
			}
		}

		// check to see if have the FacebookUserID on file but not the emailAddress.
		// that means we need to update the member with the emailAddress
		if m.MemberID == 0 && memberID_1 > 0 {
			// this means we have the FacebookUserID but not the EmailAddress
			m.MemberID = memberID_1
			sql := "UPDATE MEMBER " +
				"SET " +
				"firstName = ?, lastName = ?, email = ?, gender = ?, facebookCity = ?, facebookLocationID = ? " +
				"WHERE memberID = ?"
			app.Debug(sql)
			_, err := app.Db.Exec(sql, m.FirstName, m.LastName, m.Email, m.Gender, m.FacebookCity, m.FacebookLocationID, m.MemberID)
			if err != nil {
				app.Error("error updating member record with the email address")
				app.Error(fmt.Sprintf("%s", err))
			}

		} else {
			// this means we have the EmailAddress but not the FacebookUserID
			sql := "UPDATE MEMBER " +
				"SET " +
				"firstName = ?, lastName = ?, gender = ?, facebookUserID = ?, facebookCity = ?, facebookLocationID = ? " +
				"WHERE memberID = ?"
			app.Debug(sql)
			_, err := app.Db.Exec(sql, m.FirstName, m.LastName, m.Gender, m.FacebookID, m.FacebookCity, m.FacebookLocationID, m.MemberID)
			if err != nil {
				app.Error("error updating member record with the FacebookUserID")
				app.Error(fmt.Sprintf("%s", err))
			}

		}

		//----------------------------------------------------------
		// build the output
		//----------------------------------------------------------
		resp := Response{}
		resp.Code = 200
		resp.Status = "ok"
		resp.Data = m

		output, err := json.Marshal(resp)
		if err != nil {
			panic(err)
		}

		return 0, output
	}

	return 400, BuildErrorMessage(400, "invalid registration type", 2008)

}
