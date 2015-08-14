package kinesis

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
)

type (
	account struct {
		writeStorage core.Organization
		readStorage  core.Organization
	}
)

func (acc *account) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (acc *account) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (acc *account) Delete(ctx *context.Context) (err []errors.Error) {
	if ctx.Bag["account"].(*entity.Organization).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if err = acc.writeStorage.Delete(ctx.Bag["account"].(*entity.Organization)); err != nil {
		return err
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (acc *account) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (acc *account) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewAccount returns a new account handler tweaked specifically for Kinesis
func NewAccount(writeStorage, readStorage core.Organization) handlers.Organization {
	return &account{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
