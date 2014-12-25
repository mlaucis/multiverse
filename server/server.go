/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

/**
 * writeResponse handles the http responses and returns the data
 * @param response, response data
 * @param code, http status code
 * @param cache, response cache
 * @param w, http response writer
 * @param r, http request
 */
func writeResponse(response interface{}, code int, cache uint, w http.ResponseWriter, r *http.Request) {
	// Read response to json
	json, err := json.Marshal(response)
	if err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusInternalServerError, w)
	}

	// Set the response headers
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cache))
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
	w.Header().Set("Content-Type", "text/plain")
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
	w.Header().Set("Content-Type", "text/plain")
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
	w.Header().Set("Content-Type", "text/plain")
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
	w.Header().Set("Content-Type", "text/plain")
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