/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package main

import (
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"
	"github.com/tapglue/backend/validator"

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

	conf = config.NewConf(EnvConfigVar)

	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	storageClient := storage.Init(redis.Client())
	core.Init(storageClient)
	validator.Init(storageClient)
}

func main() {
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
	router := server.GetRouter(conf.Environment != "prod", newRelicAgent)

	log.Printf("Starting the server at port %s", conf.ListenHostPort)
	log.Fatal(http.ListenAndServe(conf.ListenHostPort, router))
}
