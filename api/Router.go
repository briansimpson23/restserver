package api

import (
	"fmt"
	"github.com/briansimpson23/restserver/app"
	"net/http"
	"strings"
)

//--------------------------------------------------------------------
// this handler receives the HTTP request and processes the request
//--------------------------------------------------------------------
func Router(r *http.Request) (int, []byte) {

	if app.Cfg["log.showHttpHeaders"] == "On" {
		logHTTPHeaders(r)
	}

	// capture the name of the API that's being requested
	url := strings.Split(strings.TrimRight(r.URL.Path[1:], "/"), "/")

	// url := strings.Split(r.URL.Path[1:], "/")
	apiName := url[0]

	app.Debug(fmt.Sprintf("apiName: %s", apiName))
	app.Debug(fmt.Sprintf("url: %s", r.URL.Path))

	//------------------------------------------------------------------
	// Figure out which API function to call
	//
	// TODO - Ultimately this can be done dynamically with interfaces.
	//        I need to learn a little more about interfaces first.
	//-------------------------------------------------------------------
	app.Debug(fmt.Sprintf("len[%d]", len(url)))

	// the members APIs
	if apiName == "customers" {
		customer := Customer{}

		if r.Method == "GET" {
			if len(url) == 2 {
				return customer.Get(r)
			}
		}
		if r.Method == "POST" {
			if len(url) == 1 {
				return customer.Create(r)
			}
		}

	}

	return 400, []byte(`"ERROR - UNKNOWN APINAME"`)

}
