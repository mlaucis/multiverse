// +build !bench

package server_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"

	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/server"

	"github.com/gorilla/mux"
	. "gopkg.in/check.v1"
)

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	payload := "{demo}"
	routeName := "createOrganization"
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
	payload := "{demo}"
	routeName := "createOrganization"
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
	routeName := "updateOrganization"
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
	c.Assert(w.Body.String(), Equals, "400 organizationId is not set or the value is incorrect")
}

// Test PUT common without CLHeader
func (s *ServerSuite) TestValidatePutCommon_NoCLHeader(c *C) {
	c.Skip("this needs a better implementation now that contexts are in place")

	payload := "{demo}"
	routeName := "updateOrganization"
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
	routeName := "deleteOrganization"
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
	c.Assert(w.Body.String(), Equals, "400 organizationId is not set or the value is incorrect")
}

// Test DELETE common without CLHeader
func (s *ServerSuite) TestValidateDeleteCommon_NoCLHeader(c *C) {
	c.Skip("skip due to context refactoring")
	payload := ""
	routeName := "deleteOrganization"
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
	organizations := CorrectDeploy(1, 0, 1, 1, 0, false, true)
	application := organizations[0].Applications[0]
	user := application.Users[0]

	routeName := "getCurrentApplicationUser"
	route := getComposedRoute(routeName)
	code, body, headers, err := runRequestWithHeaders(routeName, route, "", func(*http.Request) {}, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, http.StatusOK)
	c.Assert(body, Not(Equals), "")
	remaining, er := strconv.Atoi(headers.Get("X-RateLimit-Remaining"))
	c.Assert(er, IsNil)
	c.Assert(remaining, Equals, 999)

	receivedUser := &entity.ApplicationUser{}
	er = json.Unmarshal([]byte(body), receivedUser)
	c.Assert(er, IsNil)
	c.Assert(receivedUser.Username, Equals, user.Username)

	for i := 2; i <= 1000; i++ {
		code, body, _, err = runRequestWithHeaders(routeName, route, "", func(*http.Request) {}, signApplicationRequest(application, user, true, true))
		c.Assert(err, IsNil)
		c.Assert(code, Equals, http.StatusOK)
		c.Assert(body, Not(Equals), "")
	}

	code, body, headers, err = runRequestWithHeaders(routeName, route, "", func(*http.Request) {}, signApplicationRequest(application, user, true, true))
	c.Assert(err, IsNil)
	c.Assert(code, Equals, 429)
	remaining, er = strconv.Atoi(headers.Get("X-RateLimit-Remaining"))
	c.Assert(er, IsNil)
	c.Assert(remaining, Equals, 0)
}
