package http

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/server/response"
)

// FeedEvents returns the events of the current user driven by the social and
// interest graph.
func FeedEvents(c *controller.FeedController, users user.StrangleService) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		es, err := c.Events(app, currentUser)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(es) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := user.UsersFromIDs(users, app, es.UserIDs()...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		response.SanitizeApplicationUsers([]*v04_entity.ApplicationUser(us))

		respondJSON(w, http.StatusOK, &payloadFeedEvents{
			events: es,
			users:  us,
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

		es, ps, err := c.News(app, currentUser)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(es) == 0 && len(ps) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := user.UsersFromIDs(
			users,
			app,
			append(es.UserIDs(), ps.OwnerIDs()...)...,
		)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		response.SanitizeApplicationUsers([]*v04_entity.ApplicationUser(us))

		respondJSON(w, http.StatusOK, &payloadFeedNews{
			currentUser: currentUser,
			events:      es,
			posts:       ps,
			users:       us,
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

		ps, err := c.Posts(app, currentUser)
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

		response.SanitizeApplicationUsers([]*v04_entity.ApplicationUser(us))

		respondJSON(w, http.StatusOK, &payloadFeedPosts{posts: ps, users: us})
	}
}

type payloadFeedEvents struct {
	events event.Events
	users  user.Users
}

func (p *payloadFeedEvents) MarshalJSON() ([]byte, error) {
	es := []*v04_entity.PresentationEvent{}

	for _, e := range p.events {
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
		Users:       mapUsers(p.users),
		UsersCount:  len(p.users),
	})
}

type payloadFeedNews struct {
	currentUser *v04_entity.ApplicationUser
	events      event.Events
	posts       object.Objects
	users       user.Users
}

func (p *payloadFeedNews) MarshalJSON() ([]byte, error) {
	var (
		es           = []*v04_entity.PresentationEvent{}
		unreadEvents = 0
	)

	for _, ev := range p.events {
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

	for _, post := range p.posts {
		ps = append(ps, &payloadPost{post: post})

		if p.currentUser.LastRead != nil &&
			post.CreatedAt.After(*p.currentUser.LastRead) {
			unreadPosts++
		}
	}

	return json.Marshal(struct {
		Events            []*v04_entity.PresentationEvent                    `json:"events"`
		EventsCount       int                                                `json:"events_count"`
		EventsCountUnread int                                                `json:"events_count_unread"`
		Posts             []*payloadPost                                     `json:"posts"`
		PostsCount        int                                                `json:"posts_count"`
		PostsCountUnread  int                                                `json:"posts_count_unread"`
		Users             map[string]*v04_entity.PresentationApplicationUser `json:"users"`
		UsersCount        int                                                `json:"users_count"`
	}{
		Events:            es,
		EventsCount:       len(p.events),
		EventsCountUnread: unreadEvents,
		Posts:             ps,
		PostsCount:        len(ps),
		PostsCountUnread:  unreadPosts,
		Users:             mapUsers(p.users),
		UsersCount:        len(p.users),
	})
}

type payloadFeedPosts struct {
	posts object.Objects
	users user.Users
}

func (p *payloadFeedPosts) MarshalJSON() ([]byte, error) {
	ps := []*payloadPost{}

	for _, p := range p.posts {
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
		Users:      mapUsers(p.users),
		UsersCount: len(p.users),
	})
}
