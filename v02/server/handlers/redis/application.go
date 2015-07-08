/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package redis

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server/handlers"
	"github.com/tapglue/backend/v02/server/response"
	"github.com/tapglue/backend/v02/validator"
)

type (
	application struct {
		storage core.Application
	}
)

func (app *application) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	response.WriteResponse(ctx, ctx.Bag["application"].(*entity.Application), http.StatusOK, 10)
	return
}

func (app *application) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	application := *(ctx.Bag["application"].(*entity.Application))
	if er := json.Unmarshal(ctx.Body, &application); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to update the application (1)\n"+er.Error(), er.Error())}
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

	response.WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
	return
}

func (app *application) Delete(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	if err = app.storage.Delete(ctx.Bag["application"].(*entity.Application)); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (app *application) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		application = &entity.Application{}
	)

	if er := json.Unmarshal(ctx.Body, application); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to create the application (1)\n"+er.Error(), er.Error())}
	}

	application.AccountID = ctx.Bag["accountID"].(int64)

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
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var (
		applications []*entity.Application
	)

	if applications, err = app.storage.List(ctx.Bag["accountID"].(int64)); err != nil {
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
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	ctx.Bag["application"], err = app.storage.Read(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64))
	return
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errors.NewInternalError(0, "not implemented yet", "not implemented yet")}
}

// NewApplication returns a new application route handler
func NewApplication(storage core.Application) handlers.Application {
	return &application{
		storage: storage,
	}
}
