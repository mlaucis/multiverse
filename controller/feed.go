package controller

import (
	"sort"

	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/object"
	"github.com/tapglue/multiverse/service/user"
	"github.com/tapglue/multiverse/v04/core"
	v04_core "github.com/tapglue/multiverse/v04/core"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

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
	user *v04_entity.ApplicationUser,
	condition *v04_core.EventCondition,
) (event.Events, error) {
	ids, errs := c.connections.FriendsAndFollowingIDs(app.OrgID, app.ID, user.ID)
	if errs != nil {
		return nil, errs[0]
	}

	es, err := c.connectionEvents(app, condition, ids...)
	if err != nil {
		return nil, err
	}

	gs, err := c.globalEvents(app, condition)
	if err != nil {
		return nil, err
	}

	es = append(es, gs...)

	ts, err := c.targetUserEvents(app, user.ID, condition)
	if err != nil {
		return nil, err
	}

	es = append(es, ts...)

	es = c.distinctEvents(es)

	sort.Sort(es)

	return es, nil
}

// News returns the events and posts from the interest and social graph of the
// given user.
func (c *FeedController) News(
	app *v04_entity.Application,
	user *v04_entity.ApplicationUser,
	condition *v04_core.EventCondition,
) (event.Events, Posts, error) {
	ids, errs := c.connections.FriendsAndFollowingIDs(app.OrgID, app.ID, user.ID)
	if errs != nil {
		return nil, nil, errs[0]
	}

	es, err := c.connectionEvents(app, condition, ids...)
	if err != nil {
		return nil, nil, err
	}

	gs, err := c.globalEvents(app, condition)
	if err != nil {
		return nil, nil, err
	}

	es = append(es, gs...)

	ts, err := c.targetUserEvents(app, user.ID, condition)
	if err != nil {
		return nil, nil, err
	}

	es = append(es, ts...)

	ps, err := c.connectionPosts(app, ids...)
	if err != nil {
		return nil, nil, err
	}

	es = c.distinctEvents(es)

	errs = c.users.UpdateLastRead(app.OrgID, app.ID, user.ID)
	if errs != nil {
		// Updating the last read pointer of a user shouldn't stop the feed delivery
		// as we would accept an incorrect unread counter over a broken feed.
	}

	// Sort collection by creation time.
	sort.Sort(es)

	// FIXME(xla): The hard limit can be solved with proper pagination.
	if len(es) > 200 {
		es = es[:199]
	}

	return es, ps, nil
}

// Posts returns the posts from the interest and social graph of the given user.
func (c *FeedController) Posts(
	app *v04_entity.Application,
	user *v04_entity.ApplicationUser,
) (Posts, error) {
	ids, errs := c.connections.FriendsAndFollowingIDs(app.OrgID, app.ID, user.ID)
	if errs != nil {
		return nil, errs[0]
	}

	return c.connectionPosts(app, ids...)
}

func (c *FeedController) connectionEvents(
	app *v04_entity.Application,
	cond *v04_core.EventCondition,
	ids ...uint64,
) (event.Events, error) {
	if len(ids) == 0 {
		return event.Events{}, nil
	}

	condIDs := []interface{}{}

	for _, id := range ids {
		condIDs = append(condIDs, id)
	}

	condition := v04_core.EventCondition{}
	if cond != nil {
		condition = *cond
	}

	condition.Owned = &core.RequestCondition{
		In: []interface{}{
			true,
			false,
		},
	}
	condition.Visibility = &core.RequestCondition{
		In: []interface{}{
			int(v04_entity.EventConnections),
			int(v04_entity.EventPublic),
		},
	}
	condition.UserID = &core.RequestCondition{
		In: condIDs,
	}

	es, errs := c.events.ListAll(app.OrgID, app.ID, condition)
	if errs != nil {
		return nil, errs[0]
	}

	return es, nil
}

func (c *FeedController) connectionPosts(
	app *v04_entity.Application,
	ids ...uint64,
) (Posts, error) {
	if len(ids) == 0 {
		return Posts{}, nil
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

	return fromObjects(os), nil
}

func (c *FeedController) globalEvents(
	app *v04_entity.Application,
	cond *v04_core.EventCondition,
) (event.Events, error) {
	condition := v04_core.EventCondition{}
	if cond != nil {
		condition = *cond
	}

	condition.Visibility = &core.RequestCondition{
		Eq: int(v04_entity.EventGlobal),
	}

	gs, errs := c.events.ListAll(app.OrgID, app.ID, condition)
	if errs != nil {
		return nil, errs[0]
	}

	return gs, nil
}

func (c *FeedController) targetUserEvents(
	app *v04_entity.Application,
	targetID uint64,
	cond *v04_core.EventCondition,
) (event.Events, error) {
	condition := v04_core.EventCondition{}
	if cond != nil {
		condition = *cond
	}

	condition.Target = &core.ObjectCondition{
		ID: &core.RequestCondition{
			Eq: targetID,
		},
		Type: &core.RequestCondition{
			Eq: user.TargetType,
		},
	}

	ts, errs := c.events.ListAll(app.OrgID, app.ID, condition)
	if errs != nil {
		return nil, errs[0]
	}

	return ts, nil
}

func (c *FeedController) distinctEvents(source event.Events) event.Events {
	result := event.Events{}

	found := false
	for idx := range source {
		found = false
		for resIdx := range result {
			if source[idx].ID == result[resIdx].ID {
				found = true
				break
			}
		}
		if !found {
			result = append(result, source[idx])
		}
	}

	return result
}
