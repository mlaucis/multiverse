package app

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FindByApplicationToken(token string) (*v04_entity.Application, []errors.Error)
	FindByBackendToken(token string) (*v04_entity.Application, []errors.Error)
	FindByPublicID(publicID string) (*v04_entity.Application, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
