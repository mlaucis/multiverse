package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
)

const (
	applicationNameMin = 2
	applicationNameMax = 40

	applicationDescriptionMin = 0
	applicationDescriptionMax = 100
)

// CreateApplication validates an application on create
func CreateApplication(application *entity.Application) (errs []errors.Error) {
	if !StringLengthBetween(application.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, errmsg.ErrApplicationNameSize)
	}

	if !StringLengthBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errmsg.ErrApplicationDescriptionSize)
	}

	if application.ID != 0 {
		errs = append(errs, errmsg.ErrApplicationIDIsAlreadySet)
	}

	if application.OrgID == 0 {
		errs = append(errs, errmsg.ErrOrgIDZero)
	}

	if application.URL != "" && !IsValidURL(application.URL, true) {
		errs = append(errs, errmsg.ErrApplicationURLInvalid)
	}

	if len(application.Images) > 0 {
		if !checkImages(application.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	return
}

// UpdateApplication validates an application on update
func UpdateApplication(existingApplication, updatedApplication *entity.Application) (errs []errors.Error) {
	if !StringLengthBetween(updatedApplication.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, errmsg.ErrApplicationNameSize)
	}

	if !StringLengthBetween(updatedApplication.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errmsg.ErrApplicationDescriptionSize)
	}

	if updatedApplication.URL != "" && !IsValidURL(updatedApplication.URL, true) {
		errs = append(errs, errmsg.ErrApplicationURLInvalid)
	}

	if len(updatedApplication.Images) > 0 {
		if !checkImages(updatedApplication.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	if existingApplication.AuthToken != updatedApplication.AuthToken {
		errs = append(errs, errmsg.ErrApplicationAuthTokenUpdateNotAllowed)
	}

	return
}