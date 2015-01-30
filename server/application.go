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

// getApplication handles requests to a single application
// Request: GET /account/:AccountID/application/:ID
// Test with: curl -i localhost/account/:AccountID/application/:ID
func getApplication(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		application *entity.Application
		accountID   int64
		appID       int64
		err         error
	)

	// Read vars
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read appID
	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read application
	if application, err = core.ReadApplication(accountID, appID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(application, http.StatusOK, 10, w, r)
}

// getApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
// Test with: curl -i localhost/account/:AccountID/applications
func getApplicationList(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		account      *entity.Account
		applications []*entity.Application
		accountID    int64
		err          error
	)
	// Read vars
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read account
	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Read applications
	if applications, err = core.ReadApplicationList(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Create response
	response := &struct {
		entity.Account
		Applications []*entity.Application `json:"applications"`
	}{
		Account:      *account,
		Applications: applications,
	}

	// Write response
	writeResponse(response, http.StatusOK, 10, w, r)
}

// createApplication handles requests create an application
// Request: POST /account/:AccountID/applications
// Test with: curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App"}' localhost/account/:AccountID/applications
func createApplication(w http.ResponseWriter, r *http.Request) {
	// Validate request
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Declare vars
	var (
		application = &entity.Application{}
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

	// Parse JSON
	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(application); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// Set values
	application.AccountID = accountID
	application.Enabled = true

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	// Write resource
	if application, err = core.WriteApplication(application, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(application, http.StatusCreated, 0, w, r)
}
