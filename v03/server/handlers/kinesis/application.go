package kinesis

import (
	"net/http"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/context"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
	"github.com/tapglue/multiverse/v03/server/response"
)

type application struct {
	writeStorage core.Application
	readStorage  core.Application
}

func (app *application) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) Delete(ctx *context.Context) (err []errors.Error) {
	if err = app.writeStorage.Delete(ctx.Application); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (app *application) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

// NewApplication returns a new application route handler
func NewApplication(writeStorage, readStorage core.Application) handlers.Application {
	return &application{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
