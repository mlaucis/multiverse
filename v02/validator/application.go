/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
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
		errs = append(errs, errmsg.ApplicationNameSizeError)
	}

	if !StringLengthBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errmsg.ApplicationDescriptionSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(application.Name) {
		errs = append(errs, errmsg.ApplicationNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(application.Description) {
		errs = append(errs, errmsg.ApplicationDescriptionTypeError)
	}

	if application.ID != 0 {
		errs = append(errs, errmsg.ApplicationIDIsAlreadySetError)
	}

	if application.AccountID == 0 {
		errs = append(errs, errmsg.AccountIDZeroError)
	}

	if application.URL != "" && !IsValidURL(application.URL, true) {
		errs = append(errs, errmsg.ApplicationUserURLInvalidError)
	}

	if len(application.Images) > 0 {
		if !checkImages(application.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	return
}

// UpdateApplication validates an application on update
func UpdateApplication(existingApplication, updatedApplication *entity.Application) (errs []errors.Error) {
	if !StringLengthBetween(updatedApplication.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, errmsg.ApplicationNameSizeError)
	}

	if !StringLengthBetween(updatedApplication.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errmsg.ApplicationDescriptionSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplication.Name) {
		errs = append(errs, errmsg.ApplicationNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplication.Description) {
		errs = append(errs, errmsg.ApplicationDescriptionTypeError)
	}

	if updatedApplication.URL != "" && !IsValidURL(updatedApplication.URL, true) {
		errs = append(errs, errmsg.ApplicationUserURLInvalidError)
	}

	if len(updatedApplication.Images) > 0 {
		if !checkImages(updatedApplication.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	if existingApplication.AuthToken != updatedApplication.AuthToken {
		errs = append(errs, errmsg.ApplicationAuthTokenUpdateNotAllowedError)
	}

	return
}
