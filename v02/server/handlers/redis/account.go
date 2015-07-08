/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package redis

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/server/handlers"
	"github.com/tapglue/backend/v02/server/response"
	"github.com/tapglue/backend/v02/validator"
)

type (
	account struct {
		storage core.Account
	}
)

func (acc *account) Read(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	response.WriteResponse(ctx, ctx.Bag["account"].(*entity.Account), http.StatusOK, 10)
	return
}

func (acc *account) Update(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	account := *(ctx.Bag["account"].(*entity.Account))
	if er := json.Unmarshal(ctx.Body, &account); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to update the account (1)\n"+er.Error(), "malformed json received")}
	}

	account.ID = ctx.Bag["accountID"].(int64)

	if err := validator.UpdateAccount(ctx.Bag["account"].(*entity.Account), &account); err != nil {
		return err
	}

	updatedAccount, err := acc.storage.Update(*(ctx.Bag["account"].(*entity.Account)), account, true)
	if err != nil {
		return err
	}

	response.WriteResponse(ctx, updatedAccount, http.StatusCreated, 10)
	return nil
}

func (acc *account) Delete(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	if err = acc.storage.Delete(ctx.Bag["account"].(*entity.Account)); err != nil {
		return err
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (acc *account) Create(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	var account = &entity.Account{}

	if er := json.Unmarshal(ctx.Body, account); er != nil {
		return []errors.Error{errors.NewBadRequestError(0, "failed to create the account (1)\n"+er.Error(), er.Error())}
	}

	if err = validator.CreateAccount(account); err != nil {
		return
	}

	if account, err = acc.storage.Create(account, true); err != nil {
		return
	}

	response.WriteResponse(ctx, account, http.StatusCreated, 0)
	return
}

func (acc *account) PopulateContext(ctx *context.Context) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerDeprecatedStorage}
	ctx.Bag["account"], err = acc.storage.Read(ctx.Bag["accountID"].(int64))
	return
}

// NewAccount creates a new Account route handler
func NewAccount(storage core.Account) handlers.Account {
	return &account{
		storage: storage,
	}
}
