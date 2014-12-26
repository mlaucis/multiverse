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
 * getAccountUser handles requests to a single account user
 * Request: GET /account/:AccountID/user/:UserID
 * Test with: curl -i localhost/account/:AccountID/user/:UserID
 * @param w, response writer
 * @param r, http request
 */
func getAccountUser(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		userID    string
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// Read userID
	// TBD userID validation
	userID = vars["userId"]

	// Create mock response
	response := &struct {
		*entity.AccountUser
	}{
		AccountUser: &entity.AccountUser{
			ID:        userID,
			AccountID: accountID,
			Name:      "Demo User",
			Email:     "demouser@demo.com",
			Enabled:   true,
			LastLogin: "2014-12-20T12:10:10Z",
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	// Read account user from database

	// Query draft
	/**
	 * SELECT id, account_id, name, email, enabled, last_login, created_at, updated_at
	 * FROM account_users
	 * WHERE account_id={accountID} AND id={userID};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * getAccountUserList handles requests to list all account users
 * Request: GET /account/:AccountID/users
 * Test with: curl -i localhost/account/:AccountID/users
 * @param w, response writer
 * @param r, http request
 */
func getAccountUserList(w http.ResponseWriter, r *http.Request) {
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
		entity.Account
		AccountUser []*entity.AccountUser `json:"accountUser"`
	}{
		Account: entity.Account{
			ID:        accountID,
			Name:      "Demo Account",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
		AccountUser: []*entity.AccountUser{
			&entity.AccountUser{
				ID:        "1",
				Name:      "Demo User",
				Email:     "demouser@demo.com",
				Enabled:   true,
				LastLogin: "2014-12-20T12:10:10Z",
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.AccountUser{
				ID:        "2",
				Name:      "Demo User",
				Email:     "demouser@demo.com",
				Enabled:   true,
				LastLogin: "2014-12-20T12:10:10Z",
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.AccountUser{
				ID:        "3",
				Name:      "Demo User",
				Email:     "demouser@demo.com",
				Enabled:   true,
				LastLogin: "2014-12-20T12:10:10Z",
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
		},
	}

	// Read account users from database

	// Query draft
	/**
	 * SELECT id, account_id, name, email, enabled, last_login, created_at, updated_at
	 * FROM account_users
	 * WHERE account_id={accountID};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * createAccountUser handles requests create an account user
 * Request: POST /account/:AccountID/user
 * Test with: curl -H "Content-Type: application/json" -d '{"name":"User name"}' localhost/account/:AccountID/user
 * @param w, response writer
 * @param r, http request
 */
func createAccountUser(w http.ResponseWriter, r *http.Request) {

}