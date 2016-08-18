package controller

import (
	"sort"
	"strconv"
	"time"

	"github.com/tapglue/multiverse/platform/flake"
	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
)

// affiliations is the composite structure to map connections to users.
type affiliations map[*connection.Connection]*user.User

// connections returns only the connections of the affiliations.
func (a affiliations) connections() connection.List {
	cs := connection.List{}

	for con := range a {
		cs = append(cs, con)
	}

	return cs
}

// followers returns follow connections towards the origin.
func (a affiliations) followers(origin uint64) connection.List {
	cs := connection.List{}

	for con := range a {
		if con.Type == connection.TypeFriend {
			continue
		}

		if con.FromID == origin {
			continue
		}

		cs = append(cs, con)
	}

	return cs
}

// followers returns follow connections from the origin.
func (a affiliations) followings(origin uint64) connection.List {
	cs := connection.List{}

	for con := range a {
		if con.Type == connection.TypeFriend {
			continue
		}

		if con.ToID == origin {
			continue
		}

		cs = append(cs, con)
	}

	return cs
}

// friends returns friend connections from the origin.
func (a affiliations) friends(origin uint64) connection.List {
	cs := connection.List{}

	for con := range a {
		if con.Type == connection.TypeFollow {
			continue
		}

		if con.FromID != origin && con.ToID != origin {
			continue
		}

		cs = append(cs, con)
	}

	return cs
}

// filterFollowers return an affiliations with all follow connections towards
// the origin removed.
func (a affiliations) filterFollowers(origin uint64) affiliations {
	am := affiliations{}

	for con, user := range a {
		if con.Type == connection.TypeFollow && con.ToID == origin {
			continue
		}

		am[con] = user
	}

	return am
}

// filterFollowings returns and affiliation with all following connections
// removed.
func (a affiliations) filterFollowings(origin uint64) affiliations {
	am := affiliations{}

	for con, user := range a {
		if con.Type == connection.TypeFollow && con.FromID == origin {
			continue
		}

		am[con] = user
	}

	return am
}

// filterFriends returns an affiliation with all friend connections removed.
func (a affiliations) filterFriends(origin uint64) affiliations {
	am := affiliations{}

	for con, user := range a {
		if con.Type == connection.TypeFriend {
			continue
		}

		am[con] = user
	}

	return am
}

// userIDs returns the user ids.
func (a affiliations) userIDs() []uint64 {
	var (
		ids  = make([]uint64, 0, len(a))
		seen = map[uint64]struct{}{}
	)

	for _, user := range a {
		if _, ok := seen[user.ID]; ok {
			continue
		}

		ids = append(ids, user.ID)
		seen[user.ID] = struct{}{}
	}

	return ids
}

// users returns the list of users.
func (a affiliations) users() user.List {
	var (
		seen = map[uint64]struct{}{}
		us   = user.List{}
	)

	for _, user := range a {
		if _, ok := seen[user.ID]; ok {
			continue
		}

		seen[user.ID] = struct{}{}
		us = append(us, user)
	}

	return us
}

// condition given an index and event determines if the Event should be kept in
// the list.
type condition func(int, *event.Event) bool

// source represents an event generator of varying origin.
type source func() (event.List, error)

// Feed is the composite to transport information relevant for a feed.
type Feed struct {
	Events  event.List
	Posts   PostList
	PostMap PostMap
	UserMap user.Map
}

// FeedController bundles the business constraints for feeds.
type FeedController struct {
	connections connection.Service
	events      event.Service
	objects     object.Service
	users       user.Service
}

// NewFeedController returns a controller instance.
func NewFeedController(
	connections connection.Service,
	events event.Service,
	objects object.Service,
	users user.Service,
) *FeedController {
	return &FeedController{
		connections: connections,
		events:      events,
		objects:     objects,
		users:       users,
	}
}

