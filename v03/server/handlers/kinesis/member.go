package kinesis

import (
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
	"github.com/tapglue/backend/v03/validator"
)

type (
	accountUser struct {
		writeStorage core.Member
		readStorage  core.Member
	}
)

func (accUser *accountUser) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (accUser *accountUser) Delete(ctx *context.Context) (err []errors.Error) {
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

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
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
func NewAccountUser(writeStorage, readStorage core.Member) handlers.Member {
	return &accountUser{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
