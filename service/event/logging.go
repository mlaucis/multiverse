package event

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
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
			"namespace", namespace(orgID, appID),
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
			"id", strconv.FormatUint(eventID, 10),
			"method", "Delete",
			"namespace", namespace(orgID, appID),
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
			"namespace", namespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())
	return s.StrangleService.ListAll(orgID, appID, condition)
}

func namespace(orgID, appID int64) string {
	return fmt.Sprintf("app_%d_%d", orgID, appID)
}
