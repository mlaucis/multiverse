package connection

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// List is a collection of Connections.
type List []*v04_entity.Connection

// Service for connection interactions.
type Service interface {
	metrics.BucketByDay
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	ConnectionsByState(
		orgID, appID int64,
		id uint64,
		state v04_entity.ConnectionStateType,
	) ([]*v04_entity.Connection, []errors.Error)
	FriendsAndFollowingIDs(orgID, appID int64, id uint64) ([]uint64, []errors.Error)
	Relation(orgID, appID int64, from, to uint64) (*v04_entity.Relation, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
