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
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/validator"
	"github.com/tapglue/backend/validator/keys"
	"github.com/tapglue/backend/validator/tokens"

	"github.com/gorilla/mux"
)

const (
	apiRequestVersionString = "tg%s"

	errUserAgentNotSet           = "User-Agent header must be set (1)"
	errContentLengthNotSet       = "Content-Length header must be set (1)"
	errContentTypeNotSet         = "Content-Type header must be set (1)"
	errContentLengthNotDecodable = "Content-Length header value could not be decoded (2)"
	errContentLengthSizeNotMatch = "Content-Length header value is different from the received payload size (3)"
	errRequestBodyCannotBeEmpty  = "Request body cannot be empty (1)"
	errWrongContentType          = "Wrong Content-Type header value (1)"
)

var (
	mainLogChan  = make(chan *logger.LogMsg, 100000)
	errorLogChan = make(chan *logger.LogMsg, 100000)
	skipSecurity = false
)

// isRequestExpired checks if the request is expired or not
func isRequestExpired(ctx *context.Context) {
	if skipSecurity {
		return
	}

	// Check that the request is not older than 3 days
	// TODO check if we should lower the interval
	requestDate := ctx.R.Header.Get("x-tapglue-date")
	if requestDate == "" {
		errorHappened(ctx, "invalid request date (1)\nrequest date is empty", http.StatusBadRequest, fmt.Errorf("invalid request date"))
		return
	}

	parsedRequestDate, err := time.Parse(time.RFC3339, requestDate)
	if err != nil {
		errorHappened(ctx, "invalid request date (2)\ndate is not in RFC3339 format\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if time.Since(parsedRequestDate) > time.Duration(3*time.Minute) {
		errorHappened(ctx, "request is expired (1)", http.StatusExpectationFailed, err)
	}
}

// validateGetCommon runs a series of predefined, common, tests for GET requests
func validateGetCommon(ctx *context.Context) {
	if ctx.R.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest, nil)
		return
	}
}

// validatePutCommon runs a series of predefinied, common, tests for PUT requests
func validatePutCommon(ctx *context.Context) {
	if skipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		errorHappened(ctx, errContentLengthNotSet, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		errorHappened(ctx, errContentTypeNotSet, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		errorHappened(ctx, errWrongContentType, http.StatusBadRequest, nil)
		return
	}

	reqCL, err := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(ctx, errContentLengthNotDecodable+"\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if reqCL != ctx.R.ContentLength {
		errorHappened(ctx, errContentLengthSizeNotMatch, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Body == nil {
		errorHappened(ctx, errRequestBodyCannotBeEmpty, http.StatusBadRequest, nil)
		return
	}
}

// validateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func validateDeleteCommon(ctx *context.Context) {
	if ctx.R.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest, nil)
		return
	}
}

// validatePostCommon runs a series of predefined, common, tests for the POST requests
func validatePostCommon(ctx *context.Context) {
	if skipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		errorHappened(ctx, errUserAgentNotSet, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		errorHappened(ctx, errContentLengthNotSet, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		errorHappened(ctx, errContentTypeNotSet, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		errorHappened(ctx, errWrongContentType, http.StatusBadRequest, nil)
		return
	}

	reqCL, err := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if err != nil {
		errorHappened(ctx, errContentLengthNotDecodable+"\n"+err.Error(), http.StatusLengthRequired, err)
		return
	}

	if reqCL != ctx.R.ContentLength {
		errorHappened(ctx, errContentLengthSizeNotMatch, http.StatusBadRequest, nil)
		return
	}

	if ctx.R.Body == nil {
		errorHappened(ctx, errRequestBodyCannotBeEmpty, http.StatusBadRequest, nil)
		return
	}
}

// validateApplicationRequestToken validates that the request contains a valid request token
func validateAccountRequestToken(ctx *context.Context) {
	if skipSecurity {
		return
	}

	if errMsg, err := keys.VerifyRequest(ctx, 1); err != nil {
		errorHappened(ctx, errMsg, http.StatusUnauthorized, err)
	}
}

// validateApplicationRequestToken validates that the request contains a valid request token
func validateApplicationRequestToken(ctx *context.Context) {
	if skipSecurity {
		return
	}

	var (
		errMsg string
		err    error
	)
	if ctx.Version == "0.1" {
		errMsg, err = tokens.VerifyRequest(ctx, 3)
	} else {
		errMsg, err = keys.VerifyRequest(ctx, 2)
	}

	if err != nil {
		errorHappened(ctx, errMsg, http.StatusUnauthorized, err)
	}
}

// checkAccountSession checks if the session token is valid or not
func checkAccountSession(ctx *context.Context) {
	if skipSecurity {
		return
	}

	sessionToken, errMsg, err := validator.CheckAccountSession(ctx.R)
	if err == nil {
		ctx.SessionToken = sessionToken
		return
	}

	errorHappened(ctx, errMsg, http.StatusUnauthorized, err)
}

// checkApplicationSession checks if the session token is valid or not
func checkApplicationSession(ctx *context.Context) {
	if skipSecurity {
		return
	}

	var (
		errMsg, sessionToken string
		err                  error
	)

	if ctx.Version == "0.1" {
		sessionToken, errMsg, err = validator.CheckApplicationSimpleSession(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, ctx.R)
	} else {
		sessionToken, errMsg, err = validator.CheckApplicationSession(ctx.R)
	}

	if err == nil {
		ctx.SessionToken = sessionToken
		return
	}

	errorHappened(ctx, errMsg, http.StatusUnauthorized, err)
}

// writeCacheHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func writeCacheHeaders(cacheTime uint, ctx *context.Context) {
	if cacheTime > 0 {
		ctx.W.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		ctx.W.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		ctx.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.W.Header().Set("Pragma", "no-cache")
		ctx.W.Header().Set("Expires", "0")
	}
}

func writeCorsHeaders(ctx *context.Context) {
	ctx.W.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.W.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding, x-tapglue-id, x-tapglue-date, x-tapglue-session, x-tapglue-payload-hash, x-tapglue-signature")
	ctx.W.Header().Set("Access-Control-Allow-Credentials", "true")
}

// writeResponse handles the http responses and returns the data
func writeResponse(ctx *context.Context, response interface{}, code int, cacheTime uint) {
	// Set the response headers
	writeCacheHeaders(cacheTime, ctx)
	writeCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")

	//Check if we have a session enable and if so write it back
	if ctx.SessionToken != "" {
		ctx.W.Header().Set("x-tapglue-session", ctx.SessionToken)
	}

	ctx.StatusCode = code

	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(code)
		json.NewEncoder(ctx.W).Encode(response)
		return
	}

	ctx.W.Header().Set("Content-Encoding", "gzip")
	ctx.W.WriteHeader(code)
	gz := gzip.NewWriter(ctx.W)
	json.NewEncoder(gz).Encode(response)
	gz.Close()
}

