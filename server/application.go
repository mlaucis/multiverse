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
)

// getApplication handles requests to a single application
// Request: GET /account/:AccountID/application/:ApplicatonID
func getApplication(ctx *context) {
	var (
		application   *entity.Application
		accountID     int64
		applicationID int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if application, err = core.ReadApplication(accountID, applicationID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, application, http.StatusOK, 10)
}

// updateApplication handles requests updates an application
// Request: PUT /account/:AccountID/application/:ApplicatonID
func updateApplication(ctx *context) {
	var (
		application   = &entity.Application{}
		accountID     int64
		applicationID int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(application); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if application.ID == 0 {
		application.ID = applicationID
	}
	if application.AccountID == 0 {
		application.AccountID = accountID
	}

	if err = validator.UpdateApplication(application); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if application, err = core.UpdateApplication(application, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, application, http.StatusCreated, 0)
}

// deleteApplication handles requests to delete a single application
// Request: DELETE /account/:AccountID/application/:ApplicatonID
func deleteApplication(ctx *context) {
	var (
		accountID     int64
		applicationID int64
		err           error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if applicationID, err = strconv.ParseInt(ctx.vars["applicationId"], 10, 64); err != nil {
		errorHappened(ctx, "applicationId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if err = core.DeleteApplication(accountID, applicationID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, "", http.StatusNoContent, 10)
}

// createApplication handles requests create an application
// Request: POST /account/:AccountID/applications
func createApplication(ctx *context) {
	var (
		application = &entity.Application{}
		accountID   int64
		err         error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	decoder := json.NewDecoder(ctx.body)
	if err = decoder.Decode(application); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	application.AccountID = accountID

	if err = validator.CreateApplication(application); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusBadRequest)
		return
	}

	if application, err = core.WriteApplication(application, true); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	writeResponse(ctx, application, http.StatusCreated, 0)
}

// getApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
func getApplicationList(ctx *context) {
	var (
		account      *entity.Account
		applications []*entity.Application
		accountID    int64
		err          error
	)

	if accountID, err = strconv.ParseInt(ctx.vars["accountId"], 10, 64); err != nil {
		errorHappened(ctx, "accountId is not set or the value is incorrect", http.StatusBadRequest)
		return
	}

	if account, err = core.ReadAccount(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	if applications, err = core.ReadApplicationList(accountID); err != nil {
		errorHappened(ctx, fmt.Sprintf("%s", err), http.StatusInternalServerError)
		return
	}

	response := &struct {
		entity.Account
		Applications []*entity.Application `json:"applications"`
	}{
		Account:      *account,
		Applications: applications,
	}

	writeResponse(ctx, response, http.StatusOK, 10)
}
