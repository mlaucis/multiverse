package server

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/v02/errmsg"

	"github.com/gorilla/mux"
)

type (
	// RouteFunc defines the pattern for a route handling function
	RouteFunc func(*context.Context) []errors.Error

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

var (
	currentRevision, currentHostname string
)

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
	ctx.StatusCode = code

	// Set the response headers
	WriteCommonHeaders(cacheTime, ctx)
	WriteCorsHeaders(ctx)

	// TODO here it would be nice if we would consider the requested format when the stuff happens and deliver
	// either JSON or XML or FlatBuffers or whatever
	output := new(bytes.Buffer)
	json.NewEncoder(output).Encode(response)

	// We should only check these for valid responses, I think. Future me blame past me for this decision
	if ctx.R.Method == "GET" && ctx.StatusCode < 300 {
		// We implememt the etag check first, aka the hard check, because we don't know if something else was changed
		// in the response, not just what we calculate the etag for.
		// Example situation: we compute the etag for getFeed as being the highest LastUpdated date but maybe
		// a user was updated meanwhile which would mean that the feed might be the same, event wise, but the user wise
		// it will be different so the app should retrieve the feed and process it as maybe the display name of the user
		// was changed or something else (thumbnail or whatever)

		h := md5.New()
		io.TeeReader(output, h)
		etag := h.Sum(nil)
		ctx.W.Header().Set("ETag", fmt.Sprintf("%x", etag))

		if requestEtag := ctx.R.Header.Get("If-None-Match"); requestEtag != "" {
			if requestEtag == string(etag) {
				ctx.StatusCode = http.StatusNotModified
				ctx.W.WriteHeader(ctx.StatusCode)
				return
			}
		}

		if ifModifiedSince := ctx.R.Header.Get("If-Modified-Since"); ifModifiedSince != "" {
			if myLastModified, ok := ctx.Bag["Last-Modified"]; ok {
				ctx.W.Header().Set("Last-Modified", myLastModified.(string))
				if myLastModified.(string) == ifModifiedSince {
					ctx.StatusCode = http.StatusNotModified
					ctx.W.WriteHeader(ctx.StatusCode)
					return
				}
			}
		}
	}

	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(code)
		io.Copy(ctx.W, output)
		return
	}

	ctx.W.Header().Set("Content-Encoding", "gzip")
	ctx.W.WriteHeader(code)
	gz := gzip.NewWriter(ctx.W)
	io.Copy(gz, output)
	gz.Close()
}

// ErrorHappened handles the error message
func ErrorHappened(ctx *context.Context, err []errors.Error) {
	ctx.StatusCode = int(err[0].Type())

	WriteCommonHeaders(0, ctx)
	WriteCorsHeaders(ctx)

	// Write response
	if !strings.Contains(ctx.R.Header.Get("Accept-Encoding"), "gzip") {
		// No gzip support
		ctx.W.WriteHeader(ctx.StatusCode)
		fmt.Fprintf(ctx.W, "%d %s", err[0].Type(), err[0].Error())
	} else {
		ctx.W.Header().Set("Content-Encoding", "gzip")
		ctx.W.WriteHeader(ctx.StatusCode)
		gz := gzip.NewWriter(ctx.W)
		fmt.Fprintf(gz, "%d %s", int(err[0].Type()), err[0].Error())
		gz.Close()
	}

	go ctx.LogError(err)
}

// WriteCommonHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func WriteCommonHeaders(cacheTime uint, ctx *context.Context) {
	ctx.W.Header().Set("Strict-Transport-Security", "max-age=63072000")
	ctx.W.Header().Set("X-Content-Type-Options", "nosniff")
	ctx.W.Header().Set("X-Frame-Options", "DENY")
	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")

	ctx.W.Header().Set("X-Tapglue-Hash", currentRevision)
	ctx.W.Header().Set("X-Tapglue-Server", currentHostname)

	if cacheTime > 0 {
		ctx.W.Header().Set("Cache-Control", fmt.Sprintf(`max-age=%d, public`, cacheTime))
		ctx.W.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		ctx.W.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		ctx.W.Header().Set("Pragma", "no-cache")
		ctx.W.Header().Set("Expires", "0")
	}

	if ctx.R.Method == "GET" && ctx.StatusCode < 300 {
		if myLastModified, ok := ctx.Bag["Last-Modified"]; ok {
			ctx.W.Header().Set("Last-Modified", myLastModified.(string))
		} else {
			// This will spam the server logs for issues with missing issues but then again, it should be there...
			go ctx.LogError(errmsg.ErrMissingLastModifiedHeader.UpdateInternalMessage("missing Last-Modified from bag for route " + ctx.RouteName + " response"))
		}
	}

	if !ctx.Bag["rateLimit.enabled"].(bool) {
		return
	}
	ctx.W.Header().Set("X-RateLimit-Limit", strconv.FormatInt(1000, 10))
	ctx.W.Header().Set("X-RateLimit-Remaining", strconv.FormatInt(ctx.Bag["rateLimit.limit"].(int64), 10))
	ctx.W.Header().Set("X-RateLimit-Reset", strconv.FormatInt(ctx.Bag["rateLimit.refreshTime"].(time.Time).Unix(), 10))
}

