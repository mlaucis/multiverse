package app

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

// LogStrangleMiddleware given a Logger wraps the next StrangleService with
// logging capabilities.
func LogStrangleMiddleware(logger log.Logger, store string) StrangleMiddleware {
	return func(next StrangleService) StrangleService {
		logger = log.NewContext(logger).With(
			"service", "app",
			"store", store,
		)

		return &logStrangleService{next, logger}
	}
}

func (s *logStrangleService) FindByApplicationToken(token string) (app *v04_entity.Application, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"method", "FindByApplicationToken",
			"token", token,
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FindByApplicationToken(token)
}

func (s *logStrangleService) FindByBackendToken(token string) (app *v04_entity.Application, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration_ns", time.Since(begin).Nanoseconds(),
			"method", "FindByBackendToken",
			"token", token,
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FindByBackendToken(token)
}
