/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package kinesis

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/server"
)

type (
	connection struct {
		appUser core.ApplicationUser
		storage core.Connection
	}
)

func (conn *connection) Update(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) Delete(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) Create(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) List(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) CurrentUserList(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) FollowedByList(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) CurrentUserFollowedByList(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) Confirm(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) CreateSocial(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) Friends(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

func (conn *connection) CurrentUserFriends(ctx *context.Context) (err []errors.Error) {
	return notImplementedYet
}

// NewConnection returns a new connection handler
func NewConnection(storage core.Connection) server.Connection {
	return &connection{
		storage: storage,
	}
}

// NewConnectionWithApplicationUser initializes a new connection with an application user
func NewConnectionWithApplicationUser(storage core.Connection, appUser core.ApplicationUser) server.Connection {
	return &connection{
		storage: storage,
		appUser: appUser,
	}
}
