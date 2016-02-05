package member

import (
	"time"

	"github.com/go-kit/kit/log"
	"github.com/tapglue/multiverse/errors"

	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

type logStrangleService struct {
	StrangleService

	logger log.Logger
}

// LogStrangleMiddleware given a logger wraps the next StrangleService with
// logging capabilities.
func LogStrangleMiddleware(logger log.Logger, store string) StrangleMiddleware {
	return func(next StrangleService) StrangleService {
		logger = log.NewContext(logger).With(
			"service", "member",
			"store", store,
		)

		return &logStrangleService{next, logger}
	}
}

func (s *logStrangleService) FindBySession(
	session string,
) (member *v04_entity.Member, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"method", "FindBySession",
			"session", session,
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FindBySession(session)
}