// Events returns the events from the interest and social graph of the given user.
func (c *FeedController) Events(
	currentApp *app.App,
	origin uint64,
	opts *event.QueryOptions,
) (*Feed, error) {
	am, err := c.neighbours(currentApp, origin, 0)
	if err != nil {
		return nil, err
	}

	var (
		neighbours = am.filterFollowers(origin)
		sources    = []source{
			sourceConnection(append(am.followers(origin), am.friends(origin)...)),
			sourceGlobal(c.events, currentApp, opts),
			sourceNeighbours(
				c.events,
				currentApp,
				opts,
				am.filterFollowers(origin).userIDs()...,
			),
			sourceTarget(c.events, currentApp, origin, opts),
		}
	)

	us := am.users()

	for _, u := range neighbours {
		a, err := c.neighbours(currentApp, u.ID, origin)
		if err != nil {
			return nil, err
		}

		cs := append(a.followings(u.ID), a.friends(u.ID)...)

		sources = append(sources, sourceConnection(cs))
		us = append(us, am.users()...)
	}

	es, err := collect(sources...)
	if err != nil {
		return nil, err
	}

	ps, err := extractPosts(c.objects, currentApp, es)
	if err != nil {
		return nil, err
	}

	err = enrichCounts(c.events, c.objects, currentApp, ps)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, currentApp, origin, ps)
	if err != nil {
		return nil, err
	}

	pm := ps.toMap()

	es = filter(
		es,
		conditionDuplicate(),
		conditionPostMissing(pm),
	)

	um, err := fillupUsers(c.users, currentApp, origin, us.ToMap(), es)
	if err != nil {
		return nil, err
	}

	sort.Sort(es)

	if len(es) > 200 {
		es = es[:199]
	}

	return &Feed{
		Events:  es,
		PostMap: pm,
		UserMap: um,
	}, nil
}

// News returns the events and posts from the interest and social graph of the
// given user.
func (c *FeedController) News(
	currentApp *app.App,
	origin uint64,
	eventOpts *event.QueryOptions,
	postOpts *object.QueryOptions,
) (*Feed, error) {
	am, err := c.neighbours(currentApp, origin, 0)
	if err != nil {
		return nil, err
	}

	var (
		neighbours = am.filterFollowers(origin)
		sources    = []source{
			sourceConnection(append(am.followers(origin), am.friends(origin)...)),
			sourceGlobal(c.events, currentApp, eventOpts),
			sourceNeighbours(
				c.events,
				currentApp,
				eventOpts,
				am.filterFollowers(origin).userIDs()...,
			),
			sourceTarget(c.events, currentApp, origin, eventOpts),
		}
	)

	us := am.users()

	for _, u := range neighbours {
		a, err := c.neighbours(currentApp, u.ID, origin)
		if err != nil {
			return nil, err
		}

		cs := append(a.followings(u.ID), a.friends(u.ID)...)

		sources = append(sources, sourceConnection(cs))
		us = append(us, am.users()...)
	}

	es, err := collect(sources...)
	if err != nil {
		return nil, err
	}

	ps, err := extractPosts(c.objects, currentApp, es)
	if err != nil {
		return nil, err
	}

	err = enrichCounts(c.events, c.objects, currentApp, ps)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, currentApp, origin, ps)
	if err != nil {
		return nil, err
	}

	pm := ps.toMap()

	es = filter(
		es,
		conditionDuplicate(),
		conditionPostMissing(pm),
	)

	um, err := fillupUsers(c.users, currentApp, origin, us.ToMap(), es)
	if err != nil {
		return nil, err
	}

	sort.Sort(es)

	if len(es) > 200 {
		es = es[:199]
	}

	ps, err = c.connectionPosts(currentApp, postOpts, neighbours.userIDs()...)
	if err != nil {
		return nil, err
	}

	gs, err := c.globalPosts(currentApp, postOpts)
	if err != nil {
		return nil, err
	}

	gum, err := user.MapFromIDs(c.users, currentApp.Namespace(), gs.OwnerIDs()...)
	if err != nil {
		return nil, err
	}

	um = um.Merge(gum)

	ps = append(ps, gs...)

	sort.Sort(ps)

	err = enrichCounts(c.events, c.objects, currentApp, ps)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, currentApp, origin, ps)
	if err != nil {
		return nil, err
	}

	errs := c.users.PutLastRead(currentApp.Namespace(), origin, time.Now())
	if errs != nil {
		// Updating the last read pointer of a user shouldn't stop the feed delivery
		// as we would accept an incorrect unread counter over a broken feed.
	}

	return &Feed{
		Events:  es,
		Posts:   ps,
		PostMap: pm,
		UserMap: um,
	}, nil
}

