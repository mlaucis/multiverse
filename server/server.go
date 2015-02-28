/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/pprof"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/backend/validator"
	"github.com/tapglue/backend/validator/keys"

	"github.com/gorilla/mux"
)

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
	dbgMode      bool
	mainLogChan  = make(chan *LogMsg, 100000)
	errorLogChan = make(chan *LogMsg, 100000)
)

// isRequestExpired checks if the request is expired or not
func isRequestExpired(ctx *context) {
	// Check that the request is not older than 3 days
	// TODO check if we should lower the interval
	requestDate := ctx.r.Header.Get("x-tapglue-date")
	if requestDate == "" {
		errorHappened(ctx, "request date is invalid", http.StatusBadRequest)
		return
	}

	parsedRequestDate, err := time.Parse(time.RFC3339, requestDate)
	if err != nil {
		errorHappened(ctx, "request date is invalid", http.StatusBadRequest)
		return
	}

	if time.Since(parsedRequestDate) > time.Duration(3*24*time.Hour) {
		errorHappened(ctx, "request is expired", http.StatusExpectationFailed)
	}
}

// validateGetCommon runs a series of predefined, common, tests for GET requests
func validateGetCommon(ctx *context) {
	if ctx.r.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest)
		return
	}
}

// validatePutCommon runs a series of predefinied, common, tests for PUT requests
func validatePutCommon(ctx *context) {
	if ctx.r.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest)
		return
	}

	if ctx.r.Header.Get("Content-Length") == "" {
		errorHappened(ctx, errContentLengthNotSet, http.StatusBadRequest)
		return
	}

	if ctx.r.Header.Get("Content-Type") == "" {
		errorHappened(ctx, errContentTypeNotSet, http.StatusBadRequest)
		return
	}

	if ctx.r.Header.Get("Content-Type") != "application/json" {
		errorHappened(ctx, errWrongContentType, http.StatusBadRequest)
		return
	}

	reqCL, err := strconv.ParseInt(ctx.r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(ctx, fmt.Sprintf(errContentLengthNotDecodable, err), http.StatusBadRequest)
		return
	}

	if reqCL != ctx.r.ContentLength {
		errorHappened(ctx, errContentLengthSizeNotMatch, http.StatusBadRequest)
		return
	}

	if ctx.r.Body == nil {
		errorHappened(ctx, errRequestBodyCannotBeEmpty, http.StatusBadRequest)
		return
	}
}

// validateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func validateDeleteCommon(ctx *context) {
	if ctx.r.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest)
		return
	}
}

// validatePostCommon runs a series of predefined, common, tests for the POST requests
func validatePostCommon(ctx *context) {
	if ctx.r.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest)
		return
	}

	if ctx.r.Header.Get("Content-Length") == "" {
		errorHappened(ctx, errContentLengthNotSet, http.StatusBadRequest)
		return
	}

	if ctx.r.Header.Get("Content-Type") == "" {
		errorHappened(ctx, errContentTypeNotSet, http.StatusBadRequest)
		return
	}

	if ctx.r.Header.Get("Content-Type") != "application/json" {
		errorHappened(ctx, errWrongContentType, http.StatusBadRequest)
		return
	}

	reqCL, err := strconv.ParseInt(ctx.r.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(ctx, fmt.Sprintf(errContentLengthNotDecodable, err), http.StatusLengthRequired)
		return
	}

	if reqCL != ctx.r.ContentLength {
		errorHappened(ctx, errContentLengthSizeNotMatch, http.StatusBadRequest)
		return
	}

	if ctx.r.Body == nil {
		errorHappened(ctx, errRequestBodyCannotBeEmpty, http.StatusBadRequest)
		return
	}
}

// validateApplicationRequestToken validates that the request contains a valid request token
func validateApplicationRequestToken(ctx *context) {
	if keys.VerifyRequest(ctx.scope, ctx.version, ctx.r) {
		return
	}

	errorHappened(ctx, "request is not properly signed", http.StatusUnauthorized)
}

