package redis

import (
	"fmt"

	"github.com/tapglue/multiverse/context"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
)

type application struct {
	storage, postgresStorage core.Application
}

func (app *application) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Delete(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) PopulateContext(ctx *context.Context) (err []errors.Error) {
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidApplicationCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", user, pass))}
	}
	ctx.Bag["application"], err = app.storage.FindByKey(user)
	if err == nil {
		ctx.Bag["accountID"] = ctx.Bag["application"].(*entity.Application).OrgID
		ctx.Bag["applicationID"] = ctx.Bag["application"].(*entity.Application).ID
	} else if err[0].Code() == errmsg.ErrApplicationNotFound.Code() {
		ctx.Bag["application"], err = app.postgresStorage.FindByKey(user)
		if err == nil {
			ctx.Bag["accountID"] = ctx.Bag["application"].(*entity.Application).OrgID
			ctx.Bag["applicationID"] = ctx.Bag["application"].(*entity.Application).ID
			go func(application *entity.Application) {
				app.storage.Create(application, false)
			}(ctx.Bag["application"].(*entity.Application))
		}
	}

	return
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewApplication returns a new application route handler
func NewApplication(storage, postgresStorage core.Application) handlers.Application {
	return &application{
		storage:         storage,
		postgresStorage: postgresStorage,
	}
}
