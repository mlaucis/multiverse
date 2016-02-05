package object

import (
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tapglue/multiverse/platform/metrics"
)

type logService struct {
	Service

	logger log.Logger
}

// LogMiddleware given a Logger wraps the next Service with logging
// capabilities.
func LogMiddleware(logger log.Logger, store string) ServiceMiddleware {
	return func(next Service) Service {
		logger = log.NewContext(logger).With(
			"service", "object",
			"store", store,
		)

		return &logService{next, logger}
	}
}

func (s *logService) CreatedByDay(
	ns string,
	start, end time.Time,
) (ts metrics.Timeseries, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"datapoints", len(ts),
			"duration", time.Since(begin),
			"end", end,
			"method", "CreatedByDay",
			"namespace", ns,
			"start", start,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.CreatedByDay(ns, start, end)
}

func (s *logService) Put(ns string, object *Object) (o *Object, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"method", "put",
			"namespace", ns,
			"object", object,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.Put(ns, object)
}

func (s *logService) Query(ns string, opts QueryOptions) (os []*Object, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"method", "query",
			"namespace", ns,
			"opts", opts,
			"size", len(os),
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.Query(ns, opts)
}

func (s *logService) Remove(ns string, id uint64) (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"id", strconv.FormatUint(id, 10),
			"method", "remove",
			"namespace", ns,
			"id", id,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.Remove(ns, id)
}

func (s *logService) Setup(ns string) (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"method", "setup",
			"namespace", ns,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.Setup(ns)
}

func (s *logService) Teardown(ns string) (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"method", "teardown",
			"namespace", ns,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.Teardown(ns)
}
