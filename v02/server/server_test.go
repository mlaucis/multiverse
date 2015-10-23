package server_test

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/errors"
	ratelimiter_redis "github.com/tapglue/multiverse/limiter/redis"
	"github.com/tapglue/multiverse/logger"
	. "github.com/tapglue/multiverse/utils"
	"github.com/tapglue/multiverse/v02/core"
	v02_postgres_core "github.com/tapglue/multiverse/v02/core/postgres"
	"github.com/tapglue/multiverse/v02/entity"
	"github.com/tapglue/multiverse/v02/server"
	v02_postgres "github.com/tapglue/multiverse/v02/storage/postgres"
	v02_redis "github.com/tapglue/multiverse/v02/storage/redis"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	t.Parallel()
	TestingT(t)
}

type (
	ServerSuite          struct{}
	AccountSuite         struct{}
	AccountUserSuite     struct{}
	ApplicationSuite     struct{}
	ApplicationUserSuite struct{}
	ConnectionSuite      struct{}
	EventSuite           struct{}
)

var (
	_ = Suite(&ServerSuite{})
	_ = Suite(&AccountSuite{})
	_ = Suite(&AccountUserSuite{})
	_ = Suite(&ApplicationSuite{})
	_ = Suite(&ApplicationUserSuite{})
	_ = Suite(&ConnectionSuite{})
	_ = Suite(&EventSuite{})

	conf               *config.Config
	doLogTest          = flag.Bool("lt", false, "Set flag in order to get logs output from the tests")
	doCurlLogs         = flag.Bool("ct", false, "Set flag in order to get logs output from the tests as curl requests, sets -lt=true")
	doLogResponseTimes = flag.Bool("rt", false, "Set flag in order to get logs with response times only")
	doLogResponses     = flag.Bool("rl", false, "Set flag in order to get logs with response headers and bodies")
	quickBenchmark     = flag.Bool("qb", false, "Set flag in order to run only the benchmarks and skip all tests")
	mainLogChan        = make(chan *logger.LogMsg)
	errorLogChan       = make(chan *logger.LogMsg)

	coreAcc     core.Account
	coreAccUser core.AccountUser
	coreApp     core.Application
	coreAppUser core.ApplicationUser
	coreConn    core.Connection
	coreEvt     core.Event

	v02PostgresClient v02_postgres.Client

	nilTime             *time.Time
	redigoRateLimitPool *redigo.Pool
)

