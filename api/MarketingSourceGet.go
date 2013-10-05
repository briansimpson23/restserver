package api

import (
	"encoding/json"
	"fmt"
	"github.com/user/apiserver/app"
	"net/http"
	"strings"
)

func MarketingSourceGet(r *http.Request) (int, []byte) {

	//----------------------------------------------------------
	// validate the input and parse the memberID from the URL
	//----------------------------------------------------------
	tmp := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")
	if len(tmp) != 2 {
		e := make(map[string]string)
		e["ERROR"] = "INVALID_INPUT_FORMAT"
		b, _ := json.Marshal(e)
		return 403, b
	}
	marketingSourceID := tmp[1]

	app.Info("loading marketingSourceID[" + marketingSourceID + "]")

	//----------------------------------------------------------
	// run the query to get the records from the database
	//----------------------------------------------------------
	sql := "SELECT marketingSourceID, marketingSourceName " +
		"FROM MARKETING_SOURCE WHERE marketingSourceID = ?"
	app.Debug(sql)
	row := app.Db.QueryRow(sql, marketingSourceID)

	//----------------------------------------------------------
	// parse the database record into the Member struct
	//----------------------------------------------------------
	m := MarketingSource{}
	app.Debug("parsing the returned row")
	err := row.Scan(&m.MarketingSourceID, &m.MarketingSourceName)
	if err != nil {
		app.Error(err.Error())
		e := make(map[string]string)
		if err.Error() == "sql: no rows in result set" {
			e["ERROR"] = "UNKNOWN MARKETING_SOURCE_ID"
		} else {
			e["ERROR"] = fmt.Sprintf("%s", err)
		}

		b, _ := json.Marshal(e)
		return 404, b
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
