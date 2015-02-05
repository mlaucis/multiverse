/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strconv"
	"strings"
	"testing"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"

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
	logChan       = make(chan *LogMsg)
	storageClient *storage.Client
)

// Setup once when the suite starts running
func (s *ServerSuite) SetUpTest(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")
	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	redis.Client().FlushDb()
	storageClient = storage.Init(redis.Client())
	core.Init(storageClient)

	go TGLog(logChan)
}

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	req, err := http.NewRequest(
		"POST",
		getComposedRoute("index"),
		nil,
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test POST common with CLHeader
func (s *ServerSuite) TestValidatePostCommon_CLHeader(c *C) {
	payload := "{demo}"
	req, err := http.NewRequest(
		"POST",
		getComposedRoute("index"),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 invalid character 'd' looking for beginning of object key string")
}

// Test GET common with CLHeader
func (s *ServerSuite) TestValidateGetCommon_CLHeader(c *C) {
	req, err := http.NewRequest(
		"GET",
		getComposedRoute("getAccount", 100),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	getAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test GET common without CLHeader
func (s *ServerSuite) TestValidateGetCommon_NoCLHeader(c *C) {
	req, err := http.NewRequest(
		"GET",
		getComposedRoute("getAccount", 100),
		nil,
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	getAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test PUT common with CLHeader
func (s *ServerSuite) TestValidatePutCommon_CLHeader(c *C) {
	payload := "{demo}"
	req, err := http.NewRequest(
		"PUT",
		getComposedRoute("updateAccount", 100),
		strings.NewReader(payload),
	)
	c.Assert(err, IsNil)

	clHeader(payload, req)

	w := httptest.NewRecorder()
	updateAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test PUT common without CLHeader
func (s *ServerSuite) TestValidatePutCommon_NoCLHeader(c *C) {
	req, err := http.NewRequest(
		"PUT",
		getComposedRoute("updateAccount", 100),
		nil,
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	updateAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// Test DELETE common with CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_CLHeader(c *C) {
	req, err := http.NewRequest(
		"DELETE",
		getComposedRoute("deleteAccount", 100),
		nil,
	)
	c.Assert(err, IsNil)

	clHeader("", req)

	w := httptest.NewRecorder()
	deleteAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 accountId is not set or the value is incorrect")
}

// Test DELETE common without CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_NoCLHeader(c *C) {
	req, err := http.NewRequest(
		"DELETE",
		getComposedRoute("deleteAccount", 100),
		nil,
	)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	deleteAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 User-Agent header must be set")
}

// clHeader create a correct request header
func clHeader(payload string, req *http.Request) {
	req.Header.Add("User-Agent", "go test (+localhost)")
	if len(payload) > 0 {
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
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version%s\n", routeName, apiVersion))
	}

	return fmt.Sprintf(routes[apiVersion][routeName].composePattern(apiVersion), params...)
}

// runRequest takes a route, path, payload and token, performs a request and return a response recorder
func runRequest(routeName, routePath, payload, token string) (*httptest.ResponseRecorder, error) {
	route := getRoute(routeName)

	req, err := http.NewRequest(
		route.method,
		routePath,
		strings.NewReader(payload),
	)
	if err != nil {
		return nil, err
	}

	clHeader(payload, req)
	if token != "" {
		signRequest(token, req)
	}

	w := httptest.NewRecorder()
	m := mux.NewRouter()

	m.HandleFunc(route.routePattern(apiVersion), customHandler(routeName, route, nil, logChan)).Methods(route.method)
	m.ServeHTTP(w, req)

	return w, nil
}
