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
	if ctx.R.Header.Get("X-Jarvis-Auth") != "ZTBmZjI3MGE2M2YzYzAzOWI1MjhiYTNi" {
		return []errors.Error{errmsg.ErrServerReqMissingJarvisID}
	}

	if ctx.Bag["account"].(*entity.Account).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if err = acc.writeStorage.Delete(ctx.Bag["account"].(*entity.Account)); err != nil {
		return err
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (acc *account) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (acc *account) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewAccount returns a new account handler tweaked specifically for Kinesis
func NewAccount(writeStorage, readStorage core.Account) server.Account {
	return &account{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
