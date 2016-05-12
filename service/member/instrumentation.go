package member

import (
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const serviceName = "member"

type instrumentStrangleService struct {
	StrangleService

	component string
	errCount  kitmetrics.Counter
	opCount   kitmetrics.Counter
	opLatency *prometheus.HistogramVec
	store     string
}

// InstrumentStrangleMiddleware observes key aspects of Service operations and
// exposes Prometheus metrics.
func InstrumentStrangleMiddleware(
	component, store string,
	errCount kitmetrics.Counter,
	opCount kitmetrics.Counter,
	opLatency *prometheus.HistogramVec,
) StrangleMiddleware {
	return func(next StrangleService) StrangleService {
		return &instrumentStrangleService{
			component:       component,
			errCount:        errCount,
			opCount:         opCount,
			opLatency:       opLatency,
			store:           store,
			StrangleService: next,
		}
	}
}

func (s *instrumentStrangleService) FindBySession(
	session string,
) (member *v04_entity.Member, errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}

		s.track(begin, "FindBySession", err)
	}(time.Now())

	return s.StrangleService.FindBySession(session)
}

func (s *instrumentStrangleService) track(
	begin time.Time,
	method string,
	err error,
) {
	var (
		c = kitmetrics.Field{
			Key:   metrics.FieldComponent,
			Value: s.component,
		}
		m = kitmetrics.Field{
			Key:   metrics.FieldMethod,
			Value: method,
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
		s.errCount.With(c).With(m).With(service).With(store).Add(1)
	}

	s.opCount.With(c).With(m).With(service).With(store).Add(1)

	s.opLatency.With(prometheus.Labels{
		metrics.FieldComponent: s.component,
		metrics.FieldMethod:    method,
		metrics.FieldNamespace: "",
		metrics.FieldService:   serviceName,
		metrics.FieldStore:     s.store,
	}).Observe(time.Since(begin).Seconds())
}
