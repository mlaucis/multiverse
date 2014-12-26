/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// To be used across API demos, should be deleted when no raw examples present
var api_demo_time = time.Date(2014, time.December, 25, 12, 30, 0, 0, time.UTC)

func validatePostCommon(w http.ResponseWriter, r *http.Request) error {
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

/**
 * writeCacheHeaders will add the corresponding cache headers based on the time supplied (in seconds)
 * @param cacheTime, response cache
 * @param w, http response writer
 */
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

/**
 * writeResponse handles the http responses and returns the data
 * @param response, response data
 * @param code, http status code
 * @param cacheTime, response cache
 * @param w, http response writer
 * @param r, http request
 */
func writeResponse(response interface{}, code int, cacheTime uint, w http.ResponseWriter, r *http.Request) {
	// Read response to json
	json, err := json.Marshal(response)
	if err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusInternalServerError, w)
		return
	}

	// Set the response headers
	writeCacheHeaders(cacheTime, w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Write response
	w.Write(json)
}

/**
 * errorHappened handles the error message
 * @param message, error message
 * @param code, http status code
 * @param w, response writer
 */
func errorHappened(message string, code int, w http.ResponseWriter) {
	w.WriteHeader(code)
	writeCacheHeaders(0, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(fmt.Sprintf("%d %s", code, message)))
}

/**
 * home handles request to API root
 * Request: GET /
 * Test with: `curl -i localhost/`
 * @param w, response writer
 * @param r, http request
 */
func home(w http.ResponseWriter, r *http.Request) {
	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

/**
 * humans handles requests to humans.txt
 * Request: GET /humans.txt
 * Test with: curl -i localhost/humans.txt
 * @param w, response writer
 * @param r, http request
 */
func humans(w http.ResponseWriter, r *http.Request) {
	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`/* TEAM */
Founder: Normal Wiese, Onur Akpolat
http://gluee.co
Location: Berlin, Germany.

/* THANKS */
Name: @dlsniper

/* SITE */
Last update: 2014/12/17
Standards: HTML5
Components: None
Software: Go`))
}

/**
 * robots handles requests to robots.txt
 * Request: GET /robots.txt
 * Test with: curl -i localhost/robots.txt
 * @param w, response writer
 * @param r, http request
 */
func robots(w http.ResponseWriter, r *http.Request) {
	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}

/**
 * GetRouter creates the router
 * @return router, mux router with all routes defined
 */
func GetRouter() *mux.Router {

	// Create router
	router := mux.NewRouter().StrictSlash(true)

	// Read routes
	for _, route := range routes {
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(route.name).
			Handler(Logger(route.handlerFunc, route.name))
	}

	// Return router
	return router
}
