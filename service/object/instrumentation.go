package object

import (
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/tapglue/multiverse/platform/metrics"
)

const serviceName = "object"

type instrumentService struct {
	errCount  kitmetrics.Counter
	next      Service
	opCount   kitmetrics.Counter
	opLatency *prometheus.HistogramVec
	store     string
}

// InstrumentMiddleware observes key aspects of Service operations and exposes
// Prometheus metrics.
func InstrumentMiddleware(
	store string,
	errCount kitmetrics.Counter,
	opCount kitmetrics.Counter,
	opLatency *prometheus.HistogramVec,
) ServiceMiddleware {
	return func(next Service) Service {
		return &instrumentService{
			errCount:  errCount,
			next:      next,
			opCount:   opCount,
			opLatency: opLatency,
			store:     store,
		}
	}
}

func (s *instrumentService) Count(ns string, opts QueryOptions) (count int, err error) {
	defer func(begin time.Time) {
		s.track("Count", ns, begin, err)
	}(time.Now())

	return s.next.Count(ns, opts)
}

func (s *instrumentService) CreatedByDay(
	ns string,
	start, end time.Time,
) (ts metrics.Timeseries, err error) {
	defer func(begin time.Time) {
		s.track("CreatedByDay", ns, begin, err)
	}(time.Now())

	return s.next.CreatedByDay(ns, start, end)
}

func (s *instrumentService) Put(ns string, object *Object) (o *Object, err error) {
	defer func(begin time.Time) {
		s.track("put", ns, begin, err)
	}(time.Now())

	return s.next.Put(ns, object)
}

func (s *instrumentService) Query(ns string, opts QueryOptions) (os List, err error) {
	defer func(begin time.Time) {
		s.track("query", ns, begin, err)
	}(time.Now())

	return s.next.Query(ns, opts)
}

func (s *instrumentService) Remove(ns string, id uint64) (err error) {
	defer func(begin time.Time) {
		s.track("remove", ns, begin, err)
	}(time.Now())

	return s.next.Remove(ns, id)
}

func (s *instrumentService) Setup(ns string) (err error) {
	defer func(begin time.Time) {
		s.track("setup", ns, begin, err)
	}(time.Now())

	return s.next.Setup(ns)
}

func (s *instrumentService) Teardown(ns string) (err error) {
	defer func(begin time.Time) {
		s.track("teardown", ns, begin, err)
	}(time.Now())

	return s.next.Teardown(ns)
}

func (s *instrumentService) track(
	method string,
	namespace string,
	begin time.Time,
	err error,
) {
	var (
		m = kitmetrics.Field{
			Key:   metrics.FieldMethod,
			Value: method,
		}
		n = kitmetrics.Field{
			Key:   metrics.FieldNamespace,
			Value: namespace,
		}
		service = kitmetrics.Field{
			Key:   metrics.FieldService,
			Value: serviceName,
		}
		store = kitmetrics.Field{
			Key:   metrics.FieldStore,
			Value: s.store,
		}
	)

	if err != nil {
		s.errCount.With(m).With(n).With(service).With(store).Add(1)
	}

	s.opCount.With(m).With(n).With(service).With(store).Add(1)

	s.opLatency.With(prometheus.Labels{
		metrics.FieldMethod:    method,
		metrics.FieldNamespace: namespace,
		metrics.FieldService:   serviceName,
		metrics.FieldStore:     s.store,
	}).Observe(time.Since(begin).Seconds())
}
