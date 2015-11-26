package object

import (
	"strings"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	fieldMethod    = "method"
	fieldNamespace = "namespace"
	fieldStore     = "store"
)

type instrumentService struct {
	Service

	errCount  metrics.Counter
	opCount   metrics.Counter
	opLatency metrics.TimeHistogram
	store     string
}

// InstrumentMiddleware observes key aspects of Service operations and exposes
// Prometheus metrics.
func InstrumentMiddleware(ns string, store string) ServiceMiddleware {
	return func(next Service) Service {
		var (
			fieldKeys = []string{fieldMethod, fieldNamespace, fieldStore}
			namespace = strings.Replace(ns, "-", "_", -1)
			subsytem  = "service_object"

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

		return &instrumentService{
			errCount:  errCount,
			opCount:   opCount,
			opLatency: opLatency,
			Service:   next,
			store:     store,
		}
	}
}

func (s *instrumentService) Put(ns string, object *Object) (o *Object, err error) {
	defer func(begin time.Time) {
		s.track("put", ns, begin, err)
	}(time.Now())

	return s.Service.Put(ns, object)
}

func (s *instrumentService) Query(ns string, opts QueryOptions) (os []*Object, err error) {
	defer func(begin time.Time) {
		s.track("query", ns, begin, err)
	}(time.Now())

	return s.Service.Query(ns, opts)
}

func (s *instrumentService) Remove(ns string, id uint64) (err error) {
	defer func(begin time.Time) {
		s.track("remove", ns, begin, err)
	}(time.Now())

	return s.Service.Remove(ns, id)
}

func (s *instrumentService) Setup(ns string) (err error) {
	defer func(begin time.Time) {
		s.track("setup", ns, begin, err)
	}(time.Now())

	return s.Service.Setup(ns)
}

func (s *instrumentService) Teardown(ns string) (err error) {
	defer func(begin time.Time) {
		s.track("teardown", ns, begin, err)
	}(time.Now())

	return s.Service.Teardown(ns)
}

func (s *instrumentService) track(
	method string,
	namespace string,
	begin time.Time,
	err error,
) {
	var (
		m = metrics.Field{
			Key:   fieldMethod,
			Value: method,
		}
		ns = metrics.Field{
			Key:   fieldNamespace,
			Value: namespace,
		}
		store = metrics.Field{
			Key:   fieldStore,
			Value: s.store,
		}
	)

	if err != nil {
		s.errCount.With(m).With(ns).With(store).Add(1)
	}

	s.opCount.With(m).With(ns).With(store).Add(1)
	s.opLatency.With(m).With(ns).With(store).Observe(time.Since(begin))
}
