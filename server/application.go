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
	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/mysql"
)

// getAccountApplications handles requests to a single application
// Request: GET /app/:AppID
// Test with: curl -i localhost/app/:AppID
func getAccountApplication(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		application *entity.Application
		appID       uint64
		err         error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read appID
	if appID, err = strconv.ParseUint(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if application, err = mysql.GetApplicationByID(appID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	// Write response
	writeResponse(application, http.StatusOK, 10, w, r)
}

// getAccountApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
// Test with: curl -i localhost/account/:AccountID/applications
func getAccountApplicationList(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		account      *entity.Account
		applications []*entity.Application
		accountID    int64
		err          error
	)
	// Read variables from request
	vars := mux.Vars(r)

	// Read accountID
	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	// Read account from database
	if account, err = mysql.GetAccountByID(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if applications, err = mysql.GetAccountAllApplications(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

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

// createAccountApplication handles requests create an application
// Request: POST /account/:AccountID/app
// Test with: curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App"}' localhost/account/:AccountID/app
func createAccountApplication(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		application = &entity.Application{}
		accountID   int64
		err         error
	)
	// Read variables from request
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(application); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	// TODO validation should be added here, for example, name shouldn't be empty ;)

	if application, err = mysql.AddAccountApplication(accountID, application); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(application, http.StatusCreated, 0, w, r)
}
