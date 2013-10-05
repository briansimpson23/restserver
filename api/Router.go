package api

import (
	"../app"
	"fmt"
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

	// the mastergoals APIs
	if apiName == "mastergoals" {
		if r.Method == "GET" {
			if len(url) == 1 {
				return MasterGoalsGet(r)
			}
		}
	}

	// the members APIs
	if apiName == "members" {

		if r.Method == "GET" {

			if len(url) == 2 {
				return MemberGet(r)
			}

			if len(url) == 3 && url[2] == "goals" {
				app.Debug("calling MemberGoalsGet()")
				return MemberGoalsGet(r)
			}

			if len(url) == 3 && url[2] == "ads" {
				app.Debug("calling MemberAdsGet()")
				return MemberAdsGet(r)
			}

			if len(url) == 4 && url[2] == "goals" {
				app.Debug("calling MemberGoalGet()")
				return MemberGoalGet(r)
			}

			if len(url) == 5 && url[2] == "goals" && url[4] == "ads" {
				app.Debug("calling MemberGoalAdsGet()")
				return MemberGoalAdsGet(r)
			}

		}

		if r.Method == "POST" {

			if len(url) == 1 {
				return MemberCreate(r)
			}

			if len(url) == 3 && url[2] == "goals" {
				return MemberGoalCreate(r)
			}

		}

		if r.Method == "PUT" {
			//TODO this should be a "PUT" method.  Or something PUT /members/1/goals/23/completed
			if len(url) == 4 && url[2] == "goals" {
				return MemberGoalUpdate(r)
			}
		}

		if r.Method == "DELETE" {
			//TODO this should be a "PUT" method.  Or something PUT /members/1/goals/23/completed
			if len(url) == 4 && url[2] == "goals" {
				return MemberGoalUpdate(r)
			}
		}
	}

	// the membersessions APIs
	if apiName == "membersessions" {

		if r.Method == "GET" {
			return MemberSessionGet(r)
		}

		if r.Method == "POST" {
			return MemberSessionCreate(r)
		}

		if r.Method == "DELETE" {
			return MemberSessionTerminate(r)
		}

	}

	return 400, []byte(`"ERROR - UNKNOWN APINAME"`)

	// if r.Method == "POST" {
	// 	switch apiName {
	// 	case "membersession":
	// 		return MemberSession(r)
	// 	case "member":
	// 		return MemberCreate(r)
	// 	default:
	// 		return 400, []byte(`"ERROR - UNKNOWN APINAME"`)
	// 	}
	// } else {
	// 	switch apiName {
	// 	case "members":
	// 		return MemberGet(r)
	// 	case "membergoals":
	// 		return MemberGoalGet(r)
	// 	case "marketingsource":
	// 		return MarketingSourceGet(r)
	// 	default:
	// 		return 400, []byte(`"ERROR - UNKNOWN APINAME"`)
	// 	}
	// }
}
