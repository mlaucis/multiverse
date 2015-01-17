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

	"github.com/tapglue/backend/aerospike"
	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/server"
	"github.com/yvasiyarov/gorelic"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_BACKEND_CONFIG_PATH"
)

var cfg *config.Cfg

func init() {
	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	// Get configuration
	cfg = config.NewConf(EnvConfigVar)

	// Initialize database
	aerospike.InitAerospike(cfg.Aerospike())
}

func main() {

	newRelicAgent := gorelic.NewAgent()
	newRelicAgent.Verbose = true
	newRelicAgent.NewrelicLicense, newRelicAgent.NewrelicName = cfg.NewRelic()
	newRelicAgent.Run()

	// Get router
	router := server.GetRouter(newRelicAgent)

	// Start server
	log.Fatal(http.ListenAndServe(cfg.ListenHost(), router))
}
