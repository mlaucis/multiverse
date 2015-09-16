package kinesis

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/context"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
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
	user := *(ctx.Bag["applicationUser"].(*entity.ApplicationUser))
	var er error
	if er = json.Unmarshal(ctx.Body, &user); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
	}

	user.ID = ctx.Bag["applicationUserID"].(uint64)

	if err = validator.UpdateUser(
		appUser.readStorage,
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser),
		&user); err != nil {
		return
	}

	updatedUser, err := appUser.writeStorage.Update(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		*(ctx.Bag["applicationUser"].(*entity.ApplicationUser)),
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

	response.WriteResponse(ctx, updatedUser, http.StatusCreated, 0)
	return
}

func (appUser *applicationUser) Delete(ctx *context.Context) (err []errors.Error) {
	userID, er := strconv.ParseUint(ctx.Vars["applicationUserID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	if err = appUser.writeStorage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		userID); err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (appUser *applicationUser) DeleteCurrent(ctx *context.Context) (err []errors.Error) {
	if err = appUser.writeStorage.Delete(
		ctx.Bag["accountID"].(int64),
		ctx.Bag["applicationID"].(int64),
		ctx.Bag["applicationUser"].(*entity.ApplicationUser).ID); err != nil {
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
