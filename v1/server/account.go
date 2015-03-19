/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/server/utils"
	"github.com/tapglue/backend/v1/core"
	"github.com/tapglue/backend/v1/entity"
	"github.com/tapglue/backend/v1/validator"
)

// getAccount handles requests to a single account
// Request: GET /account/:AccountID
func getAccount(ctx *context.Context) {
	utils.WriteResponse(ctx, ctx.Account, http.StatusOK, 10)
}

// updateAccount handles requests to update a single account
// Request: PUT /account/:AccountID
func updateAccount(ctx *context.Context) {
	var err error

	account := *ctx.Account
	if err = json.NewDecoder(ctx.R.Body).Decode(&account); err != nil {
		utils.ErrorHappened(ctx, "failed to update the account (1)"+err.Error(), http.StatusBadRequest, err)
		return
	}

	account.ID = ctx.AccountID

	if err = validator.UpdateAccount(ctx.Account, &account); err != nil {
		utils.ErrorHappened(ctx, "failed to update the account (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	updatedAccount, err := core.UpdateAccount(*ctx.Account, account, true)
	if err != nil {
		utils.ErrorHappened(ctx, "failed to update the account (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, updatedAccount, http.StatusCreated, 10)
}

// deleteAccount handles requests to delete a single account
// Request: DELETE /account/:AccountID
func deleteAccount(ctx *context.Context) {
	var (
		err error
	)

	if err = core.DeleteAccount(ctx.AccountID); err != nil {
		utils.ErrorHappened(ctx, "failed to delete the account (1)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "", http.StatusNoContent, 10)
}

// createAccount handles requests create an account
// Request: POST /accounts
func createAccount(ctx *context.Context) {
	var (
		account = &entity.Account{}
		err     error
	)

	if err = json.NewDecoder(ctx.Body).Decode(account); err != nil {
		utils.ErrorHappened(ctx, "failed to create the account (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = validator.CreateAccount(account); err != nil {
		utils.ErrorHappened(ctx, "failed to create the account (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if account, err = core.WriteAccount(account, true); err != nil {
		utils.ErrorHappened(ctx, "failed to create the account (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, account, http.StatusCreated, 0)
}
