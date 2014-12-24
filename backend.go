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
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	cfg = config.GetConfig(EnvConfigVar)

	db.InitDatabases(cfg)
}

func main() {
	//http.Handle("/", server.GetRouter())

	router := server.GetRouter()

	log.Fatal(http.ListenAndServe(cfg.ListenHost, router))
}
