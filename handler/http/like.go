package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
)

// LikeCreate emits new like event for the post by the current user.
func LikeCreate(c *controller.LikeController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		postID, err := extractPostID(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		like, err := c.Create(app, currentUser.ID, postID)
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
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		postID, err := extractPostID(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		err = c.Delete(app, currentUser.ID, postID)
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
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		postID, err := extractPostID(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		opts, err := extractLikeOpts(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		opts.Before, err = extractTimeCursorBefore(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		opts.Limit, err = extractLimit(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		feed, err := c.List(app, currentUser.ID, postID, opts)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Likes) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadLikes{
			likes: feed.Likes,
			pagination: pagination(
				r,
				opts.Limit,
				eventCursorAfter(feed.Likes, opts.Limit),
				eventCursorBefore(feed.Likes, opts.Limit),
			),
			userMap: feed.UserMap,
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
	likes      event.List
	pagination *payloadPagination
	userMap    user.Map
}

func (p *payloadLikes) MarshalJSON() ([]byte, error) {
	ls := []*payloadLike{}

	for _, like := range p.likes {
		ls = append(ls, &payloadLike{like: like})
	}

	return json.Marshal(struct {
		Likes      []*payloadLike     `json:"likes"`
		LikesCount int                `json:"likes_count"`
		Pagination *payloadPagination `json:"paging"`
		UserMap    *payloadUserMap    `json:"users"`
		UserCount  int                `json:"users_count"`
	}{
		Likes:      ls,
		LikesCount: len(ls),
		Pagination: p.pagination,
		UserMap:    &payloadUserMap{userMap: p.userMap},
		UserCount:  len(p.userMap),
	})
}
