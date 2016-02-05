package event

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	"github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type logService struct {
	Service

	logger log.Logger
}

// LogMiddleware given a Logger wraps the next Service with logging capabilities.
func LogMiddleware(logger log.Logger, store string) ServiceMiddleware {
	return func(next Service) Service {
		logger = log.NewContext(logger).With(
			"service", "event",
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
			"end", end.Format(metrics.BucketFormat),
			"method", "CreatedByDay",
			"namespace", ns,
			"start", start.Format(metrics.BucketFormat),
		}

		if err != nil {
			ps = append(ps, "err", err)
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.Service.CreatedByDay(ns, start, end)
}

type logStrangleService struct {
	StrangleService

	logger log.Logger
}

// LogStrangleMiddleware given a Logger wraps the next StrangleService with
// logging capabilities.
func LogStrangleMiddleware(logger log.Logger, store string) StrangleMiddleware {
	return func(next StrangleService) StrangleService {
		logger = log.NewContext(logger).With(
			"service", "event",
			"store", store,
		)

		return &logStrangleService{next, logger}
	}
}

func (s *logStrangleService) Create(
	orgID, appID int64,
	userID uint64,
	event *v04_entity.Event,
) (errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"event", event,
			"method", "Create",
			"namespace", convertNamespace(orgID, appID),
			"user_id", strconv.FormatUint(userID, 10),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.Create(orgID, appID, userID, event)
}

func (s *logStrangleService) Delete(
	orgID, appID int64,
	userID uint64,
	eventID uint64,
) (errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"event_id", strconv.FormatUint(eventID, 10),
			"method", "Delete",
			"namespace", convertNamespace(orgID, appID),
			"user_id", strconv.FormatUint(userID, 10),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.Delete(orgID, appID, userID, eventID)
}

func (s *logStrangleService) ListAll(
	orgID, appID int64,
	condition core.EventCondition,
) (es []*v04_entity.Event, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"condition", condition,
			"duration", time.Since(begin),
			"method", "ListAll",
			"namespace", convertNamespace(orgID, appID),
			"result_size", strconv.Itoa(len(es)),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())
	return s.StrangleService.ListAll(orgID, appID, condition)
}

func convertNamespace(orgID, appID int64) string {
	return fmt.Sprintf("app_%d_%d", orgID, appID)
}
