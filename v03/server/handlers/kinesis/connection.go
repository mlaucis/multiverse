package kinesis

import (
	"net/http"
	"strconv"

	"github.com/tapglue/multiverse/context"
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/server/handlers"
	"github.com/tapglue/multiverse/v03/server/response"
)

type connection struct {
	readAppUser  core.ApplicationUser
	writeStorage core.Connection
	readStorage  core.Connection
}

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) Delete(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)
	userFromID := ctx.Bag["applicationUserID"].(uint64)

	userToID, er := strconv.ParseUint(ctx.Vars["applicationUserToID"], 10, 64)
	if er != nil {
		return []errors.Error{errmsg.ErrApplicationUserIDInvalid}
	}

	connection, err := conn.readStorage.Read(accountID, applicationID, userFromID, userToID)
	if err != nil {
		return
	}

	err = conn.writeStorage.Delete(accountID, applicationID, connection)
	if err != nil {
		return
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

func (conn *connection) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) List(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) FollowedByList(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) CurrentUserFollowedByList(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) Confirm(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) CreateSocial(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) Friends(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) CurrentUserFriends(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) CreateFriend(*context.Context) []errors.Error {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) CreateFollow(*context.Context) []errors.Error {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewConnectionWithApplicationUser returns a new connection handler
func NewConnectionWithApplicationUser(storage, permaStorage core.Connection, readAppUser core.ApplicationUser) handlers.Connection {
	return &connection{
		writeStorage: storage,
		readStorage:  permaStorage,
		readAppUser:  readAppUser,
	}
}
