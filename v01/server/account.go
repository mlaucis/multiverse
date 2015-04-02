/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/context"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
)

// getAccount handles requests to a single account
// Request: GET /account/:AccountID
func getAccount(ctx *context.Context) (err *tgerrors.TGError) {
	WriteResponse(ctx, ctx.Account, http.StatusOK, 10)
	return
}

// updateAccount handles requests to update a single account
// Request: PUT /account/:AccountID
func updateAccount(ctx *context.Context) (err *tgerrors.TGError) {
	account := *ctx.Account
	if er := json.Unmarshal(ctx.Body, &account); er != nil {
		return tgerrors.NewBadRequestError("failed to update the account (1)\n"+er.Error(), "malformed json received")
	}

	account.ID = ctx.AccountID

	if err := validator.UpdateAccount(ctx.Account, &account); err != nil {
		return err
	}

	updatedAccount, err := core.UpdateAccount(*ctx.Account, account, true)
	if err != nil {
		return err
	}

	WriteResponse(ctx, updatedAccount, http.StatusCreated, 10)
	return nil
}

// deleteAccount handles requests to delete a single account
// Request: DELETE /account/:AccountID
func deleteAccount(ctx *context.Context) (err *tgerrors.TGError) {
	if err = core.DeleteAccount(ctx.AccountID); err != nil {
		return err
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return nil
}

// createAccount handles requests create an account
// Request: POST /accounts
func createAccount(ctx *context.Context) (err *tgerrors.TGError) {
	var account = &entity.Account{}

	if er := json.Unmarshal(ctx.Body, account); er != nil {
		return tgerrors.NewBadRequestError("failed to create the account (1)\n"+er.Error(), er.Error())
	}

	if err = validator.CreateAccount(account); err != nil {
		return
	}

	if account, err = core.WriteAccount(account, true); err != nil {
		return
	}

	WriteResponse(ctx, account, http.StatusCreated, 0)
	return
}
