package validator

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
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
		errs = append(errs, errmsg.ErrApplicationNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errmsg.ErrApplicationDescriptionSize.SetCurrentLocation())
	}

	if application.ID != 0 {
		errs = append(errs, errmsg.ErrApplicationIDIsAlreadySet.SetCurrentLocation())
	}

	if application.OrgID == 0 {
		errs = append(errs, errmsg.ErrOrgIDZero.SetCurrentLocation())
	}

	return
}

// UpdateApplication validates an application on update
func UpdateApplication(existingApplication, updatedApplication *entity.Application) (errs []errors.Error) {
	if !StringLengthBetween(updatedApplication.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, errmsg.ErrApplicationNameSize.SetCurrentLocation())
	}

	if !StringLengthBetween(updatedApplication.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errmsg.ErrApplicationDescriptionSize.SetCurrentLocation())
	}

	if existingApplication.AuthToken != updatedApplication.AuthToken {
		errs = append(errs, errmsg.ErrApplicationAuthTokenUpdateNotAllowed.SetCurrentLocation())
	}

	if existingApplication.BackendToken != updatedApplication.BackendToken {
		errs = append(errs, errmsg.ErrApplicationAuthTokenUpdateNotAllowed.SetCurrentLocation())
	}

	return
}
