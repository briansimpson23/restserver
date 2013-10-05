package api

import (
	"../app"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func MemberSessionTerminate(r *http.Request) (int, []byte) {

	type Output struct {
		sessionID string
		status    string
	}
	var o Output

	//----------------------------------------------------------
	// get the sessionID from the URL
	//----------------------------------------------------------
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 2 {
		return 403, BuildErrorMessage(403, "Invalid Input Format", 2006)
	}
	sessionID := tmp[1]
	app.Debug("received logoff request for sessionID[" + sessionID + "]")

	//----------------------------------------------------------
	// run the query to update the session in the database
	//----------------------------------------------------------
	sql := "UPDATE MEMBER_SESSION SET logoffTime = now() WHERE sessionID = ? AND logoffTime IS NULL"
	app.Debug(sql)
	result, err := app.Db.Exec(sql, sessionID)
	if err != nil {
		panic(err)
	}
	app.Debug(fmt.Sprintf("SQL RESULT: %v", result))

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
