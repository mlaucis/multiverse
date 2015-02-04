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

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

const (
	apiVersion = "0.1"
)

var (
	_       = Suite(&ServerSuite{})
	conf    *config.Config
	logChan = make(chan *LogMsg)
)

// Setup once when the suite starts running
func (s *ServerSuite) SetUpTest(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")
	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	redis.Client().FlushDb()
	storageClient := storage.Init(redis.Client())
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

	clHeader("", req)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 invalid Content-Length size")
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

func clHeader(payload string, req *http.Request) {
	req.Header.Add("User-Agent", "go test (+localhost)")
	if len(payload) > 0 {
		req.Header.Add("Content-Length", strconv.FormatInt(int64(len(payload)), 10))
	}
}

func getRoute(routeName string) *route {
	if _, ok := routes[apiVersion][routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version%s\n", routeName, apiVersion))
	}

	return routes[apiVersion][routeName]
}

func getComposedRoute(routeName string, params ...interface{}) string {
	if _, ok := routes[apiVersion][routeName]; !ok {
		panic(fmt.Errorf("You requested a route, %s, that does not exists in the routing table for version%s\n", routeName, apiVersion))
	}

	return fmt.Sprintf(routes[apiVersion][routeName].composePattern(apiVersion), params...)
}
