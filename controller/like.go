package controller

import (
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	typeLike = "tg_like"
)

// LikeController bundles the business constraints for likes on posts.
type LikeController struct {
	events event.StrangleService
	posts  object.Service
}

// NewLikeController returns a controller instance.
func NewLikeController(
	events event.StrangleService,
	posts object.Service,
) *LikeController {
	return &LikeController{
		events: events,
		posts:  posts,
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

	es, errs := c.events.ListAll(app.OrgID, app.ID, core.EventCondition{
		ObjectID: &core.RequestCondition{
			Eq: postID,
		},
		Owned: &core.RequestCondition{
			Eq: true,
		},
		Type: &core.RequestCondition{
			Eq: typeLike,
		},
		UserID: &core.RequestCondition{
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

	es, errs := c.events.ListAll(app.OrgID, app.ID, core.EventCondition{
		ObjectID: &core.RequestCondition{
			Eq: postID,
		},
		Owned: &core.RequestCondition{
			Eq: true,
		},
		Type: &core.RequestCondition{
			Eq: typeLike,
		},
		UserID: &core.RequestCondition{
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
) (event.List, error) {
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

	es, errs := c.events.ListAll(app.OrgID, app.ID, core.EventCondition{
		ObjectID: &core.RequestCondition{
			Eq: postID,
		},
		Owned: &core.RequestCondition{
			Eq: true,
		},
		Type: &core.RequestCondition{
			Eq: typeLike,
		},
	})
	if errs != nil {
		return nil, errs[0]
	}

	return es, nil
}
