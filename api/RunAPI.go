package api

import (
	"fmt"
	"github.com/briansimpson23/restserver/app"
	"net/http"
	"time"
)

type Response struct {
	Code   int
	Status string
	Data   interface{}
}

//--------------------------------------------------------------------
// this handler receives the HTTP request and processes the request
//--------------------------------------------------------------------
func RunAPI(w http.ResponseWriter, r *http.Request, startTime time.Time) {

	app.StartTime = startTime
	//app.Debug(fmt.Sprintf("startTime : %s", app.StartTime))

	// set the url as the loggerName.  The loggerName is used in the
	// log output so we can track the incoming request.
	app.LoggerName = r.Method + " " + r.URL.Path

	app.Debug("+------------------------------------------------------------")
	app.Debug("|")
	//-----------------------------------------------
	// check the APIKey and make sure it is valid
	//-----------------------------------------------
	if !IsAPIKeyValid(r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(401)
		fmt.Fprintf(w, "%s", BuildErrorMessage(401, "AccessDenied", 2002))
		return
	}

	//-----------------------------------------------
	// return the API response to the requestor
	//-----------------------------------------------
	//TODO - this needs to log the actual error code we are returning
	//TODO - right now it is always logging an HTTP/200
	w.Header().Set("Content-Type", "application/json")
	err, output := Router(r)
	if err > 0 {
		w.WriteHeader(err)
	}

	app.Debug("|")
	app.Debug("+------------------------------------------------------------")
	app.Info("returned HTTP/200")
	fmt.Fprintf(w, "%s", output)
	return

}
