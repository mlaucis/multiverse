package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
	v04_core "github.com/tapglue/multiverse/v04/core"
)

// FeedEvents returns the events of the current user driven by the social and
// interest graph.
func FeedEvents(c *controller.FeedController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		opts, err := whereToOpts(r.URL.Query().Get("where"))
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		feed, err := c.Events(app, currentUser.ID, opts)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Events) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedEvents{
			events:  feed.Events,
			postMap: feed.PostMap,
			userMap: feed.UserMap,
		})
	}
}

// FeedNews returns the superset aggregration of events and posts driven by the
// social and interest graph of the current user.
func FeedNews(c *controller.FeedController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		opts, err := whereToOpts(r.URL.Query().Get("where"))
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		feed, err := c.News(app, currentUser.ID, opts)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Events) == 0 && len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedNews{
			events:   feed.Events,
			posts:    feed.Posts,
			postMap:  feed.PostMap,
			userMap:  feed.UserMap,
			lastRead: currentUser.LastRead,
		})
	}
}

// FeedPosts returns the posts of the current user driven by the social and
// interest graph.
func FeedPosts(c *controller.FeedController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		feed, err := c.Posts(app, currentUser.ID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(feed.Posts) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadFeedPosts{
			posts:   feed.Posts,
			userMap: feed.UserMap,
		})
	}
}

type payloadFeedEvents struct {
	events  event.List
	postMap controller.PostMap
	userMap user.Map
}

func (p *payloadFeedEvents) MarshalJSON() ([]byte, error) {
	es := []*payloadEvent{}

	for _, e := range p.events {
		es = append(es, &payloadEvent{event: e})
	}

	pm := map[string]*payloadPost{}

	for id, post := range p.postMap {
		pm[strconv.FormatUint(id, 10)] = &payloadPost{post: post}
	}

	return json.Marshal(struct {
		Events       []*payloadEvent         `json:"events"`
		EventsCount  int                     `json:"events_count"`
		PostMap      map[string]*payloadPost `json:"post_map"`
		PostMapCount int                     `json:"post_map_count"`
		Users        *payloadUserMap         `json:"users"`
		UsersCount   int                     `json:"users_count"`
	}{
		Events:       es,
		EventsCount:  len(es),
		PostMap:      pm,
		PostMapCount: len(pm),
		Users:        &payloadUserMap{userMap: p.userMap},
		UsersCount:   len(p.userMap),
	})
}

type payloadFeedNews struct {
	events   event.List
	posts    controller.PostList
	postMap  controller.PostMap
	userMap  user.Map
	lastRead time.Time
}

func (p *payloadFeedNews) MarshalJSON() ([]byte, error) {
	var (
		es           = []*payloadEvent{}
		unreadEvents = 0
	)

	for _, ev := range p.events {
		es = append(es, &payloadEvent{event: ev})

		if ev.CreatedAt.After(p.lastRead) {
			unreadEvents++
		}
	}

	var (
		ps          = []*payloadPost{}
		unreadPosts = 0
	)

	for _, post := range p.posts {
		ps = append(ps, &payloadPost{post: post})

		if post.CreatedAt.After(p.lastRead) {
			unreadPosts++
		}
	}

	pm := map[string]*payloadPost{}

	for id, post := range p.postMap {
		pm[strconv.FormatUint(id, 10)] = &payloadPost{post: post}
	}

	return json.Marshal(struct {
		Events            []*payloadEvent         `json:"events"`
		EventsCount       int                     `json:"events_count"`
		EventsCountUnread int                     `json:"events_count_unread"`
		Posts             []*payloadPost          `json:"posts"`
		PostsCount        int                     `json:"posts_count"`
		PostsCountUnread  int                     `json:"posts_count_unread"`
		PostMap           map[string]*payloadPost `json:"post_map"`
		PostMapCount      int                     `json:"post_map_count"`
		UserMap           *payloadUserMap         `json:"users"`
		UserCount         int                     `json:"users_count"`
	}{
		Events:            es,
		EventsCount:       len(es),
		EventsCountUnread: unreadEvents,
		Posts:             ps,
		PostsCount:        len(ps),
		PostsCountUnread:  unreadPosts,
		PostMap:           pm,
		PostMapCount:      len(pm),
		UserMap:           &payloadUserMap{userMap: p.userMap},
		UserCount:         len(p.userMap),
	})
}

type payloadFeedPosts struct {
	posts   controller.PostList
	userMap user.Map
}

func (p *payloadFeedPosts) MarshalJSON() ([]byte, error) {
	ps := []*payloadPost{}

	for _, p := range p.posts {
		ps = append(ps, &payloadPost{post: p})
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

func whereToOpts(input string) (*event.QueryOptions, error) {
	cond, errs := v04_core.NewEventFilter(input)
	if errs != nil {
		return nil, errs[0]
	}

	if cond == nil {
		return nil, nil
	}

	opts := event.QueryOptions{}

	if cond.Object != nil && cond.Object.ID != nil {
		if cond.Object.ID.Eq != nil {
			id, err := parseID(cond.Object.ID.Eq)
			if err != nil {
				return nil, err
			}

			opts.ExternalObjectIDs = []string{
				id,
			}
		}

		if cond.Object.ID.In != nil {
			opts.ExternalObjectIDs = []string{}

			for _, input := range cond.Object.ID.In {
				id, err := parseID(input)
				if err != nil {
					return nil, err
				}

				opts.ExternalObjectIDs = append(opts.ExternalObjectIDs, id)
			}
		}
	}

	if cond.Object != nil && cond.Object.Type != nil {
		if cond.Object.Type.Eq != nil {
			t, ok := cond.Object.Type.Eq.(string)
			if !ok {
				return nil, fmt.Errorf("error in where param")
			}

			opts.ExternalObjectTypes = []string{
				t,
			}
		}

		if cond.Object.Type.In != nil {
			opts.ExternalObjectTypes = []string{}

			for _, input := range cond.Object.Type.In {
				t, ok := input.(string)
				if !ok {
					return nil, fmt.Errorf("error in where param")
				}

				opts.ExternalObjectTypes = append(opts.ExternalObjectTypes, t)
			}
		}
	}

	if cond.Type != nil {
		if cond.Type.Eq != nil {
			t, ok := cond.Type.Eq.(string)
			if !ok {
				return nil, fmt.Errorf("error in where param")
			}

			opts.Types = []string{
				t,
			}
		}

		if cond.Type.In != nil {
			opts.Types = []string{}

			for _, input := range cond.Type.In {
				t, ok := input.(string)
				if !ok {
					return nil, fmt.Errorf("error in where param")
				}

				opts.Types = append(opts.Types, t)
			}
		}
	}

	return &opts, nil
}
