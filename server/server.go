/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/backend/validator"

	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
)

const (
	userAgentNotSet           = "User-Agent header must be set"
	contentLengthNotSet       = "Content-Length header must be set"
	contentLengthNotDecodable = "Content-Length header value could not be decoded. %q"
	contentLengthSizeNotMatch = "Content-Length header value is different fromt the received value"
	requestBodyCannotBeEmpty  = "request body cannot be empty"
)

var (
	dbgMode bool
)

func getReqAuthToken(r *http.Request) string {
	return r.Header.Get("Authorization")
}

// validateGetCommon runs a series of predefinied, common, tests for GET requests
func validateGetCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(userAgentNotSet, http.StatusBadRequest, r, w)
		return
	}
}

// validatePutCommon runs a series of predefinied, common, tests for PUT requests
func validatePutCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(userAgentNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Length") == "" {
		errorHappened(contentLengthNotSet, http.StatusBadRequest, r, w)
		return
	}

	reqCL, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(fmt.Sprintf(contentLengthNotDecodable, err), http.StatusBadRequest, r, w)
		return
	}

	if reqCL != r.ContentLength {
		errorHappened(contentLengthSizeNotMatch, http.StatusBadRequest, r, w)
		return
	}

	if r.Body == nil {
		errorHappened(requestBodyCannotBeEmpty, http.StatusBadRequest, r, w)
		return
	}
}

// validateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func validateDeleteCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(userAgentNotSet, http.StatusBadRequest, r, w)
		return
	}
}

// validatePostCommon runs a series of predefined, common, tests for the POST requests
func validatePostCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(userAgentNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Length") == "" {
		errorHappened(contentLengthNotSet, http.StatusBadRequest, r, w)
		return
	}

	reqCL, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(fmt.Sprintf(contentLengthNotDecodable, err), http.StatusBadRequest, r, w)
		return
	}

	if reqCL != r.ContentLength {
		errorHappened(contentLengthSizeNotMatch, http.StatusBadRequest, r, w)
		return
	}

	if r.Body == nil {
		errorHappened(requestBodyCannotBeEmpty, http.StatusBadRequest, r, w)
		return
	}
}

// validateAccountRequestToken validates that the request contains a valid request token
func validateAccountRequestToken(w http.ResponseWriter, r *http.Request) {
	var (
		accountID int64
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("invalid accountId number", http.StatusBadRequest, r, w)
		return
	}

	if !validator.ValidateAccountRequestToken(accountID, getReqAuthToken(r)) {
		errorHappened("request is not properly signed", http.StatusBadRequest, r, w)
		return
	}
}

// validateApplicationRequestToken validates that the request contains a valid request token
func validateApplicationRequestToken(w http.ResponseWriter, r *http.Request) {
	var (
		accountID     int64
		applicationID int64
		err           error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened("invalid accountId number", http.StatusBadRequest, r, w)
		return
	}

	if applicationID, err = strconv.ParseInt(vars["applicationId"], 10, 64); err != nil {
		errorHappened("invalid applicationId number", http.StatusBadRequest, r, w)
		return
	}

	if !validator.ValidateApplicationRequestToken(accountID, applicationID, getReqAuthToken(r)) {
		errorHappened("request is not properly signed", http.StatusBadRequest, r, w)
		return
	}
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
func getSanitizedHeaders(headers http.Header) http.Header {
	if !dbgMode {
		headers.Del("Authorization")
	}

	// TODO sanitize headers that shouldn't not appear in the logs

	return headers
}

// writeResponse handles the http responses and returns the data
func writeResponse(response interface{}, code int, cacheTime uint, w http.ResponseWriter, r *http.Request) {
	// Set the response headers
	writeCacheHeaders(cacheTime, w)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	// Write response
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		w.WriteHeader(code)
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Encoding", "gzip")
	w.WriteHeader(code)
	gz := gzip.NewWriter(w)
	json.NewEncoder(gz).Encode(response)
	gz.Close()
}

// errorHappened handles the error message
func errorHappened(message string, code int, r *http.Request, w http.ResponseWriter) {
	writeCacheHeaders(0, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	// Write response
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		w.WriteHeader(code)
		fmt.Fprintf(w, "%d %s", code, message)
	} else {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(code)
		gz := gzip.NewWriter(w)
		fmt.Fprintf(gz, "%d %s", code, message)
		gz.Close()
	}
	_, filename, line, ok := runtime.Caller(1)
	if !ok {
		return
	}

	headers := getSanitizedHeaders(r.Header)

	log.Printf(
		"Error %q in %s/%s:%d while %s\t%s\t%+v\n",
		message,
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
	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(w http.ResponseWriter, r *http.Request) {
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
	writeCacheHeaders(10*24*3600, w)
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.Write([]byte(`User-agent: *
Disallow: /`))
}

func signRequest(token string, req *http.Request) {
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
}

func customHandler(routeName string, r *route, newRelicAgent *gorelic.Agent, logChan chan *LogMsg) http.HandlerFunc {
	var extraHandlers []http.HandlerFunc
	switch r.method {
	case "DELETE":
		{
			extraHandlers = append(extraHandlers, validateDeleteCommon)
		}
	case "GET":
		{
			extraHandlers = append(extraHandlers, validateGetCommon)
		}
	case "PUT":
		{
			extraHandlers = append(extraHandlers, validatePutCommon)
		}
	case "POST":
		{
			extraHandlers = append(extraHandlers, validatePostCommon)
		}
	}

	r.handlers = append(extraHandlers, r.handlers...)

	handlerFunc := func(resp http.ResponseWriter, req *http.Request) {
		start := time.Now()

		for _, handler := range r.handlers {
			// Any response that happens in a handler MUST send a Content-Type header
			if resp.Header().Get("Content-Type") != "" {
				break
			}
			handler(resp, req)
		}

		logChan <- &LogMsg{
			method:     req.Method,
			requestURI: req.RequestURI,
			name:       routeName,
			headers:    req.Header,
			start:      start,
			end:        time.Now(),
		}
	}

	if newRelicAgent != nil {
		return http.HandlerFunc(newRelicAgent.WrapHTTPHandlerFunc(handlerFunc))
	}
	return handlerFunc
}

// GetRouter creates the router
func GetRouter(debugMode bool, newRelicAgent *gorelic.Agent, logChan chan *LogMsg) *mux.Router {
	dbgMode = debugMode
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.method).
				Path(route.routePattern(version)).
				Name(routeName).
				HandlerFunc(customHandler(routeName, route, newRelicAgent, logChan))
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
