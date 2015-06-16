/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server"
)

type (
	application struct {
		writeStorage core.Application
		readStorage  core.Application
	}
)

func (app *application) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) Delete(ctx *context.Context) (err []errors.Error) {
	if err = app.writeStorage.Delete(ctx.Bag["application"].(*entity.Application)); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (app *application) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (app *application) PopulateContextFromID(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewApplication returns a new application route handler
func NewApplication(writeStorage, readStorage core.Application) server.Application {
	return &application{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
