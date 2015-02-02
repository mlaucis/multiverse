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

	"github.com/gorilla/mux"
)

// getAccount handles requests to a single account
// Request: GET /account/:ID
// Test with: curl -i localhost/0.1/account/:ID
func getAccount(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID int64
		account   *entity.Account
		err       error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(account, http.StatusOK, 10, w, r)
}

// updateAccount handles requests to update a single account
// Request: PUT /account/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"token":"token_1_TmV3IEFjY291bnQ=", "name":"New Account","description":"Description of the account", "enabled": true, "created_at":"2015-02-02T19:13:18.239759449Z", "received_at":"2015-02-02T19:13:18.239759449Z", "metadata":"{}"}' -X PUT localhost/0.1/account/:ID
func updateAccount(w http.ResponseWriter, r *http.Request) {
	if err := validatePutCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID int64
		account   = &entity.Account{}
		err       error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	account.ID = accountID

	if err = validator.UpdateAccount(account); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if account, err = core.UpdateAccount(account, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(account, http.StatusOK, 10, w, r)
}

// deleteAccount handles requests to delete a single account
// Request: DELETE /account/:ID
// Test with: curl -i -X DELETE localhost/0.1/account/:ID
func deleteAccount(w http.ResponseWriter, r *http.Request) {
	if err := validateDeleteCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID int64
		result    string
		err       error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if result, err = core.DeleteAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	fmt.Println("%d", result)

	writeResponse(result, http.StatusNoContent, 10, w, r)
}

// createAccount handles requests create an account
// Request: POST /accounts
// Test with: curl -i -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/0.1/accounts
func createAccount(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		account = &entity.Account{}
		err     error
	)

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if err = validator.CreateAccount(account); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if account, err = core.WriteAccount(account, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(account, http.StatusCreated, 0, w, r)
}
