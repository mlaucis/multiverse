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

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
// Test with: curl -i localhost/account/:AccountID/user/:UserID
func getAccountUser(w http.ResponseWriter, r *http.Request) {
	var (
		accountID   uint64
		userID      uint64
		accountUser *entity.AccountUser
		err         error
	)

	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read userID
	if userID, err = strconv.ParseUint(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if accountUser, err = db.GetAccountUserByID(accountID, userID); err != nil {
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
	var (
		accountID    uint64
		account      *entity.Account
		accountUsers []*entity.AccountUser
		err          error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if account, err = db.GetAccountByID(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if accountUsers, err = db.GetAccountAllUsers(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

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
// Request: POST /account/:AccountID/user
// Test with: curl -H "Content-Type: application/json" -d '{"name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost/account/:AccountID/user
func createAccountUser(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountUser = &entity.AccountUser{}
		accountID   uint64
		err         error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	if accountUser, err = db.AddAccountUser(accountID, accountUser); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(accountUser, http.StatusCreated, 0, w, r)
}
