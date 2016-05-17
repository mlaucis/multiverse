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
	"path/filepath"
	"runtime"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	klog "github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/config"
	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/errors"
	handler "github.com/tapglue/multiverse/handler/http"
	"github.com/tapglue/multiverse/limiter/redis"
	tgLogger "github.com/tapglue/multiverse/logger"
	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/server"
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/member"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/org"
	"github.com/tapglue/multiverse/service/session"
	"github.com/tapglue/multiverse/service/user"
	v04_postgres_core "github.com/tapglue/multiverse/v04/core/postgres"
	v04_redis_core "github.com/tapglue/multiverse/v04/core/redis"
	v04_postgres "github.com/tapglue/multiverse/v04/storage/postgres"
	v04_redis "github.com/tapglue/multiverse/v04/storage/redis"
)

const (
	// EnvConfigVar holds the name of the environment variable that holds the path to the config
	EnvConfigVar     = "TAPGLUE_INTAKER_CONFIG_PATH"
	apiVersionNext   = "0.4"
	component        = "gateway-http"
	namespaceService = "service"
	namespaceSource  = "source"
	subsystemErr     = "err"
	subsystemOp      = "op"
	subsystemQueue   = "queue"
)

// Supported source types.
const (
	sourceNop = "nop"
	sourceSQS = "sqs"
)

var (
	currentRevision = "0000000-dev"

	apppath   string
	conf      *config.Config
	startTime time.Time
)

func init() {
	startTime = time.Now()

	// Use all available CPU's
	runtime.GOMAXPROCS(runtime.NumCPU())

	// Seed random generator
	mr.Seed(time.Now().UTC().UnixNano())

	cwd, _ := os.Getwd()
	apppath, _ = filepath.Abs(filepath.Join(cwd, string(filepath.Separator)))
	apppath += string(filepath.Separator)
}

