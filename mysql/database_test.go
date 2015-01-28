/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package mysql

import (
	"runtime"
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
	cfg *config.Config
)

func (dbs *DatabaseSuite) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg = config.NewConf("")
	InitDatabases(cfg.DB())

	GetMaster().Ping()
	GetSlave().Ping()

	_, err := GetMaster().Exec("DELETE FROM `accounts`")
	c.Assert(err, IsNil)
}
