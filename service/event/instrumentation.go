package event

import (
	"strings"
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/platform/metrics"
)

const (
	fieldMethod    = "method"
	fieldNamespace = "namespace"
	fieldStore     = "store"
	subsytem       = "service_event"
)

var (
	fieldKeys = []string{fieldMethod, fieldNamespace, fieldStore}

	namespace         string
	errCount, opCount kitmetrics.Counter
	opLatency         kitmetrics.TimeHistogram
)

type instrumentService struct {
	Service

	errCount  kitmetrics.Counter
	opCount   kitmetrics.Counter
	opLatency kitmetrics.TimeHistogram
	store     string
}

// InstrumentMiddleware observes key apsects of Service operations and exposes
// Prometheus metrics.
func InstrumentMiddleware(ns, store string) ServiceMiddleware {
	if namespace == "" {
		namespace = strings.Replace(ns, "-", "_", -1)
	}

	if errCount == nil {
		errCount = kitprometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsytem,
			Name:      "err_count",
			Help:      "Number of failed operations",
		}, fieldKeys)
	}

	if opCount == nil {
		opCount = kitprometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsytem,
			Name:      "op_count",
			Help:      "Number of operations performed",
		}, fieldKeys)
	}

	if opLatency == nil {
		opLatency = kitmetrics.NewTimeHistogram(
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
	}

	return func(next Service) Service {
		return &instrumentService{
			errCount:  errCount,
			opCount:   opCount,
			opLatency: opLatency,
			Service:   next,
			store:     store,
		}
	}
}

func (s *instrumentService) ActiveUserIDs(
	ns string,
	p Period,
) (ids []uint64, err error) {
	defer func(begin time.Time) {
		s.track("ActiveUserIDs", ns, begin, err)
	}(time.Now())

	return s.Service.ActiveUserIDs(ns, p)
}

func (s *instrumentService) CreatedByDay(
	ns string,
	start, end time.Time,
) (ts metrics.Timeseries, err error) {
	defer func(begin time.Time) {
		s.track("CreatedByDay", ns, begin, err)
	}(time.Now())

	return s.Service.CreatedByDay(ns, start, end)
}

func (s *instrumentService) Put(
	ns string,
	input *Event,
) (output *Event, err error) {
	defer func(begin time.Time) {
		s.track("Put", ns, begin, err)
	}(time.Now())

	return s.Service.Put(ns, input)
}

func (s *instrumentService) Query(
	ns string,
	opts QueryOptions,
) (list List, err error) {
	defer func(begin time.Time) {
		s.track("Query", ns, begin, err)
	}(time.Now())

	return s.Service.Query(ns, opts)
}

func (s *instrumentService) Setup(ns string) (err error) {
	defer func(begin time.Time) {
		s.track("Setup", ns, begin, err)
	}(time.Now())

	return s.Service.Setup(ns)
}

func (s *instrumentService) Teardown(ns string) (err error) {
	defer func(begin time.Time) {
		s.track("Teardown", ns, begin, err)
	}(time.Now())

	return s.Service.Teardown(ns)
}

func (s *instrumentService) track(
	method, namespace string,
	begin time.Time,
	err error,
) {
	var (
		m = kitmetrics.Field{
			Key:   fieldMethod,
			Value: method,
		}
		n = kitmetrics.Field{
			Key:   fieldNamespace,
			Value: namespace,
		}
		store = kitmetrics.Field{
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
