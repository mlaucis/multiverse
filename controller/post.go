package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	v04_core "github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const typePost = "tg_post"

var defaultOwned = true

// Post is the intermediate representation for posts.
type Post struct {
	IsLiked bool

	*object.Object
}

// PostList is a collection of Post.
type PostList []*Post

// PostMap is the user collection indexed by their ids.
type PostMap map[uint64]*Post

func (ps PostList) toMap() PostMap {
	pm := PostMap{}

	for _, post := range ps {
		pm[post.ID] = post
	}

	return pm
}

// OwnerIDs extracts the OwnerID of every post.
func (ps PostList) OwnerIDs() []uint64 {
	ids := []uint64{}

	for _, p := range ps {
		ids = append(ids, p.OwnerID)
	}

	return ids
}

func postsFromObjects(os object.List) PostList {
	ps := PostList{}

	for _, o := range os {
		ps = append(ps, &Post{Object: o})
	}

	return ps
}

// PostController bundles the business constraints for posts.
type PostController struct {
	connections connection.StrangleService
	events      event.StrangleService
	objects     object.Service
}

// NewPostController returns a controller instance.
func NewPostController(
	connections connection.StrangleService,
	events event.StrangleService,
	objects object.Service,
) *PostController {
	return &PostController{
		connections: connections,
		events:      events,
		objects:     objects,
	}
}

// Create associates the given Object with the owner and adds default type to it
// and stores it in the Object service.
func (c *PostController) Create(
	app *v04_entity.Application,
	post *Post,
	owner *v04_entity.ApplicationUser,
) (*Post, error) {
	post.OwnerID = owner.ID
	post.Owned = defaultOwned
	post.Type = typePost

	o, err := c.objects.Put(app.Namespace(), post.Object)
	if err != nil {
		return nil, err
	}

	return &Post{Object: o}, nil
}

// Delete marks a Post as deleted and updates it in the service.
func (c *PostController) Delete(app *v04_entity.Application, id uint64) error {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &id,
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
	})
	if err != nil {
		return err
	}

	// A delete should be idempotent and always succeed.
	if len(os) == 0 {
		return nil
	}

	o := os[0]
	o.Deleted = true

	_, err = c.objects.Put(app.Namespace(), o)
	if err != nil {
		return err
	}

	return nil
}

// ListAll returns all objects which are of type post.
func (c *PostController) ListAll(
	app *v04_entity.Application,
	user *v04_entity.ApplicationUser,
) (PostList, error) {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
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

	ps := postsFromObjects(os)

	err = enrichIsLiked(c.events, app, user.ID, ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

// ListUser returns all posts for the given user id as visible by the
// connection user id.
func (c *PostController) ListUser(
	app *v04_entity.Application,
	connectionID uint64,
	userID uint64,
) (PostList, error) {
	vs := []object.Visibility{
		object.VisibilityPublic,
		object.VisibilityGlobal,
	}

	// Check relation and include connection visibility.
	if connectionID != userID {
		r, errs := c.connections.Relation(app.OrgID, app.ID, connectionID, userID)
		if errs != nil {
			return nil, errs[0]
		}

		if (r.IsFriend != nil && *r.IsFriend) || (r.IsFollower != nil && *r.IsFollower) {
			vs = append(vs, object.VisibilityConnection)
		}
	}

	// We want all visibilities if the connection and target are the same.
	if connectionID == userID {
		vs = append(vs, object.VisibilityConnection, object.VisibilityPrivate)
	}

	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		OwnerIDs: []uint64{
			userID,
		},
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
		Visibilities: vs,
	})
	if err != nil {
		return nil, err
	}

	ps := postsFromObjects(os)

	err = enrichIsLiked(c.events, app, connectionID, ps)
	if err != nil {
		return nil, err
	}

	return ps, nil
}

// Retrieve returns the Post for the given id.
func (c *PostController) Retrieve(
	app *v04_entity.Application,
	id uint64,
) (*Post, error) {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &id,
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

	if len(os) != 1 {
		return nil, ErrNotFound
	}

	return &Post{Object: os[0]}, nil
}

// Update  stores the new post with the service.
func (c *PostController) Update(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	id uint64,
	post *Post,
) (*Post, error) {
	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &id,
		OwnerIDs: []uint64{
			owner.ID,
		},
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

	// Preserve information.
	p := ps[0]
	p.Attachments = post.Attachments
	p.Tags = post.Tags
	p.Visibility = post.Visibility

	o, err := c.objects.Put(app.Namespace(), p)
	if err != nil {
		return nil, err
	}

	return &Post{Object: o}, nil
}

func enrichIsLiked(
	events event.StrangleService,
	app *v04_entity.Application,
	userID uint64,
	ps PostList,
) error {
	for _, p := range ps {
		es, errs := events.ListAll(app.OrgID, app.ID, v04_core.EventCondition{
			ObjectID: &v04_core.RequestCondition{
				Eq: p.ID,
			},
			Type: &v04_core.RequestCondition{
				Eq: typeLike,
			},
			UserID: &v04_core.RequestCondition{
				Eq: userID,
			},
		})
		if errs != nil {
			return errs[0]
		}

		if len(es) == 1 {
			p.IsLiked = true
		}
	}

	return nil
}
