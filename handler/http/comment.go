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
)

// CommentCreate passes the Comment from the payload to the controller.
func CommentCreate(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			deviceID    = deviceIDFromContext(ctx)
			p           = &payloadComment{}
			tokenType   = tokenTypeFromContext(ctx)

			origin = createOrigin(deviceID, tokenType, currentUser.ID)
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

		comment, err := c.Create(currentApp, origin, postID, p.comment)
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
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
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

		err = c.Delete(app, currentUser.ID, postID, commentID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// CommentList returns all comments for the given a Post.
func CommentList(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		list, err := c.List(app, currentUser.ID, postID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(list.Comments) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
		}

		respondJSON(w, http.StatusOK, &payloadComments{
			comments: list.Comments,
			userMap:  list.UserMap,
		})
	}
}

// CommentRetrieve return the comment for the requested id.
func CommentRetrieve(c *controller.CommentController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
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

		comment, err := c.Retrieve(app, currentUser.ID, postID, commentID)
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
			currentApp  = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			deviceID    = deviceIDFromContext(ctx)
			p           = &payloadComment{}
			tokenType   = tokenTypeFromContext(ctx)

			origin = createOrigin(deviceID, tokenType, currentUser.ID)
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

		comment, err := c.Update(
			currentApp,
			origin,
			postID,
			commentID,
			p.comment,
		)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadComment{comment: comment})
	}
}

type payloadComment struct {
	contents object.Contents
	comment  *object.Object
}

func (p *payloadComment) MarshalJSON() ([]byte, error) {
	c := p.comment

	return json.Marshal(struct {
		Content   string          `json:"content"`
		Contents  object.Contents `json:"contents"`
		ID        string          `json:"id"`
		PostID    string          `json:"post_id"`
		Private   *object.Private `json:"private,omitempty"`
		UserID    string          `json:"user_id"`
		CreatedAt time.Time       `json:"created_at"`
		UpdatedAt time.Time       `json:"updated_at"`
	}{
		Content:   c.Attachments[0].Contents[object.DefaultLanguage],
		Contents:  c.Attachments[0].Contents,
		ID:        strconv.FormatUint(c.ID, 10),
		PostID:    strconv.FormatUint(c.ObjectID, 10),
		Private:   c.Private,
		UserID:    strconv.FormatUint(c.OwnerID, 10),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	})
}

func (p *payloadComment) UnmarshalJSON(raw []byte) error {
	f := struct {
		Content  string            `json:"content"`
		Contents map[string]string `json:"contents"`
		Private  *object.Private   `json:"private,omitempty"`
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

	p.comment = &object.Object{
		Attachments: []object.Attachment{
			{
				Contents: f.Contents,
			},
		},
		Private: f.Private,
	}

	return nil
}

type payloadComments struct {
	comments object.List
	userMap  user.Map
}

func (p *payloadComments) MarshalJSON() ([]byte, error) {
	cs := []*payloadComment{}

	for _, comment := range p.comments {
		cs = append(cs, &payloadComment{comment: comment})
	}

	return json.Marshal(struct {
		Comments      []*payloadComment `json:"comments"`
		CommentsCount int               `json:"comments_count"`
		UserMap       *payloadUserMap   `json:"users"`
		UsersCount    int               `json:"users_count"`
	}{
		Comments:      cs,
		CommentsCount: len(cs),
		UserMap:       &payloadUserMap{userMap: p.userMap},
		UsersCount:    len(p.userMap),
	})
}
