package main

import (
	"github.com/user/apiserver/api"
	"github.com/user/apiserver/app"
	"net/http"
	"time"
)

//--------------------------------------------------------------------
// this handler receives the HTTP request and processes the request
//--------------------------------------------------------------------
func handler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	api.RunAPI(w, r, start)

}

func main() {

	app.Init()
	app.Info("API server started")

	http.HandleFunc("/", handler)
	http.ListenAndServe(":"+app.Cfg["http.port"], nil)

}
