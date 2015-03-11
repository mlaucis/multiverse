/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"
	. "github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/validator"
	"github.com/tapglue/backend/validator/keys"
	"github.com/tapglue/backend/validator/tokens"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

const (
	apiVersion = "0.1"
)

var (
	_             = Suite(&ServerSuite{})
	conf          *config.Config
	storageClient *storage.Client
	doLogTest     = flag.Bool("lt", false, "Set flag in order to get logs output from the tests")
	doCurlLogs    = flag.Bool("ct", false, "Set flag in order to get logs output from the tests as curl requests, sets -lt=true")
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
	validator.Init(storageClient)

	if *doLogTest {
		// overwrite log channel to make it blocking
		mainLogChan = make(chan *LogMsg)
		errorLogChan = make(chan *LogMsg)
		if *doCurlLogs {
			go TGCurlLog(mainLogChan)
		} else {
			go TGLog(mainLogChan)
		}
		go TGLog(errorLogChan)
	} else {
		go tgSilentLog(mainLogChan)
		go tgSilentLog(errorLogChan)
	}
}

func tgSilentLog(msg chan *LogMsg) {
	for {
		select {
		case m := <-msg:
			{
				_ = m
			}
		}
	}
}

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := getRoute(routeName)
	routePath := requestRoute.routePattern(apiVersion)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test POST common with CLHeader
func (s *ServerSuite) TestValidatePostCommon_CLHeader(c *C) {
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := getRoute(routeName)
	routePath := requestRoute.routePattern(apiVersion)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 invalid character 'd' looking for beginning of object key string")
}

// Test GET common with CLHeader
func (s *ServerSuite) TestValidateGetCommon_CLHeader(c *C) {
	payload := ""
	routeName := "index"
	requestRoute := getRoute(routeName)
	routePath := requestRoute.routePattern(apiVersion)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusOK)
}

// Test GET common without CLHeader
func (s *ServerSuite) TestValidateGetCommon_NoCLHeader(c *C) {
	payload := ""
	routeName := "index"
	requestRoute := getRoute(routeName)
	routePath := requestRoute.routePattern(apiVersion)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test PUT common with CLHeader
func (s *ServerSuite) TestValidatePutCommon_CLHeader(c *C) {
	c.Skip("needs a better implementation")
	return

	payload := "{demo}"
	routeName := "updateAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test PUT common without CLHeader
func (s *ServerSuite) TestValidatePutCommon_NoCLHeader(c *C) {
	payload := "{demo}"
	routeName := "updateAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test DELETE common with CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_CLHeader(c *C) {
	c.Skip("needs a better implementation")
	return

	payload := "{demo}"
	routeName := "deleteAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	createCommonRequestHeaders(req)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test DELETE common without CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_NoCLHeader(c *C) {
	payload := ""
	routeName := "deleteAccount"
	requestRoute := getRoute(routeName)
	routePath := getComposedRoute(routeName, 0)

	req, err := http.NewRequest(
		requestRoute.method,
		routePath,
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(routePath, customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
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
	req.Header.Add("x-tapglue-date", time.Now().Format(time.RFC3339))
	req.Header.Add("User-Agent", "go test (+localhost)")
	payload := PeakBody(req).Bytes()
	if len(payload) > 0 {
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
	}
}

// getRoute takes a route name and returns the route including the version
func getRoute(routeName string) *route {
	if _, ok := routes[apiVersion][routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version%s\n", routeName, apiVersion))
	}

	return routes[apiVersion][routeName]
}

// getComposedRoute takes a routeName and parameter and returns the route including the version
func getComposedRoute(routeName string, params ...interface{}) string {
	if _, ok := routes[apiVersion][routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version %s\n", routeName, apiVersion))
	}

	pattern := routes[apiVersion][routeName].composePattern(apiVersion)

	if len(params) == 0 {
		return pattern
	}

	return fmt.Sprintf(pattern, params...)
}

// runRequest takes a route, path, payload and token, performs a request and return a response recorder
func runRequest(routeName, routePath, payload, secretKey, sessionToken string, numKeyParts int) (int, string, error) {
	requestRoute := getRoute(routeName)

	req, err := http.NewRequest(
		requestRoute.method,
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
		var err error = nil
		if numKeyParts == 3 {
			err = tokens.SignRequest(secretKey, requestRoute.scope, apiVersion, numKeyParts, req)
		} else {
			err = keys.SignRequest(secretKey, requestRoute.scope, apiVersion, numKeyParts, req)
		}

		if err != nil {
			panic(err)
		}
	}

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.
		HandleFunc(requestRoute.routePattern(apiVersion), customHandler(routeName, apiVersion, requestRoute, mainLogChan, errorLogChan)).
		Methods(requestRoute.method)
	m.ServeHTTP(w, req)

	return w.Code, w.Body.String(), nil
}

func getAccountUserSessionToken(user *entity.AccountUser) string {
	sessionToken, err := core.CreateAccountUserSession(user)
	if err != nil {
		panic(err)
	}

	return sessionToken
}

func getApplicationUserSessionToken(user *entity.User) string {
	sessionToken, err := core.CreateApplicationUserSession(correctUser)
	if err != nil {
		panic(err)
	}

	return sessionToken
}
