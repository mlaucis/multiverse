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

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:ID
// Test with: curl -i localhost/account/:AccountID/user/:ID
func getAccountUser(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		accountID   int64
		userID      int64
		accountUser *entity.AccountUser
		err         error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userID
	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read resource
	if accountUser, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(accountUser, http.StatusOK, 10, w, r)
}

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
// Test with: curl -i localhost/account/:AccountID/users
func getAccountUserList(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		accountID    int64
		account      *entity.Account
		accountUsers []*entity.AccountUser
		err          error
	)
	// Read vars
	vars := mux.Vars(r)

	// Read id
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read resource
	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Read account users
	if accountUsers, err = core.ReadAccountUserList(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create response
	response := &struct {
		entity.Account
		AccountUsers []*entity.AccountUser `json:"accountUsers"`
	}{
		Account:      *account,
		AccountUsers: accountUsers,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost/account/:AccountID/users
func createAccountUser(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		accountUser = &entity.AccountUser{}
		accountID   int64
		err         error
	)
	// Read vars
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Decode JSON
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Set values
	accountUser.AccountID = accountID
	accountUser.Enabled = true

	// Validate resource
	if err = validator.ValidateAccountUser(accountUser); err != nil {
		return
	}

	// Write resource
	if accountUser, err = core.WriteAccountUser(accountUser, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(accountUser, http.StatusCreated, 0, w, r)
}
