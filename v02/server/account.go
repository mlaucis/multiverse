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

var (
	acc core.Account
)

// GetAccount handles requests to a single account
// Request: GET /account/:AccountID
func GetAccount(ctx *context.Context) (err tgerrors.TGError) {
	WriteResponse(ctx, ctx.Bag["account"].(*entity.Account), http.StatusOK, 10)
	return
}

// UpdateAccount handles requests to update a single account
// Request: PUT /account/:AccountID
func UpdateAccount(ctx *context.Context) (err tgerrors.TGError) {
	account := *(ctx.Bag["account"].(*entity.Account))
	if er := json.Unmarshal(ctx.Body, &account); er != nil {
		return tgerrors.NewBadRequestError("failed to update the account (1)\n"+er.Error(), "malformed json received")
	}

	account.ID = ctx.Bag["accountID"].(int64)

	if err := validator.UpdateAccount(ctx.Bag["account"].(*entity.Account), &account); err != nil {
		return err
	}

	updatedAccount, err := acc.Update(*(ctx.Bag["account"].(*entity.Account)), account, true)
	if err != nil {
		return err
	}

	WriteResponse(ctx, updatedAccount, http.StatusCreated, 10)
	return nil
}

// DeleteAccount handles requests to delete a single account
// Request: DELETE /account/:AccountID
func DeleteAccount(ctx *context.Context) (err tgerrors.TGError) {
	if err = acc.Delete(ctx.Bag["accountID"].(int64)); err != nil {
		return err
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

// CreateAccount handles requests create an account
// Request: POST /accounts
func CreateAccount(ctx *context.Context) (err tgerrors.TGError) {
	var account = &entity.Account{}

	if er := json.Unmarshal(ctx.Body, account); er != nil {
		return tgerrors.NewBadRequestError("failed to create the account (1)\n"+er.Error(), er.Error())
	}

	if err = validator.CreateAccount(account); err != nil {
		return
	}

	if account, err = acc.Create(account, true); err != nil {
		return
	}

	WriteResponse(ctx, account, http.StatusCreated, 0)
	return
}