// isSessionValid checks if the session token is valid or not
func checkSession(ctx *context) {
	sessionToken, err := validator.CheckSession(ctx.r)
	if err == nil {
		ctx.sessionToken = sessionToken
		return
	}

	errorHappened(ctx, "invalid session", http.StatusUnauthorized)
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
func writeResponse(ctx *context, response interface{}, code int, cacheTime uint) {
	// Set the response headers
	writeCacheHeaders(cacheTime, ctx.w)
	ctx.w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	ctx.w.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding")
	ctx.w.Header().Set("Access-Control-Allow-Credentials", "true")

	//Check if we have a session enable and if so write it back
	if ctx.sessionToken != "" {
		ctx.w.Header().Set("x-tapglue-session", ctx.sessionToken)
	}

	// Write response
	if !strings.Contains(ctx.r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.w.WriteHeader(code)
		json.NewEncoder(ctx.w).Encode(response)
		return
	}

	ctx.w.Header().Set("Content-Encoding", "gzip")
	ctx.w.WriteHeader(code)
	gz := gzip.NewWriter(ctx.w)
	json.NewEncoder(gz).Encode(response)
	gz.Close()
}

// errorHappened handles the error message
func errorHappened(ctx *context, message string, code int) {
	writeCacheHeaders(0, ctx.w)
	ctx.w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	// Write response
	if !strings.Contains(ctx.r.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.w.WriteHeader(code)
		fmt.Fprintf(ctx.w, "%d %s", code, message)
	} else {
		ctx.w.Header().Set("Content-Encoding", "gzip")
		ctx.w.WriteHeader(code)
		gz := gzip.NewWriter(ctx.w)
		fmt.Fprintf(gz, "%d %s", code, message)
		gz.Close()
	}
	_, filename, line, ok := runtime.Caller(1)
	if !ok {
		return
	}

	ctx.errorLog <- &LogMsg{
		remoteAddr: ctx.r.RemoteAddr,
		method:     ctx.r.Method,
		requestURI: ctx.r.RequestURI,
		headers:    ctx.r.Header,
		name:       ctx.routeName,
		start:      ctx.startTime,
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
func home(ctx *context) {
	writeCacheHeaders(10*24*3600, ctx.w)
	ctx.w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.w.Write([]byte(`these aren't the droids you're looking for`))
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(ctx *context) {
	writeCacheHeaders(10*24*3600, ctx.w)
	ctx.w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.w.Write([]byte(`/* TEAM */
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
func robots(ctx *context) {
	writeCacheHeaders(10*24*3600, ctx.w)
	ctx.w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.w.Write([]byte(`User-agent: *
Disallow: /`))
}

func customHandler(routeName, version string, route *route, mainLog, errorLog chan *LogMsg) http.HandlerFunc {
	var extraHandlers []routeFunc
	switch route.method {
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

	route.handlers = append(extraHandlers, route.handlers...)

	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := NewContext(w, r, mainLog, errorLog, routeName, route.scope, version)
		if err != nil {
			errorHappened(ctx, "failed to get a request context", http.StatusInternalServerError)
			return
		}

		for _, handler := range route.handlers {
			// Any response that happens in a handler MUST send a Content-Type header
			if w.Header().Get("Content-Type") != "" {
				break
			}
			handler(ctx)
		}

		ctx.mainLog <- &LogMsg{
			remoteAddr: ctx.r.RemoteAddr,
			method:     ctx.r.Method,
			requestURI: ctx.r.RequestURI,
			headers:    ctx.r.Header,
			name:       ctx.routeName,
			start:      ctx.startTime,
			end:        time.Now(),
		}
	}
}

// GetRouter creates the router
func GetRouter(debugMode bool) (*mux.Router, chan *LogMsg, chan *LogMsg, error) {
	dbgMode = debugMode
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.method).
				Path(route.routePattern(version)).
				Name(routeName).
				HandlerFunc(customHandler(routeName, version, route, mainLogChan, errorLogChan))
		}
	}

	for _, routeName := range []string{"index", "humans", "robots"} {
		version := ""
		route := routes["0.1"][routeName]
		if route == nil {
			panic(fmt.Sprintf("route %s not found", routeName))
		}
		router.
			Methods(route.method).
			Path(route.pattern).
			Name(routeName).
			HandlerFunc(customHandler(routeName, version, route, mainLogChan, errorLogChan))
	}

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	return router, mainLogChan, errorLogChan, nil
}
