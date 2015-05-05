/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server_test

import (
	"flag"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/logger"
	. "github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

const apiVersion = "0.2"

var (
	_                  = Suite(&ServerSuite{})
	conf               *config.Config
	doLogTest          = flag.Bool("lt", false, "Set flag in order to get logs output from the tests")
	doCurlLogs         = flag.Bool("ct", false, "Set flag in order to get logs output from the tests as curl requests, sets -lt=true")
	doLogResponseTimes = flag.Bool("rt", false, "Set flag in order to get logs with response times only")
	mainLogChan        = make(chan *logger.LogMsg)
	errorLogChan       = make(chan *logger.LogMsg)

	coreAcc     core.Account
	coreAccUser core.AccountUser
	coreApp     core.Application
	coreAppUser core.ApplicationUser
	coreConn    core.Connection
	coreEvt     core.Event
)

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	payload := "{demo}"
	routeName := "createAccount"
	requestRoute := server.GetRoute(routeName)
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
	requestRoute := server.GetRoute(routeName)
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

// Test PUT common with CLHeader
func (s *ServerSuite) TestValidatePutCommon_CLHeader(c *C) {
	c.Skip("needs a better implementation")

	payload := "{demo}"
	routeName := "updateAccount"
	requestRoute := server.GetRoute(routeName)
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
	requestRoute := server.GetRoute(routeName)
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
	requestRoute := server.GetRoute(routeName)
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
	requestRoute := server.GetRoute(routeName)
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

	pattern := server.GetRoute(routeName).ComposePattern(apiVersion)

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

	pattern := server.GetRoute(routeName).ComposePattern(apiVersion)

	if len(params) == 0 {
		return pattern
	}

	pattern = strings.Replace(pattern, "%d", "%s", -1)
	pattern = strings.Replace(pattern, "%.f", "%s", -1)
	pattern = strings.Replace(pattern, "%.7f", "%s", -1)

	return fmt.Sprintf(pattern, params...)
}

// runRequest takes a route, path, payload and token, performs a request and return a response recorder
func runRequest(routeName, routePath, payload string, signFunc func(*http.Request)) (int, string, errors.Error) {
	var (
		requestRoute *server.Route
		routePattern string
	)

	if routeName == "index" {
		requestRoute = server.GetRoute(routeName)
		routePattern = "/"
	} else if routeName == "humans" {
		requestRoute = server.GetRoute(routeName)
		routePattern = "/humans.txt"
	} else if routeName == "robots" {
		requestRoute = server.GetRoute(routeName)
		routePattern = "/robots.txt"
	} else {
		requestRoute = server.GetRoute(routeName)
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

	createCommonRequestHeaders(req)

	signFunc(req)

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
	sessionToken, err := coreAccUser.CreateSession(user)
	if err != nil {
		panic(err)
	}

	return sessionToken
}

// createApplicationUserSessionToken creates an application user session and returns the token
func createApplicationUserSessionToken(accountID, applicationID int64, user *entity.ApplicationUser) string {
	sessionToken, err := coreAppUser.CreateSession(accountID, applicationID, user)
	if err != nil {
		panic(err)
	}

	return sessionToken
}

func nilSigner(*http.Request) {

}

func signAccountRequest(account *entity.Account, accountUser *entity.AccountUser, goodAccountToken, goodAccountUserToken bool) func(*http.Request) {
	return func(r *http.Request) {

	}
}

func signApplicationRequest(application *entity.Application, applicationUser *entity.ApplicationUser, goodAccountToken, goodAccountUserToken bool) func(*http.Request) {
	return func(r *http.Request) {

	}
}
