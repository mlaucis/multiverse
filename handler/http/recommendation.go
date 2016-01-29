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
			app  = appFromContext(ctx)
			user = userFromContext(ctx)
		)

		us, err := c.UsersActive(app, user, event.ByDay)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadUsers{users: us})
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

		respondJSON(w, http.StatusCreated, &payloadUsers{users: us})
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

		respondJSON(w, http.StatusCreated, &payloadUsers{users: us})
	}
}

type payloadUsers struct {
	users user.List
}

func (p *payloadUsers) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct{}{})
}
