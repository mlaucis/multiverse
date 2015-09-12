// Package server holds all the server related logic
package server

import (
	"compress/gzip"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	ratelimiter_redis "github.com/tapglue/backend/limiter/redis"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/tgflake"
	"github.com/tapglue/backend/utils"
	v02_server "github.com/tapglue/backend/v02/server"
	v02_kinesis "github.com/tapglue/backend/v02/storage/kinesis"
	v02_postgres "github.com/tapglue/backend/v02/storage/postgres"
	v03_server "github.com/tapglue/backend/v03/server"
	v03_kinesis "github.com/tapglue/backend/v03/storage/kinesis"
	v03_postgres "github.com/tapglue/backend/v03/storage/postgres"
	v03_redis "github.com/tapglue/backend/v03/storage/redis"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
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

	rawKinesisClient              v03_kinesis.Client
	rawPostgresClient             v03_postgres.Client
	rateLimiterPool, appCachePool *redigo.Pool
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
	environment string,
	debugMode, skipSecurityChecks bool,
) (*mux.Router, chan *logger.LogMsg, chan *logger.LogMsg, error) {
	router := mux.NewRouter().StrictSlash(true)

	v02_server.InitRouter(router, mainLogChan, errorLogChan, environment, skipSecurityChecks, debugMode)
	v03_server.InitRouter(router, mainLogChan, errorLogChan, environment, skipSecurityChecks, debugMode)

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
		router.PathPrefix("/debug/").Handler(http.DefaultServeMux)
	}

	router.NotFoundHandler = http.HandlerFunc(notFoundHandler)

	return router, mainLogChan, errorLogChan, nil
}

// SetupFlakes initializes the flakes for all the existing applications in the system
func SetupFlakes(db *sqlx.DB) {
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

// Setup initializes the dependencies
func Setup(conf *config.Config, revision, hostname string) {
	currentRevision = revision
	currentHostname = hostname

	rateLimiterPool = v03_redis.NewRedigoPool(conf.RateLimiter)
	applicationRateLimiter := ratelimiter_redis.NewLimiter(rateLimiterPool, "ratelimiter:app:")
	v02_server.SetupRateLimit(applicationRateLimiter)
	v03_server.SetupRateLimit(applicationRateLimiter)

	var (
		v02KinesisClient v02_kinesis.Client
		v03KinesisClient v03_kinesis.Client
	)
	if conf.Environment == "prod" {
		v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
		v03KinesisClient = v03_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
	} else {
		if conf.Kinesis.Endpoint != "" {
			v02KinesisClient = v02_kinesis.NewTest(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Kinesis.Endpoint, conf.Environment, conf.Kinesis.StreamName)
			v03KinesisClient = v03_kinesis.NewWithEndpoint(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Kinesis.Endpoint, conf.Environment, conf.Kinesis.StreamName)
		} else {
			v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
			v03KinesisClient = v03_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
		}
	}
	rawKinesisClient = v03KinesisClient

	v03KinesisClient.SetupStreams([]string{conf.Kinesis.StreamName})

	v02PostgresClient := v02_postgres.New(conf.Postgres)
	v03PostgresClient := v03_postgres.New(conf.Postgres)
	rawPostgresClient = v03PostgresClient

	SetupFlakes(v03PostgresClient.SlaveDatastore(-1))

	appCachePool = v03_redis.NewRedigoPool(conf.CacheApp)

	v02_server.Setup(v02KinesisClient, v02PostgresClient, currentRevision, currentHostname)
	v03_server.Setup(v03KinesisClient, v03PostgresClient, appCachePool, currentRevision, currentHostname)

}
