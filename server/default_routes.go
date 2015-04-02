package server

import "net/http"

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
}

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func home(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Refresh", "3; url=https://tapglue.com")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`/* TEAM */
Founder: Normal Wiese, Onur Akpolat
Lead developer: Florin Patan
http://tapglue.com
Location: Berlin, Germany.

/* SITE */
Last update: 2015/03/15
Standards: HTML5
Components: None
Software: Go, Redis`))
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robots(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}
