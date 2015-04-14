/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

type (
	// Application defines the routes for the application
	Application interface {
		// Read handles requests to a single application
		Read(*context.Context) tgerrors.TGError

		// Update handles requests updates an application
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single application
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests create an application
		Create(*context.Context) tgerrors.TGError

		// List handles requests list all account applications
		List(*context.Context) tgerrors.TGError

		// PopulateContext adds the application to the context
		PopulateContext(ctx *context.Context) tgerrors.TGError
	}

	application struct {
		storage core.Application
	}
)

func (app *application) Read(ctx *context.Context) (err tgerrors.TGError) {
	WriteResponse(ctx, ctx.Bag["application"].(*entity.Application), http.StatusOK, 10)
	return
}

func (app *application) Update(ctx *context.Context) (err tgerrors.TGError) {
	application := *(ctx.Bag["application"].(*entity.Application))
	if er := json.Unmarshal(ctx.Body, &application); er != nil {
		return tgerrors.NewBadRequestError("failed to update the application (1)\n"+er.Error(), er.Error())
	}

	application.ID = ctx.Bag["applicationID"].(int64)
	application.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.UpdateApplication(ctx.Bag["application"].(*entity.Application), &application); err != nil {
		return
	}

	updatedApplication, err := app.storage.Update(*ctx.Bag["application"].(*entity.Application), application, true)
	if err != nil {
		return
	}

	WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
	return
}

func (app *application) Delete(ctx *context.Context) (err tgerrors.TGError) {
	if err = app.storage.Delete(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64)); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (app *application) Create(ctx *context.Context) (err tgerrors.TGError) {
	var (
		application = &entity.Application{}
	)

	if er := json.Unmarshal(ctx.Body, application); er != nil {
		return tgerrors.NewBadRequestError("failed to create the application (1)\n"+er.Error(), er.Error())
	}

	application.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.CreateApplication(application); err != nil {
		return
	}

	if application, err = app.storage.Create(application, true); err != nil {
		return
	}

	WriteResponse(ctx, application, http.StatusCreated, 0)
	return
}

func (app *application) List(ctx *context.Context) (err tgerrors.TGError) {
	var (
		applications []*entity.Application
	)

	if applications, err = app.storage.List(ctx.Bag["accountID"].(int64)); err != nil {
		return
	}

	response := &struct {
		Applications []*entity.Application `json:"applications"`
	}{
		Applications: applications,
	}

	WriteResponse(ctx, response, http.StatusOK, 10)
	return
}

func (app *application) PopulateContext(ctx *context.Context) (err tgerrors.TGError) {
	ctx.Bag["application"], err = app.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64))
	return
}

// NewApplication returns a new application route handler
func NewApplication(storage core.Application) Application {
	return &application{
		storage: storage,
	}
}
