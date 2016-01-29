package controller

import (
	"math/rand"

	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
)

// RecommendationController bundles the business constriants for recommendations.
type RecommendationController struct {
	connections connection.StrangleService
	events      event.AggregateService
	users       user.StrangleService
}

// NewRecommendationController returns a controller instance.
func NewRecommendationController(
	connections connection.StrangleService,
	events event.AggregateService,
	users user.StrangleService,
) *RecommendationController {
	return &RecommendationController{
		connections: connections,
		events:      events,
		users:       users,
	}
}

// UsersActive returns the list of users with activity in the time period.
func (c *RecommendationController) UsersActive(
	app *v04_entity.Application,
	origin *v04_entity.ApplicationUser,
	period event.Period,
) (user.List, error) {
	ids, err := c.events.ActiveUserIDs(app.Namespace(), period)
	if err != nil {
		return nil, err
	}

	us, err := user.ListFromIDs(c.users, app, ids...)
	if err != nil {
		return nil, err
	}

	shuffleUsers(us)

	return us, nil
}

func (c *RecommendationController) filterConnections(
	app *v04_entity.Application,
	origin *v04_entity.ApplicationUser,
	recommends user.List,
) (user.List, error) {
	us := user.List{}

	for _, user := range recommends {
		r, errs := c.connections.Relation(app.OrgID, app.ID, origin.ID, user.ID)
		if errs != nil {
			return nil, errs[0]
		}

		if r.IsFriend != nil && *r.IsFriend || (r.IsFollowed != nil && *r.IsFollowed) {
			continue
		}

		us = append(us, user)
	}

	return us, nil
}

func shuffleUsers(us user.List) {
	for i := range us {
		j := rand.Intn(i + 1)
		us[i], us[j] = us[j], us[i]
	}
}
