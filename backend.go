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
