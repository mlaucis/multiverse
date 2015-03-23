/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/server/utils"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
)

// getApplication handles requests to a single application
// Request: GET /account/:AccountID/application/:ApplicatonID
func getApplication(ctx *context.Context) {
	utils.WriteResponse(ctx, ctx.Application, http.StatusOK, 10)
}

// updateApplication handles requests updates an application
// Request: PUT /account/:AccountID/application/:ApplicatonID
func updateApplication(ctx *context.Context) {
	var err error

	application := *ctx.Application
	if err = json.NewDecoder(ctx.Body).Decode(&application); err != nil {
		utils.ErrorHappened(ctx, "failed to update the application (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	application.ID = ctx.ApplicationID
	application.AccountID = ctx.AccountID

	if err = validator.UpdateApplication(ctx.Application, &application); err != nil {
		utils.ErrorHappened(ctx, "failed to update the application (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	updatedApplication, err := core.UpdateApplication(*ctx.Application, application, true)
	if err != nil {
		utils.ErrorHappened(ctx, "failed to update the application (3)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
}

// deleteApplication handles requests to delete a single application
// Request: DELETE /account/:AccountID/application/:ApplicatonID
func deleteApplication(ctx *context.Context) {
	if err := core.DeleteApplication(ctx.AccountID, ctx.ApplicationID); err != nil {
		utils.ErrorHappened(ctx, "failed to delete the application (1)", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, "", http.StatusNoContent, 10)
}

// createApplication handles requests create an application
// Request: POST /account/:AccountID/applications
func createApplication(ctx *context.Context) {
	var (
		application = &entity.Application{}
		err         error
	)

	if err = json.NewDecoder(ctx.Body).Decode(application); err != nil {
		utils.ErrorHappened(ctx, "failed to create an application (1)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	application.AccountID = ctx.AccountID

	if err = validator.CreateApplication(application); err != nil {
		utils.ErrorHappened(ctx, "failed to create an application (2)\n"+err.Error(), http.StatusBadRequest, err)
		return
	}

	if application, err = core.WriteApplication(application, true); err != nil {
		utils.ErrorHappened(ctx, "failed to create an application", http.StatusInternalServerError, err)
		return
	}

	utils.WriteResponse(ctx, application, http.StatusCreated, 0)
}

// getApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
func getApplicationList(ctx *context.Context) {
	var (
		applications []*entity.Application
		err          error
	)

	if applications, err = core.ReadApplicationList(ctx.AccountID); err != nil {
		utils.ErrorHappened(ctx, "failed to get the applications list (1)", http.StatusInternalServerError, err)
		return
	}

	response := &struct {
		Applications []*entity.Application `json:"applications"`
	}{
		Applications: applications,
	}

	utils.WriteResponse(ctx, response, http.StatusOK, 10)
}
