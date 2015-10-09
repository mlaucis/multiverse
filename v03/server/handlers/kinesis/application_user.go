package kinesis

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/context"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
	"github.com/tapglue/multiverse/v03/server/response"
	"github.com/tapglue/multiverse/v03/validator"
)

type applicationUser struct {
	writeStorage core.ApplicationUser
	readStorage  core.ApplicationUser
}

func (appUser *applicationUser) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) ReadCurrent(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) UpdateCurrent(ctx *context.Context) (err []errors.Error) {
	user := *ctx.ApplicationUser
	var er error
	if er = json.Unmarshal(ctx.Body, &user); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	user.ID = ctx.ApplicationUserID

	if err = validator.UpdateUser(
		appUser.readStorage,
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser,
		&user); err != nil {
		return
	}

	updatedUser, err := appUser.writeStorage.Update(
		ctx.OrganizationID,
		ctx.ApplicationID,
		*ctx.ApplicationUser,
		user,
		false)
	if err != nil {
		return
	}
	if updatedUser == nil {
		updatedUser = &user
	}

	updatedUser.Password = ""
	updatedUser.Enabled = false
	appUser.readStorage.FriendStatistics(ctx.OrganizationID, ctx.ApplicationID, updatedUser)

	response.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Delete(ctx *context.Context) (err []errors.Error) {
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	if err = appUser.writeStorage.Delete(
		ctx.OrganizationID,
		ctx.ApplicationID,
		userID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) DeleteCurrent(ctx *context.Context) (err []errors.Error) {
	if err = appUser.writeStorage.Delete(
		ctx.OrganizationID,
		ctx.ApplicationID,
		ctx.ApplicationUser.ID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) Login(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) RefreshSession(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) Logout(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) Search(*context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (appUser *applicationUser) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewApplicationUser returns a new application user routes handler
func NewApplicationUser(writeStorage, readStorage core.ApplicationUser) handlers.ApplicationUser {
	return &applicationUser{
		writeStorage: writeStorage,
		readStorage:  readStorage,
	}
}
