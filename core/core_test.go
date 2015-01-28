/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"runtime"
	"testing"

	"github.com/tapglue/backend/config"

	. "gopkg.in/check.v1"
	"github.com/tapglue/backend/redis")

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type (
	CoreSuite struct{}
)

var (
	_   = Suite(&CoreSuite{})
	cfg *config.Cfg
)

func (ass *CoreSuite) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg = config.NewConf("")
	_ = cfg
	redis.Init()
}
