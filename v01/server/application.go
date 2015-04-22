/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v01/context"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
	"github.com/tapglue/backend/v01/validator"
)

// getApplication handles requests to a single application
// Request: GET /account/:AccountID/application/:ApplicatonID
func getApplication(ctx *context.Context) (err errors.Error) {
	WriteResponse(ctx, ctx.Application, http.StatusOK, 10)
	return
}

// updateApplication handles requests updates an application
// Request: PUT /account/:AccountID/application/:ApplicatonID
func updateApplication(ctx *context.Context) (err errors.Error) {
	application := *ctx.Application
	if er := json.Unmarshal(ctx.Body, &application); er != nil {
		return errors.NewBadRequestError("failed to update the application (1)\n"+er.Error(), er.Error())
	}

	application.ID = ctx.ApplicationID
	application.AccountID = ctx.AccountID

	if err = validator.UpdateApplication(ctx.Application, &application); err != nil {
		return
	}

	updatedApplication, err := core.UpdateApplication(*ctx.Application, application, true)
	if err != nil {
		return
	}

	WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
	return
}

// deleteApplication handles requests to delete a single application
// Request: DELETE /account/:AccountID/application/:ApplicatonID
func deleteApplication(ctx *context.Context) (err errors.Error) {
	if err = core.DeleteApplication(ctx.AccountID, ctx.ApplicationID); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// createApplication handles requests create an application
// Request: POST /account/:AccountID/applications
func createApplication(ctx *context.Context) (err errors.Error) {
	var (
		application = &entity.Application{}
	)

	if er := json.Unmarshal(ctx.Body, application); er != nil {
		return errors.NewBadRequestError("failed to create the application (1)\n"+er.Error(), er.Error())
	}

	application.AccountID = ctx.AccountID

	if err = validator.CreateApplication(application); err != nil {
		return
	}

	if application, err = core.WriteApplication(application, true); err != nil {
		return
	}

	WriteResponse(ctx, application, http.StatusCreated, 0)
	return
}

// getApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
func getApplicationList(ctx *context.Context) (err errors.Error) {
	var (
		applications []*entity.Application
	)

	if applications, err = core.ReadApplicationList(ctx.AccountID); err != nil {
		return
	}

	response := &struct {
		Applications []*entity.Application `json:"applications"`
	}{
		Applications: applications,
	}

	WriteResponse(ctx, response, http.StatusOK, 10)
	return
}
