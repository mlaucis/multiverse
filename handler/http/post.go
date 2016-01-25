package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
)

// PostCreate passes the Post from the payload to the controller.
func PostCreate(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
			p    = &payloadPost{}
		)

		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		post, err := c.Create(app, p.post, user)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadPost{post: post})
	}
}

// PostDelete flags the Post as deleted.
func PostDelete(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		id, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, id)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// PostList returns all posts for a user as visible by the current user.
func PostList(c *controller.PostController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		userID, err := strconv.ParseUint(mux.Vars(r)["userID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		ps, err := c.ListUser(app, currentUser.ID, userID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := user.UsersFromIDs(users, app, ps.OwnerIDs()...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPosts{
			posts: ps,
			users: us,
		})
	}
}

// PostListAll returns all publicly visible posts.
func PostListAll(c *controller.PostController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		ps, err := c.ListAll(app, currentUser)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := user.UsersFromIDs(users, app, ps.OwnerIDs()...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPosts{
			posts: ps,
			users: us,
		})
	}
}

// PostListMe returns all posts of the current user.
func PostListMe(c *controller.PostController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		ps, err := c.ListUser(app, currentUser.ID, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := user.UsersFromIDs(users, app, ps.OwnerIDs()...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPosts{
			posts: ps,
			users: us,
		})
	}
}

// PostRetrieve returns the requested Post.
func PostRetrieve(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		id, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		post, err := c.Retrieve(app, id)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPost{post: post})
	}
}

// PostUpdate reaplces a post with new values.
func PostUpdate(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
			p    = payloadPost{}
		)

		id, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		updated, err := c.Update(app, user, id, p.post)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPost{post: updated})
	}
}

type postFields struct {
	Attachments []object.Attachment `json:"attachments"`
	CreatedAt   time.Time           `json:"created_at,omitempty"`
	ID          string              `json:"id"`
	IsLiked     bool                `json:"is_liked"`
	Tags        []string            `json:"tags,omitempty"`
	UpdatedAt   time.Time           `json:"updated_at,omitempty"`
	UserID      string              `json:"user_id"`
	Visibility  object.Visibility   `json:"visibility"`
}

type payloadPost struct {
	post *controller.Post
}

func (p *payloadPost) MarshalJSON() ([]byte, error) {
	var (
		o = p.post
		f = postFields{
			Attachments: o.Attachments,
			CreatedAt:   o.CreatedAt,
			ID:          strconv.FormatUint(o.ID, 10),
			IsLiked:     o.IsLiked,
			Tags:        o.Tags,
			UpdatedAt:   o.UpdatedAt,
			UserID:      strconv.FormatUint(o.OwnerID, 10),
			Visibility:  o.Visibility,
		}
	)

	return json.Marshal(f)
}

func (p *payloadPost) UnmarshalJSON(raw []byte) error {
	f := postFields{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.post = &controller.Post{Object: &object.Object{}}
	p.post.Attachments = f.Attachments
	p.post.Tags = f.Tags
	p.post.Visibility = f.Visibility

	return nil
}

type payloadPosts struct {
	posts controller.PostList
	users user.List
}

func (p *payloadPosts) MarshalJSON() ([]byte, error) {
	ps := []*payloadPost{}

	for _, post := range p.posts {
		ps = append(ps, &payloadPost{post: post})
	}

	return json.Marshal(struct {
		Posts      []*payloadPost `json:"posts"`
		PostsCount int            `json:"posts_count"`
		Users      payloadUserMap `json:"users"`
		UsersCount int            `json:"users_count"`
	}{
		Posts:      ps,
		PostsCount: len(ps),
		Users:      mapUsers(p.users),
		UsersCount: len(p.users),
	})
}