// NotificationsSelf returns the events which target the origin user and their
// content.
func (c *FeedController) NotificationsSelf(
	currentApp *app.App,
	origin uint64,
	opts *event.QueryOptions,
) (*Feed, error) {
	am, err := c.neighbours(currentApp, origin, 0)
	if err != nil {
		return nil, err
	}

	ps, err := userPosts(c.objects, currentApp, origin)
	if err != nil {
		return nil, err
	}

	var (
		fs      = am.filterFollowings(origin)
		sources = []source{
			sourceComment(c.objects, currentApp, ps.IDs()...),
			sourceConnection(fs.connections()),
			sourceLikes(c.events, currentApp, opts, ps.IDs()...),
			sourceTarget(c.events, currentApp, origin, opts),
		}
	)

	es, err := collect(sources...)
	if err != nil {
		return nil, err
	}

	um, err := fillupUsers(c.users, currentApp, origin, fs.users().ToMap(), es)
	if err != nil {
		return nil, err
	}

	sort.Sort(es)

	if len(es) > 200 {
		es = es[:199]
	}

	return &Feed{
		Events:  es,
		PostMap: ps.toMap(),
		UserMap: um,
	}, nil
}

// Posts returns the posts from the interest and social graph of the given user.
func (c *FeedController) Posts(
	currentApp *app.App,
	origin uint64,
	opts *object.QueryOptions,
) (*Feed, error) {
	am, err := c.neighbours(currentApp, origin, 0)
	if err != nil {
		return nil, err
	}

	neighbours := am.filterFollowers(origin)

	ps, err := c.connectionPosts(currentApp, opts, neighbours.userIDs()...)
	if err != nil {
		return nil, err
	}

	gs, err := c.globalPosts(currentApp, opts)
	if err != nil {
		return nil, err
	}

	um, err := user.MapFromIDs(c.users, currentApp.Namespace(), gs.OwnerIDs()...)
	if err != nil {
		return nil, err
	}

	ps = append(ps, gs...)

	sort.Sort(ps)

	err = enrichCounts(c.events, c.objects, currentApp, ps)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, currentApp, origin, ps)
	if err != nil {
		return nil, err
	}

	return &Feed{
		Posts:   ps,
		UserMap: am.users().ToMap().Merge(um),
	}, nil
}

func (c *FeedController) connectionPosts(
	currentApp *app.App,
	options *object.QueryOptions,
	ids ...uint64,
) (PostList, error) {
	if len(ids) == 0 {
		return PostList{}, nil
	}

	opts := object.QueryOptions{}
	if options != nil {
		opts = *options
	}

	opts.OwnerIDs = ids
	opts.Owned = &defaultOwned
	opts.Types = []string{TypePost}
	opts.Visibilities = []object.Visibility{
		object.VisibilityConnection,
		object.VisibilityPublic,
	}

	os, err := c.objects.Query(currentApp.Namespace(), opts)
	if err != nil {
		return nil, err
	}

	return postsFromObjects(os), nil
}

func (c *FeedController) globalPosts(
	currentApp *app.App,
	options *object.QueryOptions,
) (PostList, error) {
	opts := object.QueryOptions{}
	if options != nil {
		opts = *options
	}

	opts.Owned = &defaultOwned
	opts.Types = []string{TypePost}
	opts.Visibilities = []object.Visibility{
		object.VisibilityGlobal,
	}

	os, err := c.objects.Query(currentApp.Namespace(), opts)
	if err != nil {
		return nil, err
	}

	return postsFromObjects(os), nil
}

