package api

import (
	"encoding/json"
	"github.com/user/apiserver/app"
	"net/http"
	"strings"
)

func MemberGet(r *http.Request) (int, []byte) {

	//----------------------------------------------------------
	// validate the input and parse the memberID from the URL
	//----------------------------------------------------------
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 2 {
		return 400, BuildErrorMessage(400, "Invalid MemberID Format", 2006)
	}
	memberID := tmp[1]

	//----------------------------------------------------------
	// run the query to get the records from the database
	//----------------------------------------------------------
	sql := "SELECT memberID, firstName, lastName, email, zipCode, gender, facebookUserID, twitterUsername, memberStatusID, marketingSourceID " +
		"FROM MEMBER WHERE memberID=?"
	app.Debug(sql)
	row := app.Db.QueryRow(sql, memberID)

	//----------------------------------------------------------
	// parse the database record into the Member struct
	//----------------------------------------------------------
	m := Member{}
	app.Debug("parsing the returned row")
	err := row.Scan(&m.MemberID, &m.FirstName, &m.LastName, &m.Email, &m.ZipCode, &m.Gender, &m.FacebookID, &m.TwitterUsername,
		&m.MemberStatusID, &m.MarketingSourceID)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			app.Error("unknown memberID[" + memberID + "]")
			return 404, BuildErrorMessage(404, "Not Found", 2003)
		} else {
			app.Error(err.Error())
			return 500, BuildErrorMessage(500, "System Error", 2004)
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
