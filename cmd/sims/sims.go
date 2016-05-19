package main

import (
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	awsSession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/go-kit/kit/log"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/jmoiron/sqlx"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/device"
	"github.com/tapglue/multiverse/service/user"
)

const (
	component        = "sims"
	namespaceService = "service"
	namespaceSource  = "source"
	subsystemErr     = "err"
	subsystemOp      = "op"
	subsystemQueue   = "queue"

	simsTestARN = "arn:aws:sns:eu-central-1:775034650473:app/APNS_SANDBOX/simsTest"
)

var (
	// Control flow.
	ErrEndpointMissing  = errors.New("endpoint missing")
	ErrPlatformNotFound = errors.New("platform not found")

	defaultDeleted = false
	defaultEnabled = true
	// Set at build time.
	revision = "0000000-dev"
)

type ackFunc func() error
type createEndpointFunc func(namespace string, platformARN string, device *device.Device) (string, error)
type fetchUserFunc func(namespace string, id uint64) (*user.User, error)
type findEndpointARNsFunc func(namespace string, userID uint64) ([]string, error)
type getPlatformARNFunc func(namespace string) (string, error)
type pushAPNSSandboxFunc func(arn, message string) error

type message struct {
	ackFunc   ackFunc
	message   string
	namespace string
	recipient uint64
}

