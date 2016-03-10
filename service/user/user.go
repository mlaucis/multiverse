package user

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/platform/metrics"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	v04_errmsg "github.com/tapglue/multiverse/v04/errmsg"
)

// TargetType is the identifier used for events targeting a User.
const TargetType = "tg_user"

// Service for user interactions.
type Service interface {
	metrics.BucketByDay
}

// ServiceMiddleware is a chainable behaviour modifier for Service.
type ServiceMiddleware func(Service) Service

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FilterByEmail(
		orgID, appID int64,
		emails []string,
	) ([]*v04_entity.ApplicationUser, []errors.Error)
	FilterBySocialIDs(
		orgID, appID int64,
		platform string,
		ids []string,
	) ([]*v04_entity.ApplicationUser, []errors.Error)
	FindBySession(
		orgID, appID int64,
		key string,
	) (*v04_entity.ApplicationUser, []errors.Error)
	Read(
		orgID, appID int64,
		id uint64,
		stats bool,
	) (*v04_entity.ApplicationUser, []errors.Error)
	UpdateLastRead(orgID, appID int64, userID uint64) []errors.Error
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService

// Map is a user collection with their id as index.
type Map map[uint64]*v04_entity.ApplicationUser

// Merge combines two maps.
func (m Map) Merge(x Map) Map {
	for id, user := range x {
		m[id] = user
	}

	return m
}

// List is a collection of users.
type List []*v04_entity.ApplicationUser

// IDs returns the list of user ids.
func (l List) IDs() []uint64 {
	ids := []uint64{}

	for _, user := range l {
		ids = append(ids, user.ID)
	}

	return ids
}

// ToMap turns the user list into a Map.
func (l List) ToMap() Map {
	m := Map{}

	for _, user := range l {
		m[user.ID] = user
	}

	return m
}

// MapFromIDs return a populated user map for the given list of ids.
func MapFromIDs(
	s StrangleService,
	app *v04_entity.Application,
	ids ...uint64,
) (Map, error) {
	um := Map{}

	for _, id := range ids {
		if _, ok := um[id]; ok {
			continue
		}

		user, errs := s.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == v04_errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		um[user.ID] = user
	}

	return um, nil
}

// ListFromIDs gathers a user collection from the service for the given ids.
func ListFromIDs(
	s StrangleService,
	app *v04_entity.Application,
	ids ...uint64,
) (List, error) {
	var (
		seen = map[uint64]struct{}{}
		us   = List{}
	)

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		u, errs := s.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == v04_errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		us = append(us, u)
	}

	return us, nil
}