func main() {
	var (
		awsID      = flag.String("aws.id", "", "Identifier for AWS requests")
		awsRegion  = flag.String("aws.region", "us-east-1", "AWS Region to operate in")
		awsSecret  = flag.String("aws.secret", "", "Identification secret for AWS requests")
		source     = flag.String("source", sourceNop, "Source type used for state change propagations")
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

	// Setup instrumenation
	serviceFieldKeys := []string{
		metrics.FieldComponent,
		metrics.FieldMethod,
		metrics.FieldNamespace,
		metrics.FieldService,
		metrics.FieldStore,
	}

	serviceErrCount := kitprometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespaceService,
		Subsystem: subsystemErr,
		Name:      "count",
		Help:      "Number of failed service operations",
	}, serviceFieldKeys)

	serviceOpCount := kitprometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespaceService,
		Subsystem: subsystemOp,
		Name:      "count",
		Help:      "Number of service operations performed",
	}, serviceFieldKeys)

	serviceOpLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespaceService,
			Subsystem: subsystemOp,
			Name:      "latency_seconds",
			Help:      "Distribution of service op duration in seconds",
		},
		serviceFieldKeys,
	)
	prometheus.MustRegister(serviceOpLatency)

	sourceFieldKeys := []string{
		metrics.FieldComponent,
		metrics.FieldMethod,
		metrics.FieldNamespace,
		metrics.FieldSource,
		metrics.FieldStore,
	}

	sourceErrCount := kitprometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespaceSource,
		Subsystem: subsystemErr,
		Name:      "count",
		Help:      "Number of failed source operations",
	}, sourceFieldKeys)

	sourceOpCount := kitprometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespaceSource,
		Subsystem: subsystemOp,
		Name:      "count",
		Help:      "Number of source operations performed",
	}, sourceFieldKeys)

	sourceOpLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespaceSource,
			Subsystem: subsystemOp,
			Name:      "latency_seconds",
			Help:      "Distribution of source op duration in seconds",
			Buckets:   metrics.BucketsQueue,
		},
		sourceFieldKeys,
	)
	prometheus.MustRegister(sourceOpLatency)

	sourceQueueLatency := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: namespaceSource,
			Subsystem: subsystemQueue,
			Name:      "latency_seconds",
			Help:      "Distribution of message queue latency in seconds",
			Buckets:   metrics.BucketsQueue,
		},
		sourceFieldKeys,
	)
	prometheus.MustRegister(sourceQueueLatency)

	// Setup sources
	var (
		aSession = awsSession.New(&aws.Config{
			Credentials: credentials.NewStaticCredentials(*awsID, *awsSecret, ""),
			Region:      aws.String(*awsRegion),
		})
		sqsAPI = sqs.New(aSession)

		conSource connection.Source
	)

	switch *source {
	case sourceNop:
		conSource = connection.NopSource()
	case sourceSQS:
		s, err := connection.SQSSource(sqsAPI)
		if err != nil {
			logger.Log("err", err, "lifecycle", "abort")
			os.Exit(1)
		}

		conSource = s
	default:
		logger.Log(
			"err", fmt.Sprintf("unsupported Source type %s", *source),
			"lifecycle", "abort",
		)
		os.Exit(1)
	}

	conSource = connection.InstrumentSourceMiddleware(
		component,
		*source,
		sourceErrCount,
		sourceOpCount,
		sourceOpLatency,
		sourceQueueLatency,
	)(conSource)
	conSource = connection.LogSourceMiddleware(*source, logger)(conSource)

	// Setup services
	var (
		pgClient    = v04_postgres.New(conf.Postgres)
		redisClient = v04_redis.NewRedigoPool(conf.RateLimiter)
		rApps       = v04_redis_core.NewApplication(redisClient)
		rateLimiter = redis.NewLimiter(redisClient, "test:ratelimiter:app:")
	)

	var apps app.StrangleService
	apps = v04_postgres_core.NewApplication(pgClient, rApps)
	apps = app.InstrumentStrangleMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(apps)
	apps = app.LogStrangleMiddleware(logger, "postgres")(apps)

	var connections connection.Service
	connections = connection.NewPostgresService(pgClient.MainDatastore())
	connections = connection.InstrumentServiceMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(connections)
	connections = connection.LogServiceMiddleware(logger, "postgres")(connections)
	// Combine connection service and source.
	connections = connection.SourcingServiceMiddleware(conSource)(connections)

	var events event.Service
	events = event.NewPostgresService(pgClient.MainDatastore())
	events = event.InstrumentMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(events)
	events = event.LogMiddleware(logger, "postgres")(events)

	var members member.StrangleService
	members = v04_postgres_core.NewMember(pgClient)
	members = member.InstrumentStrangleMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(members)
	members = member.LogStrangleMiddleware(logger, "postgres")(members)

	var objects object.Service
	objects = object.NewPostgresService(pgClient.MainDatastore())
	objects = object.InstrumentMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(objects)
	objects = object.LogMiddleware(logger, "postgres")(objects)

	var orgs org.StrangleService
	orgs = v04_postgres_core.NewOrganization(pgClient)
	orgs = org.InstrumentStrangleMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(orgs)
	orgs = org.LogStrangleMiddleware(logger, "postgres")(orgs)

	var sessions session.Service
	sessions = session.NewPostgresService(pgClient.MainDatastore())
	sessions = session.InstrumentMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(sessions)
	sessions = session.LogMiddleware(logger, "postgres")(sessions)

	var users user.Service
	users = user.NewPostgresService(pgClient.MainDatastore())
	users = user.InstrumentMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(users)
	users = user.LogMiddleware(logger, "postgres")(users)

	// Setup controllers
	var (
		analyticsController = controller.NewAnalyticsController(
			apps,
			connections,
			events,
			objects,
			users,
		)
		commentController        = controller.NewCommentController(connections, objects, users)
		connectionController     = controller.NewConnectionController(connections, users)
		eventController          = controller.NewEventController(connections, events, objects, users)
		feedController           = controller.NewFeedController(connections, events, objects, users)
		likeController           = controller.NewLikeController(connections, events, objects, users)
		objectController         = controller.NewObjectController(connections, objects)
		postController           = controller.NewPostController(connections, events, objects, users)
		recommendationController = controller.NewRecommendationController(
			connections,
			events,
			users,
		)
		userController = controller.NewUserController(connections, sessions, users)
	)

	// Setup middlewares
	var (
		withConstraints = handler.Chain(
			handler.CtxPrepare(apiVersionNext),
			handler.Log(logger),
			handler.Instrument(component),
			handler.SecureHeaders(),
			handler.DebugHeaders(currentRevision, currentHostname),
			handler.DeviceID(),
			handler.CORS(),
			handler.Gzip(),
			handler.HasUserAgent(),
			handler.ValidateContent(),
		)
		withOrg = handler.Chain(
			withConstraints,
			handler.CtxOrg(orgs),
		)
		withMember = handler.Chain(
			withOrg,
			handler.CtxMember(members),
		)
		withApp = handler.Chain(
			withConstraints,
			handler.CtxApp(apps),
			handler.RateLimit(rateLimiter),
		)
		withUser = handler.Chain(
			withApp,
			handler.CtxUser(sessions, users),
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

	router.Methods("POST").Path(`/analytics`).Name("analytics").HandlerFunc(
		handler.Wrap(
			handler.Chain(
				handler.CtxPrepare(apiVersionNext),
				handler.Log(logger),
			),
			handler.AnalyticsDeprecated(),
		),
	)

	router.Methods("GET").Path(`/health-45016490610398192`).Name("healthcheck").HandlerFunc(
		handler.Wrap(
			handler.CtxPrepare(apiVersionNext),
			handler.Health(pgClient.MainDatastore(), redisClient),
		),
	)

	next := router.PathPrefix(fmt.Sprintf("/%s", apiVersionNext)).Subrouter()

	next.Methods("POST").Path(`/analytics`).Name("analytics").HandlerFunc(
		handler.Wrap(
			handler.Chain(
				handler.CtxPrepare(apiVersionNext),
				handler.Log(logger),
			),
			handler.Analytics(),
		),
	)

	next.Methods("GET").Path(`/me/connections/{state:[a-z]+}`).Name("getConnectionByState").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionByState(connectionController),
		),
	)

	next.Methods("DELETE").Path(`/me/connections/{type:[a-z]+}/{toID:[0-9]+}`).Name("deleteConnection").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionDelete(connectionController),
		),
	)

	next.Methods("POST").Path(`/me/connections/social`).Name("createSocialConnections").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionSocial(connectionController),
		),
	)

	next.Methods("PUT").Path(`/me/connections`).Name("updateConnection").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionUpdate(connectionController),
		),
	)

	next.Methods("GET").Path(`/me/followers`).Name("getFollowersMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionFollowersMe(connectionController),
		),
	)

	next.Methods("GET").Path(`/me/follows`).Name("getFollowingsMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionFollowingsMe(connectionController),
		),
	)

	next.Methods("GET").Path(`/me/friends`).Name("getFriendsMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionFriendsMe(connectionController),
		),
	)

	next.Methods("GET").Path(`/users/{userID:[0-9]+}/followers`).Name("getFollowers").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionFollowers(connectionController),
		),
	)

	next.Methods("GET").Path(`/users/{userID:[0-9]+}/follows`).Name("getFollowings").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionFollowings(connectionController),
		),
	)

	next.Methods("GET").Path(`/users/{userID:[0-9]+}/friends`).Name("getFriends").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ConnectionFriends(connectionController),
		),
	)

	next.Methods("DELETE").Path(`/me/events/{id:[0-9]+}`).Name("deleteEvent").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.EventDelete(eventController),
		),
	)

	next.Methods("PUT").Path(`/me/events/{id:[0-9]+}`).Name("eventUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.EventUpdate(eventController),
		),
	)

	next.Methods("GET").Path(`/me/events`).Name("eventListMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.EventListMe(eventController),
		),
	)

	next.Methods("POST").Path(`/me/events`).Name("eventCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.EventCreate(eventController),
		),
	)

	next.Methods("GET").Path(`/users/{userID:[0-9]+}/events`).Name("eventListUser").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.EventListUser(eventController),
		),
	)

	next.Methods("POST").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/comments`).Name("externalCommentCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ExternalCommentCreate(commentController),
		),
	)

	next.Methods("DELETE").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/comments/{commentID:[0-9]+}`).Name("externalCommentDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ExternalCommentDelete(commentController),
		),
	)

	next.Methods("GET").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/comments/{commentID:[0-9]+}`).Name("externalCommentRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ExternalCommentRetrieve(commentController),
		),
	)

	next.Methods("PUT").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/comments/{commentID:[0-9]+}`).Name("externalCommentUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ExternalCommentUpdate(commentController),
		),
	)

	next.Methods("GET").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/comments`).Name("externalCommentList").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.ExternalCommentList(commentController),
		),
	)

	next.Methods("POST").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/likes`).Name("externalLikeCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ExternalLikeCreate(likeController),
		),
	)

	next.Methods("DELETE").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/likes`).Name("externalLikeDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ExternalLikeDelete(likeController),
		),
	)

	next.Methods("GET").Path(`/externals/{externalID:[a-zA-Z0-9\-\_]+}/likes`).Name("externalLikeList").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.ExternalLikeList(likeController),
		),
	)

	next.Methods("OPTIONS").Path("/").Name("CORS").HandlerFunc(
		handler.Wrap(
			withMember,
			func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
		),
	)

	next.Methods("GET").Path("/me/feed/events").Name("feedEvents").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.FeedEvents(feedController),
		),
	)

	next.Methods("GET").Path("/me/feed/posts").Name("feedPosts").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.FeedPosts(feedController),
		),
	)

	next.Methods("GET").Path("/me/feed").Name("feedNews").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.FeedNews(feedController),
		),
	)

	next.Methods("POST").Path("/objects").Name("objectCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectCreate(objectController),
		),
	)

	next.Methods("DELETE").Path("/objects/{objectID:[0-9]+}").Name("objectDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectDelete(objectController),
		),
	)

	next.Methods("GET").Path("/objects/{objectID:[0-9]+}").Name("objectRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectRetrieve(objectController),
		),
	)

	next.Methods("PUT").Path("/objects/{objectID:[0-9]+}").Name("objectUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectUpdate(objectController),
		),
	)

	next.Methods("GET").Path("/objects").Name("objectListAll").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.ObjectListAll(objectController),
		),
	)

	next.Methods("GET").Path("/me/objects/connections").Name("objectListConnections").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectListConnections(objectController),
		),
	)

	next.Methods("GET").Path("/me/objects").Name("objectList").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.ObjectList(objectController),
		),
	)

	next.Methods("GET").Path(`/orgs/{orgID:[a-zA-Z0-9\-]+}/apps/{appID:[a-zA-Z0-9\-]+}/analytics`).Name("appAnalytics").HandlerFunc(
		handler.Wrap(
			withMember,
			handler.AnalyticsApp(analyticsController),
		),
	)

	next.Methods("POST").Path("/posts/{postID:[0-9]+}/comments").Name("commentCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentCreate(commentController),
		),
	)

	next.Methods("DELETE").Path("/posts/{postID:[0-9]+}/comments/{commentID:[0-9]+}").Name("commentDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentDelete(commentController),
		),
	)

	next.Methods("GET").Path("/posts/{postID:[0-9]+}/comments/{commentID:[0-9]+}").Name("commentRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentRetrieve(commentController),
		),
	)

	next.Methods("PUT").Path("/posts/{postID:[0-9]+}/comments/{commentID:[0-9]+}").Name("commentUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentUpdate(commentController),
		),
	)

	next.Methods("GET").Path("/posts/{postID:[0-9]+}/comments").Name("commentList").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.CommentList(commentController),
		),
	)

	next.Methods("POST").Path("/posts/{postID:[0-9]+}/likes").Name("likeCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.LikeCreate(likeController),
		),
	)

	next.Methods("DELETE").Path("/posts/{postID:[0-9]+}/likes").Name("likeDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.LikeDelete(likeController),
		),
	)

	next.Methods("GET").Path("/posts/{postID:[0-9]+}/likes").Name("likeList").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.LikeList(likeController),
		),
	)

	next.Methods("POST").Path("/posts").Name("postCreate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostCreate(postController),
		),
	)

	next.Methods("DELETE").Path("/posts/{postID:[0-9]+}").Name("postDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostDelete(postController),
		),
	)

	next.Methods("GET").Path("/posts/{postID:[0-9]+}").Name("postRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostRetrieve(postController),
		),
	)

	next.Methods("PUT").Path("/posts/{postID:[0-9]+}").Name("postUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostUpdate(postController),
		),
	)

	next.Methods("GET").Path("/posts").Name("postListAll").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostListAll(postController),
		),
	)

	next.Methods("GET").Path("/me/posts").Name("postListMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostListMe(postController),
		),
	)

	next.Methods("GET").Path("/users/{userID:[0-9]+}/posts").Name("postList").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.PostList(postController),
		),
	)

	next.Methods("GET").Path("/recommendations/users/active/day").Name("recommendUsersActiveDay").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.RecommendUsersActiveDay(recommendationController),
		),
	)

	next.Methods("GET").Path("/recommendations/users/active/week").Name("recommendUsersActiveWeek").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.RecommendUsersActiveWeek(recommendationController),
		),
	)

	next.Methods("GET").Path("/recommendations/users/active/month").Name("recommendUsersActiveMonth").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.RecommendUsersActiveMonth(recommendationController),
		),
	)

	next.Methods("GET").Path("/me").Name("userRetrieveMe").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserRetrieveMe(userController),
		),
	)

	next.Methods("PUT").Path("/me").Name("userUpdate").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserUpdate(userController),
		),
	)

	next.Methods("POST").Path("/me/login").Name("userMeLogin").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.UserLogin(userController),
		),
	)

	next.Methods("DELETE").Path("/me/logout").Name("userLogout").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserLogout(userController),
		),
	)

	next.Methods("DELETE").Path("/me").Name("userDelete").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserDelete(userController),
		),
	)

	next.Methods("GET").Path("/users/{userID:[0-9]+}").Name("userRetrieve").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserRetrieve(userController),
		),
	)

	next.Methods("POST").Path("/users/login").Name("userLogin").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.UserLogin(userController),
		),
	)

	next.Methods("POST").Path("/users/search/emails").Name("userSearchEmails").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserSearchEmails(userController),
		),
	)

	next.Methods("POST").Path(`/users/search/{platform:[a-z]+}`).Name("userSearchPlatform").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserSearchPlatform(userController),
		),
	)

	next.Methods("GET").Path(`/users/search`).Name("userSearch").HandlerFunc(
		handler.Wrap(
			withUser,
			handler.UserSearch(userController),
		),
	)

	next.Methods("POST").Path(`/users`).Name("userCreate").HandlerFunc(
		handler.Wrap(
			withApp,
			handler.UserCreate(userController),
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

	go func() {
		http.Handle("/metrics", prometheus.Handler())

		logger.Log(
			"duration", time.Now().Sub(startTime).Nanoseconds(),
			"lifecycle", "start",
			"listen", conf.TelemetryAddr,
			"sub", "telemetry",
		)
		log.Fatal(http.ListenAndServe(conf.TelemetryAddr, nil))
	}()

	logger.Log(
		"duration", time.Now().Sub(startTime).Nanoseconds(),
		"lifecycle", "start",
		"listen", conf.ListenHostPort,
		"sub", "api",
	)

	if conf.UseSSL {
		log.Fatal(server.ListenAndServeTLS(apppath+"self.crt", apppath+"self.key"))
	} else {
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
	pem, err := ioutil.ReadFile(apppath + "root-ca.pem")
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
	pem, err := ioutil.ReadFile(apppath + "origin-pull-ca.pem")
	if err != nil {
		panic(err)
	}

	rootCertPool := x509.NewCertPool()
	if !rootCertPool.AppendCertsFromPEM(pem) {
		panic("Failed appending certs")
	}

	return rootCertPool
}
