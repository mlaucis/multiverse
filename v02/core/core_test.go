/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"runtime"
	"testing"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/storage/redis"
	"github.com/tapglue/backend/v02/storage"

	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type (
	CoreSuite struct{}
)

var (
	_    = Suite(&CoreSuite{})
	conf *config.Config
)

func (ass *CoreSuite) SetUpSuite(c *C) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	conf = config.NewConf("")
	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	redis.Client().FlushDb()
	storageClient := storage.Init(redis.Client())
	Init(storageClient)
}
