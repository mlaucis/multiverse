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

// Users is a User collection.
type Users []*v04_entity.ApplicationUser

// UsersFromIDs gathers a user collection from the service for the given ids.
func UsersFromIDs(s StrangleService, orgID, appID int64, ids ...uint64) (Users, error) {
	var (
		seen = map[uint64]struct{}{}
		us   = Users{}
	)

	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}

		u, errs := s.Read(orgID, appID, id, false)
		if errs != nil {
			return nil, errs[0]
		}

		us = append(us, u)
	}

	return us, nil
}
