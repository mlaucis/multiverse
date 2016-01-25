package controller

import (
	"sort"
	"strconv"

	"github.com/tapglue/multiverse/platform/flake"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	v04_core "github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
)

// affiliations is the composite structure to map connections to users.
type affiliations map[*v04_entity.Connection]*v04_entity.ApplicationUser

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
		if con.Type == v04_entity.ConnectionTypeFriend {
			continue
		}

		if con.UserFromID == origin {
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
		if con.Type == v04_entity.ConnectionTypeFriend {
			continue
		}

		if con.UserToID == origin {
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
		if con.Type == v04_entity.ConnectionTypeFollow {
			continue
		}

		if con.UserFromID != origin && con.UserToID != origin {
			continue
		}

		cs = append(cs, con)
	}

	return cs
}

// filterFollowers return an affiliations with all follow connections towards
// the origin rmeoved.
func (a affiliations) filterFollowers(origin uint64) affiliations {
	am := affiliations{}

	for con, user := range a {
		if con.Type == v04_entity.ConnectionTypeFollow && con.UserFromID != origin {
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
type condition func(int, *v04_entity.Event) bool

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
	connections connection.StrangleService
	events      event.StrangleService
	objects     object.Service
	users       user.StrangleService
}

// NewFeedController returns a controller instance.
func NewFeedController(
	connections connection.StrangleService,
	events event.StrangleService,
	objects object.Service,
	users user.StrangleService,
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
	app *v04_entity.Application,
	origin *v04_entity.ApplicationUser,
	cond *v04_core.EventCondition,
) (*Feed, error) {
	am, err := c.neighbours(app, origin, nil)
	if err != nil {
		return nil, err
	}

	var (
		neighbours = am.filterFollowers(origin.ID)
		sources    = []source{
			sourceConnection(append(am.followers(origin.ID), am.friends(origin.ID)...)),
			sourceGlobal(c.events, app, cond),
			sourceNeighbours(
				c.events,
				app,
				cond,
				am.filterFollowers(origin.ID).userIDs()...,
			),
		}
	)

	us := am.users()

	for _, u := range neighbours {
		a, err := c.neighbours(app, u, origin)
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

	ps, err := extractPosts(c.objects, app, es)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin.ID, ps)
	if err != nil {
		return nil, err
	}

	pm := ps.toMap()

	es = filter(
		es,
		conditionDuplicate(),
		conditionPostMissing(pm),
	)

	um, err := fillupUsers(c.users, app, origin, us.ToMap(), es)
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
	app *v04_entity.Application,
	origin *v04_entity.ApplicationUser,
	cond *v04_core.EventCondition,
) (*Feed, error) {
	am, err := c.neighbours(app, origin, nil)
	if err != nil {
		return nil, err
	}

	var (
		neighbours = am.filterFollowers(origin.ID)
		sources    = []source{
			sourceConnection(append(am.followers(origin.ID), am.friends(origin.ID)...)),
			sourceGlobal(c.events, app, cond),
			sourceNeighbours(
				c.events,
				app,
				cond,
				am.filterFollowers(origin.ID).userIDs()...,
			),
		}
	)

	us := am.users()

	for _, u := range neighbours {
		a, err := c.neighbours(app, u, origin)
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

	ps, err := extractPosts(c.objects, app, es)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin.ID, ps)
	if err != nil {
		return nil, err
	}

	pm := ps.toMap()

	es = filter(
		es,
		conditionDuplicate(),
		conditionPostMissing(pm),
	)

	um, err := fillupUsers(c.users, app, origin, us.ToMap(), es)
	if err != nil {
		return nil, err
	}

	sort.Sort(es)

	if len(es) > 200 {
		es = es[:199]
	}

	ps, err = c.connectionPosts(app, neighbours.userIDs()...)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin.ID, ps)
	if err != nil {
		return nil, err
	}

	errs := c.users.UpdateLastRead(app.OrgID, app.ID, origin.ID)
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

// Posts returns the posts from the interest and social graph of the given user.
func (c *FeedController) Posts(
	app *v04_entity.Application,
	origin *v04_entity.ApplicationUser,
) (*Feed, error) {
	am, err := c.neighbours(app, origin, nil)
	if err != nil {
		return nil, err
	}

	ps, err := c.connectionPosts(app, am.userIDs()...)
	if err != nil {
		return nil, err
	}

	err = enrichIsLiked(c.events, app, origin.ID, ps)
	if err != nil {
		return nil, err
	}

	return &Feed{
		Posts:   ps,
		UserMap: am.users().ToMap(),
	}, nil
}

func (c *FeedController) connectionPosts(
	app *v04_entity.Application,
	ids ...uint64,
) (PostList, error) {
	if len(ids) == 0 {
		return PostList{}, nil
	}

	os, err := c.objects.Query(app.Namespace(), object.QueryOptions{
		OwnerIDs: ids,
		Owned:    &defaultOwned,
		Types: []string{
			typePost,
		},
		Visibilities: []object.Visibility{
			object.VisibilityConnection,
			object.VisibilityPublic,
		},
	})
	if err != nil {
		return nil, err
	}

	return postsFromObjects(os), nil
}

func (c *FeedController) neighbours(
	app *v04_entity.Application,
	u *v04_entity.ApplicationUser,
	origin *v04_entity.ApplicationUser,
) (affiliations, error) {
	ucs, errs := c.connections.ConnectionsByState(
		app.OrgID,
		app.ID,
		u.ID,
		v04_entity.ConnectionStateConfirmed,
	)
	if errs != nil {
		return nil, errs[0]
	}

	am := affiliations{}

	for _, con := range ucs {
		if origin != nil &&
			(con.UserToID == origin.ID || con.UserFromID == origin.ID) {
			continue
		}

		id := con.UserToID

		if con.UserToID == u.ID {
			id = con.UserFromID
		}

		user, errs := c.users.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		am[con] = user
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

	return func(idx int, event *v04_entity.Event) bool {
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
	return func(idx int, event *v04_entity.Event) bool {
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
	app *v04_entity.Application,
	es event.List,
) (PostList, error) {
	ps := PostList{}

	for _, event := range es {
		if event.ObjectID == 0 {
			continue
		}

		os, err := objects.Query(app.Namespace(), object.QueryOptions{
			ID: &event.ObjectID,
		})
		if err != nil {
			return nil, err
		}

		if len(os) == 1 && os[0].Type == typePost {
			ps = append(ps, &Post{
				Object: os[0],
			})
		}
	}

	return ps, nil
}

// fillupUsers given a map of users and events fills up all missing users.
func fillupUsers(
	users user.StrangleService,
	app *v04_entity.Application,
	origin *v04_entity.ApplicationUser,
	um user.Map,
	es event.List,
) (user.Map, error) {
	for _, id := range es.UserIDs() {
		if _, ok := um[id]; ok || id == origin.ID {
			continue
		}

		user, errs := users.Read(app.OrgID, app.ID, id, false)
		if errs != nil {
			// Check for existence.
			if errs[0].Code() == errmsg.ErrApplicationUserNotFound.Code() {
				continue
			}
			return nil, errs[0]
		}

		um[id] = user
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
			if con.State != v04_entity.ConnectionStateConfirmed {
				continue
			}

			t := v04_entity.TypeEventFollow

			if con.Type == v04_entity.ConnectionTypeFriend {
				t = v04_entity.TypeEventFriend
			}

			id, err := flake.NextID("connection-events")
			if err != nil {
				return nil, err
			}

			es = append(es, &v04_entity.Event{
				ID:    id,
				Owned: true,
				Target: &v04_entity.Object{
					ID:   strconv.FormatUint(con.UserToID, 10),
					Type: v04_entity.TypeTargetUser,
				},
				Type:       t,
				UserID:     con.UserFromID,
				Visibility: v04_entity.EventPrivate,
				Common: v04_entity.Common{
					Enabled:   true,
					CreatedAt: con.UpdatedAt,
					UpdatedAt: con.UpdatedAt,
				},
			})
		}

		sort.Sort(es)

		return es, nil
	}
}

// sourceGlobal returns all events for app with visibility EventGlobal.
func sourceGlobal(
	events event.StrangleService,
	app *v04_entity.Application,
	cond *v04_core.EventCondition,
) source {
	condition := v04_core.EventCondition{}
	if cond != nil {
		condition = *cond
	}

	condition.Visibility = &v04_core.RequestCondition{
		Eq: int(v04_entity.EventGlobal),
	}

	return func() (event.List, error) {
		es, errs := events.ListAll(app.OrgID, app.ID, condition)
		if errs != nil {
			return nil, errs[0]
		}

		return es, nil
	}
}

// connectionUsers returns all events owned by the given user ids.
func sourceNeighbours(
	events event.StrangleService,
	app *v04_entity.Application,
	cond *v04_core.EventCondition,
	ids ...uint64,
) source {
	if len(ids) == 0 {
		return func() (event.List, error) {
			return event.List{}, nil
		}
	}

	condIDs := []interface{}{}

	for _, id := range ids {
		condIDs = append(condIDs, id)
	}

	condition := v04_core.EventCondition{}
	if cond != nil {
		condition = *cond
	}

	condition.Owned = &v04_core.RequestCondition{
		In: []interface{}{
			true,
			false,
		},
	}
	condition.Visibility = &v04_core.RequestCondition{
		In: []interface{}{
			int(v04_entity.EventConnections),
			int(v04_entity.EventPublic),
		},
	}
	condition.UserID = &v04_core.RequestCondition{
		In: condIDs,
	}

	return func() (event.List, error) {
		es, errs := events.ListAll(app.OrgID, app.ID, condition)
		if errs != nil {
			return nil, errs[0]
		}

		return es, nil
	}
}
