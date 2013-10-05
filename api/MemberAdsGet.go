package api

import (
	"encoding/json"
	"fmt"
	"github.com/user/apiserver/app"
	"net/http"
	"strings"
)

func MemberAdsGet(r *http.Request) (int, []byte) {

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
	app.Debug(fmt.Sprintf("memberID[%d]", memberID))

	//-------------------------------------------------------------------
	// run the query to get the records from the database
	//-------------------------------------------------------------------
	sql := "SELECT adID, title, shortDesc, linkToURL, linkToName, imagePath FROM ADVERTISEMENT ORDER BY RAND()"
	rows, err := app.Db.Query(sql)
	if err != nil {
		app.Error(err.Error())
		return 0, nil
	}

	// AdID       int
	// Title      string
	// ShortDesc  string
	// LinkToURL  string
	// LinkToName string
	// ImagePath  string

	ads := []Advertisement{}
	for rows.Next() {
		ad := Advertisement{}
		rows.Scan(&ad.AdID, &ad.Title, &ad.ShortDesc, &ad.LinkToURL, &ad.LinkToName, &ad.ImagePath)

		ads = append(ads, ad)
	}

	//----------------------------------------------------------
	// build the output
	//----------------------------------------------------------
	resp := Response{}
	resp.Code = 200
	resp.Status = "ok"
	resp.Data = ads

	output, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}

	return 0, output

}
