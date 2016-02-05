package org

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
			"service", "org",
			"store", store,
		)

		return &logStrangleService{next, logger}
	}
}

func (s *logStrangleService) FindByKey(
	key string,
) (org *v04_entity.Organization, errs []errors.Error) {
	defer func(begin time.Time) {
		ps := []interface{}{
			"duration", time.Since(begin),
			"key", key,
			"method", "FindByKey",
		}

		if errs != nil {
			ps = append(ps, "err", errs[0])
		}

		if org != nil {
			ps = append(ps, "org_id", strconv.FormatInt(org.ID, 10))
		}

		_ = s.logger.Log(ps...)
	}(time.Now())

	return s.StrangleService.FindByKey(key)
}
