/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/tapglue/backend/core"
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/validator"
)

// getAccount handles requests to a single account
// Request: GET /account/:AccountID
// Test with: curl -i localhost/account/:AccountID
func getAccount(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		accountID int64
		account   *entity.Account
		err       error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read account from database
	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(account, http.StatusOK, 10, w, r)
}

// createAccount handles requests create an account
// Request: POST /accounts
// Test with: curl -i -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/accounts
func createAccount(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		account = &entity.Account{}
		err     error
	)

	// Decode JSON
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(account); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Set values
	account.Enabled = true

	// Validate resource
	if err = validator.ValidateAccount(account); err != nil {
		return
	}

	// Write account
	if account, err = core.WriteAccount(account, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(account, http.StatusCreated, 0, w, r)
}
