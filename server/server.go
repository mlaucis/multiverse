/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/tapglue/backend/config"

	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
)

var (
	dbgMode bool
)

func getReqAuthToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

// validateGetCommon runs a series of predefinied, common, tests for GET requests
func validateGetCommon(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("User-Agent") == "" {
		return fmt.Errorf("User-Agent header must be set")
	}

	return nil
}

// validateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func validateDeleteCommon(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("User-Agent") == "" {
		return fmt.Errorf("User-Agent header must be set")
	}

	return nil
}

// validatePostCommon runs a series of predefined, common, tests for the POST requests
func validatePostCommon(w http.ResponseWriter, r *http.Request) error {
	if r.Header.Get("User-Agent") == "" {
		return fmt.Errorf("User-Agent header must be set")
	}

	if r.ContentLength < 1 {
		return fmt.Errorf("invalid Content-Length size")
	}

	if r.Header.Get("Content-Length") == "" {
		return fmt.Errorf("Content-Length header must be set")
	}

	reqCL, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		return fmt.Errorf("Content-Length header value could not be decoded. %q", err)
	}

	if reqCL != r.ContentLength {
		fmt.Errorf("Content-Length header value is different fromt the received value")
	}

	return nil
}

// writeCacheHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func writeCacheHeaders(cacheTime uint, w http.ResponseWriter) {
	if cacheTime > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		w.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
}

// getSanitizedHeaders returns the sanitized request headers
func getSanitizedHeaders(r *http.Request) http.Header {
	headers := r.Header

	if !dbgMode {
		headers.Del("Authorization")
	}

	// TODO sanitize headers that shouldn't not appear in the logs

	return headers
}

// writeResponse handles the http responses and returns the data
func writeResponse(response interface{}, code int, cacheTime uint, w http.ResponseWriter, r *http.Request) {
	// Convert response to json
	json, err := json.Marshal(response)
	if err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Set the response headers
	writeCacheHeaders(cacheTime, w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	w.WriteHeader(code)

	// Write response
	w.Write(json)
}

// errorHappened handles the error message
func errorHappened(err error, code int, r *http.Request, w http.ResponseWriter) {
	w.WriteHeader(code)
	writeCacheHeaders(0, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(fmt.Sprintf("%d %q", code, err)))

	if config.Conf().Environment == "test" {
		return
	}

	_, filename, line, ok := runtime.Caller(1)
	if !ok {
		return
	}

	headers := getSanitizedHeaders(r)

	log.Printf(
		"Error %q in %s/%s:%d while %s\t%s\t%+v\n",
		err,
		filepath.Base(filepath.Dir(filename)),
		filepath.Base(filename),
		line,
		r.Method,
		r.RequestURI,
		headers,
	)
}

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func home(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`/* TEAM */
Founder: Normal Wiese, Onur Akpolat
http://tapglue.co
Location: Berlin, Germany.

/* THANKS */
Name: @dlsniper

/* SITE */
Last update: 2014/12/17
Standards: HTML5
Components: None
Software: Go`))
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robots(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}

// GetRouter creates the router
func GetRouter(debugMode bool, newRelicAgent *gorelic.Agent) *mux.Router {
	dbgMode = debugMode
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.method).
				Path(route.routePattern(version)).
				Name(routeName).
				Handler(Logger(route.handlerFunc, routeName, newRelicAgent))
		}
	}

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	return router
}
