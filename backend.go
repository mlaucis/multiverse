/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package main

import (
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"

	"github.com/yvasiyarov/gorelic"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_BACKEND_CONFIG_PATH"
)

var (
	conf          *config.Config
	newRelicAgent *gorelic.Agent
)

func init() {
	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	rand.Seed(time.Now().UTC().UnixNano())

	// Get configuration
	conf = config.NewConf(EnvConfigVar)

	// Initialize components
	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	storageClient := storage.Init(redis.Client())
	core.Init(storageClient)
}

func main() {

	// Setup newrelic
	if conf.Newrelic.Enabled {
		newRelicAgent = gorelic.NewAgent()
		newRelicAgent.Verbose = true
		newRelicAgent.NewrelicLicense = conf.Newrelic.Key
		newRelicAgent.NewrelicName = conf.Newrelic.Name
		newRelicAgent.Run()
	} else {
		newRelicAgent = nil
	}

	// Get router
	router := server.GetRouter(newRelicAgent)

	// Start server
	log.Fatal(http.ListenAndServe(conf.ListenHostPort, router))
}

c.Assert(getAccountUser, DeepEquals, savedAccountUser)
... obtained *entity.AccountUser = &entity.AccountUser{ID:3, AccountID:0, Role:(*entity.AccountRole)(nil), UserCommon:entity.UserCommon{Username:"", Password:"iamsecure..not", DisplayName:"Demo User", FirstName:"", LastName:"", Email:"d@m.o", URL:"", Activated:"", LastLogin:time.Time{sec:0, nsec:0, loc:(*time.Location)(0x465220)}}, Common:entity.Common{Image:[]*entity.Image(nil), Metadata:"", Enabled:false, CreatedAt:time.Time{sec:0, nsec:0, loc:(*time.Location)(0x465220)}, UpdatedAt:time.Time{sec:0, nsec:0, loc:(*time.Location)(0x465220)}, ReceivedAt:0}}
... expected *entity.AccountUser = &entity.AccountUser{ID:3, AccountID:0, Role:(*entity.AccountRole)(nil), UserCommon:entity.UserCommon{Username:"", Password:"iamsecure..not", DisplayName:"Demo User", FirstName:"", LastName:"", Email:"d@m.o", URL:"", Activated:"", LastLogin:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}}, Common:entity.Common{Image:[]*entity.Image(nil), Metadata:"", Enabled:false, CreatedAt:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, UpdatedAt:time.Time{sec:0, nsec:0, loc:(*time.Location)(nil)}, ReceivedAt:0}}
