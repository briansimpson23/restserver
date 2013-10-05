package api

import (
	"../app"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func MemberGoalCreate(r *http.Request) (int, []byte) {

	//----------------------------------------------------------
	// get the payload from the HTTP request and parse it
	// into the needed variables
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
	//TODO - handle this error
	memberID, _ := strconv.ParseInt(tmp[1], 0, 64)
	masterGoalID := r.FormValue("masterGoalID")
	app.Debug(fmt.Sprintf("    memberID[%d]", memberID))
	app.Debug(fmt.Sprintf("masterGoalID[%d]", masterGoalID))

	//-------------------------------------------------------
	// populate the MemberGoal struct
	//-------------------------------------------------------
	g := MemberGoal{}
	if masterGoalID != "" {
		g.MasterGoalID, _ = strconv.Atoi(masterGoalID)
	}
	g.GoalName = strings.TrimSpace(r.FormValue("goalName"))
	app.Debug(fmt.Sprintf("MemberGoal %s", g))

	//-------------------------------------------------------
	// validate the input
	//-------------------------------------------------------
	if memberID < 1 {
		app.Error(fmt.Sprintf("Invalid memberID of [%d]", memberID))
		return 400, BuildErrorMessage(400, "Invalid memberID", 2008)
	}

	if g.GoalName == "" && g.MasterGoalID == 0 {
		app.Error("Missing goalName")
		return 400, BuildErrorMessage(400, "Missing Goal Name", 2008)
	}

	// get the highest goalID for this member
	sql := "SELECT max(goalID) as goalID FROM MEMBER_GOAL WHERE memberID = ?"
	var maxGoalID int = 0
	row := app.Db.QueryRow(sql, memberID)
	err := row.Scan(&maxGoalID)
	if err != nil {
		app.Debug(fmt.Sprintf("ERROR CHECKING MAX GOAL ID - %v", err))
		//panic(err)
	}
	g.GoalID = maxGoalID + 1
	app.Debug(fmt.Sprintf("g.GoalID %d", g.GoalID))

	//----------------------------------------------------------
	// run the query to get the records from the database
	//----------------------------------------------------------
	sql = "INSERT INTO MEMBER_GOAL " +
		"SET " +
		"memberID = ?, goalID = ?, masterGoalID = ?, goalName = ?, dateAdded = now() "
	//s, _ := app.Db.Prepare(sql)
	app.Debug(sql)
	_, err = app.Db.Exec(sql, memberID, g.GoalID, g.MasterGoalID, g.GoalName)
	if err != nil {
		panic(err)
	}

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = 200
	resp.Status = "ok"
	resp.Data = g

	output, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return 0, output

}
