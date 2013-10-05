package api

import (
	"fmt"
	"github.com/briansimpson23/restserver/app"
	"net/http"
	"strings"
)

func IsAPIKeyValid(r *http.Request) bool {

	//TODO - ultimately we want to keep the active keys in memory so
	//       don't have to hit the database on every request
	//     - do something like check the active ones on a request
	//       if no match then check the DB.
	//     - put a maximum limit on how long the keys are stored in RAM
	//       maybe expire them after 2 minutes.  This way if we invalidate
	//       a key in the database it will get killed off in a timely manner

	apiKey := strings.Trim(fmt.Sprintf("%s", r.Header["Api-Key"]), "[]")

	// keyStatus := strings.ToLower(app.Cfg["security.apiKey"])

	if strings.EqualFold(app.Cfg["security.apiKey"], "off") {
		app.Warn(fmt.Sprintf("security.apiKey is turned off"))
		return true
	}

	// parse the Api-Key from the HTTP Header
	// apiKey := strings.Trim(fmt.Sprintf("%s", r.Header["Api-Key"]), "[]")
	app.Debug(fmt.Sprintf("IsAPIKeyValid() - apiKey = %s", apiKey))

	// check and see if the Api-Key is in the Database
	app.Debug(fmt.Sprintf("IsAPIKeyValid() - running the query to find the apiKey[%s]", apiKey))
	row := app.Db.QueryRow("SELECT appID FROM API_KEY WHERE apiKey = ?", apiKey)

	err := row.Scan(&app.AppID)
	if err != nil {
		app.Warn(fmt.Sprintf("unable to verify apiKey - %s", err.Error))
		return false
	}

	app.Debug("IsAPIKeyValid() - returning true")
	return true
}
