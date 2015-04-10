/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"encoding/json"
	"net/http"

	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/validator"
)

var (
	app core.Application
)

// GetApplication handles requests to a single application
// Request: GET /account/:AccountID/application/:ApplicatonID
func GetApplication(ctx *context.Context) (err tgerrors.TGError) {
	WriteResponse(ctx, ctx.Bag["application"].(*entity.Application), http.StatusOK, 10)
	return
}

// UpdateApplication handles requests updates an application
// Request: PUT /account/:AccountID/application/:ApplicatonID
func UpdateApplication(ctx *context.Context) (err tgerrors.TGError) {
	application := *(ctx.Bag["application"].(*entity.Application))
	if er := json.Unmarshal(ctx.Body, &application); er != nil {
		return tgerrors.NewBadRequestError("failed to update the application (1)\n"+er.Error(), er.Error())
	}

	application.ID = ctx.Bag["applicationID"].(int64)
	application.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.UpdateApplication(ctx.Bag["application"].(*entity.Application), &application); err != nil {
		return
	}

	updatedApplication, err := app.Update(*ctx.Bag["application"].(*entity.Application), application, true)
	if err != nil {
		return
	}

	WriteResponse(ctx, updatedApplication, http.StatusCreated, 0)
	return
}

// DeleteApplication handles requests to delete a single application
// Request: DELETE /account/:AccountID/application/:ApplicatonID
func DeleteApplication(ctx *context.Context) (err tgerrors.TGError) {
	if err = app.Delete(ctx.Bag["accountID"].(int64), ctx.Bag["applicationID"].(int64)); err != nil {
		return
	}

	WriteResponse(ctx, "", http.StatusNoContent, 10)
	return
}

// CreateApplication handles requests create an application
// Request: POST /account/:AccountID/applications
func CreateApplication(ctx *context.Context) (err tgerrors.TGError) {
	var (
		application = &entity.Application{}
	)

	if er := json.Unmarshal(ctx.Body, application); er != nil {
		return tgerrors.NewBadRequestError("failed to create the application (1)\n"+er.Error(), er.Error())
	}

	application.AccountID = ctx.Bag["accountID"].(int64)

	if err = validator.CreateApplication(application); err != nil {
		return
	}

	if application, err = app.Create(application, true); err != nil {
		return
	}

	WriteResponse(ctx, application, http.StatusCreated, 0)
	return
}

// GetApplicationList handles requests list all account applications
// Request: GET /account/:AccountID/applications
func GetApplicationList(ctx *context.Context) (err tgerrors.TGError) {
	var (
		applications []*entity.Application
	)

	if applications, err = app.List(ctx.Bag["accountID"].(int64)); err != nil {
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
