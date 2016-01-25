package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/user"
	v04_core "github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// FeedEvents returns the events of the current user driven by the social and
// interest graph.
func FeedEvents(c *controller.FeedController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		where, errs := v04_core.NewEventFilter(r.URL.Query().Get("where"))
		if errs != nil {
			respondError(w, 0, wrapError(ErrBadRequest, errs[0].Error()))
			return
		}

		feed, err := c.Events(app, currentUser, where)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Events) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedEvents{
			feed: feed,
		})
	}
}

// FeedNews returns the superset aggregration of events and posts driven by the
// social and interest graph of the current user.
func FeedNews(c *controller.FeedController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		where, errs := v04_core.NewEventFilter(r.URL.Query().Get("where"))
		if errs != nil {
			respondError(w, 0, wrapError(ErrBadRequest, errs[0].Error()))
			return
		}

		feed, err := c.News(app, currentUser, where)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Events) == 0 && len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedNews{
			currentUser: currentUser,
			feed:        feed,
		})
	}
}

// FeedPosts returns the posts of the current user driven by the social and
// interest graph.
func FeedPosts(c *controller.FeedController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		feed, err := c.Posts(app, currentUser)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedPosts{feed: feed})
	}
}

type payloadFeedEvents struct {
	feed *controller.Feed
}

func (p *payloadFeedEvents) MarshalJSON() ([]byte, error) {
	es := []*v04_entity.PresentationEvent{}

	for _, e := range p.feed.Events {
		es = append(es, &v04_entity.PresentationEvent{Event: e})
	}

	return json.Marshal(struct {
		Events      []*v04_entity.PresentationEvent                    `json:"events"`
		EventsCount int                                                `json:"events_count"`
		Users       map[string]*v04_entity.PresentationApplicationUser `json:"users"`
		UsersCount  int                                                `json:"users_count"`
	}{
		Events:      es,
		EventsCount: len(es),
		Users:       mapUserPresentation(p.feed.UserMap),
		UsersCount:  len(p.feed.UserMap),
	})
}

type payloadFeedNews struct {
	currentUser *v04_entity.ApplicationUser
	feed        *controller.Feed
}

func (p *payloadFeedNews) MarshalJSON() ([]byte, error) {
	var (
		es           = []*v04_entity.PresentationEvent{}
		unreadEvents = 0
	)

	for _, ev := range p.feed.Events {
		es = append(es, &v04_entity.PresentationEvent{Event: ev})

		if p.currentUser.LastRead != nil &&
			ev.CreatedAt.After(*p.currentUser.LastRead) {
			unreadEvents++
		}
	}

	var (
		ps          = []*payloadPost{}
		unreadPosts = 0
	)

	for _, post := range p.feed.Posts {
		ps = append(ps, &payloadPost{post: post})

		if p.currentUser.LastRead != nil &&
			post.CreatedAt.After(*p.currentUser.LastRead) {
			unreadPosts++
		}
	}

	pm := map[string]*payloadPost{}

	for id, post := range p.feed.PostMap {
		pm[strconv.FormatUint(id, 10)] = &payloadPost{post: post}
	}

	return json.Marshal(struct {
		Events            []*v04_entity.PresentationEvent                    `json:"events"`
		EventsCount       int                                                `json:"events_count"`
		EventsCountUnread int                                                `json:"events_count_unread"`
		Posts             []*payloadPost                                     `json:"posts"`
		PostsCount        int                                                `json:"posts_count"`
		PostsCountUnread  int                                                `json:"posts_count_unread"`
		PostMap           map[string]*payloadPost                            `json:"post_map"`
		PostMapCount      int                                                `json:"post_map_count"`
		Users             map[string]*v04_entity.PresentationApplicationUser `json:"users"`
		UsersCount        int                                                `json:"users_count"`
	}{
		Events:            es,
		EventsCount:       len(es),
		EventsCountUnread: unreadEvents,
		Posts:             ps,
		PostsCount:        len(ps),
		PostsCountUnread:  unreadPosts,
		PostMap:           pm,
		PostMapCount:      len(pm),
		Users:             mapUserPresentation(p.feed.UserMap),
		UsersCount:        len(p.feed.UserMap),
	})
}

type payloadFeedPosts struct {
	feed *controller.Feed
}

func (p *payloadFeedPosts) MarshalJSON() ([]byte, error) {
	ps := []*payloadPost{}

	for _, p := range p.feed.Posts {
		ps = append(ps, &payloadPost{post: p})
	}

	return json.Marshal(struct {
		Posts      []*payloadPost                                     `json:"posts"`
		PostsCount int                                                `json:"posts_count"`
		Users      map[string]*v04_entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                                `json:"users_count"`
	}{
		Posts:      ps,
		PostsCount: len(ps),
		Users:      mapUserPresentation(p.feed.UserMap),
		UsersCount: len(p.feed.UserMap),
	})
}
