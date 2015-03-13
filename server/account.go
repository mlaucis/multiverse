/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getAccount handles requests to a single account
// Request: GET /account/:AccountID
func getAccount(ctx *context.Context) {
	writeResponse(ctx, ctx.Account, http.StatusOK, 10)
}

// updateAccount handles requests to update a single account
// Request: PUT /account/:AccountID
func updateAccount(ctx *context.Context) {
	var (
		account = &entity.Account{}
		err     error
	)

	decoder := json.NewDecoder(ctx.R.Body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	account.ID = ctx.AccountID

	if err = validator.UpdateAccount(account); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if account, err = core.UpdateAccount(account, true); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, account, http.StatusCreated, 10)
}

// deleteAccount handles requests to delete a single account
// Request: DELETE /account/:AccountID
func deleteAccount(ctx *context.Context) {
	var (
		err error
	)

	if err = core.DeleteAccount(ctx.AccountID); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createAccount handles requests create an account
// Request: POST /accounts
func createAccount(ctx *context.Context) {
	var (
		account = &entity.Account{}
		err     error
	)

	decoder := json.NewDecoder(ctx.Body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if err = validator.CreateAccount(account); err != nil {
		errorHappened(ctx, err.Error(), http.StatusBadRequest, err)
		return
	}

	if account, err = core.WriteAccount(account, true); err != nil {
		errorHappened(ctx, err.Error(), http.StatusInternalServerError, err)
		return
	}

	writeResponse(ctx, account, http.StatusCreated, 0)
}
