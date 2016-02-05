package member

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// StrangleService is an intermediate interface to understand the
// dependencies of new middlewares and controllers.
type StrangleService interface {
	FindBySession(string) (*v04_entity.Member, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for
// MemberStrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
