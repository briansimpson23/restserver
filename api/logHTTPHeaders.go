package api

import (
	"../app"
	"fmt"
	"net/http"
	"net/http/httputil"
)

func logHTTPHeaders(r *http.Request) {

	// &{Method:GET URL:/members/1 Proto:HTTP/1.0 ProtoMajor:1 ProtoMinor:0 Header:map[User-Agent:[Wget/1.12 (linux-gnu)] Accept:[*/*] Connection:[Keep-Alive]] Body:0xc2000b09c0 ContentLength:0 TransferEncoding:[] Close:false Host:host1.local.venovix.com:8080 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:127.0.0.1:55724 RequestURI:/members/1 TLS:<nil>}

	// &{Method:POST URL:/membersession Proto:HTTP/1.1 ProtoMajor:1 ProtoMinor:1 Header:map[User-Agent:[serviceSDK/0.0.1] Accept:[*/*] Api-Key:[LKJLK123KLJ123JLKJL12] Api-Session:[43gg84] Content-Length:[188] Content-Type:[multipart/form-data; boundary=----------------------------6e7bc3ffb675]] Body:0xc2001041c0 ContentLength:188 TransferEncoding:[] Close:false Host:localhost:8080 Form:map[] PostForm:map[] MultipartForm:<nil> Trailer:map[] RemoteAddr:[::1]:43235 RequestURI:/membersession TLS:<nil>}

	dump, _ := httputil.DumpRequest(r, true)
	app.Debug(fmt.Sprintf("ZZZZZ %s", dump))

	app.Debug(fmt.Sprintf("HTTP HEADER => %+v\n\n", r))
	app.Debug(fmt.Sprintf("HttpRequest: Method[%s]", r.Method))
	app.Debug(fmt.Sprintf("HttpRequest: URL[%s]", r.URL))
	app.Debug(fmt.Sprintf("HttpRequest: Protocol[%s]", r.Proto))
	app.Debug(fmt.Sprintf("HttpRequest: RemoteAddr[%s]", r.RemoteAddr))
	app.Debug(fmt.Sprintf("HttpRequest: RequestURI[%s]", r.RequestURI))
	app.Debug("HttpRequest: Headers")
	for name, value := range r.Header {
		app.Debug(fmt.Sprintf("HttpRequest:    %s: %s", name, value))
	}
	//app.Debug(fmt.Sprintf("HttpRequest: Body[%s]", r.Body.Reader()))
	app.Debug(fmt.Sprintf("HttpRequest: ContentLength[%d]", r.ContentLength))
	app.Debug("HttpRequest:   Form")
	for name, value := range r.Form {
		app.Debug(fmt.Sprintf("HttpRequest:    %s: %s", name, value))
	}
	app.Debug("HttpRequest:   PostForm")
	for name, value := range r.PostForm {
		app.Debug(fmt.Sprintf("HttpRequest:    %s: %s", name, value))
	}
	app.Debug("HttpRequest:   Trailer")
	for name, value := range r.Trailer {
		app.Debug(fmt.Sprintf("HttpRequest:    %s: %s", name, value))
	}
	app.Debug(fmt.Sprintf("HttpRequest:POSTFORM %s", r.PostForm))
	app.Debug(fmt.Sprintf("HttpRequest:FORM %s", r.Form))
	app.Debug(fmt.Sprintf("HttpRequest:payload %s", r.FormValue("payload")))

}
