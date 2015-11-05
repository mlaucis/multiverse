package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	"github.com/tapglue/multiverse/v04/server/handlers"
	"github.com/tapglue/multiverse/v04/server/response"
	"github.com/tapglue/multiverse/v04/validator"
)

type application struct {
	storage core.Application
}

func (app *application) Read(ctx *context.Context) (err []errors.Error) {
	// TODO This one read only the current application maybe we want to have something to read any application?
	response.WriteResponse(ctx, ctx.Application, http.StatusOK, 10)
	return
}

func (app *application) Update(ctx *context.Context) (err []errors.Error) {
	application := *ctx.Application
	if er := json.Unmarshal(ctx.Body, &application); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	application.ID = ctx.ApplicationID
	application.OrgID = ctx.OrganizationID
	application.PublicID = ctx.Application.PublicID
	application.PublicOrgID = ctx.Organization.PublicID

	if err = validator.UpdateApplication(ctx.Application, &application); err != nil {
		return
	}

	updatedApplication, err := app.storage.Update(*ctx.Application, application, true)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
	return
}

func (app *application) Delete(ctx *context.Context) (err []errors.Error) {
	if err = app.storage.Delete(ctx.Application); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (app *application) Create(ctx *context.Context) (err []errors.Error) {
	var (
		application = &entity.Application{}
	)

	if er := json.Unmarshal(ctx.Body, application); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error()).SetCurrentLocation()}
	}

	application.OrgID = ctx.OrganizationID
	application.PublicOrgID = ctx.Organization.PublicID

	if err = validator.CreateApplication(application); err != nil {
		return
	}

	if application, err = app.storage.Create(application, true); err != nil {
		return
	}

	response.WriteResponse(ctx, application, http.StatusCreated, 0)
	return
}

func (app *application) List(ctx *context.Context) (err []errors.Error) {
	var (
		applications []*entity.Application
	)

	if applications, err = app.storage.List(ctx.OrganizationID); err != nil {
		return
	}

	resp := &struct {
		Applications []*entity.Application `json:"applications"`
	}{
		Applications: applications,
	}

	response.WriteResponse(ctx, resp, http.StatusOK, 10)
	return
}

func (app *application) PopulateContext(ctx *context.Context) (err []errors.Error) {
	appToken, userToken, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidApplicationCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", appToken, userToken)).SetCurrentLocation()}
	}

	if len(appToken) == 32 {
		ctx.TokenType = context.TokenTypeApplication
		ctx.Application, err = app.storage.FindByApplicationToken(appToken)
	} else if len(appToken) == 44 {
		ctx.TokenType = context.TokenTypeBackend
		ctx.Application, err = app.storage.FindByBackendToken(appToken)
	} else {
		ctx.TokenType = context.TokenTypeUnknown
		return []errors.Error{errmsg.ErrAuthInvalidApplicationCredentials.UpdateInternalMessage(fmt.Sprintf("got unexpected token size %s:%s", appToken, userToken)).SetCurrentLocation()}
	}

	if err == nil {
		ctx.OrganizationID = ctx.Application.OrgID
		ctx.ApplicationID = ctx.Application.ID
	}

	if ctx.Application == nil {
		return []errors.Error{errmsg.ErrApplicationNotFound.SetCurrentLocation()}
	}

	return
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err []errors.Error) {
	applicationID := ctx.Vars["applicationID"]
	if !validator.IsValidUUID5(applicationID) {
		return []errors.Error{errmsg.ErrApplicationIDInvalid.SetCurrentLocation()}
	}

	ctx.Application, err = app.storage.FindByPublicID(applicationID)
	if err == nil {
		if ctx.Application == nil {
			return []errors.Error{errmsg.ErrApplicationNotFound.SetCurrentLocation()}
		}
		ctx.OrganizationID = ctx.Application.OrgID
		ctx.ApplicationID = ctx.Application.ID
	}

	if ctx.Application == nil {
		return []errors.Error{errmsg.ErrApplicationNotFound.SetCurrentLocation()}
	}

	return
}

// NewApplication returns a new application route handler
func NewApplication(storage core.Application) handlers.Application {
	return &application{
		storage: storage,
	}
}
