/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gluee/backend/entity"
	"github.com/gorilla/mux"
)

// getAccountApplications handles requests to a single application
// Request: GET /app/:AppID
// Test with: curl -i localhost/app/:AppID
func getAccountApplication(w http.ResponseWriter, r *http.Request) {
	var (
		appID uint64
		err   error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened("appId is not set or the value is incorrect", http.StatusBadRequest, r, w)
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
			CreatedAt: apiDemoTime,
			UpdatedAt: apiDemoTime,
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

// getAccountApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
// Test with: curl -i localhost/account/:AccountID/applications
func getAccountApplicationList(w http.ResponseWriter, r *http.Request) {
	var (
		accountID uint64
		err       error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
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
			CreatedAt: apiDemoTime,
			UpdatedAt: apiDemoTime,
		},
		Application: []*entity.Application{
			&entity.Application{
				ID:        1,
				Key:       "abc123def",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: apiDemoTime,
				UpdatedAt: apiDemoTime,
			},
			&entity.Application{
				ID:        2,
				Key:       "abc345def",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: apiDemoTime,
				UpdatedAt: apiDemoTime,
			},
			&entity.Application{
				ID:        3,
				Key:       "abc678ef",
				Name:      "Demo App",
				Enabled:   true,
				CreatedAt: apiDemoTime,
				UpdatedAt: apiDemoTime,
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

// createAccountApplication handles requests create an application
// Request: POST /account/:AccountID/app
// Test with: curl -H "Content-Type: application/json" -d '{"name":"New App"}' localhost/account/:AccountID/app
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
		errorHappened(fmt.Sprintf("%q", err), http.StatusRequestEntityTooLarge, r, w)
	}

	if err := r.Body.Close(); err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusBadRequest, r, w)
	}

	if err := json.Unmarshal(body, &app); err != nil {
		errorHappened(fmt.Sprintf("%q", err), http.StatusBadRequest, r, w)
	}

	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseUint(vars["accountId"], 10, 64); err != nil {
		errorHappened("accountId is not set or the value is incorrect", http.StatusBadRequest, r, w)
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
			CreatedAt: apiDemoTime,
			UpdatedAt: apiDemoTime,
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
