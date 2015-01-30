/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/storage/redis"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

var (
	_    = Suite(&ServerSuite{})
	conf *config.Config
)

// Setup once when the suite starts running
func (s *ServerSuite) SetUpTest(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")
	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
}

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	req, err := http.NewRequest(
		"POST",
		"http://localhost:8089/",
		nil,
	)
	clHeader("", req)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 \"invalid Content-Length size\"")
}

// Test POST common with CLHeader
func (s *ServerSuite) TestValidatePostCommon_CLHeader(c *C) {
	payload := "{demo}"
	req, err := http.NewRequest(
		"POST",
		"http://localhost:8089/",
		strings.NewReader(payload),
	)
	clHeader(payload, req)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 \"invalid character 'd' looking for beginning of object key string\"")
}
