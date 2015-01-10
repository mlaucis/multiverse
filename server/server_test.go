/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/db"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type ServerSuite struct{}

var (
	_   = Suite(&ServerSuite{})
	cfg *config.Cfg
)

// Setup once when the suite starts running
func (s *ServerSuite) SetUpTest(c *C) {
	cfg = config.NewConf("")
	db.InitDatabases(cfg.DB())
	db.GetMaster().Ping()
	db.GetSlave().Ping()
	_, err := db.GetMaster().Exec("DELETE FROM `accounts`")
	c.Assert(err, IsNil)
}

// Test POST common without CLHeader
func (s *ServerSuite) TestValidatePostCommon_NoCLHeader(c *C) {
	req, err := http.NewRequest(
		"POST",
		"http://localhost:8089/",
		nil,
	)
	test_CLHeader("", req)
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
	test_CLHeader(payload, req)
	c.Assert(err, IsNil)

	w := httptest.NewRecorder()
	createAccount(w, req)

	c.Assert(w.Code, Equals, http.StatusBadRequest)
	c.Assert(w.Body.String(), Equals, "400 \"invalid character 'd' looking for beginning of object key string\"")
}
