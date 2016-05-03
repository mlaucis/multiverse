package object

import (
	"strconv"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tapglue/multiverse/platform/metrics"
)

type logService struct {
	logger log.Logger
	next   Service
}

// LogMiddleware given a Logger wraps the next Service with logging
// capabilities.
func LogMiddleware(logger log.Logger, store string) ServiceMiddleware {
	return func(next Service) Service {
		logger = log.NewContext(logger).With(
			"service", "object",
			"store", store,
		)

		return &logService{logger: logger, next: next}
	}
}

func (s *logService) Count(ns string, opts QueryOptions) (count int, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"count", count,
			"duration_ns", time.Since(begin).Nanoseconds(),
			"method", "Query",
			"namespace", ns,
			"opts", opts,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.next.Count(ns, opts)
}

func (s *logService) CreatedByDay(
	ns string,
	start, end time.Time,
) (ts metrics.Timeseries, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"datapoints", len(ts),
			"duration_ns", time.Since(begin).Nanoseconds(),
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

	return s.next.CreatedByDay(ns, start, end)
}

func (s *logService) Put(ns string, input *Object) (output *Object, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"input", input,
			"method", "Put",
			"namespace", ns,
			"output", output,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.next.Put(ns, input)
}

func (s *logService) Query(ns string, opts QueryOptions) (os List, err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"method", "Query",
			"namespace", ns,
			"opts", opts,
			"size", len(os),
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.next.Query(ns, opts)
}

func (s *logService) Remove(ns string, id uint64) (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"id", strconv.FormatUint(id, 10),
			"method", "Remove",
			"namespace", ns,
			"id", id,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.next.Remove(ns, id)
}

func (s *logService) Setup(ns string) (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"method", "Setup",
			"namespace", ns,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.next.Setup(ns)
}

func (s *logService) Teardown(ns string) (err error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"method", "teardown",
			"namespace", ns,
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.next.Teardown(ns)
}
