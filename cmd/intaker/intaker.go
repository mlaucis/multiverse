// Command intaker will launch a specified frontend for Tapglue
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

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/logger"
	"github.com/tapglue/multiverse/server"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar = "TAPGLUE_INTAKER_CONFIG_PATH"
)

var (
	conf            *config.Config
	startTime       time.Time
	currentRevision string
	forceNoSec      = flag.Bool("force-no-sec", false, "Force no sec enables launching the backend in production without security checks")
)

func init() {
	startTime = time.Now()

	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	mr.Seed(time.Now().UTC().UnixNano())

}

func main() {
	flag.Parse()

	conf = config.NewConf(EnvConfigVar)

	if conf.SkipSecurity && conf.Environment == "prod" {
		if !*forceNoSec {
			panic("attempted to launch in production with no security checks enabled")
		}
	}

	// rsyslog will create the timestamps for us
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

	currentHostname, err := os.Hostname()
	if err != nil {
		panic(fmt.Sprintf("failed to retrieve the current hostname. Error: %q", err))
	}
	if currentHostname == "" {
		panic("hostname is empty")
	}

	server.Setup(conf, currentRevision, currentHostname)

	// Get router
	router, mainLogChan, errorLogChan, err := server.GetRouter(conf.Environment, conf.Environment != "prod", conf.SkipSecurity)
	if err != nil {
		panic(err)
	}

	go logger.JSONLog(mainLogChan)
	go logger.JSONLog(errorLogChan)

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

	go func() {
		http.Handle("/metrics", prometheus.Handler())

		log.Fatal(http.ListenAndServe(conf.TelemetryAddr, nil))
	}()

	if conf.UseSSL {
		log.Printf("Starting SSL server at \"%s\" in %s", conf.ListenHostPort, time.Now().Sub(startTime))
		log.Fatal(server.ListenAndServeTLS("./self.crt", "./self.key"))
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
	TLSConfig.ClientAuth = tls.NoClientCert
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
