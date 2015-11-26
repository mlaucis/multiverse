package user

import (
	"strconv"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

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

func (s *logStrangleService) FindBySession(orgID, appID int64, key string) (user *v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"app", strconv.FormatInt(appID, 10),
			"duration", time.Since(begin),
			"key", key,
			"method", "FindBySession",
			"org", strconv.FormatInt(orgID, 10),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FindBySession(orgID, appID, key)
}

func (s *logStrangleService) Read(orgID, appID int64, id uint64, stats bool) (user *v04_entity.ApplicationUser, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"app", strconv.FormatInt(appID, 10),
			"duration", time.Since(begin),
			"id", strconv.FormatUint(id, 10),
			"method", "Read",
			"org", strconv.FormatInt(orgID, 10),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.Read(orgID, appID, id, stats)
}
