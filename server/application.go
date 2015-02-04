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

// getApplication handles requests to a single application
// Request: GET /account/:AccountID/application/:ID
// Test with: curl -i localhost/0.1/account/:AccountID/application/:ID
func getApplication(w http.ResponseWriter, r *http.Request) {
	if err := validateGetCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		application *entity.Application
		accountID   int64
		appID       int64
		err         error
	)

	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if application, err = core.ReadApplication(accountID, appID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(application, http.StatusOK, 10, w, r)
}

// updateApplication handles requests updates an application
// Request: PUT /account/:AccountID/application/:ID
// Test with: curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App"}' -X PUT localhost/0.1/account/:AccountID/application/:ID
func updateApplication(w http.ResponseWriter, r *http.Request) {
	if err := validatePutCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		application = &entity.Application{}
		accountID   int64
		appID       int64
		err         error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(application); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if application.ID == 0 {
		application.ID = appID
	}
	if application.AccountID == 0 {
		application.AccountID = accountID
	}

	if err = validator.UpdateApplication(application); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if application, err = core.UpdateApplication(application, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(application, http.StatusCreated, 0, w, r)
}

// deleteApplication handles requests to delete a single application
// Request: DELETE /account/:AccountID/application/:ID
// Test with: curl -i -X DELETE localhost/0.1/account/:AccountID/application/:ID
func deleteApplication(w http.ResponseWriter, r *http.Request) {
	if err := validateDeleteCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		accountID int64
		appID     int64
		err       error
	)
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if appID, err = strconv.ParseInt(vars["appId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("appId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if err = core.DeleteApplication(accountID, appID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse("", http.StatusNoContent, 10, w, r)
}

// getApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
// Test with: curl -i localhost/0.1/account/:AccountID/applications
func getApplicationList(w http.ResponseWriter, r *http.Request) {
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
	vars := mux.Vars(r)

	if accountID, err = strconv.ParseInt(vars["accountId"], 10, 64); err != nil {
		errorHappened(fmt.Errorf("accountId is not set or the value is incorrect"), http.StatusBadRequest, r, w)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	if applications, err = core.ReadApplicationList(accountID); err != nil {
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

	writeResponse(response, http.StatusOK, 10, w, r)
}

// createApplication handles requests create an application
// Request: POST /account/:AccountID/applications
// Test with: curl -i -H "Content-Type: application/json" -d '{"key": "hmac(256)", "name":"New App"}' localhost/0.1/account/:AccountID/applications
func createApplication(w http.ResponseWriter, r *http.Request) {
	if err := validatePostCommon(w, r); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	var (
		application = &entity.Application{}
		accountID   int64
		err         error
	)
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

	application.AccountID = accountID

	if err = validator.CreateApplication(application); err != nil {
		errorHappened(err, http.StatusBadRequest, r, w)
		return
	}

	if application, err = core.WriteApplication(application, true); err != nil {
		errorHappened(err, http.StatusInternalServerError, r, w)
		return
	}

	writeResponse(application, http.StatusCreated, 0, w, r)
}
