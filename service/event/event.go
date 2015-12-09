package event

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// Events is an Event collection.
type Events []*v04_entity.Event

// IDs returns ID for every Event.
func (es Events) IDs() []uint64 {
	ids := []uint64{}

	for _, e := range es {
		ids = append(ids, e.ID)
	}

	return ids
}

// UserIDs returns UserID for every Event.
func (es Events) UserIDs() []uint64 {
	ids := []uint64{}

	for _, e := range es {
		ids = append(ids, e.UserID)
	}

	return ids
}

// StrangleService is an intermediate interface to understand the dependencies
// of new middlewares and controllers.
type StrangleService interface {
	Create(orgID, appID int64, userID uint64, event *v04_entity.Event) []errors.Error
	Delete(orgID, appID int64, userID, eventID uint64) []errors.Error
	ListAll(
		orgID, appID int64,
		condition core.EventCondition,
	) ([]*v04_entity.Event, []errors.Error)
}

// StrangleMiddleware is a chainable behaviour modifier for StrangleService.
type StrangleMiddleware func(StrangleService) StrangleService
