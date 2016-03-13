package org

import (
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const serviceName = "org"

type instrumentStrangleService struct {
	StrangleService

	errCount  kitmetrics.Counter
	opCount   kitmetrics.Counter
	opLatency *prometheus.HistogramVec
	store     string
}

// InstrumentStrangleMiddleware observes key aspects of Service operations and
// exposes Prometheus metrics.
func InstrumentStrangleMiddleware(
	store string,
	errCount kitmetrics.Counter,
	opCount kitmetrics.Counter,
	opLatency *prometheus.HistogramVec,
) StrangleMiddleware {
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

func (s *instrumentStrangleService) FindByKey(
	key string,
) (org *v04_entity.Organization, errs []errors.Error) {
	defer func(begin time.Time) {
		var (
			m = kitmetrics.Field{
				Key:   metrics.FieldMethod,
				Value: "FindByKey",
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

		if errs != nil {
			s.errCount.With(m).With(service).With(store).Add(1)
		}

		s.opCount.With(m).With(service).With(store).Add(1)

		s.opLatency.With(prometheus.Labels{
			metrics.FieldMethod:    "FindByKey",
			metrics.FieldNamespace: "",
			metrics.FieldService:   serviceName,
			metrics.FieldStore:     s.store,
		}).Observe(time.Since(begin).Seconds())
	}(time.Now())

	return s.StrangleService.FindByKey(key)
}
