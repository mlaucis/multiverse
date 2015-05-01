package server

import (
	"encoding/json"
	"net/http"
)

var generalRoutes = map[string]struct {
	path    string
	method  string
	handler http.HandlerFunc
}{
	"index": {
		path:    "/",
		method:  "GET",
		handler: home,
	},
	"humans": {
		path:    "/humans.txt",
		method:  "GET",
		handler: humans,
	},
	"robots": {
		path:    "/robots.txt",
		method:  "GET",
		handler: robots,
	},
	"versions": {
		path:    "/versions",
		method:  "GET",
		handler: versions,
	},
}

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func home(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Refresh", "3; url=https://tapglue.com")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`/* TEAM */
Founders: Normal Wiese, Onur Akpolat
Lead developer: Florin Patan
http://tapglue.com
Location: Berlin, Germany.

/* SITE */
Last update: 2015/03/15
Standards: HTML5
Components: None
Software: Go, AWS Kinesis, PostgreSQL, REDIS, Docker`))
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robots(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}

// versions endpoint handles the api status for each api version
// Request: GET /versions
// Test with: curl -i localhost/versions
func versions(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(24*3600, w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(map[string]string{
		"0.1": "deprecated",
		"0.2": "current",
	})
}
