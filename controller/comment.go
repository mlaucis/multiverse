package controller

import (
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
)

const (
	// TypeComment identifies a comment object.
	TypeComment = "tg_comment"

	attachmentContent = "content"
)

// CommentFeed is a collection of comments with their referneced users.
type CommentFeed struct {
	Comments object.List
	UserMap  user.Map
}

// CommentController bundles the business constraints for comemnts on posts.
type CommentController struct {
	connections connection.Service
	objects     object.Service
	users       user.Service
}

// NewCommentController returns a controller instance.
func NewCommentController(
	connections connection.Service,
	objects object.Service,
	users user.Service,
) *CommentController {
	return &CommentController{
		connections: connections,
		objects:     objects,
		users:       users,
	}
}

// Create associates the given Comment with the Post passed by id.
func (c *CommentController) Create(
	currentApp *app.App,
	origin Origin,
	postID uint64,
	input *object.Object,
) (*object.Object, error) {
	err := constrainCommentPrivate(origin, input.Private)
	if err != nil {
		return nil, err
	}

	ps, err := c.objects.Query(currentApp.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{TypePost},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, ErrNotFound
	}

	comment := &object.Object{
		Attachments: []object.Attachment{
			object.NewTextAttachment(
				attachmentContent,
				input.Attachments[0].Contents,
			),
		},
		ObjectID:   postID,
		OwnerID:    origin.UserID,
		Owned:      true,
		Private:    input.Private,
		Type:       TypeComment,
		Visibility: ps[0].Visibility,
	}

	if err := comment.Validate(); err != nil {
		return nil, wrapError(ErrInvalidEntity, "invalid Comment: %s", err)
	}

	if err := isPostVisible(c.connections, currentApp, ps[0], origin.UserID); err != nil {
		return nil, err
	}

	return c.objects.Put(currentApp.Namespace(), comment)
}

// Delete flags the Comment as deleted.
func (c *CommentController) Delete(
	currentApp *app.App,
	origin uint64,
	postID uint64,
	commentID uint64,
) error {
	cs, err := c.objects.Query(currentApp.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			origin,
		},
		Types: []string{
			TypeComment,
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

	_, err = c.objects.Put(currentApp.Namespace(), cs[0])
	if err != nil {
		return err
	}

	return nil
}

// List returns all comemnts for the given post id.
func (c *CommentController) List(
	currentApp *app.App,
	origin uint64,
	postID uint64,
) (*CommentFeed, error) {
	ps, err := c.objects.Query(currentApp.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{TypePost},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, ErrNotFound
	}

	if err := isPostVisible(c.connections, currentApp, ps[0], origin); err != nil {
		return nil, err
	}

	cs, err := c.objects.Query(currentApp.Namespace(), object.QueryOptions{
		ObjectIDs: []uint64{
			postID,
		},
		Types: []string{
			TypeComment,
		},
		Owned: &defaultOwned,
	})
	if err != nil {
		return nil, err
	}

	um, err := user.MapFromIDs(c.users, currentApp.Namespace(), cs.OwnerIDs()...)
	if err != nil {
		return nil, err
	}

	return &CommentFeed{Comments: cs, UserMap: um}, nil
}

// Retrieve returns the comment given id.
func (c *CommentController) Retrieve(
	currentApp *app.App,
	origin uint64,
	postID, commentID uint64,
) (*object.Object, error) {
	cs, err := c.objects.Query(currentApp.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			origin,
		},
		Types: []string{
			TypeComment,
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
	currentApp *app.App,
	origin Origin,
	postID, commentID uint64,
	new *object.Object,
) (*object.Object, error) {
	err := constrainCommentPrivate(origin, new.Private)
	if err != nil {
		return nil, err
	}

	cs, err := c.objects.Query(currentApp.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			origin.UserID,
		},
		Owned: &defaultOwned,
		Types: []string{
			TypeComment,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(cs) != 1 {
		return nil, ErrNotFound
	}

	old := cs[0]

	old.Attachments = []object.Attachment{
		object.NewTextAttachment(
			attachmentContent,
			new.Attachments[0].Contents,
		),
	}

	if origin.IsBackend() && new.Private != nil {
		old.Private = new.Private
	}

	return c.objects.Put(currentApp.Namespace(), old)
}

func constrainCommentPrivate(origin Origin, private *object.Private) error {
	if !origin.IsBackend() && private != nil {
		return wrapError(ErrUnauthorized,
			"private can only be set by backend integration",
		)
	}

	return nil
}
