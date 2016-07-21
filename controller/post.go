package controller

import (
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// TypePost identifies an object as a Post.
const TypePost = "tg_post"

var defaultOwned = true

// Post is the intermediate representation for posts.
type Post struct {
	Counts  PostCounts
	IsLiked bool

	*object.Object
}

// PostCounts bundles all connected entity counts.
type PostCounts struct {
	Comments int
	Likes    int
}

// PostFeed is the composite answer for post list methods.
type PostFeed struct {
	Posts   PostList
	UserMap user.Map
}

// PostMap is the user collection indexed by their ids.
type PostMap map[uint64]*Post

// PostList is a collection of Post.
type PostList []*Post

func (ps PostList) toMap() PostMap {
	pm := PostMap{}

	for _, post := range ps {
		pm[post.ID] = post
	}

	return pm
}

func (ps PostList) Len() int {
	return len(ps)
}

func (ps PostList) Less(i, j int) bool {
	return ps[i].CreatedAt.After(ps[j].CreatedAt)
}

func (ps PostList) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
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
	connections connection.Service
	events      event.Service
	objects     object.Service
	users       user.Service
}

// NewPostController returns a controller instance.
func NewPostController(
	connections connection.Service,
	events event.Service,
	objects object.Service,
	users user.Service,
) *PostController {
	return &PostController{
		connections: connections,
		events:      events,
		objects:     objects,
		users:       users,
	}
}

// Create associates the given Object with the owner and adds default type to it
// and stores it in the Object service.
func (c *PostController) Create(
	app *v04_entity.Application,
	origin Origin,
	post *Post,
) (*Post, error) {
	post.OwnerID = origin.UserID
	post.Owned = defaultOwned
	post.Type = TypePost
	// TypePost identifies an object as a Post.

	if err := post.Validate(); err != nil {
		return nil, wrapError(ErrInvalidEntity, "invalid Post: %s", err)
	}

	if err := constrainPostVisibility(origin, post.Visibility); err != nil {
		return nil, err
	}

	if err := post.Object.Validate(); err != nil {
		return nil, wrapError(ErrInvalidEntity, "%s", err)
	}

	o, err := c.objects.Put(app.Namespace(), post.Object)
	if err != nil {
		return nil, err
	}

	return &Post{Object: o}, nil
}

// Delete marks a Post as deleted and updates it in the service.
func (c *PostController) Delete(
	app *v04_entity.Application,
	origin uint64,
	id uint64,
) error {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &id,
		Owned: &defaultOwned,
		Types: []string{
			TypePost,
		},
	})
	if err != nil {
		return err
	}

	// A delete should be idempotent and always succeed.
	if len(os) == 0 {
		return nil
	}

	post := os[0]

	if post.OwnerID != origin {
		return wrapError(ErrUnauthorized, "not allowed to delete post")
	}

	post.Deleted = true

	_, err = c.objects.Put(app.Namespace(), post)
	if err != nil {
		return err
	}

	return nil
}

// ListAll returns all objects which are of type post.
func (c *PostController) ListAll(
	app *v04_entity.Application,
	origin uint64,
) (*PostFeed, error) {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		Owned: &defaultOwned,
		Types: []string{
			TypePost,
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

	err = enrichCounts(c.events, c.objects, app, ps)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin, ps)
	if err != nil {
		return nil, err
	}

	um, err := user.MapFromIDs(c.users, app.Namespace(), ps.OwnerIDs()...)
	if err != nil {
		return nil, err
	}

	return &PostFeed{
		Posts:   ps,
		UserMap: um,
	}, nil
}

// ListUser returns all posts for the given user id as visible by the
// connection user id.
func (c *PostController) ListUser(
	app *v04_entity.Application,
	origin uint64,
	userID uint64,
) (*PostFeed, error) {
	vs := []object.Visibility{
		object.VisibilityPublic,
		object.VisibilityGlobal,
	}

	// Check relation and include connection visibility.
	if origin != userID {
		r, err := queryRelation(c.connections, app, origin, userID)
		if err != nil {
			return nil, err
		}

		if r.isFriend || r.isFollowing {
			vs = append(vs, object.VisibilityConnection)
		}
	}

	// We want all visibilities if the connection and target are the same.
	if origin == userID {
		vs = append(vs, object.VisibilityConnection, object.VisibilityPrivate)
	}

	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		OwnerIDs: []uint64{
			userID,
		},
		Owned: &defaultOwned,
		Types: []string{
			TypePost,
		},
		Visibilities: vs,
	})
	if err != nil {
		return nil, err
	}

	ps := postsFromObjects(os)

	err = enrichCounts(c.events, c.objects, app, ps)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin, ps)
	if err != nil {
		return nil, err
	}

	um, err := user.MapFromIDs(c.users, app.Namespace(), ps.OwnerIDs()...)
	if err != nil {
		return nil, err
	}

	return &PostFeed{
		Posts:   ps,
		UserMap: um,
	}, nil
}

