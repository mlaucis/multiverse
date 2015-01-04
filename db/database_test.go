/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package db provides configuration for database connection
package db

import (
	"testing"

	"github.com/tapglue/backend/config"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type (
	DatabaseSuite struct{}
)

var (
	_   = Suite(&DatabaseSuite{})
	cfg *config.Cfg
)

func (dbs *DatabaseSuite) SetUpTest(c *C) {
	cfg = config.NewConf("")
}

// Test InitDatabases
func (dbs *DatabaseSuite) TestInitDatabases(c *C) {
	// This will panic if things don't work
	InitDatabases(cfg.DB())
}

// Test GetMaster
func (dbs *DatabaseSuite) TestGetMaster(c *C) {
	// This will panic if things don't work
	InitDatabases(cfg.DB())
	GetMaster().Ping()
}

// Test GetSlave
func (dbs *DatabaseSuite) TestGetSlave(c *C) {
	// This will panic if things don't work
	InitDatabases(cfg.DB())
	GetSlave().Ping()
}
