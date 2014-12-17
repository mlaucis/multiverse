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
	CONFIG_PATH = "GLUEMOBILE_BACKEND_CONFIG_PATH"
)

var cfg *config.Config

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	rand.Seed(time.Now().UTC().UnixNano())
	cfg = config.GetConfig(CONFIG_PATH)
}

func main() {
	http.Handle("/", server.GetRouter())

	http.ListenAndServe(cfg.ListenHost, nil)
}
