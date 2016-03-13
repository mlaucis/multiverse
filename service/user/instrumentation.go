package user

import (
	"time"

	kitmetrics "github.com/go-kit/kit/metrics"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const serviceName = "user"

type instrumentService struct {
	Service

	errCount  kitmetrics.Counter
	opCount   kitmetrics.Counter
	opLatency *prometheus.HistogramVec
	store     string
}

// InstrumentMiddleware observes key apsects of Service operations and exposes
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
			opCount:   opCount,
			opLatency: opLatency,
			Service:   next,
			store:     store,
		}
	}
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

func (s *instrumentService) track(
	method, namespace string,
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

type instrumentStrangleService struct {
	StrangleService

	errCount  kitmetrics.Counter
	opCount   kitmetrics.Counter
	opLatency *prometheus.HistogramVec
	store     string
}

// InstrumentStrangleMiddleware observes key aspects of Service operations and exposes
// Prometheus metrics.
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

func (s *instrumentStrangleService) FilterByEmail(
	orgID, appID int64,
	emails []string,
) (users []*v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "FilterByEmail", err)
	}(time.Now())

	return s.StrangleService.FilterByEmail(orgID, appID, emails)
}

func (s *instrumentStrangleService) FilterBySocialIDs(
	orgID, appID int64,
	platform string,
	ids []string,
) (users []*v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "FilterBySocialIDs", err)
	}(time.Now())

	return s.StrangleService.FilterBySocialIDs(orgID, appID, platform, ids)
}

func (s *instrumentStrangleService) FindBySession(
	orgID, appID int64,
	key string,
) (uesr *v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "FindBySession", err)
	}(time.Now())

	return s.StrangleService.FindBySession(orgID, appID, key)
}

func (s *instrumentStrangleService) Read(
	orgID, appID int64,
	id uint64,
	stats bool,
) (user *v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "Read", err)
	}(time.Now())

	return s.StrangleService.Read(orgID, appID, id, stats)
}

func (s *instrumentStrangleService) UpdateLastRead(
	orgID, appID int64,
	id uint64,
) (errs []errors.Error) {
	defer func(begin time.Time) {
		var err error
		if errs != nil {
			err = errs[0]
		}
		s.track(orgID, appID, begin, "UpdateLastRead", err)
	}(time.Now())
	return s.StrangleService.UpdateLastRead(orgID, appID, id)
}

func (s *instrumentStrangleService) track(
	orgID, appID int64,
	begin time.Time,
	method string,
	err error,
) {
	var (
		namespace = convertNamespace(orgID, appID)
		m         = kitmetrics.Field{
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
