package kinesis

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server/handlers"
	"github.com/tapglue/backend/v02/server/response"
)

type (
	account struct {
		writeStorage core.Account
		readStorage  core.Account
	}
)

func (acc *account) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (acc *account) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (acc *account) Delete(ctx *context.Context) (err []errors.Error) {
	if ctx.Bag["account"].(*entity.Account).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if err = acc.writeStorage.Delete(ctx.Bag["account"].(*entity.Account)); err != nil {
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
func NewAccount(writeStorage, readStorage core.Account) handlers.Account {
	return &account{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
