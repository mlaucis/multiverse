/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Command backend is the heavy lifting part of the tapglue backend
package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/storage"
	"github.com/tapglue/backend/storage/redis"
	"github.com/tapglue/backend/validator"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_BACKEND_CONFIG_PATH"
)

var (
	conf       *config.Config
	startTime  time.Time
	forceNoSec = flag.Bool("force-no-sec", false, "Force no sec enables launching the backend in production without security checks")
)

func init() {
	startTime = time.Now()

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
}

func main() {
	// Get router
	router, mainLogChan, errorLogChan, err := server.GetRouter(conf.Environment, conf.Environment != "prod", conf.SkipSecurity)
	if err != nil {
		panic(err)
	}

	if conf.JSONLogs {
		go logger.JSONLog(mainLogChan)
		go logger.JSONLog(errorLogChan)
	} else {
		go logger.TGLog(mainLogChan)
		go logger.TGLog(errorLogChan)
	}


	server := &http.Server{
		Addr:           conf.ListenHostPort,
		Handler:        router,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if conf.UseSSL {
		server.TLSConfig = configTLS()
	}

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
	if conf.UseSSL {
		log.Printf("Starting SSL server at \"%s\" in %s", conf.ListenHostPort, time.Now().Sub(startTime))
		log.Fatal(server.ListenAndServeTLS("./cert/STAR_tapglue_com.pem", "./cert/STAR_tapglue_com.key"))
	} else {
		log.Printf("Starting NORMAL server at \"%s\" in %s", conf.ListenHostPort, time.Now().Sub(startTime))
		log.Fatal(server.ListenAndServe())
	}
}

func configTLS() *tls.Config {
	TLSConfig := &tls.Config{}
	TLSConfig.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}
	TLSConfig.MinVersion = tls.VersionTLS12
	TLSConfig.SessionTicketsDisabled = true
	TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	/*TLSConfig.Certificates = make([]tls.Certificate, 1)
	TLSConfig.Certificates[0], TLSConfig.RootCAs, TLSConfig.ClientCAs = loadCertificates()*/

	return TLSConfig
}

func loadCertificates() (tls.Certificate, *x509.CertPool, *x509.CertPool) {
	mycert, err := tls.LoadX509KeyPair("./cert/STAR_tapglue_com.crt", "./cert/STAR_taplue_com.key")
	if err != nil {
		panic(err)
	}

	pem, err := ioutil.ReadFile("./cert/STAR_tapglue_com.ca-bundle")
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return mycert, certPool, certPool
}
