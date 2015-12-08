// Command intaker will launch a specified frontend for Tapglue
package main

import (
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	_ "expvar"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/syslog"
	mr "math/rand"
	"net/http"
	_ "net/http/pprof"
	"os"
	"runtime"
	"time"

	klog "github.com/go-kit/kit/log"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/errors"
	handler "github.com/tapglue/multiverse/handler/http"
	"github.com/tapglue/multiverse/limiter/redis"
	tgLogger "github.com/tapglue/multiverse/logger"
	"github.com/tapglue/multiverse/server"
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_postgres_core "github.com/tapglue/multiverse/v04/core/postgres"
	v04_redis_core "github.com/tapglue/multiverse/v04/core/redis"
	v04_postgres "github.com/tapglue/multiverse/v04/storage/postgres"
	v04_redis "github.com/tapglue/multiverse/v04/storage/redis"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar   = "TAPGLUE_INTAKER_CONFIG_PATH"
	apiVersionNext = "0.4"
	component      = "intaker"
)

var (
	currentRevision = "0000000-dev"

	conf      *config.Config
	startTime time.Time
)

func init() {
	startTime = time.Now()

	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	mr.Seed(time.Now().UTC().UnixNano())

}

func main() {
	var (
		forceNoSec = flag.Bool("force-no-sec", false, "Force no sec enables launching the backend in production without security checks")
	)
	flag.Parse()

	conf = config.NewConf(EnvConfigVar)

	if conf.SkipSecurity && conf.Environment == "prod" {
		if !*forceNoSec {
			panic("attempted to launch in production with no security checks enabled")
		}
	}

	// Setup logging
	var out io.Writer = os.Stdout

	if conf.UseSysLog {
		syslogWriter, err := syslog.New(syslog.LOG_INFO, "intaker")
		if err == nil {
			log.Printf("logging to syslog is enabled. Please tail your syslog for intaker app for further logs\n")
			log.SetFlags(0) // rsyslog will create the timestamps for us
			log.SetOutput(syslogWriter)
			out = syslogWriter
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

	logger := klog.NewContext(
		klog.NewJSONLogger(out),
	).With(
		"caller", klog.Caller(3),
		"component", component,
		"host", currentHostname,
		"revision", currentRevision,
	)

	// Setup services
	var (
		pgClient    = v04_postgres.New(conf.Postgres)
		redisClient = v04_redis.NewRedigoPool(conf.RateLimiter)
		rApps       = v04_redis_core.NewApplication(redisClient)
		rateLimiter = redis.NewLimiter(redisClient, "test:ratelimiter:app:")
	)

	var apps app.StrangleService
	apps = v04_postgres_core.NewApplication(pgClient, rApps)
	apps = app.InstrumentStrangleMiddleware(component, "postgres")(apps)
	apps = app.LogStrangleMiddleware(logger, "postgres")(apps)

	var connections connection.StrangleService
	connections = v04_postgres_core.NewConnection(pgClient)
	connections = connection.InstrumentMiddleware(component, "postgres")(connections)
	connections = connection.LogStrangleMiddleware(logger, "postgres")(connections)

	var objects object.Service
	objects = object.NewPostgresService(pgClient.MainDatastore())
	objects = object.InstrumentMiddleware(component, "postgres")(objects)
	objects = object.LogMiddleware(logger, "postgres")(objects)

	var users user.StrangleService
	users = v04_postgres_core.NewApplicationUser(pgClient)
	users = user.InstrumentMiddleware(component, "postgres")(users)
	users = user.LogStrangleMiddleware(logger, "postgres")(users)

	// Setup controllers
	var (
		commentController = controller.NewCommentController(objects)
		objectController  = controller.NewObjectController(connections, objects)
		postController    = controller.NewPostController(connections, objects)
	)

	// Setup middlewares
	var (
		withApp = handler.Chain(
			handler.CtxPrepare(apiVersionNext),
			handler.Log(logger),
			handler.Instrument(component),
			handler.SecureHeaders(),
			handler.DebugHeaders(currentRevision, currentHostname),
			handler.Gzip(),
			handler.HasUserAgent(),
			handler.ValidateContent(),
			handler.CtxApp(apps),
			handler.RateLimit(rateLimiter),
		)
		withUser = handler.Chain(
			withApp,
			handler.CtxUser(users),
		)
	)

	// Setup Server
	server.Setup(conf, currentRevision, currentHostname)

	// Setup Router
	router, mainLogChan, errorLogChan, err := server.GetRouter(conf.Environment, conf.Environment != "prod", conf.SkipSecurity)
	if err != nil {
		panic(err)
	}

	go tgLogger.JSONLog(mainLogChan)
	go tgLogger.JSONLog(errorLogChan)

	next := router.PathPrefix(fmt.Sprintf("/%s", apiVersionNext)).Subrouter()

	next.Methods("POST").PathPrefix("/objects").Name("objectCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectCreate(objectController),
		),
	)

	next.Methods("DELETE").PathPrefix("/objects/{objectID:[0-9]+}").Name("objectDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectDelete(objectController),
		),
	)

	next.Methods("GET").PathPrefix("/objects/{objectID:[0-9]+}").Name("objectRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectRetrieve(objectController),
		),
	)

	next.Methods("PUT").PathPrefix("/objects/{objectID:[0-9]+}").Name("objectUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectUpdate(objectController),
		),
	)

	next.Methods("GET").PathPrefix("/objects").Name("objectListAll").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.ObjectListAll(objectController),
		),
	)

	next.Methods("GET").PathPrefix("/me/objects/connections").Name("objectListConnections").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectListConnections(objectController),
		),
	)

	next.Methods("GET").PathPrefix("/me/objects").Name("objectList").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectList(objectController),
		),
	)

	next.Methods("POST").PathPrefix("/posts/{postID:[0-9]+}/comments").Name("commentCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentCreate(commentController),
		),
	)

	next.Methods("DELETE").PathPrefix("/posts/{postID:[0-9]+}/comments/{commentID:[0-9]+}").Name("commentDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentDelete(commentController),
		),
	)

	next.Methods("GET").PathPrefix("/posts/{postID:[0-9]+}/comments/{commentID:[0-9]+}").Name("commentRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentRetrieve(commentController),
		),
	)

	next.Methods("PUT").PathPrefix("/posts/{postID:[0-9]+}/comments/{commentID:[0-9]+}").Name("commentUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentUpdate(commentController),
		),
	)

	next.Methods("GET").PathPrefix("/posts/{postID:[0-9]+}/comments").Name("commentList").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.CommentList(commentController, users),
		),
	)

	next.Methods("POST").PathPrefix("/posts").Name("postCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostCreate(postController),
		),
	)

	next.Methods("DELETE").PathPrefix("/posts/{postID:[0-9]+}").Name("postDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostDelete(postController),
		),
	)

	next.Methods("GET").PathPrefix("/posts/{postID:[0-9]+}").Name("postRetrieve").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.PostRetrieve(postController),
		),
	)

	next.Methods("PUT").PathPrefix("/posts/{postID:[0-9]+}").Name("postUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostUpdate(postController),
		),
	)

	next.Methods("GET").PathPrefix("/posts").Name("postListAll").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.PostListAll(postController, users),
		),
	)

	next.Methods("GET").PathPrefix("/me/posts/connections").Name("postListConnections").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostListMeConnections(postController, users),
		),
	)

	next.Methods("GET").PathPrefix("/me/posts").Name("postListMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostListMe(postController, users),
		),
	)

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
