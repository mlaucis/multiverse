/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package redis

import (
	"runtime"
	"testing"

	"github.com/tapglue/backend/config"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type (
	AerospikeSuite struct{}
)

var (
	_   = Suite(&AerospikeSuite{})
	cfg *config.Config
)

func (ass *AerospikeSuite) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cfg = config.NewConf("")
	InitAerospike(cfg.Aerospike())
}
