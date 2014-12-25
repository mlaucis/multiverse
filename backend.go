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

	"github.com/gluee/backend/config"
	"github.com/gluee/backend/db"
	"github.com/gluee/backend/server"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "GLUEE_BACKEND_CONFIG_PATH"
)

var cfg *config.Config

func init() {
	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	// Get configuration
	cfg = config.GetConfig(EnvConfigVar)

	// Initialize database
	db.InitDatabases(cfg)
}

func main() {

	// Get router
	router := server.GetRouter()

	// Start server
	log.Fatal(http.ListenAndServe(cfg.ListenHost, router))
}
