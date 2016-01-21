package user

import (
	"github.com/tapglue/multiverse/errors"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// TargetType is the identifier used for events targeting a User.
const TargetType = "tg_user"

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	FindBySession(orgID, appID int64, key string) (*v04_entity.ApplicationUser, []errors.Error)
	Read(orgID, appID int64, id uint64, stats bool) (*v04_entity.ApplicationUser, []errors.Error)
	UpdateLastRead(orgID, appID int64, userID uint64) []errors.Error
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService

// Map is a user collection with their id as index.
type Map map[uint64]*v04_entity.ApplicationUser

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

// UsersFromIDs gathers a user collection from the service for the given ids.
func UsersFromIDs(
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
			// FIXME(xla): We can ignore returned errors for as this method is only
			// used to construct user maps in responses and the logging wrapper of the
			// user service takes care of error reporting. Yet it needs proper
			// addressing as it is a dangerous assumption to believe the usage of this
			// method will only be in one context.
			continue
		}

		us = append(us, u)
	}

	return us, nil
}
