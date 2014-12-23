/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package main

import (
	"math/rand"
	"net/http"
	"runtime"
	"time"

	"github.com/gluee/backend/config"
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
}

func main() {
	http.Handle("/", server.GetRouter())

	http.ListenAndServe(cfg.ListenHost, nil)
}
