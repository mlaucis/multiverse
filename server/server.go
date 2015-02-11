/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/validator"
	"github.com/tapglue/backend/validator/keys"
)

// We have to have our own type so that we can break what go forces us to do
type noCloseReaderCloser struct {
	*bytes.Buffer
}

const (
	apiRequestVersionString = "tg%s"

	errUserAgentNotSet           = "User-Agent header must be set"
	errContentLengthNotSet       = "Content-Length header must be set"
	errContentTypeNotSet         = "Content-Type header must be set"
	errContentLengthNotDecodable = "Content-Length header value could not be decoded. %q"
	errContentLengthSizeNotMatch = "Content-Length header value is different fromt the received value"
	errRequestBodyCannotBeEmpty  = "Request body cannot be empty"
	errWrongContentType          = "Wrong Content-Type header value"
)

var (
	dbgMode bool
	logChan = make(chan *LogMsg, 100000)
)

// We should do some closing here but then again, that's what we want to prevent
func (m noCloseReaderCloser) Close() error {
	return nil
}

// peakBody allows us to look at the request body and get the values without closing the body
func peakBody(r *http.Request) *bytes.Buffer {
	buf, _ := ioutil.ReadAll(r.Body)
	buff := noCloseReaderCloser{bytes.NewBuffer(buf)}
	r.Body = noCloseReaderCloser{bytes.NewBuffer(buf)}
	return buff.Buffer
}

// isRequestExpired checks if the request is expired or not
func isRequestExpired(w http.ResponseWriter, r *http.Request) {
	// Check that the request is not older than 3 days
	// TODO check if we should lower the interval
	requestDate := r.Header.Get("x-tapglue-date")
	if requestDate == "" {
		errorHappened("request date is invalid", http.StatusBadRequest, r, w)
		return
	}

	parsedRequestDate, err := time.Parse(time.RFC3339, requestDate)
	if err != nil {
		errorHappened("request date is invalid", http.StatusBadRequest, r, w)
		return
	}

	if time.Since(parsedRequestDate) > time.Duration(3*24*time.Hour) {
		errorHappened("request is expired", http.StatusExpectationFailed, r, w)
	}
}

// validateGetCommon runs a series of predefined, common, tests for GET requests
func validateGetCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(errUserAgentNotSet, http.StatusBadRequest, r, w)
		return
	}
}

// validatePutCommon runs a series of predefinied, common, tests for PUT requests
func validatePutCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(errUserAgentNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Length") == "" {
		errorHappened(errContentLengthNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Type") == "" {
		errorHappened(errContentTypeNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		errorHappened(errWrongContentType, http.StatusBadRequest, r, w)
		return
	}

	reqCL, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(fmt.Sprintf(errContentLengthNotDecodable, err), http.StatusBadRequest, r, w)
		return
	}

	if reqCL != r.ContentLength {
		errorHappened(errContentLengthSizeNotMatch, http.StatusBadRequest, r, w)
		return
	}

	if r.Body == nil {
		errorHappened(errRequestBodyCannotBeEmpty, http.StatusBadRequest, r, w)
		return
	}
}

// validateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func validateDeleteCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(errUserAgentNotSet, http.StatusBadRequest, r, w)
		return
	}
}

// validatePostCommon runs a series of predefined, common, tests for the POST requests
func validatePostCommon(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("User-Agent") == "" {
		errorHappened(errUserAgentNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Length") == "" {
		errorHappened(errContentLengthNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Type") == "" {
		errorHappened(errContentTypeNotSet, http.StatusBadRequest, r, w)
		return
	}

	if r.Header.Get("Content-Type") != "application/json" {
		errorHappened(errWrongContentType, http.StatusBadRequest, r, w)
		return
	}

	reqCL, err := strconv.ParseInt(r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(fmt.Sprintf(errContentLengthNotDecodable, err), http.StatusLengthRequired, r, w)
		return
	}

	if reqCL != r.ContentLength {
		errorHappened(errContentLengthSizeNotMatch, http.StatusBadRequest, r, w)
		return
	}

	if r.Body == nil {
		errorHappened(errRequestBodyCannotBeEmpty, http.StatusBadRequest, r, w)
		return
	}
}

// validateApplicationRequestToken validates that the request contains a valid request token
func validateApplicationRequestToken(requestScope, requestVersion string, w http.ResponseWriter, r *http.Request) {
	if keys.VerifyRequest(requestScope, requestVersion, r) {
		return
	}

	errorHappened("request is not properly signed", http.StatusUnauthorized, r, w)
}

// isSessionValid checks if the session token is valid or not
func isSessionValid(w http.ResponseWriter, r *http.Request) {
	if err := validator.IsSessionValid(r); err == nil {
		return
	}

	errorHappened("invalid session", http.StatusUnauthorized, r, w)
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

	logChan <- &LogMsg{
		method:     r.Method,
		requestURI: r.RequestURI,
		headers:    r.Header,
		name:       "-",
		start:      time.Now(),
		end:        time.Now(),
		message: fmt.Sprintf(
			"Error %q in %s/%s:%d",
			message,
			filepath.Base(filepath.Dir(filename)),
			filepath.Base(filename),
			line,
		),
	}
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

func customHandler(routeName string, r *route) http.HandlerFunc {
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

	return func(resp http.ResponseWriter, req *http.Request) {
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
}

// GetRouter creates the router
func GetRouter(debugMode bool) (*mux.Router, chan *LogMsg, error) {
	dbgMode = debugMode
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.method).
				Path(route.routePattern(version)).
				Name(routeName).
				HandlerFunc(customHandler(routeName, route))
		}
	}

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	return router, logChan, nil
}
