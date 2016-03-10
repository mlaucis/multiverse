package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// EventController bundles the business constraints of Events.
type EventController struct {
	connections connection.StrangleService
	events      event.Service
	objects     object.Service
	users       user.StrangleService
}

// NewEventController returns a controller instance.
func NewEventController(
	connections connection.StrangleService,
	events event.Service,
	objects object.Service,
	users user.StrangleService,
) *EventController {
	return &EventController{
		connections: connections,
		events:      events,
		objects:     objects,
		users:       users,
	}
}

// ListUser returns the events of a user as seen by the origin user.
func (c *EventController) ListUser(
	app *v04_entity.Application,
	originID uint64,
	userID uint64,
) (*Feed, error) {
	var (
		enabled = true
		opts    = event.QueryOptions{
			Enabled: &enabled,
			UserIDs: []uint64{
				userID,
			},
			Visibilities: []event.Visibility{
				event.VisibilityGlobal,
				event.VisibilityPublic,
			},
		}
	)

	r, errs := c.connections.Relation(app.OrgID, app.ID, originID, userID)
	if errs != nil {
		return nil, errs[0]
	}

	if (r.IsFriend != nil && *r.IsFriend) || (r.IsFollowed != nil && *r.IsFollowed) {
		opts.Visibilities = append(opts.Visibilities, event.VisibilityConnection)
	}

	es, err := c.events.Query(app.Namespace(), opts)
	if err != nil {
		return nil, err
	}

	um, err := fillupUsers(c.users, app, originID, user.Map{}, es)
	if err != nil {
		return nil, err
	}

	ps, err := extractPosts(c.objects, app, es)
	if err != nil {
		return nil, err
	}

	pum, err := extractUsersFromPosts(c.users, app, ps)
	if err != nil {
		return nil, err
	}

	return &Feed{
		Events:  es,
		PostMap: ps.toMap(),
		UserMap: um.Merge(pum),
	}, nil
}
