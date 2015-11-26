package user

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FindBySession(orgID, appID int64, key string) (*v04_entity.ApplicationUser, []errors.Error)
	Read(orgID, appID int64, id uint64, stats bool) (*v04_entity.ApplicationUser, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
