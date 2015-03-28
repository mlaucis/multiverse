/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"
	. "github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
	"github.com/tapglue/backend/v01/validator/keys"
	"github.com/tapglue/backend/v01/validator/tokens"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/server/utils"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

const (
	apiVersion = "0.1"
)

var (
	_                  = Suite(&ServerSuite{})
	conf               *config.Config
	storageClient      *storage.Client
	doLogTest          = flag.Bool("lt", false, "Set flag in order to get logs output from the tests")
	doCurlLogs         = flag.Bool("ct", false, "Set flag in order to get logs output from the tests as curl requests, sets -lt=true")
	doLogResponseTimes = flag.Bool("rt", false, "Set flag in order to get logs with response times only")
	mainLogChan        = make(chan *logger.LogMsg)
	errorLogChan       = make(chan *logger.LogMsg)
)

// Setup once when the suite starts running
func (s *ServerSuite) SetUpTest(c *C) {
	flag.Parse()

	if *doCurlLogs {
		*doLogTest = true
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")
	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	redis.Client().FlushDb()
	storageClient = storage.Init(redis.Client())
	core.Init(storageClient)
	server.Init()
	validator.Init(storageClient)

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
}

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := server.GetRoute(routeName, apiVersion)
	routePath := requestRoute.RoutePattern(apiVersion)

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set (1)")
}

// Test POST common with CLHeader
func (s *ServerSuite) TestValidatePostCommon_CLHeader(c *C) {
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := server.GetRoute(routeName, apiVersion)
	routePath := requestRoute.RoutePattern(apiVersion)

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
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 failed to create the account (1)\n"+"invalid character 'd' looking for beginning of object key string")
}

// Test GET common with CLHeader
func (s *ServerSuite) TestValidateGetCommon_CLHeader(c *C) {
	payload := ""
	routeName := "index"
	requestRoute := server.GetRoute(routeName, "")
	routePath := requestRoute.RoutePattern("")

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
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusOK)
}

// Test GET common without CLHeader
func (s *ServerSuite) TestValidateGetCommon_NoCLHeader(c *C) {
	payload := ""
	routeName := "index"
	requestRoute := server.GetRoute(routeName, "")
	routePath := requestRoute.RoutePattern("")

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set (1)")
}

// Test PUT common with CLHeader
func (s *ServerSuite) TestValidatePutCommon_CLHeader(c *C) {
	c.Skip("needs a better implementation")

	payload := "{demo}"
	routeName := "updateAccount"
	requestRoute := server.GetRoute(routeName, apiVersion)
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
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
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
	requestRoute := server.GetRoute(routeName, apiVersion)
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
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
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
	requestRoute := server.GetRoute(routeName, apiVersion)
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
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
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
	requestRoute := server.GetRoute(routeName, apiVersion)
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
		HandleFunc(routePath, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test a correct humans request
func (s *ServerSuite) TestHumans_OK(c *C) {
	routeName := "humans"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, "", "", "", 0)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
}

// Test a correct robots request
func (s *ServerSuite) TestRobots_OK(c *C) {
	routeName := "robots"
	route := getComposedRoute(routeName)
	code, body, err := runRequest(routeName, route, "", "", "", 0)
	c.Assert(err, IsNil)

	c.Assert(code, Equals, http.StatusOK)
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

	pattern := server.GetRoute(routeName, apiVersion).ComposePattern(apiVersion)

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

	pattern := server.GetRoute(routeName, apiVersion).ComposePattern(apiVersion)

	if len(params) == 0 {
		return pattern
	}

	pattern = strings.Replace(pattern, "%d", "%s", -1)
	pattern = strings.Replace(pattern, "%.f", "%s", -1)
	pattern = strings.Replace(pattern, "%.7f", "%s", -1)

	return fmt.Sprintf(pattern, params...)
}

// runRequest takes a route, path, payload and token, performs a request and return a response recorder
func runRequest(routeName, routePath, payload, secretKey, sessionToken string, numKeyParts int) (int, string, error) {
	var (
		requestRoute *utils.Route
		routePattern string
	)

	if routeName == "index" {
		requestRoute = server.GetRoute(routeName, "")
		routePattern = "/"
	} else if routeName == "humans" {
		requestRoute = server.GetRoute(routeName, "")
		routePattern = "/humans.txt"
	} else if routeName == "robots" {
		requestRoute = server.GetRoute(routeName, "")
		routePattern = "/robots.txt"
	} else {
		requestRoute = server.GetRoute(routeName, apiVersion)
		routePattern = requestRoute.RoutePattern(apiVersion)
	}

	req, err := http.NewRequest(
		requestRoute.Method,
		routePath,
		strings.NewReader(payload),
	)
	if err != nil {
		panic(err)
	}

	if sessionToken != "" {
		req.Header.Set("x-tapglue-session", sessionToken)
	}

	createCommonRequestHeaders(req)
	if secretKey != "" {
		var err error
		if numKeyParts == 3 {
			err = tokens.SignRequest(secretKey, requestRoute.Scope, apiVersion, numKeyParts, req)
		} else {
			err = keys.SignRequest(secretKey, requestRoute.Scope, apiVersion, numKeyParts, req)
		}

		if err != nil {
			panic(err)
		}
	}

	w := httptest.NewRecorder()
	m := mux.NewRouter()
	m.
		HandleFunc(routePattern, server.CustomHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan, "test", true, false)).
		Methods(requestRoute.Method)
	m.ServeHTTP(w, req)

	return w.Code, w.Body.String(), nil
}

// getAccountUserSessionToken retrieves the session token for a certain user
func getAccountUserSessionToken(user *entity.AccountUser) string {
	sessionToken, err := core.CreateAccountUserSession(user)
	if err != nil {
		panic(err)
	}

	return sessionToken
}

// createApplicationUserSessionToken creates an application user session and returns the token
func createApplicationUserSessionToken(user *entity.User) string {
	sessionToken, err := core.CreateApplicationUserSession(user)
	if err != nil {
		panic(err)
	}

	return sessionToken
}
