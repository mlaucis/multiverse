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
)

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:UserID
// Test with: curl -i localhost/account/:AccountID/user/:UserID
func getAccountUser(w http.ResponseWriter, r *http.Request) {
	var (
		accountID   int64
		userID      int64
		accountUser *entity.AccountUser
		err         error
	)

	// Read variables from request
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

	if accountUser, err = core.GetAccountUserByID(accountID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(accountUser, http.StatusOK, 10, w, r)
}

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
// Test with: curl -i localhost/account/:AccountID/users
// func getAccountUserList(w http.ResponseWriter, r *http.Request) {
// 	var (
// 		accountID    int64
// 		account      *entity.Account
// 		accountUsers []*entity.AccountUser
// 		err          error
// 	)
// 	// Read variables from request
// 	vars := mux.Vars(r)

// 	// Read accountID
// 	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
// 		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
// 		return
// 	}

// 	if account, err = core.GetAccountByID(accountID); err != nil {
// 		errorHappened(err, http.StatusInternalServerError, r, w)
// 		return
// 	}

// 	if accountUsers, err = core.GetAccountAllUsers(accountID); err != nil {
// 		errorHappened(err, http.StatusInternalServerError, r, w)
// 		return
// 	}

// 	response := &struct {
// 		entity.Account
// 		AccountUsers []*entity.AccountUser `json:"accountUsers"`
// 	}{
// 		Account:      *account,
// 		AccountUsers: accountUsers,
// 	}

// 	// Write response
// 	writeResponse(response, http.StatusOK, 10, w, r)
// }

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost/account/:AccountID/users
func createAccountUser(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountUser = &entity.AccountUser{}
		accountID   int64
		err         error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	accountUser.Enabled = true

	if accountUser, err = core.AddAccountUser(accountID, accountUser, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(accountUser, http.StatusCreated, 0, w, r)
}
