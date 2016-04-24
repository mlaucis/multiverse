package http

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/event"
	"github.com/tapglue/multiverse/service/user"
	v04_entity "github.com/tapglue/multiverse/v04/entity"
	v04_response "github.com/tapglue/multiverse/v04/server/response"
)

// RecommendUsersActiveDay returns a list of active users in the last day.
func RecommendUsersActiveDay(c *controller.RecommendationController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, user, event.ByDay)
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
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, user, event.ByWeek)
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
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, user, event.ByMonth)
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
	users user.StrangleList
}

func (p *payloadUsers) MarshalJSON() ([]byte, error) {
	ps := []*v04_entity.PresentationApplicationUser{}

	for _, user := range p.users {
		v04_response.SanitizeApplicationUser(user)

		ps = append(ps, &v04_entity.PresentationApplicationUser{
			ApplicationUser: user,
		})
	}

	return json.Marshal(struct {
		Users      []*v04_entity.PresentationApplicationUser `json:"users"`
		UsersCount int                                       `json:"users_count"`
	}{
		Users:      ps,
		UsersCount: len(ps),
	})
}
