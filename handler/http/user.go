package http

import (
	"encoding/json"
	"net/http"

	"golang.org/x/net/context"

	"github.com/gorilla/mux"
	"github.com/tapglue/multiverse/controller"
)

// UserSearchEmails returns all Users for the emails of the payload.
func UserSearchEmails(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadSearchEmails{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		if len(p.Emails) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := c.ListByEmails(app, currentUser.ID, p.Emails...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{
			users: us,
		})
	}
}

// UserSearchPlatform returns all users for the given ids and platform.
func UserSearchPlatform(c *controller.UserController) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			platform    = mux.Vars(r)["platform"]
			app         = appFromContext(ctx)
			currentUser = userFromContext(ctx)
			p           = payloadSearchPlatform{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		if len(p.IDs) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		us, err := c.ListByPlatformIDs(app, currentUser.ID, platform, p.IDs...)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(us) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadUsers{
			users: us,
		})
	}
}

type payloadSearchEmails struct {
	Emails []string `json:"emails"`
}

type payloadSearchPlatform struct {
	IDs []string `json:"ids"`
}