func main() {
	var (
		begin = time.Now()

		awsID         = flag.String("aws.id", "", "Identifier for AWS requests")
		awsRegion     = flag.String("aws.region", "us-east-1", "AWS region to operate in")
		awsSecret     = flag.String("aws.secret", "", "Identification secret for AWS requests")
		postgresURL   = flag.String("postgres.url", "", "Postgres URL to connect to")
		telemetryAddr = flag.String("telemetry.addr", ":9001", "Address to expose telemetry on")
	)
	flag.Parse()

	logger := log.NewContext(
		log.NewJSONLogger(os.Stdout),
	).With(
		"caller", log.Caller(3),
		"component", component,
		"revision", revision,
	)

	hostname, err := os.Hostname()
	if err != nil {
		logger.Log("err", err, "lifecycle", "abort")
	}

	logger = log.NewContext(logger).With("host", hostname)

	// Setup instrumenation
	go func(addr string, begin time.Time, logger log.Logger) {
		http.Handle("/metrics", prometheus.Handler())

		logger = log.NewContext(logger).With(
			"listen", addr,
			"sub", "telemetry",
		)

		logger.Log(
			"duration", time.Now().Sub(begin).Nanoseconds(),
			"lifecycle", "start",
		)

		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logger.Log(
				"err", err,
				"lifecycle", "abort",
			)
		}
	}(*telemetryAddr, begin, logger)

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

	aSession := awsSession.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(*awsID, *awsSecret, ""),
		Region:      aws.String(*awsRegion),
	})

	db, err := sqlx.Connect("postgres", *postgresURL)
	if err != nil {
		logger.Log("err", err, "lifecycle", "abort")
		os.Exit(1)
	}

	var connections connection.Service
	connections = connection.NewPostgresService(db)
	connections = connection.InstrumentServiceMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(connections)
	connections = connection.LogServiceMiddleware(logger, "postgres")(connections)

	var conSource connection.Source

	s, err := connection.SQSSource(sqs.New(aSession))
	if err != nil {
		logger.Log("err", err, "lifecycle", "abort")
		os.Exit(1)
	}

	conSource = s
	conSource = connection.InstrumentSourceMiddleware(
		component,
		"sqs",
		sourceErrCount,
		sourceOpCount,
		sourceOpLatency,
		sourceQueueLatency,
	)(conSource)
	conSource = connection.LogSourceMiddleware("sqs", logger)(conSource)

	var devices device.Service
	devices = device.PostgresService(db)
	devices = device.InstrumentServiceMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(devices)
	devices = device.LogServiceMiddleware(logger, "postgres")(devices)

	var users user.Service
	users = user.NewPostgresService(db)
	users = user.InstrumentMiddleware(component, "postgres", serviceErrCount, serviceOpCount, serviceOpLatency)(users)
	users = user.LogMiddleware(logger, "postgres")(users)

	snsService := sns.New(aSession)

	var createEndpoint createEndpointFunc
	createEndpoint = func(ns, pARN string, device *device.Device) (string, error) {
		r, err := snsService.CreatePlatformEndpoint(&sns.CreatePlatformEndpointInput{
			PlatformApplicationArn: aws.String(pARN),
			Token: aws.String(device.Token),
		})
		if err != nil {
			return "", err
		}

		device.EndpointARN = *r.EndpointArn

		d, err := devices.Put(ns, device)
		if err != nil {
			return "", err
		}

		return d.EndpointARN, nil
	}

	var fetchUser fetchUserFunc
	fetchUser = func(ns string, id uint64) (*user.User, error) {
		us, err := users.Query(ns, user.QueryOptions{
			Enabled: &defaultEnabled,
			IDs: []uint64{
				id,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(us) == 0 {
			return nil, fmt.Errorf("user '%d' not found", id)
		}

		return us[0], nil
	}

	var getPlatformARN getPlatformARNFunc
	getPlatformARN = func(ns string) (string, error) {
		if ns == "app_1_610" {
			return simsTestARN, nil
		}

		return "", ErrPlatformNotFound
	}

	var findEndpointARNs findEndpointARNsFunc
	findEndpointARNs = func(ns string, userID uint64) ([]string, error) {
		as := []string{}

		pARN, err := getPlatformARN(ns)
		if err != nil {
			if err == ErrPlatformNotFound {
				return as, nil
			}
			return nil, err
		}

		ds, err := devices.Query(ns, device.QueryOptions{
			Deleted: &defaultDeleted,
			Platforms: []device.Platform{
				device.PlatformIOS,
			},
			UserIDs: []uint64{
				userID,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, d := range ds {
			if d.EndpointARN != "" {
				as = append(as, d.EndpointARN)
				continue
			}

			arn, err := createEndpoint(ns, pARN, d)
			if err != nil {
				return nil, err
			}

			as = append(as, arn)
		}

		return as, nil
	}

	var pushAPNSSandbox pushAPNSSandboxFunc
	pushAPNSSandbox = func(arn, msg string) error {
		_, err := snsService.Publish(&sns.PublishInput{
			Message: aws.String(
				fmt.Sprintf(
					`{"APNS_SANDBOX":"{\"aps\":{\"alert\":\"%s\"}}"}`,
					msg,
				),
			),
			MessageStructure: aws.String("json"),
			TargetArn:        aws.String(arn),
		})
		return err
	}

	logger.Log(
		"duration", time.Now().Sub(begin).Nanoseconds(),
		"lifecycle", "start",
		"sub", "worker",
	)

	msgc := make(chan message)

	go func() {
		err := connectionConsumer(msgc, conSource, fetchUser)
		if err != nil {
			logger.Log("err", err, "lifecycle", "abort")
			os.Exit(1)
		}
	}()

	err = pushChannel(
		msgc,
		findEndpointARNs,
		getPlatformARN,
		pushAPNSSandbox,
	)
	if err != nil {
		logger.Log("err", err, "lifecycle", "abort")
		os.Exit(1)
	}
}

func connectionConsumer(
	msgc chan<- message,
	conSource connection.Source,
	fetchUser fetchUserFunc,
) error {
	for {
		c, err := conSource.Consume()
		if err != nil {
			if connection.IsEmptySource(err) {
				continue
			}
			return err
		}

		// CONSUMER
		// filter
		if c.Old != nil {
			continue
		}

		if c.New.State != connection.StateConfirmed {
			continue
		}

		if c.New.Type != connection.TypeFollow {
			continue
		}

		// fetch recipients
		origin, err := fetchUser(c.Namespace, c.New.FromID)
		if err != nil {
			return fmt.Errorf("origin fetch: %s", err)
		}

		target, err := fetchUser(c.Namespace, c.New.ToID)
		if err != nil {
			return fmt.Errorf("target fetch: %s", err)
		}

		// send message
		msgc <- message{
			ackFunc: func() error {
				acked := false

				if acked {
					return nil
				}

				err := conSource.Ack(c.AckID)
				if err == nil {
					acked = true
				}
				return err
			},
			message: fmt.Sprintf(
				"%s %s (%s) started following you",
				origin.Firstname,
				origin.Lastname,
				origin.Username,
			),
			namespace: c.Namespace,
			recipient: target.ID,
		}
	}
}

func pushChannel(
	msgc <-chan message,
	findEndpointARNs findEndpointARNsFunc,
	getPlatformARN getPlatformARNFunc,
	pushAPNSSandbox pushAPNSSandboxFunc,
) error {
	// CHANNEL
	for msg := range msgc {
		// check if platform is enabled
		_, err := getPlatformARN(msg.namespace)
		if err != nil {
			if err == ErrPlatformNotFound {
				continue
			}
			return err
		}

		// find arns
		as, err := findEndpointARNs(msg.namespace, msg.recipient)
		if err != nil {
			return err
		}
		if len(as) == 0 {
			continue
		}

		for _, arn := range as {
			err := pushAPNSSandbox(arn, msg.message)
			if err != nil {
				return err
			}
		}

		// publish to endpoint
		err = msg.ackFunc()
		if err != nil {
			return err
		}
	}

	return nil
}
