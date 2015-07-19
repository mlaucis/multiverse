// Package server provides handling for all the requests towards this module
package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/limiter"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server/response"
	"github.com/tapglue/backend/v02/storage/postgres"

	"strings"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/tgflake"
)

type (
	errorResponse struct {
		Code             int    `json:"code"`
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url,omitempty"`
	}
)

const (
	// APIVersion holds which API Version does this module holds
	APIVersion = "0.2"

	appRateLimit        int64 = 1000
	appRateLimitSeconds int64 = 60
)

var (
	postgresAccount, kinesisAccount                 core.Account
	postgresAccountUser, kinesisAccountUser         core.AccountUser
	postgresApplication, kinesisApplication         core.Application
	postgresApplicationUser, kinesisApplicationUser core.ApplicationUser
	postgresConnection, kinesisConnection           core.Connection
	postgresEvent, kinesisEvent                     core.Event

	appRateLimiter limiter.Limiter
)

// ValidateGetCommon runs a series of predefined, common, tests for GET requests
func ValidateGetCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.R.Header.Get("User-Agent") != "" {
		return
	}
	return []errors.Error{errmsg.ErrServerReqBadUserAgent}
}

// ValidatePutCommon runs a series of predefinied, common, tests for PUT requests
func ValidatePutCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		err = append(err, errmsg.ErrServerReqBadUserAgent)
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		err = append(err, errmsg.ErrServerReqContentLengthMissing)
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		err = append(err, errmsg.ErrServerReqContentTypeMissing)
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		err = append(err, errmsg.ErrServerReqContentTypeMismatch)
	}

	reqCL, er := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if er != nil {
		err = append(err, errmsg.ErrServerReqContentLengthInvalid)
	}

	if reqCL != ctx.R.ContentLength {
		err = append(err, errmsg.ErrServerReqContentLengthSizeMismatch)
	} else {
		// TODO better handling here for limits, maybe make them customizable
		if reqCL > 2048 {
			err = append(err, errmsg.ErrServerReqPayloadTooBig)
		}
	}

	if ctx.R.Body == nil {
		err = append(err, errmsg.ErrServerReqBodyEmpty)
	}
	return
}

// ValidateDeleteCommon runs a series of predefinied, common, tests for DELETE requests
func ValidateDeleteCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.R.Header.Get("User-Agent") == "" {
		err = append(err, errmsg.ErrServerReqBadUserAgent)
	}

	return
}

// ValidatePostCommon runs a series of predefined, common, tests for the POST requests
func ValidatePostCommon(ctx *context.Context) (err []errors.Error) {
	if ctx.SkipSecurity {
		return
	}

	if ctx.R.Header.Get("User-Agent") == "" {
		err = append(err, errmsg.ErrServerReqBadUserAgent)
	}

	if ctx.R.Header.Get("Content-Length") == "" {
		err = append(err, errmsg.ErrServerReqContentLengthMissing)
	}

	if ctx.R.Header.Get("Content-Type") == "" {
		err = append(err, errmsg.ErrServerReqContentTypeMissing)
	}

	if ctx.R.Header.Get("Content-Type") != "application/json" &&
		ctx.R.Header.Get("Content-Type") != "application/json; charset=UTF-8" {
		err = append(err, errmsg.ErrServerReqContentTypeMismatch)
	}

	reqCL, er := strconv.ParseInt(ctx.R.Header.Get("Content-Length"), 10, 64)
	if er != nil {
		err = append(err, errmsg.ErrServerReqContentLengthInvalid)
	}

	if reqCL != ctx.R.ContentLength {
		err = append(err, errmsg.ErrServerReqContentLengthSizeMismatch)
	} else {
		// TODO better handling here for limits, maybe make them customizable
		if reqCL > 2048 {
			err = append(err, errmsg.ErrServerReqPayloadTooBig)
		}
	}

	if ctx.R.Body == nil {
		err = append(err, errmsg.ErrServerReqBodyEmpty)
	}
	return
}

// GetRoute takes a route name and returns the route including the version
func GetRoute(routeName string) *Route {
	for idx := range Routes {
		if Routes[idx].Name == routeName {
			return Routes[idx]
		}
	}

	panic(fmt.Sprintf("route %q not found", routeName))
}

// RateLimitApplication takes care of appling the rate limits for the application
func RateLimitApplication(ctx *context.Context) []errors.Error {
	if ctx.SkipSecurity {
		return nil
	}

	hash, _, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errors.NewBadRequestError(2300, "something went wrong with the authentication", "something went wrong with the authentication")}
	}

	hash = fmt.Sprintf("%s.%s", hash, ctx.R.Method)

	limit, refreshTime, err := appRateLimiter.Request(&limiter.Limitee{hash, appRateLimit, appRateLimitSeconds})
	if err != nil {
		return []errors.Error{errors.NewInternalError(0, "something went wrong", err.Error())}
	}

	ctx.Bag["rateLimit.enabled"] = true
	ctx.Bag["rateLimit.limit"] = limit
	ctx.Bag["rateLimit.refreshTime"] = refreshTime

	if limit == 0 {
		return []errors.Error{errors.New(429, 0, "Too Many Requests", "over quota", false)}
	}

	return nil
}

