/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Command backend is the heavy lifting part of the tapglue backend
package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"io/ioutil"
	"log"
	mr "math/rand"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/server"
	"github.com/tapglue/backend/tgerrors"
	v01_core "github.com/tapglue/backend/v01/core"
	v01_storage "github.com/tapglue/backend/v01/storage"
	v01_redis "github.com/tapglue/backend/v01/storage/redis"
	v01_validator "github.com/tapglue/backend/v01/validator"
	v02_core "github.com/tapglue/backend/v02/core"
	v02_redis_core "github.com/tapglue/backend/v02/core/redis"
	v02_server "github.com/tapglue/backend/v02/server"
	v02_storage "github.com/tapglue/backend/v02/storage"
	v02_kinesis "github.com/tapglue/backend/v02/storage/kinesis"
	v02_redis "github.com/tapglue/backend/v02/storage/redis"
	v02_validator "github.com/tapglue/backend/v02/validator"
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
	mr.Seed(time.Now().UTC().UnixNano())

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

	tgerrors.Init(conf.Environment != "prod")
	v01_redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	v02_redis.Init(conf.Redis.Hosts[0], conf.Redis.Password, conf.Redis.DB, conf.Redis.PoolSize)
	v02_kinesis.Init(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region)
	v01StorageClient := v01_storage.Init(v01_redis.Client())
	v02StorageClient := v02_storage.Init(v02_redis.Client(), v02_kinesis.Client())

	account := v02_redis_core.NewAccount(v02StorageClient, v02_redis.Client())
	accountUser := v02_redis_core.NewAccountUser(v02StorageClient, v02_redis.Client())
	application := v02_redis_core.NewApplication(v02StorageClient, v02_redis.Client())
	applicationUser := v02_redis_core.NewApplicationUser(v02StorageClient, v02_redis.Client())
	connection := v02_redis_core.NewConnection(v02StorageClient, v02_redis.Client())
	event := v02_redis_core.NewEvent(v02StorageClient, v02_redis.Client())

	v02_server.InitCores(account, accountUser, application, applicationUser, connection, event)

	v01_core.Init(v01StorageClient)
	v02_core.Init(v02StorageClient)
	v01_validator.Init(v01StorageClient)
	v02_validator.Init(v02StorageClient, account, accountUser, application, applicationUser)
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
	TLSConfig.CipherSuites = []uint16{
		tls.TLS_FALLBACK_SCSV,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_ECDHE_RSA_WITH_3DES_EDE_CBC_SHA,
		tls.TLS_RSA_WITH_AES_256_CBC_SHA,
		tls.TLS_RSA_WITH_AES_128_CBC_SHA,
		tls.TLS_RSA_WITH_3DES_EDE_CBC_SHA,
	}

	TLSConfig.Rand = rand.Reader
	TLSConfig.MinVersion = tls.VersionTLS10
	TLSConfig.SessionTicketsDisabled = false
	TLSConfig.InsecureSkipVerify = false
	TLSConfig.ClientAuth = tls.VerifyClientCertIfGiven
	TLSConfig.PreferServerCipherSuites = true
	TLSConfig.ClientSessionCache = tls.NewLRUClientSessionCache(1000)
	TLSConfig.RootCAs = loadCertificates()

	return TLSConfig
}

func loadCertificates() *x509.CertPool {
	pem, err := ioutil.ReadFile("./cert/STAR_tapglue_com.ca-bundle")
	if err != nil {
		panic(err)
	}

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return rootCertPool
}
