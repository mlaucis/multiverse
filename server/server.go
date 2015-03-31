/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"strconv"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/server/utils"
	"github.com/tapglue/backend/tgerrors"
	v1_server "github.com/tapglue/backend/v01/server"

	"github.com/gorilla/mux"
)

const (
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
	routes       = make(map[string]map[string]*utils.Route)
)

// isRequestExpired checks if the request is expired or not
func isRequestExpired(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	// Check that the request is not older than 3 days
	// TODO check if we should lower the interval
	requestDate := ctx.R.Header.Get("x-tapglue-date")
	if requestDate == "" {
		return tgerrors.NewBadRequestError("invalid request date(1)\nrequest date is empty", "empty request date")
	}

	parsedRequestDate, er := time.Parse(time.RFC3339, requestDate)
	if er != nil {
		return tgerrors.NewBadRequestError("invalid request date(2)\nrequest date is not a valid format", "not rfc3339 request date")
	}

	if time.Since(parsedRequestDate) > time.Duration(3*time.Minute) {
		return tgerrors.NewBadRequestError("invalid request date(1)\nrequest expired", "request expired")
	}
	return
}

// validateGetCommon runs a series of predefined, common, tests for GET requests
func validateGetCommon(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.R.Header.Get("User-Agent") != "" {
		return
	}
	return tgerrors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
}

// validatePutCommon runs a series of predefinied, common, tests for PUT requests
func validatePutCommon(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		return tgerrors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		return tgerrors.NewBadRequestError("Content-Length header missing", "missing content-length header")
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		return tgerrors.NewBadRequestError("Content-Type header empty", "missing content-type header")
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		return tgerrors.NewBadRequestError("Content-Type header is empty", "content-type header mismatch")
	}

	reqCL, er := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if er != nil {
		return tgerrors.NewBadRequestError("Content-Length header is invalid", "content-length header is not an int")
	}

	if reqCL != ctx.R.ContentLength {
		return tgerrors.NewBadRequestError("Content-Length header size mismatch", "content-length header size mismatch")
	}

	if ctx.R.Body == nil {
		return tgerrors.NewBadRequestError("Empty request body", "empty request body")
	}
	return
}

// validateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func validateDeleteCommon(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.R.Header.Get("User-Agent") != "" {
		return
	}

	return tgerrors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
}

// validatePostCommon runs a series of predefined, common, tests for the POST requests
func validatePostCommon(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		return tgerrors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		return tgerrors.NewBadRequestError("Content-Length header missing", "missing content-length header")
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		return tgerrors.NewBadRequestError("Content-Type header empty", "missing content-type header")
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		return tgerrors.NewBadRequestError("Content-Type header is empty", "content-type header mismatch")
	}

	reqCL, er := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if er != nil {
		return tgerrors.NewBadRequestError("Content-Length header is invalid", "content-length header is not an int")
	}

	if reqCL != ctx.R.ContentLength {
		return tgerrors.NewBadRequestError("Content-Length header size mismatch", "content-length header size mismatch")
	}

	if ctx.R.Body == nil {
		return tgerrors.NewBadRequestError("Empty request body", "empty request body")
	}
	return
}

// GetRoute takes a route name and returns the route including the version
func GetRoute(routeName, apiVersion string) *utils.Route {
	if _, ok := routes[apiVersion][routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version \"%s\"\n", routeName, apiVersion))
	}

	return routes[apiVersion][routeName]
}

// CustomHandler composes the http handling function according to its definition and returns it
func CustomHandler(routeName, version string, route *utils.Route, mainLog, errorLog chan *logger.LogMsg, environment string, debugMode, skipSecurity bool) http.HandlerFunc {
	extraHandlers := []utils.RouteFunc{utils.CorsHandler}
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

	if version != "0.1" && version != "" {
		extraHandlers = append(extraHandlers, isRequestExpired)
	}

	route.Handlers = append(extraHandlers, route.Handlers...)

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, err := context.NewContext(w, r, mainLog, errorLog, routeName, route.Scope, version, route.Filters, environment, debugMode)
		if err != nil {
			utils.ErrorHappened(ctx, err)
			return
		}

		ctx.SkipSecurity = skipSecurity

		for _, handler := range route.Handlers {
			// Any response that happens in a handler MUST send a Content-Type header
			if err := handler(ctx); err != nil {
				utils.ErrorHappened(ctx, err)
				break
			}

		}

		ctx.LogRequest(ctx.StatusCode, -1)
	}
}

// GetRouter creates the router
func GetRouter(environment string, debugMode, skipSecurityChecks bool) (*mux.Router, chan *logger.LogMsg, chan *logger.LogMsg, error) {
	router := mux.NewRouter().StrictSlash(true)

	for version, innerRoutes := range routes {
		for routeName, route := range innerRoutes {
			router.
				Methods(route.Method, "OPTIONS").
				Path(route.RoutePattern(version)).
				Name(routeName).
				HandlerFunc(CustomHandler(routeName, version, route, mainLogChan, errorLogChan, environment, debugMode, skipSecurityChecks))
		}
	}

	if debugMode {
		router.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		router.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		router.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		router.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	}

	router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	return router, mainLogChan, errorLogChan, nil
}

// Init ensures that the package is properly initialized
func Init() {
	routes[""] = defaultRoutes
	routes["0.1"] = v1_server.Routes
}
