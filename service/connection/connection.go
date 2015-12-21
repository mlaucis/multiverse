package connection

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FriendsAndFollowingIDs(orgID, appID int64, id uint64) ([]uint64, []errors.Error)
	Relation(orgID, appID int64, from, to uint64) (*v04_entity.Relation, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