// WriteCorsHeaders will write the needed CORS headers
func WriteCorsHeaders(ctx *context.Context) {
	ctx.W.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.W.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	ctx.W.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding, Authorization")
	ctx.W.Header().Set("Access-Control-Allow-Credentials", "true")
}

// CorsHandler will handle the CORS requests
func CorsHandler(ctx *context.Context) (err []errors.Error) {
	WriteCommonHeaders(100, ctx)
	WriteCorsHeaders(ctx)
	return
}

// VersionHandler handlers the requests for the version status
func VersionHandler(ctx *context.Context) []errors.Error {
	response := struct {
		Version string `json:"version"`
		Status  string `json:"status"`
	}{"v" + version, "current"}
	WriteResponse(ctx, response, 200, 86400)
	return nil
}

// GetRoute takes a route name and returns the route including the version
func GetRoute(routeName string) *Route {
	if _, ok := Routes[routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table\n", routeName))
	}

	return Routes[routeName]
}

// ValidateGetCommon runs a series of predefined, common, tests for GET requests
func ValidateGetCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.R.Header.Get("User-Agent") != "" {
		return
	}
	return []errors.Error{errmsg.ErrBadUserAgent}
}

// ValidatePutCommon runs a series of predefinied, common, tests for PUT requests
func ValidatePutCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		err = append(err, errmsg.ErrBadUserAgent)
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		err = append(err, errmsg.ErrContentLengthMissing)
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		err = append(err, errmsg.ErrContentTypeMissing)
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		err = append(err, errmsg.ErrContentTypeMismatch)
	}

	reqCL, er := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if er != nil {
		err = append(err, errmsg.ErrContentLengthInvalid)
	}

	if reqCL != ctx.R.ContentLength {
		err = append(err, errmsg.ErrContentLengthSizeMismatch)
	} else {
		// TODO better handling here for limits, maybe make them customizable
		if reqCL > 2048 {
			err = append(err, errmsg.ErrPayloadTooBig)
		}
	}

	if ctx.R.Body == nil {
		err = append(err, errmsg.ErrRequestBodyEmpty)
	}
	return
}

// ValidateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func ValidateDeleteCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.R.Header.Get("User-Agent") == "" {
		err = append(err, errmsg.ErrBadUserAgent)
	}

	return
}

// ValidatePostCommon runs a series of predefined, common, tests for the POST requests
func ValidatePostCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		err = append(err, errmsg.ErrBadUserAgent)
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		err = append(err, errmsg.ErrContentLengthMissing)
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		err = append(err, errmsg.ErrContentTypeMissing)
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		err = append(err, errmsg.ErrContentTypeMismatch)
	}

	reqCL, er := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if er != nil {
		err = append(err, errmsg.ErrContentLengthInvalid)
	}

	if reqCL != ctx.R.ContentLength {
		err = append(err, errmsg.ErrContentLengthSizeMismatch)
	} else {
		// TODO better handling here for limits, maybe make them customizable
		if reqCL > 2048 {
			err = append(err, errmsg.ErrPayloadTooBig)
		}
	}

	if ctx.R.Body == nil {
		err = append(err, errmsg.ErrRequestBodyEmpty)
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
	route.Filters = append([]context.Filter{func(ctx *context.Context) []errors.Error {
		ctx.Vars = mux.Vars(ctx.R)
		return nil
	}}, route.Filters...)

	return func(w http.ResponseWriter, r *http.Request) {

		ctx, err := context.New(w, r, mainLog, errorLog, routeName, route.Scope, version, route.Filters, environment, debugMode)
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

		go ctx.LogRequest(ctx.StatusCode, -1)
	}
}

// Init takes care of initializing the router with the requests needed
func Init(revision, hostname string) {
	currentRevision = revision
	currentHostname = hostname
}
