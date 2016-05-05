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
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = &payloadPost{}
			tokenType   = tokenTypeFromContext(ctx)
		)

		err := json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		post, err := c.Create(
			currentApp,
			createOrigin(tokenType, currentUser.ID),
			p.post,
		)
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
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		id, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, currentUser.ID, id)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// PostList returns all posts for a user as visible by the current user.
func PostList(c *controller.PostController) Handler {
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

		feed, err := c.ListUser(app, currentUser.ID, userID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPosts{
			posts:   feed.Posts,
			userMap: feed.UserMap,
		})
	}
}

// PostListAll returns all publicly visible posts.
func PostListAll(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		feed, err := c.ListAll(app, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPosts{
			posts:   feed.Posts,
			userMap: feed.UserMap,
		})
	}
}

// PostListMe returns all posts of the current user.
func PostListMe(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		feed, err := c.ListUser(app, currentUser.ID, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPosts{
			posts:   feed.Posts,
			userMap: feed.UserMap,
		})
	}
}

// PostRetrieve returns the requested Post.
func PostRetrieve(c *controller.PostController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		id, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		post, err := c.Retrieve(app, currentUser.ID, id)
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
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadPost{}
			tokenType   = tokenTypeFromContext(ctx)
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

		updated, err := c.Update(
			currentApp,
			createOrigin(tokenType, currentUser.ID),
			id,
			p.post,
		)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadPost{post: updated})
	}
}

type payloadAttachment struct {
	attachment object.Attachment
}

func (p *payloadAttachment) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Content  string          `json:"content"`
		Contents object.Contents `json:"contents"`
		Name     string          `json:"name"`
		Type     string          `json:"type"`
	}{
		Content:  p.attachment.Contents[object.DefaultLanguage],
		Contents: p.attachment.Contents,
		Name:     p.attachment.Name,
		Type:     p.attachment.Type,
	})
}

func (p *payloadAttachment) UnmarshalJSON(raw []byte) error {
	f := struct {
		Content  string          `json:"content"`
		Contents object.Contents `json:"contents"`
		Name     string          `json:"name"`
		Type     string          `json:"type"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	if f.Contents == nil {
		if f.Content == "" {
			return ErrBadRequest
		}

		f.Contents = object.Contents{
			object.DefaultLanguage: f.Content,
		}
	}

	p.attachment = object.Attachment{
		Contents: f.Contents,
		Name:     f.Name,
		Type:     f.Type,
	}

	return nil
}

type payloadPost struct {
	post *controller.Post
}

func (p *payloadPost) MarshalJSON() ([]byte, error) {
	ps := []*payloadAttachment{}

	for _, a := range p.post.Attachments {
		ps = append(ps, &payloadAttachment{attachment: a})
	}

	return json.Marshal(struct {
		Attachments []*payloadAttachment `json:"attachments"`
		Counts      postCounts           `json:"counts"`
		CreatedAt   time.Time            `json:"created_at,omitempty"`
		ID          string               `json:"id"`
		IsLiked     bool                 `json:"is_liked"`
		Tags        []string             `json:"tags,omitempty"`
		UpdatedAt   time.Time            `json:"updated_at,omitempty"`
		UserID      string               `json:"user_id"`
		Visibility  object.Visibility    `json:"visibility"`
	}{
		Attachments: ps,
		Counts: postCounts{
			Comments: p.post.Counts.Comments,
			Likes:    p.post.Counts.Likes,
		},
		CreatedAt:  p.post.CreatedAt,
		ID:         strconv.FormatUint(p.post.ID, 10),
		IsLiked:    p.post.IsLiked,
		Tags:       p.post.Tags,
		UpdatedAt:  p.post.UpdatedAt,
		UserID:     strconv.FormatUint(p.post.OwnerID, 10),
		Visibility: p.post.Visibility,
	})
}

func (p *payloadPost) UnmarshalJSON(raw []byte) error {
	f := struct {
		Attachments []*payloadAttachment `json:"attachments"`
		Tags        []string             `json:"tags,omitempty"`
		Visibility  object.Visibility    `json:"visibility"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	as := []object.Attachment{}

	for _, a := range f.Attachments {
		as = append(as, a.attachment)
	}

	p.post = &controller.Post{Object: &object.Object{}}
	p.post.Attachments = as
	p.post.Tags = f.Tags
	p.post.Visibility = f.Visibility

	return nil
}

type payloadPosts struct {
	posts   controller.PostList
	userMap user.Map
}

func (p *payloadPosts) MarshalJSON() ([]byte, error) {
	ps := []*payloadPost{}

	for _, post := range p.posts {
		ps = append(ps, &payloadPost{post: post})
	}

	return json.Marshal(struct {
		Posts      []*payloadPost  `json:"posts"`
		PostsCount int             `json:"posts_count"`
		UserMap    *payloadUserMap `json:"users"`
		UserCount  int             `json:"users_count"`
	}{
		Posts:      ps,
		PostsCount: len(ps),
		UserMap:    &payloadUserMap{userMap: p.userMap},
		UserCount:  len(p.userMap),
	})
}

type postCounts struct {
	Comments int `json:"comments"`
	Likes    int `json:"likes"`
}

type postFields struct {
	Attachments []object.Attachment `json:"attachments"`
	Counts      postCounts          `json:"counts"`
	CreatedAt   time.Time           `json:"created_at,omitempty"`
	ID          string              `json:"id"`
	IsLiked     bool                `json:"is_liked"`
	Tags        []string            `json:"tags,omitempty"`
	UpdatedAt   time.Time           `json:"updated_at,omitempty"`
	UserID      string              `json:"user_id"`
	Visibility  object.Visibility   `json:"visibility"`
}
