package kinesis

import (
	"net/http"
	"strconv"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
)

type (
	connection struct {
		readAppUser  core.ApplicationUser
		writeStorage core.Connection
		readStorage  core.Connection
	}
)

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerNotImplementedYet}
}

func (conn *connection) Delete(ctx *context.Context) (err []errors.Error) {
	accountID := ctx.Bag["accountID"].(int64)
	applicationID := ctx.Bag["applicationID"].(int64)

	userFromID := ctx.Bag["applicationUserID"].(uint64)
	userToCustomID := ctx.Vars["applicationUserToID"]

	userToID, err := conn.determineTGUserID(accountID, applicationID, userToCustomID)
	if err != nil {
		return
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

func (conn *connection) determineTGUserID(accountID, applicationID int64, userID string) (uint64, []errors.Error) {
	id, er := strconv.ParseUint(userID, 10, 64)
	if er == nil {
		// TODO There has to be a better way to do this, no? no? But otherwise, how should we detect if the incoming id is a custom ID or not??
		if id > 27246450442288181 {
			return id, nil
		}
	}

	user, err := conn.readAppUser.FindByCustomID(accountID, applicationID, userID)
	if err != nil {
		return 0, err
	}

	return user.ID, nil
}

// NewConnectionWithApplicationUser returns a new connection handler
func NewConnectionWithApplicationUser(storage, permaStorage core.Connection, readAppUser core.ApplicationUser) handlers.Connection {
	return &connection{
		writeStorage: storage,
		readStorage:  permaStorage,
		readAppUser:  readAppUser,
	}
}
