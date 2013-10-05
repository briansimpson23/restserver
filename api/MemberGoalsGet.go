package api

import (
	"encoding/json"
	"fmt"
	"github.com/user/apiserver/app"
	"net/http"
	"strconv"
	"strings"
)

func MemberGoalsGet(r *http.Request) (int, []byte) {

	//----------------------------------------------------------
	// validate the input and parse the memberID from the URL
	//----------------------------------------------------------
	//TODO we are splitting this URL in the router and in each API
	//TODO need to change this process and centralize it
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 3 {
		e := make(map[string]string)
		e["ERROR"] = "INVALID_INPUT_FORMAT"
		b, _ := json.Marshal(e)
		return 400, b
	}
	memberID := tmp[1]
	filter := r.FormValue("filter")

	app.Debug(fmt.Sprintf("memberID[%d]", memberID))
	app.Debug(fmt.Sprintf("filter[%s]", filter))

	//-------------------------------------------------------------------
	// run the query to get the records from the database
	//-------------------------------------------------------------------
	sql := "SELECT m.goalID, IF(m.goalName = '', g.name, m.goalName) as goalName, m.masterGoalID, m.goalStatus " +
		"FROM MEMBER_GOAL m " +
		"LEFT JOIN MASTER_GOAL g ON m.masterGoalID = g.masterGoalID " +
		"WHERE m.memberID=? AND m.goalStatus > 0 " +
		"ORDER BY m.goalStatus ASC, m.dateCompleted DESC, goalName"
	app.Debug(sql)
	rows, err := app.Db.Query(sql, memberID)
	if err != nil {
		app.Error(err.Error())
		return 0, nil
	}

	goals := []MemberGoal{}
	for rows.Next() {
		g := MemberGoal{}
		g.MemberID, _ = strconv.ParseInt(memberID, 0, 64)
		rows.Scan(&g.GoalID, &g.GoalName, &g.MasterGoalID, &g.GoalStatus)

		//goal := MemberGoal{goalID, goalName}
		goals = append(goals, g)
	}

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = 200
	resp.Status = "ok"
	resp.Data = goals

	output, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return 0, output

}
