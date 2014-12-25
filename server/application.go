/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

// Package server holds all the server related logic
package server

import (
 	"encoding/json"
 	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

/**
 * getAccountApplications handles requests to a single application
 * Request: GET /app/:AppID
 * Test with: curl -i localhost/app/:AppID
 * @param w, response writer
 * @param r, http request
 */
func getAccountApplication(w http.ResponseWriter, r *http.Request) {
	var (
		appID uint64
		err   error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// Create mock response
	response := &struct {
		*entity.Application
	}{
		Application: &entity.Application{
			ID:        appID,
			Key:       "abc123def",
			AccountID: 123456,
			Name:      "Demo App",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	// Read account application from database

	// Query draft
	/**
	 * SELECT id, key, name, enabled, created_at, updated_at
	 * FROM applications
	 * WHERE id={apptID};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * getAccountApplicationList handles requests list all account applications
 * Request: GET /account/:AccountID/applications
 * Test with: curl -i localhost/account/:AccountID/applications
 * @param w, response writer
 * @param r, http request
 */
func getAccountApplicationList(w http.ResponseWriter, r *http.Request) {
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
		Application []*entity.Application `json:"application"`
	}{
		Account: entity.Account{
			ID:        accountID,
			Name:      "Demo Account",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
		Application: []*entity.Application{
			&entity.Application{
				ID:        1,
				Key:       "abc123def",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.Application{
				ID:        2,
				Key:       "abc345def",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
			&entity.Application{
				ID:        3,
				Key:       "abc678ef",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: "2014-12-15T10:10:10Z",
				UpdatedAt: "2014-12-20T12:10:10Z",
			},
		},
	}

	// Read account applications from database

	// Query drafts

	/** Query Account
	 * SELECT id, name, enabled, created_at, updated_at
	 * FROM accounts
	 * WHERE id={accountID};
	 */

	/** Query Applications
	 * SELECT id, key, name, enabled, created_at, updated_at
	 * FROM applications
	 * WHERE account_id={accountID};
	 */

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

/**
 * createAccountApplication handles requests create an application
 * Request: POST /account/:AccountID/app
 * Test with: curl -H "Content-Type: application/json" -d '{"name":"New App"}' localhost/account/:AccountID/app
 * @param w, response writer
 * @param r, http request
 */
func createAccountApplication(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		//appName string
		app entity.Application
		err error
	)

	// Validate request
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &app); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, w)
		return
	}

	// TBD Read and validate appname

	// Create mock response
	response := &struct {
		*entity.Application
	}{
		Application: &entity.Application{
			ID:        1,
			Key:       "abc123def",
			AccountID: accountID,
			Name:      "Demo App",
			Enabled:   true,
			CreatedAt: "2014-12-15T10:10:10Z",
			UpdatedAt: "2014-12-20T12:10:10Z",
		},
	}

	// Write account applications to database

	// Query drafts

	/**
	 * INSERT INTO applications (account_id, name)
	 * VALUES ({accountID}, {appName});
	 */

	// Write response
	writeResponse(response, http.StatusCreated, 10, w, r)
}