// Retrieve returns the Post for the given id.
func (c *PostController) Retrieve(
	app *v04_entity.Application,
	origin uint64,
	id uint64,
) (*Post, error) {
	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID:    &id,
		Owned: &defaultOwned,
		Types: []string{
			TypePost,
		},
	})
	if err != nil {
		return nil, err
	}

	if len(os) != 1 {
		return nil, ErrNotFound
	}

	if err := isPostVisible(c.connections, app, os[0], origin); err != nil {
		return nil, err
	}

	post := &Post{Object: os[0]}

	err = enrichCounts(c.events, c.objects, app, PostList{post})
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin, PostList{post})
	if err != nil {
		return nil, err
	}

	return post, nil
}

// Update stores a post with the new values.
func (c *PostController) Update(
	app *v04_entity.Application,
	origin Origin,
	id uint64,
	post *Post,
) (*Post, error) {
	ps, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		ID: &id,
		OwnerIDs: []uint64{
			origin.UserID,
		},
		Owned: &defaultOwned,
		Types: []string{
			TypePost,
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

	err = constrainPostVisibility(origin, p.Visibility)
	if err != nil {
		return nil, err
	}

	if err := p.Validate(); err != nil {
		return nil, wrapError(ErrInvalidEntity, "%s", err)
	}

	o, err := c.objects.Put(app.Namespace(), p)
	if err != nil {
		return nil, err
	}

	return &Post{Object: o}, nil
}

func constrainPostVisibility(origin Origin, visibility object.Visibility) error {
	if !origin.IsBackend() && visibility == object.VisibilityGlobal {
		return wrapError(
			ErrUnauthorized,
			"global visibility can only set by backend integration",
		)
	}

	return nil
}

func enrichCounts(
	events event.Service,
	objects object.Service,
	app *v04_entity.Application,
	ps PostList,
) error {
	for _, p := range ps {
		comments, err := objects.Count(app.Namespace(), object.QueryOptions{
			ObjectIDs: []uint64{
				p.ID,
			},
			Types: []string{
				TypeComment,
			},
		})
		if err != nil {
			return err
		}

		likes, err := events.Count(app.Namespace(), event.QueryOptions{
			Enabled: &defaultEnabled,
			ObjectIDs: []uint64{
				p.ID,
			},
			Types: []string{
				TypeLike,
			},
		})
		if err != nil {
			return err
		}

		p.Counts = PostCounts{
			Comments: comments,
			Likes:    likes,
		}
	}

	return nil
}

func enrichIsLiked(
	events event.Service,
	app *v04_entity.Application,
	userID uint64,
	ps PostList,
) error {
	for _, p := range ps {
		es, err := events.Query(app.Namespace(), event.QueryOptions{
			Enabled: &defaultEnabled,
			ObjectIDs: []uint64{
				p.ID,
			},
			Types: []string{
				TypeLike,
			},
			UserIDs: []uint64{
				userID,
			},
		})
		if err != nil {
			return err
		}

		if len(es) == 1 {
			p.IsLiked = true
		}
	}

	return nil
}

// isPostVisible given a post validates that the origin is allowed to see the
// post.
func isPostVisible(
	connections connection.Service,
	app *v04_entity.Application,
	post *object.Object,
	origin uint64,
) error {
	if origin == post.OwnerID {
		return nil
	}

	switch post.Visibility {
	case object.VisibilityGlobal, object.VisibilityPublic:
		return nil
	case object.VisibilityPrivate:
		return ErrNotFound
	}

	r, err := queryRelation(connections, app, origin, post.OwnerID)
	if err != nil {
		return err
	}

	if !r.isFriend && !r.isFollowing {
		return ErrNotFound
	}

	return nil
}
