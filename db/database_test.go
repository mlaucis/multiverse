/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

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

func (dbs *DatabaseSuite) SetUpSuite(c *C) {
	cfg = config.NewConf("")
	InitDatabases(cfg.DB())

	GetMaster().Ping()
	GetSlave().Ping()

	_, err := GetMaster().Exec("DELETE FROM `accounts`")
	c.Assert(err, IsNil)
}
