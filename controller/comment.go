package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	v04_errmsg "github.com/tapglue/multiverse/v04/errmsg"
)

const (
	attachmentContent = "content"
	typeComment       = "tg_comment"
)

// CommentFeed is a collection of comments with their referneced users.
type CommentFeed struct {
	Comments object.List
	UserMap  user.StrangleMap
}

// CommentController bundles the business constraints for comemnts on posts.
type CommentController struct {
	connections connection.Service
	objects     object.Service
	users       user.StrangleService
}

// NewCommentController returns a controller instance.
func NewCommentController(
	connections connection.Service,
	objects object.Service,
	users user.StrangleService,
) *CommentController {
	return &CommentController{
		connections: connections,
		objects:     objects,
		users:       users,
	}
}

// Create associates the given Comment with the Post passed by id.
func (c *CommentController) Create(
	app *v04_entity.Application,
	origin uint64,
	postID uint64,
	content string,
) (*object.Object, error) {
	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{typePost},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, ErrNotFound
	}

	if err := isPostVisible(c.connections, app, ps[0], origin); err != nil {
		return nil, err
	}

	return c.objects.Put(app.Namespace(), &object.Object{
		Attachments: []object.Attachment{
			object.NewTextAttachment(attachmentContent, content),
		},
		ObjectID:   postID,
		OwnerID:    origin,
		Owned:      true,
		Type:       typeComment,
		Visibility: ps[0].Visibility,
	})
}

// Delete flags the Comment as deleted.
func (c *CommentController) Delete(
	app *v04_entity.Application,
	origin uint64,
	postID uint64,
	commentID uint64,
) error {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			origin,
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
	origin uint64,
	postID uint64,
) (*CommentFeed, error) {
	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &postID,
		Owned: &defaultOwned,
		Types: []string{typePost},
	})
	if err != nil {
		return nil, err
	}

	if len(ps) == 0 {
		return nil, ErrNotFound
	}

	if err := isPostVisible(c.connections, app, ps[0], origin); err != nil {
		return nil, err
	}

	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ObjectIDs: []uint64{
			postID,
		},
		Types: []string{
			typeComment,
		},
		Owned: &defaultOwned,
	})
	if err != nil {
		return nil, err
	}

	um, err := extractUsersFromObjects(c.users, app, cs)
	if err != nil {
		return nil, err
	}

	return &CommentFeed{Comments: cs, UserMap: um}, nil
}

// Retrieve returns the comment given id.
func (c *CommentController) Retrieve(
	app *v04_entity.Application,
	origin uint64,
	postID, commentID uint64,
) (*object.Object, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			origin,
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
	origin uint64,
	postID, commentID uint64,
	content string,
) (*object.Object, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ObjectIDs: []uint64{
			postID,
		},
		OwnerIDs: []uint64{
			origin,
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

// ExternalCreate stores a comment with the given content associated with the
// external object.
func (c *CommentController) ExternalCreate(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
	content string,
) (*object.Object, error) {
	return c.objects.Put(app.Namespace(), &object.Object{
		Attachments: []object.Attachment{
			object.NewTextAttachment(attachmentContent, content),
		},
		ExternalID: externalID,
		OwnerID:    owner.ID,
		Owned:      true,
		Type:       typeComment,
		Visibility: object.VisibilityPublic,
	})
}

// ExternalDelete flags the Comment as deleted.
func (c *CommentController) ExternalDelete(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
	commentID uint64,
) error {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &commentID,
		ExternalIDs: []string{
			externalID,
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

// ExternalList returns all comemnts for the given external id.
func (c *CommentController) ExternalList(
	app *v04_entity.Application,
	externalID string,
) (*CommentFeed, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ExternalIDs: []string{
			externalID,
		},
		Types: []string{
			typeComment,
		},
		Owned: &defaultOwned,
	})
	if err != nil {
		return nil, err
	}

	um, err := extractUsersFromObjects(c.users, app, cs)
	if err != nil {
		return nil, err
	}

	return &CommentFeed{Comments: cs, UserMap: um}, nil
}

// ExternalRetrieve returns the comment given id.
func (c *CommentController) ExternalRetrieve(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
	commentID uint64,
) (*object.Object, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ExternalIDs: []string{
			externalID,
		},
		ID: &commentID,
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

// ExternalUpdate replaces the given comment with new values.
func (c *CommentController) ExternalUpdate(
	app *v04_entity.Application,
	owner *v04_entity.ApplicationUser,
	externalID string,
	commentID uint64,
	content string,
) (*object.Object, error) {
	cs, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ExternalIDs: []string{
			externalID,
		},
		ID: &commentID,
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

func extractUsersFromObjects(
	users user.StrangleService,
	app *v04_entity.Application,
	os object.List,
) (user.StrangleMap, error) {
	um := user.StrangleMap{}

	for _, id := range os.OwnerIDs() {
		if _, ok := um[id]; ok {
			continue
		}

		user, errs := users.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == v04_errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		um[id] = user
	}

	return um, nil
}
