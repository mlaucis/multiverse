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

// getAccountUser handles requests to a single account user
// Request: GET /account/:AccountID/user/:ID
// Test with: curl -i localhost/0.1/account/:AccountID/user/:ID
func getAccountUser(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID   int64
		userID      int64
		accountUser *entity.AccountUser
		err         error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if !validator.ValidateAccountRequestToken(accountID, getReqAuthToken(r)) {
		errorHappened(fmt.Errorf("request is not properly signed"), http.StatusBadRequest, r, w)
		return
	}

	if accountUser, err = core.ReadAccountUser(accountID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(accountUser, http.StatusOK, 10, w, r)
}

// updateAccountUser handles requests update an account user
// Request: PUT /account/:AccountID/user/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost/0.1/account/:AccountID/user/:ID
func updateAccountUser(w http.ResponseWriter, r *http.Request) {
	if err := validatePutCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountUser = &entity.AccountUser{}
		accountID   int64
		userID      int64
		err         error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if !validator.ValidateAccountRequestToken(accountID, getReqAuthToken(r)) {
		errorHappened(fmt.Errorf("request is not properly signed"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if accountUser.ID == 0 {
		accountUser.ID = userID
	}
	if accountUser.AccountID == 0 {
		accountUser.AccountID = accountID
	}

	if err = validator.UpdateAccountUser(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if accountUser, err = core.UpdateAccountUser(accountUser, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(accountUser, http.StatusCreated, 0, w, r)
}

// deleteAccountUser handles requests to delete a single account user
// Request: DELETE /account/:AccountID/user/:ID
// Test with: curl -i -X DELETE localhost/0.1/account/:AccountID/user/:ID
func deleteAccountUser(w http.ResponseWriter, r *http.Request) {
	if err := validateDeleteCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID int64
		userID    int64
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if userID, err = strconv.ParseInt(vars["userId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("userId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if err = core.DeleteAccountUser(accountID, userID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("", http.StatusNoContent, 10, w, r)
}

// getAccountUserList handles requests to list all account users
// Request: GET /account/:AccountID/users
// Test with: curl -i localhost/0.1/account/:AccountID/users
func getAccountUserList(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID    int64
		account      *entity.Account
		accountUsers []*entity.AccountUser
		err          error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if !validator.ValidateAccountRequestToken(accountID, getReqAuthToken(r)) {
		errorHappened(fmt.Errorf("request is not properly signed"), http.StatusBadRequest, r, w)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if accountUsers, err = core.ReadAccountUserList(accountID); err != nil {
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

	writeResponse(response, http.StatusOK, 10, w, r)
}

// createAccountUser handles requests create an account user
// Request: POST /account/:AccountID/users
// Test with: curl -i -H "Content-Type: application/json" -d '{"user_name":"User name", "password":"hmac(256)", "email":"de@m.o"}' localhost/0.1/account/:AccountID/users
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
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if !validator.ValidateAccountRequestToken(accountID, getReqAuthToken(r)) {
		errorHappened(fmt.Errorf("request is not properly signed"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	accountUser.AccountID = accountID

	if err = validator.CreateAccountUser(accountUser); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if accountUser, err = core.WriteAccountUser(accountUser, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(accountUser, http.StatusCreated, 0, w, r)
}
