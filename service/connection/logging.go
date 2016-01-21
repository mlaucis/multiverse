package connection

import (
	"fmt"
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
			"service", "connection",
			"store", store,
		)

		return &logStrangleService{next, logger}
	}
}

func (s *logStrangleService) ConnectionsByState(
	orgID, appID int64,
	id uint64,
	state v04_entity.ConnectionStateType,
) (cs []*v04_entity.Connection, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"connections", strconv.Itoa(len(cs)),
			"duration", time.Since(begin),
			"id", strconv.FormatUint(id, 10),
			"method", "ConnectionsByState",
			"namespace", namespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.ConnectionsByState(orgID, appID, id, state)
}

func (s *logStrangleService) FriendsAndFollowingIDs(orgID, appID int64, id uint64) (ids []uint64, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"id", strconv.FormatUint(id, 10),
			"duration", time.Since(begin),
			"method", "FriendsAndFollowingIDs",
			"namespace", namespace(orgID, appID),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FriendsAndFollowingIDs(orgID, appID, id)
}

func (s *logStrangleService) Relation(
	orgID, appID int64,
	from, to uint64,
) (r *v04_entity.Relation, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"from", strconv.FormatUint(from, 10),
			"method", "Relation",
			"namespace", namespace(orgID, appID),
			"relation", r,
			"to", strconv.FormatUint(to, 10),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.Relation(orgID, appID, from, to)
}

func namespace(orgID, appID int64) string {
	return fmt.Sprintf("app_%d_%d", orgID, appID)
}
