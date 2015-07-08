package server

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strings"
)

type generalRoute struct {
	name    string
	path    string
	method  string
	handler http.HandlerFunc
}

var generalRoutes = map[string]generalRoute{
	"humans": {
		name:    "humans",
		path:    "/humans.txt",
		method:  "GET",
		handler: humansHandler,
	},
	"robots": {
		name:    "robots",
		path:    "/robots.txt",
		method:  "GET",
		handler: robotsHandler,
	},
	"versions": {
		name:    "versions",
		path:    "/versions",
		method:  "GET",
		handler: versionsHandler,
	},
	"analytics": {
		name:    "analytics",
		path:    "/analytics",
		method:  "POST",
		handler: analyticsHandler,
	},
	"index": {
		name:    "home",
		path:    "/",
		method:  "GET",
		handler: homeHandler,
	},
}

var (
	notFoundResponse    = []byte("{\"errors\":[{\"code\":0,\"message\":\"requested resource was not found\"}]}")
	analyticsOKResponse = []byte("ok")
)

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func homeHandler(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Header().Set("Refresh", "3; url=https://tapglue.com")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humansHandler(w http.ResponseWriter, r *http.Request) {
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
func robotsHandler(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(10*24*3600, w, r)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}

// versionsHandler endpoint handles the api status for each api version
// Request: GET /versions
// Test with: curl -i localhost/versions
func versionsHandler(w http.ResponseWriter, r *http.Request) {
	response := struct {
		Version map[string]struct {
			Version string `json:"version"`
			Status  string `json:"status"`
		} `json:"versions"`
		Revision string `json:"revision"`
	}{
		Version: map[string]struct {
			Version string `json:"version"`
			Status  string `json:"status"`
		}{
			"0.1": {"0.1", "disabled"},
			"0.2": {"0.2", "current"},
		},
		Revision: currentRevision,
	}

	WriteCommonHeaders(7200, w, r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Write response
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		w.WriteHeader(200)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(200)
	gz := gzip.NewWriter(w)
	json.NewEncoder(gz).Encode(response)
	gz.Close()
}

func analyticsHandler(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(0, w, r)
	w.WriteHeader(200)
	w.Write(analyticsOKResponse)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	WriteCommonHeaders(5, w, r)

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Write response
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		w.WriteHeader(404)
		w.Write(notFoundResponse)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(404)
	gz := gzip.NewWriter(w)
	gz.Write(notFoundResponse)
	gz.Close()
}
