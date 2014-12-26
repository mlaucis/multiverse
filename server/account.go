/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gluee/backend/db"
	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
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
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
		return
	}

	// Read account from database
	account, err = db.GetAccountByID(accountID)
	if err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(account, http.StatusOK, 10, w, r)
}

// createAccount handles requests create an account
// Request: POST /account
// Test with: curl -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/account
func createAccount(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusBadRequest, r, w)
		return
	}

	var (
		account          = &entity.Account{}
		createdAccountID int64
	)

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&account); err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusBadRequest, r, w)
		return
	}

	query := "INSERT INTO `gluee`.`accounts` (`name`) VALUES (?)"
	result, err := db.GetMaster().Exec(query, account.Name)
	if err != nil {
		errorHappened("Error while saving to database", http.StatusInternalServerError, r, w)
		return
	}

	createdAccountID, err = result.LastInsertId()
	if err != nil {
		errorHappened("error while processing the request", http.StatusInternalServerError, r, w)
	}

	if account, err = db.GetAccountByID(uint64(createdAccountID)); err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(account, http.StatusOK, 0, w, r)
}
