package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// CommentCreate passes the Comment from the payload to the controller.
func CommentCreate(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
			p    = &payloadComment{}
		)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		comment, err := c.Create(app, user, postID, p.content)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadComment{comment: comment})
	}
}

// CommentDelete flags the comment as deleted.
func CommentDelete(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		commentID, err := strconv.ParseUint(mux.Vars(r)["commentID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, user, postID, commentID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// CommentList returns all comments for the given a Post.
func CommentList(
	c *controller.CommentController,
	users user.StrangleService,
) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		cs, err := c.List(app, postID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(cs) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
		}

		us, err := user.UsersFromIDs(users, app, cs.OwnerIDs()...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadComments{
			comments: cs,
			users:    us,
		})
	}
}

// CommentRetrieve return the comment for the requested id.
func CommentRetrieve(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		commentID, err := strconv.ParseUint(mux.Vars(r)["commentID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		comment, err := c.Retrieve(app, user, postID, commentID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadComment{comment: comment})
	}
}

// CommentUpdate replaces the value for a comment with the new values.
func CommentUpdate(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
			p    = &payloadComment{}
		)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		commentID, err := strconv.ParseUint(mux.Vars(r)["commentID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = json.NewDecoder(r.Body).Decode(p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		comment, err := c.Update(app, user, postID, commentID, p.content)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadComment{comment: comment})
	}
}

type commentFields struct {
	Content   string    `json:"content"`
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type payloadComment struct {
	content string
	comment *object.Object
}

func (p *payloadComment) MarshalJSON() ([]byte, error) {
	var (
		c = p.comment
		f = commentFields{
			Content:   c.Attachments[0].Content,
			ID:        strconv.FormatUint(c.ID, 10),
			PostID:    strconv.FormatUint(c.ObjectID, 10),
			UserID:    strconv.FormatUint(c.OwnerID, 10),
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
		}
	)

	return json.Marshal(f)
}

func (p *payloadComment) UnmarshalJSON(raw []byte) error {
	f := commentFields{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.content = f.Content

	return nil
}

type payloadComments struct {
	comments object.Objects
	users    user.List
}

func (p *payloadComments) MarshalJSON() ([]byte, error) {
	cs := []*payloadComment{}

	for _, comment := range p.comments {
		cs = append(cs, &payloadComment{comment: comment})
	}

	return json.Marshal(struct {
		Comments      []*payloadComment                                  `json:"comments"`
		CommentsCount int                                                `json:"comments_count"`
		Users         map[string]*v04_entity.PresentationApplicationUser `json:"users"`
		UsersCount    int                                                `json:"users_count"`
	}{
		Comments:      cs,
		CommentsCount: len(cs),
		Users:         mapUsers(p.users),
		UsersCount:    len(p.users),
	})
}
