package connection

import (
	"strconv"
	"time"

	"github.com/go-kit/kit/log"

	"github.com/tapglue/multiverse/errors"
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

func (s *logStrangleService) FriendsAndFollowingIDs(orgID, appID int64, id uint64) (ids []uint64, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"app", strconv.FormatInt(appID, 10),
			"id", strconv.FormatUint(id, 10),
			"duration", time.Since(begin),
			"method", "FriendsAndFollowingIDs",
			"org", strconv.FormatInt(orgID, 10),
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FriendsAndFollowingIDs(orgID, appID, id)
}
