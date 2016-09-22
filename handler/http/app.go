package http

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"golang.org/x/net/context"

	"github.com/tapglue/multiverse/controller"
	"github.com/tapglue/multiverse/service/app"
)

// AppCreate creates an application for the current Org.
func AppCreate(fn controller.AppCreateFunc) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentOrg = orgFromContext(ctx)
			p          = payloadApp{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		app, err := fn(currentOrg, p.name, p.description)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusCreated, &payloadApp{app: app})
	}
}

// AppDelete disables the App and renders it unusable.
func AppDelete(fn controller.AppDeleteFunc) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentOrg = orgFromContext(ctx)
			publicID   = mux.Vars(r)["appID"]
		)

		err := fn(currentOrg, publicID)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusNoContent, nil)
	}
}

// AppList returns all Apps for the current Org.
func AppList(fn controller.AppListFunc) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		currentOrg := orgFromContext(ctx)

		opts, err := extractAppOpts(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		opts.Before, err = extractTimeCursorBefore(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		opts.Limit, err = extractLimit(r)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		as, err := fn(currentOrg, opts)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		if len(as) == 0 {
			respondJSON(w, http.StatusNoContent, nil)
			return
		}

		respondJSON(w, http.StatusOK, &payloadApps{
			apps: as,
			pagination: pagination(
				r,
				opts.Limit,
				appCursorAfter(as, opts.Limit),
				appCursorBefore(as, opts.Limit),
			),
		})
	}
}

// AppUpdate updates the values of an App.
func AppUpdate(fn controller.AppUpdateFunc) Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		var (
			currentOrg = orgFromContext(ctx)
			publicID   = mux.Vars(r)["appID"]
			p          = payloadApp{}
		)

		err := json.NewDecoder(r.Body).Decode(&p)
		if err != nil {
			respondError(w, 0, wrapError(ErrBadRequest, err.Error()))
			return
		}

		app, err := fn(currentOrg, publicID, p.name, p.description)
		if err != nil {
			respondError(w, 0, err)
			return
		}

		respondJSON(w, http.StatusOK, &payloadApp{app: app})
	}
}

type payloadApp struct {
	app         *app.App
	description string
	name        string
}

func (p *payloadApp) MarshalJSON() ([]byte, error) {
	f := struct {
		BackendToken string    `json:"backend_token"`
		Description  string    `json:"description"`
		Enabled      bool      `json:"enabled"`
		InProduction bool      `json:"in_production"`
		Name         string    `json:"name"`
		OrgID        string    `json:"account_id"`
		PublicID     string    `json:"id"`
		Token        string    `json:"token"`
		URL          string    `json:"url"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
	}{
		BackendToken: p.app.BackendToken,
		Description:  p.app.Description,
		Enabled:      p.app.Enabled,
		InProduction: p.app.InProduction,
		Name:         p.app.Name,
		OrgID:        p.app.PublicOrgID,
		PublicID:     p.app.PublicID,
		Token:        p.app.Token,
		URL:          p.app.URL,
		CreatedAt:    p.app.CreatedAt,
		UpdatedAt:    p.app.UpdatedAt,
	}

	return json.Marshal(&f)
}

func (p *payloadApp) UnmarshalJSON(raw []byte) error {
	f := struct {
		Description string `json:"description"`
		Name        string `json:"name"`
	}{}

	err := json.Unmarshal(raw, &f)
	if err != nil {
		return err
	}

	p.description = f.Description
	p.name = f.Name

	return nil
}

type payloadApps struct {
	apps       app.List
	pagination *payloadPagination
}

func (p *payloadApps) MarshalJSON() ([]byte, error) {
	as := []*payloadApp{}

	for _, app := range p.apps {
		as = append(as, &payloadApp{app: app})
	}

	f := struct {
		Apps       []*payloadApp      `json:"applications"`
		Pagination *payloadPagination `json:"paging"`
	}{
		Apps:       as,
		Pagination: p.pagination,
	}

	return json.Marshal(&f)
}

func appCursorAfter(as app.List, limit int) string {
	var after string

	if len(as) > 0 {
		after = toTimeCursor(as[0].CreatedAt)
	}

	return after
}

func appCursorBefore(as app.List, limit int) string {
	var before string

	if len(as) > 0 {
		before = toTimeCursor(as[len(as)-1].CreatedAt)
	}

	return before
}
