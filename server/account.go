/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

/**
 * getAccount handles requests to a single account
 * Request: GET /account/:AccountID
 * Test with: curl -i localhost/account/:AccountID
 * @param w, response writer
 * @param r, http request
 */
func getAccount(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// Create mock response
	response := &struct {
		*entity.Account
	}{
		Account: &entity.Account{
			ID:        accountID,
			Name:      "Demo Account",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	// Read account from database

	// Query draft
	/**
	 * SELECT id, name, enabled, created_at, updated_at
	 * FROM accounts
	 * WHERE id={accountID};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * createAccount handles requests create an account
 * Request: POST /account
 * Test with: curl -H "Content-Type: application/json" -d '{"name":"New Account"}' localhost/account
 * @param w, response writer
 * @param r, http request
 */
func createAccount(w http.ResponseWriter, r *http.Request) {

}