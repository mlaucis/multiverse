package event

import (
	"strconv"

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

func (es Events) Len() int {
	return len(es)
}

func (es Events) Less(i, j int) bool {
	return es[i].CreatedAt.After(*es[j].CreatedAt)
}

func (es Events) Swap(i, j int) {
	es[i], es[j] = es[j], es[i]
}

// UserIDs returns UserID for every Event.
func (es Events) UserIDs() []uint64 {
	ids := []uint64{}

	for _, e := range es {
		ids = append(ids, e.UserID)

		// Extract user ids from target as well.
		if e.Target != nil && e.Target.Type == v04_entity.TypeTargetUser {
			id, err := strconv.ParseUint(e.Target.ID.(string), 10, 64)
			if err != nil {
				// We fail silently here for now until we find a way to log this. As the
				// only effect is that we don't add a potential user to the map
				continue
			}

			ids = append(ids, id)
		}
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
