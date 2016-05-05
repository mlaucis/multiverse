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
	connections connection.Service
	events      event.Service
	objects     object.Service
	users       user.Service
}

// NewEventController returns a controller instance.
func NewEventController(
	connections connection.Service,
	events event.Service,
	objects object.Service,
	users user.Service,
) *EventController {
	return &EventController{
		connections: connections,
		events:      events,
		objects:     objects,
		users:       users,
	}
}

// Create stores a new event for the origin user.
func (c *EventController) Create(
	app *v04_entity.Application,
	origin Origin,
	input *event.Event,
) (*event.Event, error) {
	input.Enabled = true
	input.UserID = origin.UserID

	err := constrainEventVisibility(origin, input.Visibility)
	if err != nil {
		return nil, err
	}

	event, err := c.events.Put(app.Namespace(), input)
	if err != nil {
		return nil, err
	}

	return event, nil
}

// Delete marks an event as disabled.
func (c *EventController) Delete(
	app *v04_entity.Application,
	originID uint64,
	id uint64,
) error {
	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		Enabled: &defaultEnabled,
		IDs: []uint64{
			id,
		},
		UserIDs: []uint64{
			originID,
		},
	})
	if err != nil {
		return err
	}

	if len(es) == 0 {
		return ErrNotFound
	}

	event := es[0]
	event.Enabled = false

	_, err = c.events.Put(app.Namespace(), event)
	if err != nil {
		return err
	}

	return nil
}

// List returns the events of a user as seen by the origin user.
func (c *EventController) List(
	app *v04_entity.Application,
	originID uint64,
	userID uint64,
	options *event.QueryOptions,
) (*Feed, error) {
	opts := event.QueryOptions{}
	if options != nil {
		opts = *options
	}

	opts.Enabled = &defaultEnabled
	opts.UserIDs = []uint64{
		userID,
	}
	opts.Visibilities = []event.Visibility{
		event.VisibilityGlobal,
		event.VisibilityPublic,
	}

	if originID == userID {
		opts.Visibilities = append(
			opts.Visibilities,
			event.VisibilityConnection,
			event.VisibilityPrivate,
		)
	} else {
		r, err := queryRelation(c.connections, app, originID, userID)
		if err != nil {
			return nil, err
		}

		if r.isFriend || r.isFollowing {
			opts.Visibilities = append(opts.Visibilities, event.VisibilityConnection)
		}
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

	pum, err := user.MapFromIDs(c.users, app.Namespace(), ps.OwnerIDs()...)
	if err != nil {
		return nil, err
	}

	return &Feed{
		Events:  es,
		PostMap: ps.toMap(),
		UserMap: um.Merge(pum),
	}, nil
}

// Update stores an event with new values.
func (c *EventController) Update(
	app *v04_entity.Application,
	origin Origin,
	id uint64,
	input *event.Event,
) (*event.Event, error) {
	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		Enabled: &defaultEnabled,
		IDs: []uint64{
			id,
		},
		UserIDs: []uint64{
			origin.UserID,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(es) == 0 {
		return nil, ErrNotFound
	}

	e := es[0]
	e.Language = input.Language
	e.Object = input.Object
	e.Target = input.Target
	e.Type = input.Type
	e.Visibility = input.Visibility

	err = constrainEventVisibility(origin, e.Visibility)
	if err != nil {
		return nil, err
	}

	event, err := c.events.Put(app.Namespace(), e)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func constrainEventVisibility(
	origin Origin,
	visiblity event.Visibility,
) error {
	if !origin.IsBackend() && visiblity == event.VisibilityGlobal {
		return wrapError(
			ErrUnauthorized,
			"global visibility can only be set by backend integration",
		)
	}
	return nil
}
