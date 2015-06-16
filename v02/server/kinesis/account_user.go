/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server"
	"github.com/tapglue/backend/v02/validator"
)

type (
	accountUser struct {
		writeStorage core.AccountUser
		readStorage  core.AccountUser
	}
)

func (accUser *accountUser) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) Delete(ctx *context.Context) (err []errors.Error) {
	if ctx.R.Header.Get("X-Jarvis-Auth") != "ZTBmZjI3MGE2M2YzYzAzOWI1MjhiYTNi" {
		return []errors.Error{errmsg.ErrServerReqMissingJarvisID}
	}

	accountUserID := ctx.Vars["accountUserID"]
	if !validator.IsValidUUID5(accountUserID) {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}
	accountUser, err := accUser.readStorage.FindByPublicID(ctx.Bag["accountID"].(int64), accountUserID)
	if err != nil {
		return
	}

	if err = accUser.writeStorage.Delete(accountUser); err != nil {
		return
	}

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (accUser *accountUser) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) Login(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) RefreshSession(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) Logout(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// PopulateContext adds the accountUser to the context
func (accUser *accountUser) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewAccountUser creates a new Account Route handler
func NewAccountUser(writeStorage, readStorage core.AccountUser) server.AccountUser {
	return &accountUser{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
