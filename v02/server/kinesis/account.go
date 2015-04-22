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
	account struct {
		storage core.Account
	}
)

func (acc *account) Read(ctx *context.Context) (err errors.Error) {
	return nil
}

func (acc *account) Update(ctx *context.Context) (err errors.Error) {
	return nil
}

func (acc *account) Delete(ctx *context.Context) (err errors.Error) {
	return nil
}

func (acc *account) Create(ctx *context.Context) (err errors.Error) {
	return nil
}

func (acc *account) PopulateContext(ctx *context.Context) (err errors.Error) {
	return nil
}

// NewAccount returns a new account handler tweaked specifically for Kinesis
func NewAccount(datastore core.Account) server.Account {
	return &account{
		storage: datastore,
	}
}