func init() {
	flag.Parse()

	if *doCurlLogs {
		*doLogTest = true
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")

	errors.Init(conf.Environment != "prod")

	if *doLogResponseTimes {
		go logger.TGLogResponseTimes(mainLogChan)
		go logger.TGLogResponseTimes(errorLogChan)
	} else if *doLogTest {
		if *doCurlLogs {
			go logger.TGCurlLog(mainLogChan)
		} else {
			go logger.TGLog(mainLogChan)
		}
		go logger.TGLog(errorLogChan)
	} else {
		go logger.TGSilentLog(mainLogChan)
		go logger.TGSilentLog(errorLogChan)
	}

	v02PostgresClient = v02_postgres.New(conf.Postgres)

	redigoRateLimitPool = v02_redis.NewRedigoPool(conf.RateLimiter.Hosts[0], "")
	applicationRateLimiter := ratelimiter_redis.NewLimiter(redigoRateLimitPool, "ratelimiter:app:")

	coreAcc = v02_postgres_core.NewAccount(v02PostgresClient)
	coreAccUser = v02_postgres_core.NewAccountUser(v02PostgresClient)
	coreApp = v02_postgres_core.NewApplication(v02PostgresClient)
	coreAppUser = v02_postgres_core.NewApplicationUser(v02PostgresClient)
	coreConn = v02_postgres_core.NewConnection(v02PostgresClient)
	coreEvt = v02_postgres_core.NewEvent(v02PostgresClient)

	server.SetupRateLimit(applicationRateLimiter)
	server.Setup(v02PostgresClient, "HEAD", "CI-Machine")

	testBootup(conf.Postgres)

	createdAt := struct {
		CreatedAt *time.Time
	}{}
	er := json.Unmarshal([]byte(`{"created_at": null}`), &createdAt)
	if er != nil {
		panic(er)
	}
	nilTime = createdAt.CreatedAt
}

func (s *ServerSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

func (s *AccountSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

func (s *AccountUserSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

func (s *ApplicationSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

func (s *ApplicationUserSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

func (s *ConnectionSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

func (s *EventSuite) SetUpSuite(c *C) {
	if *quickBenchmark {
		c.Skip("Running in quick benchmark mode")
	}
}

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := getRoute(routeName)
	routePath := requestRoute.RoutePattern()

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "{\"errors\":[{\"code\":5002,\"message\":\"User-Agent header must be set (1)\"},{\"code\":5004,\"message\":\"Content-Length header missing\"},{\"code\":5007,\"message\":\"Content-Type header empty\"},{\"code\":5006,\"message\":\"Content-Type header mismatch\"},{\"code\":5003,\"message\":\"Content-Length header is invalid\"},{\"code\":5005,\"message\":\"Content-Length header size mismatch\"}]}\n")
}

// Test POST common with CLHeader
func (s *ServerSuite) TestValidatePostCommon_CLHeader(c *C) {
	c.Skip("deprecated api version")
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := getRoute(routeName)
	routePath := requestRoute.RoutePattern()

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, `{"errors":[{"code":5001,"message":"invalid character 'd' looking for beginning of object key string"}]}`+"\n")
}

// Test PUT common with CLHeader
func (s *ServerSuite) TestValidatePutCommon_CLHeader(c *C) {
	c.Skip("needs a better implementation")

	payload := "{demo}"
	routeName := "updateAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test PUT common without CLHeader
func (s *ServerSuite) TestValidatePutCommon_NoCLHeader(c *C) {
	c.Skip("this needs a better implementation now that contexts are in place")

	payload := "{demo}"
	routeName := "updateAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test DELETE common with CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_CLHeader(c *C) {
	c.Skip("needs a better implementation")

	payload := "{demo}"
	routeName := "deleteAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test DELETE common without CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_NoCLHeader(c *C) {
	c.Skip("skip due to context refactoring")
	payload := ""
	routeName := "deleteAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 1)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

func (s *ServerSuite) TestRateLimit(c *C) {
	var (
		limit       = 10
		accounts    = CorrectDeploy(1, 0, 1, 1, 0, false, true)
		application = accounts[0].Applications[0]
		user        = application.Users[0]
		routeName   = "getCurrentApplicationUser"
		route       = getComposedRoute(routeName)
	)

	for i := 1; i < limit; i++ {
		code, body, headers, err := runRequestWithHeaders(routeName, route, "", func(*http.Request) {}, signApplicationRequest(application, user, true, true))
		c.Assert(err, IsNil)
		c.Assert(code, Equals, http.StatusOK)
		c.Assert(body, Not(Equals), "")
		remaining, er := strconv.Atoi(headers.Get("X-RateLimit-Remaining"))
		c.Assert(er, IsNil)
		c.Assert(remaining, Equals, limit-i)
	}

	code, _, headers, err := runRequestWithHeaders(routeName, route, "", func(*http.Request) {}, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, 429)
	remaining, er := strconv.Atoi(headers.Get("X-RateLimit-Remaining"))
	c.Assert(er, IsNil)
	c.Assert(remaining, Equals, 0)
}

func (s *ServerSuite) TestRateLimitProduction(c *C) {
	var (
		account     = CorrectDeploy(1, 1, 1, 1, 0, false, true)[0]
		member      = account.Users[0]
		application = account.Applications[0]
		routeName   = "updateApplication"
		route       = getComposedRoute(routeName, account.PublicID, application.PublicID)

		payload = fmt.Sprintf(
			`{"name":"%s", "description":"i changed the description", "url": "%s", "enabled": true, "in_production": true}`,
			application.Name,
			application.URL,
		)
	)

	code, body, err := runRequest(routeName, route, payload, signAccountRequest(account, member, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusGone)
	c.Assert(body, Not(Equals), "")
}

// createCommonRequestHeaders create a correct request header
func createCommonRequestHeaders(req *http.Request) {
	payload := PeakBody(req).Bytes()

	//req.Header.Add("x-tapglue-date", time.Now().Format(time.RFC3339))
	req.Header.Add("User-Agent", "Tapglue Test UA")

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
}

// getComposedRoute takes a routeName and parameter and returns the route including the version
func getComposedRoute(routeName string, params ...interface{}) string {
	if routeName == "index" {
		return "/"
	} else if routeName == "humans" {
		return "/humans.txt"
	} else if routeName == "robots" {
		return "/robots.txt"
	}

	pattern := getRoute(routeName).TestPattern()

	if len(params) == 0 {
		return pattern
	}

	return fmt.Sprintf(pattern, params...)
}

// getQueryRoute does something
func getQueryRoute(routeName, query string, params ...interface{}) string {
	if routeName == "index" {
		return "/"
	} else if routeName == "humans" {
		return "/humans.txt"
	} else if routeName == "robots" {
		return "/robots.txt"
	}

	pattern := getRoute(routeName).TestPattern()
	pattern += "?" + query

	if len(params) == 0 {
		return pattern
	}

	return fmt.Sprintf(pattern, params...)
}

// getComposedRouteString takes the route and stringyfies all the params
func getComposedRouteString(routeName string, params ...interface{}) string {
	if routeName == "index" {
		return "/"
	} else if routeName == "humans" {
		return "/humans.txt"
	} else if routeName == "robots" {
		return "/robots.txt"
	}

	pattern := getRoute(routeName).TestPattern()

	if len(params) == 0 {
		return pattern
	}

	return fmt.Sprintf(pattern, params...)
}

// getComposedRouteString does something
func getQueryRouteString(routeName, query string, params ...interface{}) string {
	if routeName == "index" {
		return "/"
	} else if routeName == "humans" {
		return "/humans.txt"
	} else if routeName == "robots" {
		return "/robots.txt"
	}

	pattern := getRoute(routeName).TestPattern()

	if len(params) == 0 {
		return pattern
	}

	pattern = strings.Replace(pattern, "%d", "%s", -1)
	pattern = strings.Replace(pattern, "%.f", "%s", -1)
	pattern = strings.Replace(pattern, "%.7f", "%s", -1)
	pattern += "?" + query
	return fmt.Sprintf(pattern, params...)
}

// runRequest takes a route, path, payload and token, performs a request and return a response recorder
func runRequest(routeName, routePath, payload string, signFunc func(*http.Request)) (int, string, errors.Error) {
	var (
		requestRoute *server.Route
		routePattern string
	)

	if routeName == "index" {
		requestRoute = getRoute(routeName)
		routePattern = "/"
	} else if routeName == "humans" {
		requestRoute = getRoute(routeName)
		routePattern = "/humans.txt"
	} else if routeName == "robots" {
		requestRoute = getRoute(routeName)
		routePattern = "/robots.txt"
	} else {
		requestRoute = getRoute(routeName)
		routePattern = requestRoute.RoutePattern()
	}

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	if err != nil {
		panic(err)
	}

	createCommonRequestHeaders(req)

	signFunc(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	m.
		HandleFunc(routePattern, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	body := w.Body.String()

	if *doLogResponses {
		fmt.Printf("Got response: %#v with body %s\n", w, body)
	}

	return w.Code, body, nil
}

// runRequestWithHeaders is like runRequest but with Headerzz!!!
func runRequestWithHeaders(routeName, routePath, payload string, headerz, signFunc func(*http.Request)) (int, string, http.Header, errors.Error) {
	var (
		requestRoute *server.Route
		routePattern string
	)

	if routeName == "index" {
		requestRoute = getRoute(routeName)
		routePattern = "/"
	} else if routeName == "humans" {
		requestRoute = getRoute(routeName)
		routePattern = "/humans.txt"
	} else if routeName == "robots" {
		requestRoute = getRoute(routeName)
		routePattern = "/robots.txt"
	} else {
		requestRoute = getRoute(routeName)
		routePattern = requestRoute.RoutePattern()
	}

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	if err != nil {
		panic(err)
	}

	createCommonRequestHeaders(req)
	signFunc(req)
	headerz(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	m.
		HandleFunc(routePattern, server.CustomHandler(requestRoute, mainLogChan, errorLogChan, "test", false, true)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	body := w.Body.String()

	if *doLogResponses {
		fmt.Printf("Got response: %#v with body %s\n", w, body)
	}

	return w.Code, body, w.Header(), nil
}

// getAccountUserSessionToken retrieves the session token for a certain user
func getAccountUserSessionToken(user *entity.AccountUser) string {
	sessionToken, err := coreAccUser.CreateSession(user)
	if err != nil {
		panic(err)
	}

	return sessionToken
}

func nilSigner(*http.Request) {

}

func signAccountRequest(account *entity.Account, accountUser *entity.AccountUser, goodAccountToken, goodAccountUserToken bool) func(*http.Request) {
	return func(r *http.Request) {
		user := ""
		pass := ""

		if goodAccountToken && account != nil {
			user = account.AuthToken
		}
		if goodAccountToken && account == nil {
			user = ""
		}
		if !goodAccountToken && account != nil {
			user = account.AuthToken + "a"
		}
		if !goodAccountToken && account == nil {
			user = "a"
		}

		if goodAccountUserToken && accountUser != nil {
			pass = accountUser.SessionToken
		}
		if goodAccountUserToken && accountUser == nil {
			pass = ""
		}
		if !goodAccountUserToken && accountUser != nil {
			pass = accountUser.SessionToken + "a"
		}
		if !goodAccountUserToken && accountUser == nil {
			pass = "a"
		}

		if user == "" && pass == "" {
			return
		}

		encodedAuth := Base64Encode(user + ":" + pass)

		r.Header.Add("Authorization", "Basic "+encodedAuth)
	}
}

func signApplicationRequest(application *entity.Application, applicationUser *entity.ApplicationUser, goodApplicationToken, goodApplicationUserToken bool) func(*http.Request) {
	return func(r *http.Request) {
		user := ""
		pass := ""

		if goodApplicationToken && application != nil {
			user = application.AuthToken
		}
		if goodApplicationToken && application == nil {
			user = ""
		}
		if !goodApplicationToken && application != nil {
			user = application.AuthToken + "a"
		}
		if !goodApplicationToken && application == nil {
			user = "a"
		}

		if goodApplicationUserToken && applicationUser != nil {
			pass = applicationUser.SessionToken
		}
		if goodApplicationUserToken && applicationUser == nil {
			pass = ""
		}
		if !goodApplicationUserToken && applicationUser != nil {
			pass = applicationUser.SessionToken + "a"
		}
		if !goodApplicationUserToken && applicationUser == nil {
			pass = "a"
		}

		encodedAuth := Base64Encode(user + ":" + pass)

		r.Header.Add("Authorization", "Basic "+encodedAuth)
	}
}

func getRoute(routeName string) *server.Route {
	routes := server.SetupRoutes()
	for idx := range routes {
		if routes[idx].Name == routeName {
			return routes[idx]
		}
	}

	panic(fmt.Sprintf("route %q not found", routeName))
}
