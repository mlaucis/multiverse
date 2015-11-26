package connection

import "github.com/tapglue/multiverse/errors"

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FriendsAndFollowingIDs(orgID, appID int64, id uint64) ([]uint64, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
