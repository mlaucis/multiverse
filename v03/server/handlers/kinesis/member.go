package kinesis

import (
	"net/http"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/context"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
	"github.com/tapglue/multiverse/v03/server/response"
	"github.com/tapglue/multiverse/v03/validator"
)

type member struct {
	writeStorage core.Member
	readStorage  core.Member
}

func (user *member) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (user *member) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (user *member) Delete(ctx *context.Context) (err []errors.Error) {
	accountUserID := ctx.Vars["accountUserID"]
	if !validator.IsValidUUID5(accountUserID) {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid.SetCurrentLocation()}
	}
	accountUser, err := user.readStorage.FindByPublicID(ctx.OrganizationID, accountUserID)
	if err != nil {
		return
	}

	if err = user.writeStorage.Delete(accountUser); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (user *member) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (user *member) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (user *member) Login(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (user *member) RefreshSession(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

func (user *member) Logout(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

// PopulateContext adds the accountUser to the context
func (user *member) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet.SetCurrentLocation()}
}

// NewMember creates a new member route handler
func NewMember(writeStorage, readStorage core.Member) handlers.Member {
	return &member{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
