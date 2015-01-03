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
	"github.com/tapglue/backend/db"
	"github.com/tapglue/backend/server"
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
	db.InitDatabases(cfg.DB())
}

func main() {

	// Get router
	router := server.GetRouter()

	// Start server
	log.Fatal(http.ListenAndServe(cfg.ListenHost(), router))
}
