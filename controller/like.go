package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	typeExternal = "tg_external"
	typeLike     = "tg_like"
)

var defaultEnabled = true

// LikeFeed is a collection of likes with their referenced users.
type LikeFeed struct {
	Likes   event.List
	UserMap user.StrangleMap
}

// LikeController bundles the business constraints for likes on posts.
type LikeController struct {
	connections connection.Service
	events      event.Service
	posts       object.Service
	users       user.StrangleService
}

// NewLikeController returns a controller instance.
func NewLikeController(
	connections connection.Service,
	events event.Service,
	posts object.Service,
	users user.StrangleService,
) *LikeController {
	return &LikeController{
		connections: connections,
		events:      events,
		posts:       posts,
		users:       users,
	}
}

// Create checks if a like for the owner on the post exists and if not creates
// a new event for it.
func (c *LikeController) Create(
	app *v04_entity.Application,
	origin uint64,
	postID uint64,
) (*event.Event, error) {
	ps, err := c.posts.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, ErrNotFound
	}

	post := ps[0]

	if err := isPostVisible(c.connections, app, post, origin); err != nil {
		return nil, err
	}

	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		ObjectIDs: []uint64{
			postID,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeLike,
		},
		UserIDs: []uint64{
			origin,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(es) > 0 && es[0].Enabled == true {
		return es[0], nil
	}

	var like *event.Event

	if len(es) == 0 {
		like = &event.Event{
			Enabled:    true,
			ObjectID:   postID,
			Owned:      true,
			Type:       typeLike,
			UserID:     origin,
			Visibility: event.Visibility(post.Visibility),
		}
	} else {
		like = es[0]
		like.Enabled = true
	}

	like, err = c.events.Put(app.Namespace(), like)
	if err != nil {
		return nil, err
	}

	return like, nil
}

// Delete removes an existing like event for the given user on the post.
func (c *LikeController) Delete(
	app *v04_entity.Application,
	origin uint64,
	postID uint64,
) error {
	ps, err := c.posts.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
	})
	if err != nil {
		return err
	}

	if len(ps) != 1 {
		return ErrNotFound
	}

	if err := isPostVisible(c.connections, app, ps[0], origin); err != nil {
		return err
	}

	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		Enabled: &defaultEnabled,
		ObjectIDs: []uint64{
			postID,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeLike,
		},
		UserIDs: []uint64{
			origin,
		},
	})
	if err != nil {
		return err
	}

	if len(es) == 0 {
		return nil
	}

	like := es[0]
	like.Enabled = false

	like, err = c.events.Put(app.Namespace(), like)
	if err != nil {
		return err
	}

	return nil
}

// List returns all likes for the given post.
func (c *LikeController) List(
	app *v04_entity.Application,
	origin uint64,
	postID uint64,
) (*LikeFeed, error) {
	ps, err := c.posts.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) != 1 {
		return nil, ErrNotFound
	}

	if err := isPostVisible(c.connections, app, ps[0], origin); err != nil {
		return nil, err
	}

	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		Enabled: &defaultEnabled,
		ObjectIDs: []uint64{
			postID,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeLike,
		},
	})

	um, err := user.StrangleMapFromIDs(c.users, app, es.UserIDs()...)
	if err != nil {
		return nil, err
	}

	return &LikeFeed{
		Likes:   es,
		UserMap: um,
	}, nil
}

// ExternalCreate checks if a like for the owner on the external entity exists
// and if not creates a new event for it.
func (c *LikeController) ExternalCreate(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
) (*event.Event, error) {
	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		ExternalObjectIDs: []string{
			externalID,
		},
		ExternalObjectTypes: []string{
			typeExternal,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeLike,
		},
		UserIDs: []uint64{
			owner.ID,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(es) > 0 && es[0].Enabled == true {
		return es[0], nil
	}

	var like *event.Event

	if len(es) == 0 {
		like = &event.Event{
			Enabled: true,
			Object: &event.Object{
				ID:   externalID,
				Type: typeExternal,
			},
			Owned:      true,
			Type:       typeLike,
			UserID:     owner.ID,
			Visibility: event.VisibilityConnection,
		}
	} else {
		like = es[0]
		like.Enabled = true
	}

	like, err = c.events.Put(app.Namespace(), like)
	if err != nil {
		return nil, err
	}

	return like, nil
}

// ExternalDelete removes an existing like event for the given user on the
// external entity.
func (c *LikeController) ExternalDelete(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
) error {
	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		Enabled: &defaultEnabled,
		ExternalObjectIDs: []string{
			externalID,
		},
		ExternalObjectTypes: []string{
			typeExternal,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeLike,
		},
		UserIDs: []uint64{
			owner.ID,
		},
	})
	if err != nil {
		return err
	}

	if len(es) == 0 {
		return nil
	}

	like := es[0]
	like.Enabled = false

	like, err = c.events.Put(app.Namespace(), like)
	if err != nil {
		return err
	}

	return nil
}

// ExternalList returns all likes for the external entity.
func (c *LikeController) ExternalList(
	app *v04_entity.Application,
	externalID string,
) (*LikeFeed, error) {
	es, err := c.events.Query(app.Namespace(), event.QueryOptions{
		Enabled: &defaultEnabled,
		ExternalObjectIDs: []string{
			externalID,
		},
		ExternalObjectTypes: []string{
			typeExternal,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeLike,
		},
	})
	if err != nil {
		return nil, err
	}

	um, err := user.StrangleMapFromIDs(c.users, app, es.UserIDs()...)
	if err != nil {
		return nil, err
	}

	return &LikeFeed{
		Likes:   es,
		UserMap: um,
	}, nil
}
