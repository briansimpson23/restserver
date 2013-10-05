package api

import (
	"../app"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func MemberGoalUpdate(r *http.Request) (int, []byte) {

	app.Debug("entered MemberGoalUpdate()")
	//----------------------------------------------------------
	// validate the input and parse the input from the URL
	// expecting: /members/<memberID>/goals/<goalID>
	//            status=complete
	//----------------------------------------------------------
	//TODO we are splitting this URL in the router and in each API
	//TODO need to change this process and centralize it
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 4 {
		app.Error(fmt.Sprintf("MemberGoalGet() - INVALID_INPUT_FORMAT"))
		e := make(map[string]string)
		e["ERROR"] = "INVALID_INPUT_FORMAT"
		b, _ := json.Marshal(e)
		return 400, b
	}
	memberID, _ := strconv.ParseInt(tmp[1], 0, 64)
	goalID, _ := strconv.ParseInt(tmp[3], 0, 64)
	status := strings.TrimSpace(r.FormValue("status"))

	app.Debug(fmt.Sprintf("memberID[%d]", memberID))
	app.Debug(fmt.Sprintf("  goalID[%d]", goalID))
	app.Debug(fmt.Sprintf("  status[%s]", status))

	//----------------------------------------------------------
	// RIGHT NOW THIS ONLY COMPLETES
	//----------------------------------------------------------
	if status == "complete" {
		sql := "UPDATE MEMBER_GOAL SET goalStatus = 2, dateCompleted = NOW() WHERE memberID = ? AND goalID = ? AND goalStatus = 1"
		//s, _ := app.Db.Prepare(sql)
		app.Debug(sql)
		_, err := app.Db.Exec(sql, memberID, goalID)
		if err != nil {
			panic(err)
		}
	} else if r.Method == "DELETE" {
		sql := "UPDATE MEMBER_GOAL SET goalStatus = 0, dateDeleted = NOW() WHERE memberID = ? AND goalID = ? AND goalStatus > 0"
		//s, _ := app.Db.Prepare(sql)
		app.Debug(sql)
		_, err := app.Db.Exec(sql, memberID, goalID)
		if err != nil {
			panic(err)
		}
	}

	//---------------------------------------------------
	// convert the output type into json and return it
	//---------------------------------------------------
	o := make(map[string]string)
	o["MemberGoal"] = "ok"

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
