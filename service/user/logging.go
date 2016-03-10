package user

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type logService struct {
	Service

	logger log.Logger
}

// LogMiddleware gien a Logger wraps the next Service with logging capabilities.
func LogMiddleware(logger log.Logger, store string) ServiceMiddleware {
	return func(next Service) Service {
		logger = log.NewContext(logger).With(
			"service", "user",
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
			"service", "user",
			"store", store,
		)

		return &logStrangleService{next, logger}
	}
}

func (s *logStrangleService) FilterByEmail(
	orgID, appID int64,
	emails []string,
) (users []*v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"datapoints", len(users),
			"duration", time.Since(begin),
			"inputs", len(emails),
			"method", "FitlerByEmail",
			"namespace", convertNamespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FilterByEmail(orgID, appID, emails)
}

func (s *logStrangleService) FilterBySocialIDs(
	orgID, appID int64,
	platform string,
	ids []string,
) (users []*v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"datapoints", len(users),
			"duration", time.Since(begin),
			"inputs", len(ids),
			"method", "FitlerBySocialIDs",
			"namespace", convertNamespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FilterBySocialIDs(orgID, appID, platform, ids)
}

func (s *logStrangleService) FindBySession(
	orgID, appID int64,
	key string,
) (user *v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"key", key,
			"method", "FindBySession",
			"namespace", convertNamespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FindBySession(orgID, appID, key)
}

func (s *logStrangleService) Read(
	orgID, appID int64,
	id uint64,
	stats bool,
) (user *v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"id", strconv.FormatUint(id, 10),
			"method", "Read",
			"namespace", convertNamespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.Read(orgID, appID, id, stats)
}

func (s *logStrangleService) UpdateLastRead(
	orgID, appID int64,
	id uint64,
) (errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"id", strconv.FormatUint(id, 10),
			"method", "UpdateLastRead",
			"namespace", convertNamespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.UpdateLastRead(orgID, appID, id)
}

func convertNamespace(orgID, appID int64) string {
	return fmt.Sprintf("app_%d_%d", orgID, appID)
}