// errorHappened handles the error message
func errorHappened(ctx *context.Context, message string, code int, internalError error) {
	writeCacheHeaders(0, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(code)
		fmt.Fprintf(ctx.W, "%d %s", code, message)
	} else {
		ctx.W.Header().Set("Content-Encoding", "gzip")
		ctx.W.WriteHeader(code)
		gz := gzip.NewWriter(ctx.W)
		fmt.Fprintf(gz, "%d %s", code, message)
		gz.Close()
	}

	ctx.StatusCode = code
	ctx.LogErrorWithMessage(internalError, message, 1)
}

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func home(ctx *context.Context) {
	writeCacheHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`these aren't the droids you're looking for`))
	ctx.StatusCode = 200
}

// humans handles requests to humans.txt
// Request: GET /humans.txt
// Test with: curl -i localhost/humans.txt
func humans(ctx *context.Context) {
	writeCacheHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`/* TEAM */
Founder: Normal Wiese, Onur Akpolat
Lead developer: Florin Patan
http://tapglue.com
Location: Berlin, Germany.

/* SITE */
Last update: 2015/03/15
Standards: HTML5
Components: None
Software: Go, Redis`))
	ctx.StatusCode = 200
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robots(ctx *context.Context) {
	writeCacheHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`User-agent: *
Disallow: /`))
	ctx.StatusCode = 200
}

func corsHandler(ctx *context.Context) {
	if ctx.R.Method != "OPTIONS" {
		return
	}

	writeCacheHeaders(100, ctx)
	writeCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

// CustomHandler composes the http handling function according to its definition and returns it
func CustomHandler(routeName, version string, route *Route, mainLog, errorLog chan *logger.LogMsg, environment string, debugMode bool) http.HandlerFunc {
	extraHandlers := []RouteFunc{corsHandler}
	switch route.Method {
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

	if version != "0.1" && routeName != "index" && routeName != "humans" && routeName != "robots" {
		extraHandlers = append(extraHandlers, isRequestExpired)
	}

	route.Handlers = append(extraHandlers, route.Handlers...)

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, err := context.NewContext(w, r, mainLog, errorLog, routeName, route.Scope, version, route.Filters, environment, debugMode)
		if err != nil {
			errorHappened(ctx, "failed to get a request context (1)", http.StatusInternalServerError, err)
			return
		}

		for _, handler := range route.Handlers {
			// Any response that happens in a handler MUST send a Content-Type header
			if w.Header().Get("Content-Type") != "" {
				break
			}
			handler(ctx)
		}

		ctx.LogRequest(ctx.StatusCode, -1)
	}
}

// GetRouter creates the router
func GetRouter(environment string, debugMode, skipSecurityChecks bool) (*mux.Router, chan *logger.LogMsg, chan *logger.LogMsg, error) {
	skipSecurity = skipSecurityChecks
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.Method, "OPTIONS").
				Path(route.RoutePattern(version)).
				Name(routeName).
				HandlerFunc(CustomHandler(routeName, version, route, mainLogChan, errorLogChan, environment, debugMode))
		}
	}

	for _, routeName := range []string{"index", "humans", "robots"} {
		version := ""
		route := routes["0.1"][routeName]
		if route == nil {
			panic(fmt.Sprintf("route %s not found", routeName))
		}
		router.
			Methods(route.Method, "OPTIONS").
			Path(route.Pattern).
			Name(routeName).
			HandlerFunc(CustomHandler(routeName, version, route, mainLogChan, errorLogChan, environment, debugMode))
	}

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	return router, mainLogChan, errorLogChan, nil
}
