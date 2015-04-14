/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

type (
	// Account holds the account routes
	Account interface {
		// Read handles requests to a single account
		Read(*context.Context) tgerrors.TGError

		// Update handles requests to update a single account
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single account
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests create an account
		Create(*context.Context) tgerrors.TGError

		// PopulateContext adds the account to the context
		PopulateContext(*context.Context) tgerrors.TGError
	}

	account struct {
		storage core.Account
	}
)

func (acc *account) Read(ctx *context.Context) (err tgerrors.TGError) {
	WriteResponse(ctx, ctx.Bag["account"].(*entity.Account), http.StatusOK, 10)
	return
}

func (acc *account) Update(ctx *context.Context) (err tgerrors.TGError) {
	account := *(ctx.Bag["account"].(*entity.Account))
	if er := json.Unmarshal(ctx.Body, &account); er != nil {
		return tgerrors.NewBadRequestError("failed to update the account (1)\n"+er.Error(), "malformed json received")
	}

	account.ID = ctx.Bag["accountID"].(int64)

	if err := validator.UpdateAccount(ctx.Bag["account"].(*entity.Account), &account); err != nil {
		return err
	}

	updatedAccount, err := acc.storage.Update(*(ctx.Bag["account"].(*entity.Account)), account, true)
	if err != nil {
		return err
	}

	WriteResponse(ctx, updatedAccount, http.StatusCreated, 10)
	return nil
}

func (acc *account) Delete(ctx *context.Context) (err tgerrors.TGError) {
	if err = acc.storage.Delete(ctx.Bag["accountID"].(int64)); err != nil {
		return err
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

func (acc *account) Create(ctx *context.Context) (err tgerrors.TGError) {
	var account = &entity.Account{}

	if er := json.Unmarshal(ctx.Body, account); er != nil {
		return tgerrors.NewBadRequestError("failed to create the account (1)\n"+er.Error(), er.Error())
	}

	if err = validator.CreateAccount(account); err != nil {
		return
	}

	if account, err = acc.storage.Create(account, true); err != nil {
		return
	}

	WriteResponse(ctx, account, http.StatusCreated, 0)
	return
}

func (acc *account) PopulateContext(ctx *context.Context) (err tgerrors.TGError) {
	ctx.Bag["account"], err = acc.storage.Read(ctx.Bag["accountID"].(int64))
	return
}

// NewAccount creates a new Account route handler
func NewAccount(storage core.Account) Account {
	return &account{
		storage: storage,
	}
}
