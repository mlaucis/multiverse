/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Command backend is the heavy lifting part of the tapglue backend
package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"flag"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"
	"github.com/tapglue/backend/validator"
	"github.com/tapglue/backend/worker/channel"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_BACKEND_CONFIG_PATH"
)

var (
	conf       *config.Config
	forceNoSec = flag.Bool("force-no-sec", false, "Force no sec enables launching the backend in production without security checks")
)

func init() {
	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	rand.Seed(time.Now().UTC().UnixNano())

	flag.Parse()

	conf = config.NewConf(EnvConfigVar)

	if conf.SkipSecurity && conf.Environment == "prod" {
		if !*forceNoSec {
			panic("attempted to launch in production with no security checks enabled")
		}
	}

	if conf.SkipSecurity {
		log.Printf("launching with no security checks enabled\n")
	}

	redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	storageClient := storage.Init(redis.Client())
	core.Init(storageClient)
	validator.Init(storageClient)

	queue := channel.NewQueue()
	worker := channel.NewWorker(queue)
	_ = worker
}

func main() {
	// Get router
	router, mainLogChan, errorLogChan, err := server.GetRouter(conf.Environment != "prod", conf.SkipSecurity)
	if err != nil {
		panic(err)
	}
	go server.TGLog(mainLogChan)
	go server.TGLog(errorLogChan)

	if conf.UseArtwork {
		log.Printf(`

88888888888                         888                          .d8888b.
    888                             888                         d88P  Y88b
    888                             888                         Y88b.
    888   8888b.  88888b.   .d88b.  888 888  888  .d88b.         "Y888b.    .d88b.  888d888 888  888  .d88b.  888d888
    888      "88b 888 "88b d88P"88b 888 888  888 d8P  Y8b           "Y88b. d8P  Y8b 888P"   888  888 d8P  Y8b 888P"
    888  .d888888 888  888 888  888 888 888  888 88888888             "888 88888888 888     Y88  88P 88888888 888
    888  888  888 888 d88P Y88b 888 888 Y88b 888 Y8b.           Y88b  d88P Y8b.     888      Y8bd8P  Y8b.     888
    888  "Y888888 88888P"   "Y88888 888  "Y88888  "Y8888         "Y8888P"   "Y8888  888       Y88P    "Y8888  888
                  888           888
                  888      Y8b d88P
                  888       "Y88P"

  	`)
	}

	// Get IP Address
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
	}

	var localIP string
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				localIP = ipnet.IP.String()
			}
		}
	}

	log.Printf("Starting the server at %s%s", localIP, conf.ListenHostPort)
	log.Fatal(http.ListenAndServe(conf.ListenHostPort, router))
}
