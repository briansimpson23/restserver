package api

import (
	"encoding/json"
	"fmt"
	"github.com/user/apiserver/app"
	"net/http"
	"strconv"
	"strings"
)

func MemberGoalGet(r *http.Request) (int, []byte) {

	app.Debug("entered MemberGoalGet()")
	//----------------------------------------------------------
	// validate the input and parse the memberID from the URL
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

	app.Debug(fmt.Sprintf("memberID[%d]", memberID))
	app.Debug(fmt.Sprintf("goalID[%d]", goalID))

	//-------------------------------------------------------------------
	// run the query to get the records from the database
	//-------------------------------------------------------------------
	sql := "SELECT m.goalID, IF(m.goalName = '', g.name, m.goalName) as goalName, m.goalStatus, m.dateAdded, m.dateCompleted " +
		"FROM MEMBER_GOAL m " +
		"LEFT JOIN MASTER_GOAL g ON m.masterGoalID = g.masterGoalID " +
		"WHERE m.memberID=? AND m.goalID = ? AND m.goalStatus > 0"

	// sql := "SELECT goalID, goalName, goalStatus, dateAdded, dateCompleted " +
	// 	" FROM MEMBER_GOAL WHERE memberID = ? AND goalID = ? AND goalStatus <> 0"
	app.Debug(sql)
	row := app.Db.QueryRow(sql, memberID, goalID)

	goal := MemberGoal{}
	goal.MemberID = memberID
	err := row.Scan(&goal.GoalID, &goal.GoalName, &goal.GoalStatus, &goal.DateAdded, &goal.DateCompleted)
	if err != nil {
		app.Error(fmt.Sprintf("unable to load memberGoal memberID[%d]  goalID[%d]", memberID, goalID))
		app.Error(err.Error())
		return 0, nil
	}

	// //----------------------------------------------------------
	// // build the output
	// //----------------------------------------------------------
	// wrapper := make(map[string]MemberGoal)
	// wrapper["MemberGoal"] = goal

	// output, err := json.Marshal(wrapper)
	// if err != nil {
	// 	panic(err)
	// }

	// // app.Info(fmt.Sprintf("processed memberGoal[%d][%d]", memberID, goalID))
	// return 0, output

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = 200
	resp.Status = "ok"
	resp.Data = goal

	output, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return 0, output

}
