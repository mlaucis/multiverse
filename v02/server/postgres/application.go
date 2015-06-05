/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	application struct {
		storage core.Application
	}
)

func (app *application) Read(ctx *context.Context) (err errors.Error) {
	// TODO This one read only the current application maybe we want to have something to read any application?
	computeApplicationLastModified(ctx, ctx.Bag["application"].(*entity.Application))
	server.WriteResponse(ctx, ctx.Bag["application"].(*entity.Application), http.StatusOK, 10)
	return
}

func (app *application) Update(ctx *context.Context) (err errors.Error) {
	application := *(ctx.Bag["application"].(*entity.Application))
	if er := json.Unmarshal(ctx.Body, &application); er != nil {
		return errors.NewBadRequestError("failed to update the application (1)\n"+er.Error(), er.Error())
	}

	application.ID = ctx.Bag["applicationID"].(int64)
	application.AccountID = ctx.Bag["accountID"].(int64)
	application.PublicID = ctx.Bag["account"].(*entity.Account).PublicID

	if err = validator.UpdateApplication(ctx.Bag["application"].(*entity.Application), &application); err != nil {
		return
	}

	updatedApplication, err := app.storage.Update(*ctx.Bag["application"].(*entity.Application), application, true)
	if err != nil {
		return
	}

	server.WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
	return
}

func (app *application) Delete(ctx *context.Context) (err errors.Error) {
	if err = app.storage.Delete(ctx.Bag["application"].(*entity.Application)); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (app *application) Create(ctx *context.Context) (err errors.Error) {
	var (
		application = &entity.Application{}
	)

	if er := json.Unmarshal(ctx.Body, application); er != nil {
		return errors.NewBadRequestError("failed to create the application (1)\n"+er.Error(), er.Error())
	}

	application.AccountID = ctx.Bag["accountID"].(int64)
	application.PublicAccountID = ctx.Bag["account"].(*entity.Account).PublicID

	if err = validator.CreateApplication(application); err != nil {
		return
	}

	if application, err = app.storage.Create(application, true); err != nil {
		return
	}

	server.WriteResponse(ctx, application, http.StatusCreated, 0)
	return
}

func (app *application) List(ctx *context.Context) (err errors.Error) {
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

	computeApplicationsLastModified(ctx, response.Applications)

	server.WriteResponse(ctx, response, http.StatusOK, 10)
	return
}

func (app *application) PopulateContext(ctx *context.Context) (err errors.Error) {
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return errors.NewBadRequestError("error while reading application credentials", fmt.Sprintf("got %s:%s", user, pass))
	}
	ctx.Bag["application"], err = app.storage.FindByKey(user)
	if err == nil {
		ctx.Bag["accountID"] = ctx.Bag["application"].(*entity.Application).AccountID
		ctx.Bag["applicationID"] = ctx.Bag["application"].(*entity.Application).ID
	}
	return
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err errors.Error) {
	applicationID := ctx.Vars["applicationID"]
	if !validator.IsValidUUID5(applicationID) {
		return invalidAppIDError
	}

	ctx.Bag["application"], err = app.storage.FindByPublicID(applicationID)
	if err == nil {
		if ctx.Bag["application"].(*entity.Application) == nil {
			return errors.NewNotFoundError("application not found", "application not found")
		}

		ctx.Bag["accountID"] = ctx.Bag["application"].(*entity.Application).AccountID
		ctx.Bag["applicationID"] = ctx.Bag["application"].(*entity.Application).ID
	}
	return
}

// NewApplication returns a new application route handler
func NewApplication(storage core.Application) server.Application {
	return &application{
		storage: storage,
	}
}
