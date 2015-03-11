/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getAccount handles requests to a single account
// Request: GET /account/:AccountID
func getAccount(ctx *context) {
	var (
		accountID int64
		account   *entity.Account
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, account, http.StatusOK, 10)
}

// updateAccount handles requests to update a single account
// Request: PUT /account/:AccountID
func updateAccount(ctx *context) {
	var (
		accountID int64
		account   = &entity.Account{}
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.r.Body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if account.ID == 0 {
		account.ID = accountID
	}

	if err = validator.UpdateAccount(account); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if account, err = core.UpdateAccount(account, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, account, http.StatusCreated, 10)
}

// deleteAccount handles requests to delete a single account
// Request: DELETE /account/:AccountID
func deleteAccount(ctx *context) {
	var (
		accountID int64
		err       error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if err = core.DeleteAccount(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createAccount handles requests create an account
// Request: POST /accounts
func createAccount(ctx *context) {
	var (
		account = &entity.Account{}
		err     error
	)

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if err = validator.CreateAccount(account); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if account, err = core.WriteAccount(account, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, account, http.StatusCreated, 0)
}
