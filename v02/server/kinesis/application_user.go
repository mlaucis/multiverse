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
	applicationUser struct {
		storage core.ApplicationUser
	}
)

func (appUser *applicationUser) Read(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) Update(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) Delete(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) Create(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) Login(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) RefreshSession(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) Logout(ctx *context.Context) (err errors.Error) {
	return nil
}

func (appUser *applicationUser) PopulateContext(ctx *context.Context) (err errors.Error) {
	return nil
}

// NewApplicationUser returns a new application user routes handler
func NewApplicationUser(storage core.ApplicationUser) server.ApplicationUser {
	return &applicationUser{
		storage: storage,
	}
}
