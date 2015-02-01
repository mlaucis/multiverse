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
// Test with: curl -i localhost/account/:ID
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

// createAccount handles requests create an account
// Request: POST /accounts
// Test with: curl -i -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/accounts
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

	account.Enabled = true

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