func (c *FeedController) neighbours(
	currentApp *app.App,
	origin uint64,
	root uint64,
) (affiliations, error) {
	fs, err := c.connections.Query(currentApp.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{
			origin,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		Types: []connection.Type{
			connection.TypeFollow,
		},
	})
	if err != nil {
		return nil, err
	}

	is, err := c.connections.Query(currentApp.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		FromIDs: []uint64{
			origin,
		},
		States: []connection.State{
			connection.StateConfirmed,
		},
		Types: []connection.Type{
			connection.TypeFriend,
		},
	})
	if err != nil {
		return nil, err
	}

	os, err := c.connections.Query(currentApp.Namespace(), connection.QueryOptions{
		Enabled: &defaultEnabled,
		States: []connection.State{
			connection.StateConfirmed,
		},
		ToIDs: []uint64{
			origin,
		},
		Types: []connection.Type{
			connection.TypeFriend,
		},
	})
	if err != nil {
		return nil, err
	}

	cs := append(fs, is...)
	cs = append(cs, os...)

	if root == 0 {
		fs, err := c.connections.Query(currentApp.Namespace(), connection.QueryOptions{
			Enabled: &defaultEnabled,
			States: []connection.State{
				connection.StateConfirmed,
			},
			ToIDs: []uint64{
				origin,
			},
			Types: []connection.Type{
				connection.TypeFollow,
			},
		})
		if err != nil {
			return nil, err
		}

		cs = append(cs, fs...)
	}

	am := affiliations{}

	for _, con := range cs {
		if con.ToID == root || con.FromID == root {
			continue
		}

		id := con.ToID

		if con.ToID == origin {
			id = con.FromID
		}

		us, err := c.users.Query(currentApp.Namespace(), user.QueryOptions{
			Enabled: &defaultEnabled,
			IDs: []uint64{
				id,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(us) != 1 {
			continue
		}

		am[con] = us[0]
	}

	return am, nil
}

// collect combines multiple soures into a single list of events.
func collect(sources ...source) (event.List, error) {
	events := event.List{}

	for _, s := range sources {
		es, err := s()
		if err != nil {
			return nil, err
		}

		events = append(events, es...)
	}

	return events, nil
}

// conditionDuplicate reports true if it encounters an Event with an ID already
// seen.
func conditionDuplicate() condition {
	seen := map[uint64]struct{}{}

	return func(idx int, event *event.Event) bool {
		if event.ID == 0 {
			return false
		}

		if _, ok := seen[event.ID]; ok {
			return true
		}

		seen[event.ID] = struct{}{}

		return false
	}
}

// conditionPostMissing reports true when the ObjectID of the event can't be
// found in the given ids.
func conditionPostMissing(pm PostMap) condition {
	return func(idx int, event *event.Event) bool {
		if event.ObjectID == 0 {
			return false
		}

		_, ok := pm[event.ObjectID]

		return !ok
	}
}

// extractPosts retrieves referenced post objects from a list of events.
func extractPosts(
	objects object.Service,
	currentApp *app.App,
	es event.List,
) (PostList, error) {
	ps := PostList{}

	for _, event := range es {
		if event.ObjectID == 0 {
			continue
		}

		os, err := objects.Query(currentApp.Namespace(), object.QueryOptions{
			ID: &event.ObjectID,
		})
		if err != nil {
			return nil, err
		}

		if len(os) == 1 && os[0].Type == TypePost {
			ps = append(ps, &Post{
				Object: os[0],
			})
		}
	}

	return ps, nil
}

// fillupUsers given a map of users and events fills up all missing users.
func fillupUsers(
	users user.Service,
	currentApp *app.App,
	originID uint64,
	um user.Map,
	es event.List,
) (user.Map, error) {
	for _, id := range es.UserIDs() {
		if _, ok := um[id]; ok {
			continue
		}

		us, err := users.Query(currentApp.Namespace(), user.QueryOptions{
			Enabled: &defaultEnabled,
			IDs: []uint64{
				id,
			},
		})
		if err != nil {
			return nil, err
		}

		if len(us) != 1 {
			continue
		}

		um[id] = us[0]
	}

	return um, nil
}

// filter filters out event for which one of the conditions is true.
func filter(events event.List, cs ...condition) event.List {
	es := event.List{}

	for idx, event := range events {
		keep := true

		for _, c := range cs {
			if c(idx, event) {
				keep = false
				break
			}
		}

		if !keep {
			continue
		}

		es = append(es, event)
	}

	return es
}

// sourceComment creates comment events for the given posts.
func sourceComment(
	objects object.Service,
	currentApp *app.App,
	postIDs ...uint64,
) source {
	if len(postIDs) == 0 {
		return func() (event.List, error) {
			return event.List{}, nil
		}
	}

	return func() (event.List, error) {
		es := event.List{}

		cs, err := objects.Query(currentApp.Namespace(), object.QueryOptions{
			ObjectIDs: postIDs,
			Owned:     &defaultOwned,
			Types: []string{
				TypeComment,
			},
		})
		if err != nil {
			return nil, err
		}

		for _, comment := range cs {
			id, err := flake.NextID("comment-events")
			if err != nil {
				return nil, err
			}

			es = append(es, &event.Event{
				Enabled:    true,
				ID:         id,
				ObjectID:   comment.ObjectID,
				Owned:      true,
				Type:       TypeComment,
				UserID:     comment.OwnerID,
				Visibility: event.VisibilityPrivate,
				CreatedAt:  comment.CreatedAt,
				UpdatedAt:  comment.UpdatedAt,
			})
		}

		return es, nil
	}
}

// sourceConnection creates follow events for the given connections.
func sourceConnection(cs connection.List) source {
	if len(cs) == 0 {
		return func() (event.List, error) {
			return event.List{}, nil
		}
	}

	return func() (event.List, error) {
		es := event.List{}

		for _, con := range cs {
			if con.State != connection.StateConfirmed {
				continue
			}

			t := event.TypeFollow

			if con.Type == connection.TypeFriend {
				t = event.TypeFriend
			}

			id, err := flake.NextID("connection-events")
			if err != nil {
				return nil, err
			}

			es = append(es, &event.Event{
				Enabled: true,
				ID:      id,
				Owned:   true,
				Target: &event.Target{
					ID:   strconv.FormatUint(con.ToID, 10),
					Type: event.TargetUser,
				},
				Type:       t,
				UserID:     con.FromID,
				Visibility: event.VisibilityPrivate,
				CreatedAt:  con.CreatedAt,
				UpdatedAt:  con.UpdatedAt,
			})
		}

		sort.Sort(es)

		return es, nil
	}
}

// sourceGlobal returns all events for app with visibility EventGlobal.
func sourceGlobal(
	events event.Service,
	currentApp *app.App,
	options *event.QueryOptions,
) source {
	opts := event.QueryOptions{}
	if options != nil {
		opts = *options
	}

	opts.Visibilities = []event.Visibility{
		event.VisibilityGlobal,
	}

	return func() (event.List, error) {
		es, err := events.Query(currentApp.Namespace(), opts)
		if err != nil {
			return nil, err
		}

		return es, nil
	}
}

func sourceLikes(
	events event.Service,
	currentApp *app.App,
	options *event.QueryOptions,
	postIDs ...uint64,
) source {
	if len(postIDs) == 0 {
		return func() (event.List, error) {
			return event.List{}, nil
		}
	}

	opts := event.QueryOptions{}
	if options != nil {
		opts = *options
	}
	opts.Enabled = &defaultEnabled
	opts.ObjectIDs = postIDs
	opts.Owned = &defaultOwned
	opts.Types = []string{
		TypeLike,
	}

	return func() (event.List, error) {
		return events.Query(currentApp.Namespace(), opts)
	}
}

// connectionUsers returns all events owned by the given user ids.
func sourceNeighbours(
	events event.Service,
	currentApp *app.App,
	options *event.QueryOptions,
	ids ...uint64,
) source {
	if len(ids) == 0 {
		return func() (event.List, error) {
			return event.List{}, nil
		}
	}

	opts := event.QueryOptions{}
	if options != nil {
		opts = *options
	}

	opts.Visibilities = []event.Visibility{
		event.VisibilityConnection,
		event.VisibilityPublic,
	}
	opts.UserIDs = ids

	return func() (event.List, error) {
		return events.Query(currentApp.Namespace(), opts)
	}
}

// sourceTarget returns all events where the origin is the target.
func sourceTarget(
	events event.Service,
	currentApp *app.App,
	origin uint64,
	options *event.QueryOptions,
) source {
	opts := event.QueryOptions{}
	if options != nil {
		opts = *options
	}

	opts.TargetIDs = []string{
		strconv.FormatUint(origin, 10),
	}
	opts.TargetTypes = []string{
		event.TargetUser,
	}
	opts.Visibilities = []event.Visibility{
		event.VisibilityPrivate,
	}

	return func() (event.List, error) {
		es, err := events.Query(currentApp.Namespace(), opts)
		if err != nil {
			return nil, err
		}

		return es, nil
	}
}

func userPosts(
	objects object.Service,
	currentApp *app.App,
	origin uint64,
) (PostList, error) {
	os, err := objects.Query(currentApp.Namespace(), object.QueryOptions{
		Owned: &defaultOwned,
		OwnerIDs: []uint64{
			origin,
		},
		Types: []string{
			TypePost,
		},
	})
	if err != nil {
		return nil, err
	}

	ps := PostList{}
	for _, o := range os {
		ps = append(ps, &Post{Object: o})
	}

	return ps, nil
}
