package api

import (
	"encoding/json"
	"github.com/user/apiserver/app"
	"net/http"
)

// MasterGoalID int
// Name         string
// Description  string
// ImageURL     string
// Status       int
// DateAdded    time.Time

func MasterGoalsGet(r *http.Request) (int, []byte) {

	//---------------------------------------------------------------------------------
	// run the query to get the records from the database
	//---------------------------------------------------------------------------------
	sql := "SELECT masterGoalID, name, description, imageURL, status, dateAdded " +
		"FROM MASTER_GOAL " +
		"WHERE status = 1 " +
		"ORDER BY name"
	rows, err := app.Db.Query(sql)
	if err != nil {
		app.Error(err.Error())
		return 0, nil
	}

	goals := []MasterGoal{}
	for rows.Next() {
		g := MasterGoal{}
		rows.Scan(&g.MasterGoalID, &g.Name, &g.Description, &g.ImageURL, &g.Status, &g.DateAdded)

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
