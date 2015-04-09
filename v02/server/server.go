package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/context"
)

type (
	// RouteFunc defines the pattern for a route handling function
	RouteFunc func(*context.Context) *tgerrors.TGError

	// Route holds the route pattern
	Route struct {
		Method   string
		Pattern  string
		CPattern string
		Scope    string
		Handlers []RouteFunc
		Filters  []context.Filter
	}
)

const version = "0.2"

// RoutePattern returns the route pattern for a certain version
func (r *Route) RoutePattern(version string) string {
	if version == "" {
		return r.Pattern
	}
	return "/" + version + r.Pattern
}

// ComposePattern returns the composed pattern for a route
func (r *Route) ComposePattern(version string) string {
	return "/" + version + r.CPattern
}

// WriteResponse handles the http responses and returns the data
func WriteResponse(ctx *context.Context, response interface{}, code int, cacheTime uint) {
	// Set the response headers
	WriteCommonHeaders(cacheTime, ctx)
	WriteCorsHeaders(ctx)
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

// ErrorHappened handles the error message
func ErrorHappened(ctx *context.Context, err *tgerrors.TGError) {
	WriteCommonHeaders(0, ctx)
	ctx.W.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(int(err.Type))
		fmt.Fprintf(ctx.W, "%d %s", err.Type, err.Error())
	} else {
		ctx.W.Header().Set("Content-Encoding", "gzip")
		ctx.W.WriteHeader(int(err.Type))
		gz := gzip.NewWriter(ctx.W)
		fmt.Fprintf(gz, "%d %s", int(err.Type), err.Error())
		gz.Close()
	}

	ctx.StatusCode = int(err.Type)
	ctx.LogError(err)
}

// WriteCommonHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func WriteCommonHeaders(cacheTime uint, ctx *context.Context) {
	ctx.W.Header().Set("Strict-Transport-Security", "max-age=63072000")
	ctx.W.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.W.Header().Set("X-Frame-Options", "DENY")

	if cacheTime > 0 {
		ctx.W.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		ctx.W.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		ctx.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.W.Header().Set("Pragma", "no-cache")
		ctx.W.Header().Set("Expires", "0")
	}
}

// WriteCorsHeaders will write the needed CORS headers
func WriteCorsHeaders(ctx *context.Context) {
	ctx.W.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.W.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding, x-tapglue-id, x-tapglue-date, x-tapglue-session, x-tapglue-payload-hash, x-tapglue-signature")
	ctx.W.Header().Set("Access-Control-Allow-Credentials", "true")
}

// CorsHandler will handle the CORS requests
func CorsHandler(ctx *context.Context) (err *tgerrors.TGError) {
	WriteCommonHeaders(100, ctx)
	WriteCorsHeaders(ctx)
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
	return
}

// GetRoute takes a route name and returns the route including the version
func GetRoute(routeName string) *Route {
	if _, ok := Routes[routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table\n", routeName))
	}

	return Routes[routeName]
}

// ValidateGetCommon runs a series of predefined, common, tests for GET requests
func ValidateGetCommon(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.R.Header.Get("User-Agent") != "" {
		return
	}
	return tgerrors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
}

// ValidatePutCommon runs a series of predefinied, common, tests for PUT requests
func ValidatePutCommon(ctx *context.Context) (err *tgerrors.TGError) {
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

// ValidateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func ValidateDeleteCommon(ctx *context.Context) (err *tgerrors.TGError) {
	if ctx.R.Header.Get("User-Agent") != "" {
		return
	}

	return tgerrors.NewBadRequestError("User-Agent header must be set (1)", "missing ua header")
}

// ValidatePostCommon runs a series of predefined, common, tests for the POST requests
func ValidatePostCommon(ctx *context.Context) (err *tgerrors.TGError) {
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

// CustomHandler composes the http handling function according to its definition and returns it
func CustomHandler(routeName, version string, route *Route, mainLog, errorLog chan *logger.LogMsg, environment string, debugMode, skipSecurity bool) http.HandlerFunc {
	extraHandlers := []RouteFunc{CorsHandler}
	switch route.Method {
	case "DELETE":
		{
			extraHandlers = append(extraHandlers, ValidateDeleteCommon)
		}
	case "GET":
		{
			extraHandlers = append(extraHandlers, ValidateGetCommon)
		}
	case "PUT":
		{
			extraHandlers = append(extraHandlers, ValidatePutCommon)
		}
	case "POST":
		{
			extraHandlers = append(extraHandlers, ValidatePostCommon)
		}
	}

	route.Handlers = append(extraHandlers, route.Handlers...)

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, err := context.NewContext(w, r, mainLog, errorLog, routeName, route.Scope, version, route.Filters, environment, debugMode)
		if err != nil {
			ErrorHappened(ctx, err)
			return
		}

		ctx.SkipSecurity = skipSecurity

		for _, handler := range route.Handlers {
			// Any response that happens in a handler MUST send a Content-Type header
			if err := handler(ctx); err != nil {
				ErrorHappened(ctx, err)
				break
			}

		}

		ctx.LogRequest(ctx.StatusCode, -1)
	}
}

// Init takes care of initializing the router with the requests needed
func Init(router *mux.Router, mainLogChan, errorLogChan chan *logger.LogMsg, environment string, debugMode, skipSecurityChecks bool) {
	for routeName, route := range Routes {
		router.
			Methods(route.Method, "OPTIONS").
			Path(route.RoutePattern(version)).
			Name(routeName).
			HandlerFunc(CustomHandler(routeName, version, route, mainLogChan, errorLogChan, environment, debugMode, skipSecurityChecks))
	}
}
