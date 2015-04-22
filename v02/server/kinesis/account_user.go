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
	accountUser struct {
		storage core.AccountUser
	}
)

func (accUser *accountUser) Read(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) Update(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) Delete(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) Create(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) List(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) Login(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) RefreshSession(ctx *context.Context) (err errors.Error) {
	return nil
}

func (accUser *accountUser) Logout(ctx *context.Context) (err errors.Error) {
	return nil
}

// PopulateContext adds the accountUser to the context
func (accUser *accountUser) PopulateContext(ctx *context.Context) (err errors.Error) {
	return nil
}

// NewAccountUser creates a new Account Route handler
func NewAccountUser(storage core.AccountUser) server.AccountUser {
	return &accountUser{
		storage: storage,
	}
}
