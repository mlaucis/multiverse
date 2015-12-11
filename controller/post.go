package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const typePost = "tg_post"

var defaultOwned = true

// PostController bundles the business constraints for posts.
type PostController struct {
	connections connection.StrangleService
	objects     object.Service
}

// NewPostController returns a controller instance.
func NewPostController(
	connections connection.StrangleService,
	objects object.Service,
) *PostController {
	return &PostController{
		connections: connections,
		objects:     objects,
	}
}

// Create associates the given Object with the owner and adds default type to it
// and stores it in the Object service.
func (c *PostController) Create(
	app *v04_entity.Application,
	post *object.Object,
	owner *v04_entity.ApplicationUser,
) (*object.Object, error) {
	post.OwnerID = owner.ID
	post.Owned = defaultOwned
	post.Type = typePost

	return c.objects.Put(app.Namespace(), post)
}

// Delete marks a Post as deleted and updates it in the service.
func (c *PostController) Delete(app *v04_entity.Application, id uint64) error {
	var o *object.Object

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

	o = os[0]
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
) (object.Objects, error) {
	return c.objects.Query(app.Namespace(), object.QueryOptions{
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
		Visibilities: []object.Visibility{
			object.VisibilityPublic,
			object.VisibilityGlobal,
		},
	})
}

// ListUser returns all posts for the given user id.
func (c *PostController) ListUser(
	app *v04_entity.Application,
	userID uint64,
) (object.Objects, error) {
	return c.objects.Query(app.Namespace(), object.QueryOptions{
		OwnerIDs: []uint64{
			userID,
		},
		Owned: &defaultOwned,
		Types: []string{
			typePost,
		},
	})
}

// Retrieve returns the Post for the given id.
func (c *PostController) Retrieve(
	app *v04_entity.Application,
	id uint64,
) (*object.Object, error) {
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

	return os[0], nil
}

// Update  stores the new post with the service.
func (c *PostController) Update(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	id uint64,
	post *object.Object,
) (*object.Object, error) {
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

	return c.objects.Put(app.Namespace(), p)
}
