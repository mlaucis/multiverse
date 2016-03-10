package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
)

// LikeCreate emits new like event for the post by the current user.
func LikeCreate(c *controller.LikeController) Handler {
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

		like, err := c.Create(app, user, postID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadLike{like: like})
	}
}

// LikeDelete removes an existing like event for the currentuser on the post.
func LikeDelete(c *controller.LikeController) Handler {
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

		err = c.Delete(app, user, postID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// LikeList returns all Likes for a post.
func LikeList(c *controller.LikeController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		app := appFromContext(ctx)

		postID, err := strconv.ParseUint(mux.Vars(r)["postID"], 10, 64)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		ls, err := c.List(app, postID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ls.Likes) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadLikes{
			likes: ls.Likes,
			users: ls.UserMap,
		})
	}
}

// ExternalLikeCreate emits new like event for the external entity by the
// current user.
func ExternalLikeCreate(c *controller.LikeController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			externalID = mux.Vars(r)["externalID"]
			app        = appFromContext(ctx)
			user       = userFromContext(ctx)
		)

		like, err := c.ExternalCreate(app, user, externalID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadLike{like: like})
	}
}

// ExternalLikeDelete removes an existing like event for the currentuser on the
// external entity.
func ExternalLikeDelete(c *controller.LikeController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			externalID = mux.Vars(r)["externalID"]
			app        = appFromContext(ctx)
			user       = userFromContext(ctx)
		)

		err := c.ExternalDelete(app, user, externalID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// ExternalLikeList returns all Likes for an external entity.
func ExternalLikeList(c *controller.LikeController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			externalID = mux.Vars(r)["externalID"]
			app        = appFromContext(ctx)
		)

		ls, err := c.ExternalList(app, externalID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(ls.Likes) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadLikes{
			likes: ls.Likes,
			users: ls.UserMap,
		})
	}
}

type likeFields struct {
	ID        string    `json:"id"`
	PostID    string    `json:"post_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type payloadLike struct {
	like *event.Event
}

func (p *payloadLike) MarshalJSON() ([]byte, error) {
	var (
		l = p.like
		f = likeFields{
			ID:        strconv.FormatUint(l.ID, 10),
			PostID:    strconv.FormatUint(l.ObjectID, 10),
			UserID:    strconv.FormatUint(l.UserID, 10),
			CreatedAt: l.CreatedAt,
			UpdatedAt: l.UpdatedAt,
		}
	)

	return json.Marshal(f)
}

type payloadLikes struct {
	likes event.List
	users user.Map
}

func (p *payloadLikes) MarshalJSON() ([]byte, error) {
	ls := []*payloadLike{}

	for _, like := range p.likes {
		ls = append(ls, &payloadLike{like: like})
	}

	return json.Marshal(struct {
		Likes      []*payloadLike `json:"likes"`
		LikesCount int            `json:"likes_count"`
		Users      payloadUserMap `json:"users"`
		UsersCount int            `json:"users_count"`
	}{
		Likes:      ls,
		LikesCount: len(ls),
		Users:      mapUserPresentation(p.users),
		UsersCount: len(p.users),
	})
}
