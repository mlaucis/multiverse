package controller

import (
	"github.com/tapglue/multiverse/service/object"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

const (
	attachmentContent = "content"
	typeComment       = "tg_comment"
)

// CommentController bundles the business constraints for comemnts on posts.
type CommentController struct {
	objects object.Service
}

// NewCommentController returns a controller instance.
func NewCommentController(
	objects object.Service,
) *CommentController {
	return &CommentController{
		objects: objects,
	}
}

// Create associates the given Comment with the Post passed by id.
func (c *CommentController) Create(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	postID uint64,
	content string,
) (*object.Object, error) {
	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, ErrNotFound
	}

	return c.objects.Put(app.Namespace(), &object.Object{
		Attachments: []object.Attachment{
			object.NewTextAttachment(attachmentContent, content),
		},
		ObjectID:   postID,
		OwnerID:    owner.ID,
		Owned:      true,
		Type:       typeComment,
		Visibility: ps[0].Visibility,
	})
}

// Delete flags the Comment as deleted.
func (c *CommentController) Delete(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	postID uint64,
	commentID uint64,
) error {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			owner.ID,
		},
		Types: []string{
			typeComment,
		},
		Owned: &defaultOwned,
	})
	if err != nil {
		return err
	}

	// A delete should be idempotent and always succeed.
	if len(cs) != 1 {
		return nil
	}

	cs[0].Deleted = true

	_, err = c.objects.Put(app.Namespace(), cs[0])
	if err != nil {
		return err
	}

	return nil
}

// List returns all comemnts for the given post id.
func (c *CommentController) List(
	app *v04_entity.Application,
	postID uint64,
) (object.Objects, error) {
	return c.objects.Query(app.Namespace(), object.QueryOptions{
		ObjectIDs: []uint64{
			postID,
		},
		Types: []string{
			typeComment,
		},
		Owned: &defaultOwned,
	})
}

// Retrieve returns the comment given id.
func (c *CommentController) Retrieve(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	postID, commentID uint64,
) (*object.Object, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			owner.ID,
		},
		Types: []string{
			typeComment,
		},
		Owned: &defaultOwned,
	})
	if err != nil {
		return nil, err
	}

	if len(cs) != 1 {
		return nil, ErrNotFound
	}

	return cs[0], nil
}

// Update replaces the given comment with new values.
func (c *CommentController) Update(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	postID, commentID uint64,
	content string,
) (*object.Object, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			owner.ID,
		},
		Owned: &defaultOwned,
		Types: []string{
			typeComment,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(cs) != 1 {
		return nil, ErrNotFound
	}

	cs[0].Attachments = []object.Attachment{
		object.NewTextAttachment(attachmentContent, content),
	}

	return c.objects.Put(app.Namespace(), cs[0])
}
