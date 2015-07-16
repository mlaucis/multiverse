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
	"strings"
	"time"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/limiter"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/utils"
	v02_core "github.com/tapglue/backend/v02/core"
	v02_server "github.com/tapglue/backend/v02/server"
	v02_kinesis "github.com/tapglue/backend/v02/storage/kinesis"
	v02_postgres "github.com/tapglue/backend/v02/storage/postgres"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/yvasiyarov/gorelic"
)

type (
	errorResponse struct {
		Code             int    `json:"code"`
		Message          string `json:"message"`
		DocumentationURL string `json:"documentation_url,omitempty"`
	}
)

var (
	mainLogChan  = make(chan *logger.LogMsg, 100000)
	errorLogChan = make(chan *logger.LogMsg, 100000)

	currentRevision = ""
	currentHostname = ""

	rawKinesisClient   v02_kinesis.Client
	rawPostgresClient  v02_postgres.Client
	rawRateLimiterPool *redigo.Pool
)

// WriteCommonHeaders will add the corresponding cache headers based on the time supplied (in seconds)
func WriteCommonHeaders(cacheTime uint, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "User-Agent, Content-Type, Content-Length, Accept-Encoding")
	w.Header().Set("Access-Control-Allow-Credentials", "true")

	w.Header().Set("Strict-Transport-Security", "max-age=63072000")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")

	w.Header().Set("X-Tapglue-Hash", currentRevision)
	w.Header().Set("X-Tapglue-Server", currentHostname)

	if cacheTime > 0 {
		w.Header().Set("Cache-Control", fmt.Sprintf(`"max-age=%d, public"`, cacheTime))
		w.Header().Set("Expires", time.Now().Add(time.Duration(cacheTime)*time.Second).Format(http.TimeFormat))
	} else {
		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")
	}
}

// WriteResponse handles the http responses and returns the data
func WriteResponse(ctx *context.Context, response interface{}, code int, cacheTime uint) {
	// Set the response headers
	WriteCommonHeaders(cacheTime, ctx.W, ctx.R)

	ctx.W.Header().Set("Content-Type", "application/json; charset=UTF-8")
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
func ErrorHappened(ctx *context.Context, err []errors.Error) {
	errorMessage := []errorResponse{}
	for idx := range err {
		errorMessage = append(errorMessage, errorResponse{Code: err[idx].Code(), Message: err[idx].Error()})
	}

	WriteResponse(ctx, errorMessage, int(err[0].Type()), 0)
	go ctx.LogError(err)
}

// NewContext creates a new Context
func NewContext(
	w http.ResponseWriter,
	r *http.Request,
	p map[string]string,
	mainLog, errorLog chan *logger.LogMsg,
	routeName, environment string,
	debugMode bool,
) (*context.Context, []errors.Error) {

	ctx := new(context.Context)
	ctx.StartTime = time.Now()
	ctx.R = r
	ctx.W = w
	if p != nil {
		ctx.Vars = p
	} else {
		ctx.Vars = map[string]string{}
	}
	ctx.MainLog = mainLog
	ctx.ErrorLog = errorLog
	if r.Method != "GET" {
		ctx.Body = utils.PeakBody(r).Bytes()
	}
	ctx.RouteName = routeName
	ctx.Environment = environment
	ctx.DebugMode = debugMode
	ctx.Bag = map[string]interface{}{}
	ctx.Bag["rateLimit.enabled"] = false
	ctx.AuthUsername, ctx.AuthPassword, ctx.AuthOk = r.BasicAuth()
	ctx.Query = r.URL.Query()

	return ctx, nil
}

// GetRouter creates the router
func GetRouter(
	agent *gorelic.Agent,
	environment string,
	debugMode, skipSecurityChecks bool,
) (*mux.Router, chan *logger.LogMsg, chan *logger.LogMsg, error) {
	router := mux.NewRouter().StrictSlash(true)

	v02_server.InitRouter(agent, router, mainLogChan, errorLogChan, environment, skipSecurityChecks, debugMode)

	for idx := range generalRoutes {
		router.
			Methods(generalRoutes[idx].method).
			Path(generalRoutes[idx].path).
			HandlerFunc(func(route generalRoute) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				ctx, err := NewContext(w, r, mux.Vars(r), mainLogChan, errorLogChan, route.name, "", false)
				ctx.SkipSecurity = true
				if err != nil {
					go ErrorHappened(ctx, err)
					return
				}

				route.handler(ctx)
				if ctx.R.Header.Get("User-Agent") == "ELB-HealthChecker/1.0" {
					return
				} else if ctx.R.Header.Get("User-Agent") == "updown.io bot 2.0" {
					return
				}
				go ctx.LogRequest(ctx.StatusCode, -1)
			}
		}(generalRoutes[idx]))
	}

	router.
		Methods("GET").
		Path("/favicon.ico").
		HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./favicon.ico")
	})

	if debugMode {
		router.Methods("GET").Path("/debug/pprof").Handler(http.HandlerFunc(pprof.Index))
		router.Methods("GET").Path("/debug/pprof/cmdline").Handler(http.HandlerFunc(pprof.Cmdline))
		router.Methods("GET").Path("/debug/pprof/profile").Handler(http.HandlerFunc(pprof.Profile))
		router.Methods("GET").Path("/debug/pprof/symbol").Handler(http.HandlerFunc(pprof.Symbol))
	}

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router, mainLogChan, errorLogChan, nil
}

// SetupRateLimit will allow for proper rate limits to be configured
func SetupRateLimit(applicationRateLimiter limiter.Limiter) {
	v02_server.SetupRateLimit(applicationRateLimiter)
}

func SetupRawConnections(
	kinesisClient v02_kinesis.Client,
	postgresClient v02_postgres.Client,
	rateLimiterPool *redigo.Pool) {

	rawKinesisClient = kinesisClient
	rawPostgresClient = postgresClient
	rawRateLimiterPool = rateLimiterPool
}

// SetupRedisCores takes care of initializing the core
func SetupRedisCores(
	account v02_core.Account,
	accountUser v02_core.AccountUser,
	application v02_core.Application,
	applicationUser v02_core.ApplicationUser,
	connection v02_core.Connection,
	event v02_core.Event) {

	v02_server.SetupRedisCores(account, accountUser, application, applicationUser, connection, event)
}

// SetupKinesisCores takes care of initializing the core
func SetupKinesisCores(
	account v02_core.Account,

	accountUser v02_core.AccountUser,
	application v02_core.Application,
	applicationUser v02_core.ApplicationUser,
	connection v02_core.Connection,
	event v02_core.Event) {

	v02_server.SetupKinesisCores(account, accountUser, application, applicationUser, connection, event)
}

// SetupPostgresCores takes care of initializing the PostgresSQL core
func SetupPostgresCores(
	account v02_core.Account,
	accountUser v02_core.AccountUser,
	application v02_core.Application,
	applicationUser v02_core.ApplicationUser,
	connection v02_core.Connection,
	event v02_core.Event) {

	v02_server.SetupPostgresCores(account, accountUser, application, applicationUser, connection, event)
}

func SetupFlakes() {
	v02_server.SetupFlakes(rawPostgresClient)
}

// Setup initializes the dependencies
// Must be called after initializing the cores
func Setup(revision, hostname string) {
	currentRevision = revision
	currentHostname = hostname
	v02_server.Setup(currentRevision, currentHostname)
}
