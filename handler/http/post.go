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
	v04_entity "github.com/tapglue/multiverse/v04/entity"
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

// PostListAll returns all publicly visible posts.
func PostListAll(c *controller.PostController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		ps, err := c.ListAll(app)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := extractUsers(app, ps, users)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, payloadPosts{
			posts: ps,
			users: us,
		})
	}
}

// PostListMe returns all posts of the current user.
func PostListMe(c *controller.PostController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		ps, err := c.ListUser(app, user.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := extractUsers(app, ps, users)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, payloadPosts{
			posts: ps,
			users: us,
		})
	}
}

// PostListMeConnections returns all posts from a users social graph.
func PostListMeConnections(
	c *controller.PostController,
	users user.StrangleService,
) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		ps, err := c.ListUserConnections(app, user.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := extractUsers(app, ps, users)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, payloadPosts{
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

func extractUsers(
	app *v04_entity.Application,
	os []*object.Object,
	users user.StrangleService,
) (map[string]*v04_entity.ApplicationUser, error) {
	us := map[string]*v04_entity.ApplicationUser{}

	for _, o := range os {
		user, errs := users.Read(app.OrgID, app.ID, o.OwnerID, false)
		if errs != nil {
			return nil, errs[0]
		}

		us[strconv.FormatUint(user.ID, 10)] = user
	}

	return us, nil
}

type postFields struct {
	Attachments []object.Attachment `json:"attachments"`
	CreatedAt   time.Time           `json:"created_at,omitempty"`
	ID          string              `json:"id"`
	Tags        []string            `json:"tags,omitempty"`
	UpdatedAt   time.Time           `json:"updated_at,omitempty"`
	UserID      string              `json:"user_id"`
	Visibility  object.Visibility   `json:"visibility"`
}

type payloadPost struct {
	post *object.Object
}

func (p *payloadPost) MarshalJSON() ([]byte, error) {
	var (
		o = p.post
		f = postFields{
			Attachments: o.Attachments,
			CreatedAt:   o.CreatedAt,
			ID:          strconv.FormatUint(o.ID, 10),
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

	p.post = &object.Object{
		Attachments: f.Attachments,
		Tags:        f.Tags,
		Visibility:  f.Visibility,
	}

	return nil
}

type payloadPosts struct {
	posts []*object.Object
	users map[string]*v04_entity.ApplicationUser
}

func (p *payloadPosts) MarshalJSON() ([]byte, error) {
	ps := []*payloadPost{}

	for _, post := range p.posts {
		ps = append(ps, &payloadPost{post: post})
	}

	return json.Marshal(struct {
		Posts      []*payloadPost                         `json:"posts"`
		PostsCount int                                    `json:"posts_count"`
		Users      map[string]*v04_entity.ApplicationUser `json:"users"`
		UsersCount int                                    `json:"users_count"`
	}{
		Posts:      ps,
		PostsCount: len(ps),
		Users:      p.users,
		UsersCount: len(p.users),
	})
}
