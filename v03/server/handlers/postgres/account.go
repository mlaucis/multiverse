package postgres

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/core"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
	"github.com/tapglue/backend/v03/server/handlers"
	"github.com/tapglue/backend/v03/server/response"
	"github.com/tapglue/backend/v03/validator"
)

type (
	account struct {
		storage core.Account
	}
)

func (acc *account) Read(ctx *context.Context) (err []errors.Error) {
	if ctx.Bag["account"] == nil {
		return []errors.Error{errmsg.ErrAccountMissingInContext}
	}

	if ctx.Bag["account"].(*entity.Account).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	response.ComputeAccountLastModified(ctx, ctx.Bag["account"].(*entity.Account))

	response.WriteResponse(ctx, ctx.Bag["account"].(*entity.Account), http.StatusOK, 10)
	return
}

func (acc *account) Update(ctx *context.Context) (err []errors.Error) {
	account := *(ctx.Bag["account"].(*entity.Account))

	if account.PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if er := json.Unmarshal(ctx.Body, &account); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
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
	if ctx.Bag["account"].(*entity.Account).PublicID != ctx.Vars["accountID"] {
		return []errors.Error{errmsg.ErrAccountMismatch}
	}

	if err = acc.storage.Delete(ctx.Bag["account"].(*entity.Account)); err != nil {
		return err
	}

	response.WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (acc *account) Create(ctx *context.Context) (err []errors.Error) {
	var account = &entity.Account{}

	if er := json.Unmarshal(ctx.Body, account); er != nil {
		return []errors.Error{errmsg.ErrServerReqBadJSONReceived.UpdateMessage(er.Error())}
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
	user, pass, ok := ctx.BasicAuth()
	if !ok {
		return []errors.Error{errmsg.ErrAuthInvalidAccountCredentials.UpdateInternalMessage(fmt.Sprintf("got %s:%s", user, pass))}
	}
	account, err := acc.storage.FindByKey(user)
	if account == nil {
		return []errors.Error{errmsg.ErrAccountNotFound}
	}
	if err == nil {
		ctx.Bag["account"] = account
		ctx.Bag["accountID"] = account.ID
	}
	return
}

// NewAccount returns a new account handler tweaked specifically for Kinesis
func NewAccount(datastore core.Account) handlers.Account {
	return &account{
		storage: datastore,
	}
}
