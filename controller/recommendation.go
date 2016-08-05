package controller

import (
	"math/rand"

	"github.com/tapglue/multiverse/service/app"
	"github.com/tapglue/multiverse/service/connection"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
)

// conditionUser determines if the user should be filtered.
type conditionUser func(*user.User) (bool, error)

// RecommendationController bundles the business constriants for recommendations.
type RecommendationController struct {
	connections connection.Service
	events      event.AggregateService
	users       user.Service
}

// NewRecommendationController returns a controller instance.
func NewRecommendationController(
	connections connection.Service,
	events event.Service,
	users user.Service,
) *RecommendationController {
	return &RecommendationController{
		connections: connections,
		events:      events,
		users:       users,
	}
}

// UsersActive returns the list of users with activity in the time period.
func (c *RecommendationController) UsersActive(
	currentApp *app.App,
	origin uint64,
	period event.Period,
) (user.List, error) {
	ids, err := c.events.ActiveUserIDs(currentApp.Namespace(), period)
	if err != nil {
		return nil, err
	}

	us, err := user.ListFromIDs(c.users, currentApp.Namespace(), ids...)
	if err != nil {
		return nil, err
	}

	us, err = filterUsers(
		us,
		conditionOrigin(origin),
		conditionConnection(c.connections, currentApp, origin),
	)
	if err != nil {
		return nil, err
	}

	shuffleUsers(us)

	return us, nil
}

func conditionConnection(
	connections connection.Service,
	currentApp *app.App,
	origin uint64,
) conditionUser {
	return func(user *user.User) (bool, error) {
		r, err := queryRelation(connections, currentApp, origin, user.ID)
		if err != nil {
			return false, err
		}

		if r.isFriend || r.isFollowing {
			return true, nil
		}

		return false, nil
	}
}

func conditionOrigin(
	origin uint64,
) conditionUser {
	return func(user *user.User) (bool, error) {
		return origin == user.ID, nil
	}
}

func filterUsers(users user.List, cs ...conditionUser) (user.List, error) {
	us := user.List{}

	for _, user := range users {
		keep := true

		for _, condition := range cs {
			drop, err := condition(user)
			if err != nil {
				return nil, err
			}

			if drop {
				keep = false
				break
			}
		}

		if keep {
			us = append(us, user)
		}
	}

	return us, nil
}

func shuffleUsers(us user.List) {
	for i := range us {
		j := rand.Intn(i + 1)
		us[i], us[j] = us[j], us[i]
	}
}
