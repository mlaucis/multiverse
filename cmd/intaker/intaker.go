/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"log/syslog"
	mr "math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	"github.com/tapglue/backend/config"
	"github.com/tapglue/backend/errors"
	ratelimiter_redis "github.com/tapglue/backend/limiter/redis"
	"github.com/tapglue/backend/logger"
	"github.com/tapglue/backend/server"
	v02_kinesis_core "github.com/tapglue/backend/v02/core/kinesis"
	v02_postgres_core "github.com/tapglue/backend/v02/core/postgres"
	v02_kinesis "github.com/tapglue/backend/v02/storage/kinesis"
	v02_postgres "github.com/tapglue/backend/v02/storage/postgres"
	v02_redis "github.com/tapglue/backend/v02/storage/redis"

	redigo "github.com/garyburd/redigo/redis"
	"github.com/yvasiyarov/gorelic"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_INTAKER_CONFIG_PATH"
)

var (
	conf                *config.Config
	startTime           time.Time
	redigoRateLimitPool *redigo.Pool
	currentRevision     string
	forceNoSec          = flag.Bool("force-no-sec", false, "Force no sec enables launching the backend in production without security checks")
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

	log.SetFlags(0)

	if conf.UseSysLog {
		syslogWriter, err := syslog.New(syslog.LOG_INFO, "intaker")
		if err == nil {
			log.Printf("logging to syslog is enabled. Please tail your syslog for intaker app for further logs\n")
			log.SetOutput(syslogWriter)
		} else {
			log.Printf("%v\n", err)
			log.Printf("logging to syslog failed reverting to stdout logging\n")
		}
		conf.UseArtwork = false
	}

	if conf.SkipSecurity {
		log.Printf("launching with no security checks enabled\n")
	}

	errors.Init(true)

	var v02KinesisClient v02_kinesis.Client
	if conf.Environment == "prod" {
		v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment)
	} else {
		if conf.Kinesis.Endpoint != "" {
			v02KinesisClient = v02_kinesis.NewTest(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Kinesis.Endpoint, conf.Environment)
		} else {
			v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment)
		}
	}

	switch conf.Environment {
	case "dev":
		v02KinesisClient.SetupStreams([]string{v02_kinesis.PackedStreamNameDev})
	case "test":
		v02KinesisClient.SetupStreams([]string{v02_kinesis.PackedStreamNameTest})
	case "prod":
		v02KinesisClient.SetupStreams([]string{v02_kinesis.PackedStreamNameProduction})
	}

	v02PostgresClient := v02_postgres.New(conf.Postgres)

	redigoRateLimitPool = v02_redis.NewRedigoPool(conf.Redis.Hosts[0], "")

	applicationRateLimiter := ratelimiter_redis.NewLimiter(redigoRateLimitPool, "ratelimiter.app.")

	kinesisAccount := v02_kinesis_core.NewAccount(v02KinesisClient)
	kinesisAccountUser := v02_kinesis_core.NewAccountUser(v02KinesisClient)
	kinesisApplication := v02_kinesis_core.NewApplication(v02KinesisClient)
	kinesisApplicationUser := v02_kinesis_core.NewApplicationUser(v02KinesisClient)
	kinesisConnection := v02_kinesis_core.NewConnection(v02KinesisClient)
	kinesisEvent := v02_kinesis_core.NewEvent(v02KinesisClient)

	postgresAccount := v02_postgres_core.NewAccount(v02PostgresClient)
	postgresAccountUser := v02_postgres_core.NewAccountUser(v02PostgresClient)
	postgresApplication := v02_postgres_core.NewApplication(v02PostgresClient)
	postgresApplicationUser := v02_postgres_core.NewApplicationUser(v02PostgresClient)
	postgresConnection := v02_postgres_core.NewConnection(v02PostgresClient)
	postgresEvent := v02_postgres_core.NewEvent(v02PostgresClient)

	server.SetupRawConnections(v02KinesisClient, v02PostgresClient, redigoRateLimitPool)
	server.SetupRateLimit(applicationRateLimiter)
	server.SetupKinesisCores(kinesisAccount, kinesisAccountUser, kinesisApplication, kinesisApplicationUser, kinesisConnection, kinesisEvent)
	server.SetupPostgresCores(postgresAccount, postgresAccountUser, postgresApplication, postgresApplicationUser, postgresConnection, postgresEvent)
	server.SetupFlakes()

	currentHostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve the current hostname. Error: %q", err))
	}
	if currentHostname == "" {
		panic("hostname is empty")
	}

	server.Setup(currentRevision, currentHostname)
}

func main() {
	agent := gorelic.NewAgent()
	agent.NewrelicName = "Intaker"
	agent.NewrelicLicense = "24f345545c02907b32909bd9f818b29c63bbc5c1"

	// Get router
	router, mainLogChan, errorLogChan, err := server.GetRouter(agent, conf.Environment, conf.Environment != "prod", conf.SkipSecurity)
	if err != nil {
		panic(err)
	}

	if conf.Environment == "prod" {
		agent.Run()
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

88888888888                         888                        8888888          888             888
    888                             888                          888            888             888
    888                             888                          888            888             888
    888   8888b.  88888b.   .d88b.  888 888  888  .d88b.         888   88888b.  888888  8888b.  888  888  .d88b.  888d888
    888      "88b 888 "88b d88P"88b 888 888  888 d8P  Y8b        888   888 "88b 888        "88b 888 .88P d8P  Y8b 888P"
    888  .d888888 888  888 888  888 888 888  888 88888888        888   888  888 888    .d888888 888888K  88888888 888
    888  888  888 888 d88P Y88b 888 888 Y88b 888 Y8b.            888   888  888 Y88b.  888  888 888 "88b Y8b.     888
    888  "Y888888 88888P"   "Y88888 888  "Y88888  "Y8888       8888888 888  888  "Y888 "Y888888 888  888  "Y8888  888
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
