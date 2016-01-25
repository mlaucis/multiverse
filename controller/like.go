package controller

import (
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_core "github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	typeExternal = "tg_external"
	typeLike     = "tg_like"
)

// LikeFeed is a collection of likes with their referenced users.
type LikeFeed struct {
	Likes   event.List
	UserMap user.Map
}

// LikeController bundles the business constraints for likes on posts.
type LikeController struct {
	events event.StrangleService
	posts  object.Service
	users  user.StrangleService
}

// NewLikeController returns a controller instance.
func NewLikeController(
	events event.StrangleService,
	posts object.Service,
	users user.StrangleService,
) *LikeController {
	return &LikeController{
		events: events,
		posts:  posts,
		users:  users,
	}
}

// Create checks if a like for the owner on the post exists and if not creates
// a new event for it.
func (c *LikeController) Create(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	postID uint64,
) (*v04_entity.Event, error) {
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

	post := ps[0]

	es, errs := c.events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
		ObjectID: &v04_core.RequestCondition{
			Eq: postID,
		},
		Owned: &v04_core.RequestCondition{
			Eq: true,
		},
		Type: &v04_core.RequestCondition{
			Eq: typeLike,
		},
		UserID: &v04_core.RequestCondition{
			Eq: owner.ID,
		},
	})
	if errs != nil {
		return nil, errs[0]
	}

	if len(es) == 1 {
		return es[0], nil
	}

	ev := &v04_entity.Event{
		ObjectID:   postID,
		Owned:      true,
		Type:       typeLike,
		UserID:     owner.ID,
		Visibility: uint8(post.Visibility),
	}

	errs = c.events.Create(app.OrgID, app.ID, owner.ID, ev)
	if errs != nil {
		return nil, errs[0]
	}

	return ev, nil
}

// Delete removes an existing like event for the given user on the post.
func (c *LikeController) Delete(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
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

	es, errs := c.events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
		ObjectID: &v04_core.RequestCondition{
			Eq: postID,
		},
		Owned: &v04_core.RequestCondition{
			Eq: true,
		},
		Type: &v04_core.RequestCondition{
			Eq: typeLike,
		},
		UserID: &v04_core.RequestCondition{
			Eq: owner.ID,
		},
	})
	if errs != nil {
		return errs[0]
	}

	if len(es) == 0 {
		return nil
	}

	errs = c.events.Delete(app.OrgID, app.ID, owner.ID, es[0].ID)
	if errs != nil {
		return errs[0]
	}

	return nil
}

// List returns all likes for the given post.
func (c *LikeController) List(
	app *v04_entity.Application,
	postID uint64,
) (*LikeFeed, error) {
	ps, err := c.posts.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
		Visibilities: []object.Visibility{
			object.VisibilityPublic,
			object.VisibilityGlobal,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) != 1 {
		return nil, ErrNotFound
	}

	es, errs := c.events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
		ObjectID: &v04_core.RequestCondition{
			Eq: postID,
		},
		Owned: &v04_core.RequestCondition{
			Eq: true,
		},
		Type: &v04_core.RequestCondition{
			Eq: typeLike,
		},
	})
	if errs != nil {
		return nil, errs[0]
	}

	um, err := user.MapFromIDs(c.users, app, event.List(es).UserIDs()...)
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
) (*v04_entity.Event, error) {
	es, errs := c.events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
		Object: &v04_core.ObjectCondition{
			ID: &v04_core.RequestCondition{
				Eq: externalID,
			},
			Type: &v04_core.RequestCondition{
				Eq: typeExternal,
			},
		},
		Owned: &v04_core.RequestCondition{
			Eq: true,
		},
		Type: &v04_core.RequestCondition{
			Eq: typeLike,
		},
		UserID: &v04_core.RequestCondition{
			Eq: owner.ID,
		},
	})
	if errs != nil {
		return nil, errs[0]
	}

	if len(es) == 1 {
		return es[0], nil
	}

	ev := &v04_entity.Event{
		Object: &v04_entity.Object{
			ID:   externalID,
			Type: typeExternal,
		},
		Owned:      true,
		Type:       typeLike,
		UserID:     owner.ID,
		Visibility: v04_entity.EventConnections,
	}

	errs = c.events.Create(app.OrgID, app.ID, owner.ID, ev)
	if errs != nil {
		return nil, errs[0]
	}

	return ev, nil
}

// ExternalDelete removes an existing like event for the given user on the
// external entity.
func (c *LikeController) ExternalDelete(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
) error {
	es, errs := c.events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
		Object: &v04_core.ObjectCondition{
			ID: &v04_core.RequestCondition{
				Eq: externalID,
			},
			Type: &v04_core.RequestCondition{
				Eq: typeExternal,
			},
		},
		Owned: &v04_core.RequestCondition{
			Eq: true,
		},
		Type: &v04_core.RequestCondition{
			Eq: typeLike,
		},
		UserID: &v04_core.RequestCondition{
			Eq: owner.ID,
		},
	})
	if errs != nil {
		return errs[0]
	}

	if len(es) == 0 {
		return nil
	}

	errs = c.events.Delete(app.OrgID, app.ID, owner.ID, es[0].ID)
	if errs != nil {
		return errs[0]
	}

	return nil
}

// ExternalList returns all likes for the external entity.
func (c *LikeController) ExternalList(
	app *v04_entity.Application,
	externalID string,
) (*LikeFeed, error) {
	es, errs := c.events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
		Object: &v04_core.ObjectCondition{
			ID: &v04_core.RequestCondition{
				Eq: externalID,
			},
			Type: &v04_core.RequestCondition{
				Eq: typeExternal,
			},
		},
		Owned: &v04_core.RequestCondition{
			Eq: true,
		},
		Type: &v04_core.RequestCondition{
			Eq: typeLike,
		},
	})
	if errs != nil {
		return nil, errs[0]
	}

	um, err := user.MapFromIDs(c.users, app, event.List(es).UserIDs()...)
	if err != nil {
		return nil, err
	}

	return &LikeFeed{
		Likes:   es,
		UserMap: um,
	}, nil
}