// CustomHandler generates the handler for a certain route
func CustomHandler(route *Route, mainLogChan, errorLogChan chan *logger.LogMsg, env string, skipSecurity, debug bool) http.HandlerFunc {
	extraHandlers := []RouteFunc{}
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
		ctx, err := NewContext(w, r, mux.Vars(r), mainLogChan, errorLogChan, route, env, debug)
		if err != nil {
			response.ErrorHappened(ctx, err)
			return
		}
		ctx.SkipSecurity = skipSecurity

		for idx := range route.Filters {
			if err = route.Filters[idx](ctx); err != nil {
				response.ErrorHappened(ctx, err)
				return
			}
		}

		for idx := range route.Handlers {
			if err = route.Handlers[idx](ctx); err != nil {
				response.ErrorHappened(ctx, err)
				return
			}
		}

		go ctx.LogRequest(ctx.StatusCode, -1)
	}
}

// CustomOptionsHandler handles all the OPTIONS requests for us
func CustomOptionsHandler(route *Route, mainLogChan, errorLogChan chan *logger.LogMsg, env string, skipSecurity, debug bool) http.HandlerFunc {
	// Override the route method to what we need
	route.Method = "OPTIONS"
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, err := NewContext(w, r, mux.Vars(r), mainLogChan, errorLogChan, route, env, debug)
		ctx.SkipSecurity = skipSecurity
		if err != nil {
			go response.ErrorHappened(ctx, err)
			return
		}
		ctx.StatusCode = 200
		if err = response.CORSHandler(ctx); err != nil {
			go response.ErrorHappened(ctx, err)
			return
		}

		if ctx.R.Header.Get("User-Agent") == "ELB-HealthChecker/1.0" {
			return
		} else if ctx.R.Header.Get("User-Agent") == "updown.io bot 2.0" {
			return
		}

		go ctx.LogRequest(ctx.StatusCode, -1)
	}
}

// SetupRateLimit initializes the rate limiters
func SetupRateLimit(applicationRateLimiter limiter.Limiter) {
	appRateLimiter = applicationRateLimiter
}

// SetupKinesisCores takes care of initializing the redis core
func SetupKinesisCores(
	account core.Account,
	accountUser core.AccountUser,
	application core.Application,
	applicationUser core.ApplicationUser,
	connection core.Connection,
	event core.Event) {
	kinesisAccount = account
	kinesisAccountUser = accountUser
	kinesisApplication = application
	kinesisApplicationUser = applicationUser
	kinesisConnection = connection
	kinesisEvent = event
}

// SetupPostgresCores takes care of initializing the PostgresSQL core
func SetupPostgresCores(
	account core.Account,
	accountUser core.AccountUser,
	application core.Application,
	applicationUser core.ApplicationUser,
	connection core.Connection,
	event core.Event) {
	postgresAccount = account
	postgresAccountUser = accountUser
	postgresApplication = application
	postgresApplicationUser = applicationUser
	postgresConnection = connection
	postgresEvent = event
}

// SetupFlakes initializes the flakes for all the existing applications in the system
func SetupFlakes(storageClient postgres.Client) {
	db := storageClient.MainDatastore()

	existingSchemas, err := db.Query(`SELECT nspname FROM pg_catalog.pg_namespace WHERE nspname ILIKE 'app_%_%'`)
	if err != nil {
		panic(err)
	}
	defer existingSchemas.Close()
	for existingSchemas.Next() {
		schemaName := ""
		err := existingSchemas.Scan(&schemaName)
		if err != nil {
			panic(err)
		}
		details := strings.Split(schemaName, "_")
		if len(details) != 3 || details[0] != "app" {
			continue
		}

		appID, err := strconv.ParseInt(details[2], 10, 64)
		if err != nil {
			panic(err)
		}
		_ = tgflake.Flake(appID, "users")
		_ = tgflake.Flake(appID, "events")
	}
}

// Setup initializes the route handlers
// Must be called after initializing the cores
func Setup(revision, hostname string) {
	if appRateLimiter == nil {
		panic("You must first initialize the rate limiter")
	}

	if kinesisAccount == nil || postgresAccount == nil {
		panic("You must initialize the kinesis and postgres cores first")
	}

	if revision == "" {
		panic("omfg missing revision")
	}

	response.Setup(revision, hostname)
	InitHandlers()

	Routes = SetupRoutes()
}
