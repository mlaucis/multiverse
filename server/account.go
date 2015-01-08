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
	"github.com/tapglue/backend/db"
	"github.com/tapglue/backend/entity"
)

// getAccount handles requests to a single account
// Request: GET /account/:AccountID
// Test with: curl -i localhost/account/:AccountID
func getAccount(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		account   *entity.Account
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read account from database
	if account, err = db.GetAccountByID(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(account, http.StatusOK, 10, w, r)
}

// createAccount handles requests create an account
// Request: POST /account
// Test with: curl -i -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/account
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

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	if account, err = db.AddAccount(account); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(account, http.StatusCreated, 0, w, r)
}
