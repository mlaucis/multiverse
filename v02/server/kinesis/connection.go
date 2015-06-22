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

	userFromID := ctx.Bag["applicationUserID"].(string)
	userToID := ctx.Vars["applicationUserToID"]

	userToID, err = conn.determineTGUserID(accountID, applicationID, userToID)
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

	server.WriteResponse(ctx, "", http.StatusNoContent, 10)
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

func (conn *connection) determineTGUserID(accountID, applicationID int64, userID string) (string, []errors.Error) {
	if validator.IsValidUUID5(userID) {
		return userID, nil
	}

	user, err := conn.readAppUser.FindByCustomID(accountID, applicationID, userID)
	if err != nil {
		return "", err
	}

	return user.ID, nil
}

// NewConnection returns a new connection handler
func NewConnectionWithApplicationUser(storage, permaStorage core.Connection, readAppUser core.ApplicationUser) server.Connection {
	return &connection{
		writeStorage: storage,
		readStorage:  permaStorage,
		readAppUser:  readAppUser,
	}
}
