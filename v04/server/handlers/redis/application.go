package redis

import (
	"fmt"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/context"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
)

type application struct {
	storage, postgresStorage core.Application
}

func (app *application) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) Delete(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) PopulateContext(ctx *context.Context) (err []errors.Error) {
	appToken, userToken, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidApplicationCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", appToken, userToken)).SetCurrentLocation()}
	}

	if len(appToken) == 32 {
		ctx.Application, err = app.storage.FindByApplicationToken(appToken)
		if err == nil {
			ctx.OrganizationID = ctx.Application.OrgID
			ctx.ApplicationID = ctx.Application.ID
			ctx.TokenType = context.TokenTypeApplication
		} else if err[0].Code() == errmsg.ErrApplicationNotFound.Code() {
			ctx.Application, err = app.postgresStorage.FindByApplicationToken(appToken)
			if err == nil {
				ctx.OrganizationID = ctx.Application.OrgID
				ctx.ApplicationID = ctx.Application.ID
				ctx.TokenType = context.TokenTypeApplication
				_, err := app.storage.Create(ctx.Application, false)
				if err != nil {
					ctx.LogError(err)
				}
			}
		}
	} else if len(appToken) == 44 {
		ctx.Application, err = app.storage.FindByBackendToken(appToken)
		if err == nil {
			ctx.OrganizationID = ctx.Application.OrgID
			ctx.ApplicationID = ctx.Application.ID
			ctx.TokenType = context.TokenTypeBackend
		} else if err[0].Code() == errmsg.ErrApplicationNotFound.Code() {
			ctx.Application, err = app.postgresStorage.FindByBackendToken(appToken)
			if err == nil {
				ctx.OrganizationID = ctx.Application.OrgID
				ctx.ApplicationID = ctx.Application.ID
				ctx.TokenType = context.TokenTypeBackend
				_, err := app.storage.Create(ctx.Application, false)
				if err != nil {
					ctx.LogError(err)
				}
			}
		}
	} else {
		ctx.TokenType = context.TokenTypeUnknown
		return []errors.Error{errmsg.ErrAuthInvalidApplicationCredentials.UpdateInternalMessage(fmt.Sprintf("got unexpected token size %s:%s", appToken, userToken)).SetCurrentLocation()}
	}

	if ctx.Application == nil {
		return []errors.Error{errmsg.ErrApplicationNotFound.SetCurrentLocation()}
	}

	return
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

// NewApplication returns a new application route handler
func NewApplication(storage, postgresStorage core.Application) handlers.Application {
	return &application{
		storage:         storage,
		postgresStorage: postgresStorage,
	}
}
