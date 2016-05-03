package http

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
)

// RecommendUsersActiveDay returns a list of active users in the last day.
func RecommendUsersActiveDay(c *controller.RecommendationController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, currentUser.ID, event.ByDay)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

// RecommendUsersActiveWeek returns a list of active users in the last week.
func RecommendUsersActiveWeek(c *controller.RecommendationController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, currentUser.ID, event.ByWeek)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

// RecommendUsersActiveMonth returns a list of active users in the last month.
func RecommendUsersActiveMonth(c *controller.RecommendationController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, currentUser.ID, event.ByMonth)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{users: us})
	}
}

type payloadUsers struct {
	users user.List
}

func (p *payloadUsers) MarshalJSON() ([]byte, error) {
	ps := []*payloadUser{}

	for _, u := range p.users {
		ps = append(ps, &payloadUser{
			user: u,
		})
	}

	return json.Marshal(struct {
		Users      []*payloadUser `json:"users"`
		UsersCount int            `json:"users_count"`
	}{
		Users:      ps,
		UsersCount: len(ps),
	})
}
