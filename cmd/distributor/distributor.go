// Command distributor will launch a specified consumer for Kinesis and write the received information to its target
//
// Currently it supports:
// - postgres
package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	_ "expvar"
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
	"github.com/tapglue/backend/logger"
	v02_kinesis "github.com/tapglue/backend/v02/storage/kinesis"
	v02_postgres "github.com/tapglue/backend/v02/storage/postgres"
	v02_writer "github.com/tapglue/backend/v02/writer"
	v02_writer_postgres "github.com/tapglue/backend/v02/writer/postgres"
	v03_kinesis "github.com/tapglue/backend/v03/storage/kinesis"
	v03_postgres "github.com/tapglue/backend/v03/storage/postgres"
	v03_redis "github.com/tapglue/backend/v03/storage/redis"
	v03_writer "github.com/tapglue/backend/v03/writer"
	v03_writer_postgres "github.com/tapglue/backend/v03/writer/postgres"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_DISTRIBUTOR_CONFIG_PATH"
)

var (
	startTime       time.Time
	conf            *config.Config
	currentRevision string

	consumerTarget = flag.String("target", "", "Select the target of the consumer to be launched. Currently supported: postgres")
	v02PgConsumer  v02_writer.Writer
	v03PgConsumer  v03_writer.Writer

	mainLogChan  = make(chan *logger.LogMsg, 100000)
	errorLogChan = make(chan *logger.LogMsg, 100000)

	hostname, hostnameErr = os.Hostname()

	pg   v03_postgres.Client
	ksis v03_kinesis.Client
)

func init() {
	if hostnameErr != nil {
		fmt.Println("failed to fetch the hostname")
		panic(hostnameErr)
	}

	startTime = time.Now()

	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	mr.Seed(time.Now().UTC().UnixNano())

	flag.Parse()

	conf = config.NewConf(EnvConfigVar)

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

	errors.Init(conf.Environment != "prod")

	var v02KinesisClient v02_kinesis.Client
	if conf.Environment == "prod" {
		v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
	} else {
		if conf.Kinesis.Endpoint != "" {
			v02KinesisClient = v02_kinesis.NewTest(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Kinesis.Endpoint, conf.Environment, conf.Kinesis.StreamName)
		} else {
			v02KinesisClient = v02_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
		}
	}

	v02PgClient := v02_postgres.New(conf.Postgres)
	v02PgConsumer = v02_writer_postgres.New(v02KinesisClient, v02PgClient)

	var v03KinesisClient v03_kinesis.Client
	if conf.Environment == "prod" {
		v03KinesisClient = v03_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
	} else {
		if conf.Kinesis.Endpoint != "" {
			v03KinesisClient = v03_kinesis.NewWithEndpoint(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Kinesis.Endpoint, conf.Environment, conf.Kinesis.StreamName)
		} else {
			v03KinesisClient = v03_kinesis.New(conf.Kinesis.AuthKey, conf.Kinesis.SecretKey, conf.Kinesis.Region, conf.Environment, conf.Kinesis.StreamName)
		}
	}

	v03PgClient := v03_postgres.New(conf.Postgres)
	v03RedisClient := v03_redis.NewRedigoPool(conf.CacheApp)
	v03PgConsumer = v03_writer_postgres.New(v03KinesisClient, v03PgClient, v03RedisClient)

	pg = v03PgClient
	ksis = v03KinesisClient
}

func main() {
	flag.Parse()

	if *consumerTarget == "" {
		flag.PrintDefaults()
		os.Exit(64)
	}

	if *consumerTarget != "postgres" {
		flag.PrintDefaults()
		os.Exit(64)
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
		Handler:        http.DefaultServeMux,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	if conf.UseSSL {
		server.TLSConfig = configTLS()
	}

	if conf.UseArtwork {
		log.Printf(`

88888888888                         888                        8888888b.  d8b          888            d8b 888               888
    888                             888                        888  "Y88b Y8P          888            Y8P 888               888
    888                             888                        888    888              888                888               888
    888   8888b.  88888b.   .d88b.  888 888  888  .d88b.       888    888 888 .d8888b  888888 888d888 888 88888b.  888  888 888888 .d88b.  888d888
    888      "88b 888 "88b d88P"88b 888 888  888 d8P  Y8b      888    888 888 88K      888    888P"   888 888 "88b 888  888 888   d88""88b 888P"
    888  .d888888 888  888 888  888 888 888  888 88888888      888    888 888 "Y8888b. 888    888     888 888  888 888  888 888   888  888 888
    888  888  888 888 d88P Y88b 888 888 Y88b 888 Y8b.          888  .d88P 888      X88 Y88b.  888     888 888 d88P Y88b 888 Y88b. Y88..88P 888
    888  "Y888888 88888P"   "Y88888 888  "Y88888  "Y8888       8888888P"  888  88888P'  "Y888 888     888 88888P"   "Y88888  "Y888 "Y88P"  888
                  888           888
                  888      Y8b d88P
                  888       "Y88P"

`)
	}

	go func() {
		if conf.UseSSL {
			log.Printf("Starting SSL server at \"%s\" in %s", conf.ListenHostPort, time.Now().Sub(startTime))
			log.Fatal(server.ListenAndServeTLS("./self.crt", "./self.key"))
		} else {
			log.Printf("Starting NORMAL server at \"%s\" in %s", conf.ListenHostPort, time.Now().Sub(startTime))
			log.Fatal(server.ListenAndServe())
		}
	}()

	for {
		execute(conf.Kinesis.StreamName, mainLogChan, errorLogChan)
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
	TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
	TLSConfig.PreferServerCipherSuites = true
	TLSConfig.ClientSessionCache = tls.NewLRUClientSessionCache(1000)
	//TLSConfig.RootCAs = loadCertificates()
	TLSConfig.ClientCAs = loadClientCertificates()

	return TLSConfig
}

func loadCertificates() *x509.CertPool {
	pem, err := ioutil.ReadFile("./root-ca.pem")
	if err != nil {
		panic(err)
	}

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return rootCertPool
}

func loadClientCertificates() *x509.CertPool {
	pem, err := ioutil.ReadFile("./origin-pull-ca.pem")
	if err != nil {
		panic(err)
	}

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return rootCertPool
}
