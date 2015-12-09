package event

import (
	"strings"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	fieldMethod    = "method"
	fieldNamespace = "namespace"
	fieldStore     = "store"
)

type instrumentStrangleService struct {
	StrangleService

	errCount  metrics.Counter
	opCount   metrics.Counter
	opLatency metrics.TimeHistogram
	store     string
}

// InstrumentStrangleMiddleware observes key aspects of Service operations and exposes
// Prometheus metrics.
func InstrumentStrangleMiddleware(ns string, store string) StrangleMiddleware {
	var (
		fieldKeys = []string{fieldMethod, fieldNamespace, fieldStore}
		namespace = strings.Replace(ns, "-", "_", -1)
		subsytem  = "service_event"

		errCount = kitprometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsytem,
			Name:      "err_count",
			Help:      "Number of failed operations",
		}, fieldKeys)
		opCount = kitprometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsytem,
			Name:      "op_count",
			Help:      "Number of operations performed",
		}, fieldKeys)
		opLatency = metrics.NewTimeHistogram(
			time.Second,
			kitprometheus.NewHistogram(
				prometheus.HistogramOpts{
					Namespace: namespace,
					Subsystem: subsytem,
					Name:      "op_latency_seconds",
					Help:      "Distribution of op duration in seconds",
				},
				fieldKeys,
			),
		)
	)

	return func(next StrangleService) StrangleService {
		return &instrumentStrangleService{
			errCount:        errCount,
			opCount:         opCount,
			opLatency:       opLatency,
			store:           store,
			StrangleService: next,
		}
	}
}

func (s *instrumentStrangleService) Create(
	orgID, appID int64,
	userID uint64,
	event *v04_entity.Event,
) (errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "Create", err)
	}(time.Now())

	return s.StrangleService.Create(orgID, appID, userID, event)
}

func (s *instrumentStrangleService) Delete(
	orgID, appID int64,
	userID, eventID uint64,
) (errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "Delete", err)
	}(time.Now())

	return s.StrangleService.Delete(orgID, appID, userID, eventID)
}

func (s *instrumentStrangleService) ListAll(
	orgID, appID int64,
	condition core.EventCondition,
) (es []*v04_entity.Event, errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "ListAll", err)
	}(time.Now())

	return s.StrangleService.ListAll(orgID, appID, condition)
}

func (s *instrumentStrangleService) track(
	orgID, appID int64,
	begin time.Time,
	method string,
	err error,
) {
	var (
		m = metrics.Field{
			Key:   fieldMethod,
			Value: "FriendsAndFollowingIDs",
		}
		n = metrics.Field{
			Key:   fieldNamespace,
			Value: namespace(orgID, appID),
		}
		store = metrics.Field{
			Key:   fieldStore,
			Value: s.store,
		}
	)

	if err != nil {
		s.errCount.With(m).With(n).With(store).Add(1)
	}

	s.opCount.With(m).With(n).With(store).Add(1)
	s.opLatency.With(m).With(n).With(store).Observe(time.Since(begin))
}
