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

	errUserAgentNotSet           = "User-Agent header must be set"
	errContentLengthNotSet       = "Content-Length header must be set"
	errContentTypeNotSet         = "Content-Type header must be set"
	errContentLengthNotDecodable = "Content-Length header value could not be decoded. %q"
	errContentLengthSizeNotMatch = "Content-Length header value is different fromt the received value"
	errRequestBodyCannotBeEmpty  = "Request body cannot be empty"
	errWrongContentType          = "Wrong Content-Type header value"
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
		errorHappened(ctx, "request date is invalid", http.StatusBadRequest, nil)
		return
	}

	parsedRequestDate, err := time.Parse(time.RFC3339, requestDate)
	if err != nil {
		errorHappened(ctx, "request date is invalid", http.StatusBadRequest, err)
		return
	}

	if time.Since(parsedRequestDate) > time.Duration(3*24*time.Hour) {
		errorHappened(ctx, "request is expired", http.StatusExpectationFailed, err)
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
		errorHappened(ctx, fmt.Sprintf(errContentLengthNotDecodable, err), http.StatusBadRequest, err)
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
		errorHappened(ctx, fmt.Sprintf(errContentLengthNotDecodable, err), http.StatusLengthRequired, err)
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

	if err := keys.VerifyRequest(ctx, 1); err != nil {
		errorHappened(ctx, "request is not properly signed", http.StatusUnauthorized, err)
	}
}

// validateApplicationRequestToken validates that the request contains a valid request token
func validateApplicationRequestToken(ctx *context.Context) {
	if skipSecurity {
		return
	}

	var err error
	if ctx.Version == "0.1" {
		err = tokens.VerifyRequest(ctx, 3)
	} else {
		err = keys.VerifyRequest(ctx, 2)
	}

	if err != nil {
		errorHappened(ctx, "request is not properly signed", http.StatusUnauthorized, err)
	}
}

// checkAccountSession checks if the session token is valid or not
func checkAccountSession(ctx *context.Context) {
	if skipSecurity {
		return
	}

	sessionToken, err := validator.CheckAccountSession(ctx.R)
	if err == nil {
		ctx.SessionToken = sessionToken
		return
	}

	errorHappened(ctx, "invalid session", http.StatusUnauthorized, err)
}

// checkApplicationSession checks if the session token is valid or not
func checkApplicationSession(ctx *context.Context) {
	if skipSecurity {
		return
	}

	var (
		sessionToken string
		err          error
	)

	if ctx.Version == "0.1" {
		sessionToken, err = validator.CheckApplicationSimpleSession(ctx.AccountID, ctx.ApplicationID, ctx.ApplicationUserID, ctx.R)
	} else {
		sessionToken, err = validator.CheckApplicationSession(ctx.R)
	}

	if err == nil {
		ctx.SessionToken = sessionToken
		return
	}

	errorHappened(ctx, "invalid session", http.StatusUnauthorized, err)
}

// writeCacheHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func writeCacheHeaders(cacheTime uint, ctx *context.Context) {
	/*if cacheTime > 0 {
		ctx.W.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		ctx.W.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {*/
	ctx.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
	ctx.W.Header().Set("Pragma", "no-cache")
	ctx.W.Header().Set("Expires", "0")
	//}
}

// getSanitizedHeaders returns the sanitized request headers
func getSanitizedHeaders(headers http.Header) http.Header {
	// TODO sanitize headers that shouldn't not appear in the logs

	return headers
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

	ctx.LogErrorWithMessage(internalError, message, 1)
}

// home handles request to API root
// Request: GET /
// Test with: `curl -i localhost/`
func home(ctx *context.Context) {
	writeCacheHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`these aren't the droids you're looking for`))
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
}

// robots handles requests to robots.txt
// Request: GET /robots.txt
// Test with: curl -i localhost/robots.txt
func robots(ctx *context.Context) {
	writeCacheHeaders(10*24*3600, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	ctx.W.Write([]byte(`User-agent: *
Disallow: /`))
}

func corsHandler(ctx *context.Context) {
	if ctx.R.Method != "OPTIONS" {
		return
	}

	writeCacheHeaders(100, ctx)
	writeCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
}

func customHandler(routeName, version string, route *route, mainLog, errorLog chan *logger.LogMsg, environment string, debugMode bool) http.HandlerFunc {
	var extraHandlers []routeFunc = []routeFunc{corsHandler}
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

	if version != "0.1" && routeName != "index" && routeName != "humans" && routeName != "robots" {
		extraHandlers = append(extraHandlers, isRequestExpired)
	}

	route.handlers = append(extraHandlers, route.handlers...)

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, err := context.NewContext(w, r, mainLog, errorLog, routeName, route.scope, version, route.contextFilters, environment, debugMode)
		if err != nil {
			errorHappened(ctx, "failed to get a request context", http.StatusInternalServerError, err)
			return
		}

		for _, handler := range route.handlers {
			// Any response that happens in a handler MUST send a Content-Type header
			if w.Header().Get("Content-Type") != "" {
				break
			}
			handler(ctx)
		}

		ctx.LogMessage(w.Header().Get("status-code"), -1)
	}
}

// GetRouter creates the router
func GetRouter(environment string, debugMode, skipSecurityChecks bool) (*mux.Router, chan *logger.LogMsg, chan *logger.LogMsg, error) {
	skipSecurity = skipSecurityChecks
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.method, "OPTIONS").
				Path(route.routePattern(version)).
				Name(routeName).
				HandlerFunc(customHandler(routeName, version, route, mainLogChan, errorLogChan, environment, debugMode))
		}
	}

	for _, routeName := range []string{"index", "humans", "robots"} {
		version := ""
		route := routes["0.1"][routeName]
		if route == nil {
			panic(fmt.Sprintf("route %s not found", routeName))
		}
		router.
			Methods(route.method, "OPTIONS").
			Path(route.pattern).
			Name(routeName).
			HandlerFunc(customHandler(routeName, version, route, mainLogChan, errorLogChan, environment, debugMode))
	}

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	return router, mainLogChan, errorLogChan, nil
}
