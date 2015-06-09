/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

const (
	applicationNameMin = 2
	applicationNameMax = 40

	applicationDescriptionMin = 0
	applicationDescriptionMax = 100
)

var (
	errorApplicationNameSize = errors.NewBadRequestError(fmt.Sprintf("application name must be between %d and %d characters", applicationNameMin, applicationNameMax), "")
	errorApplicationNameType = errors.NewBadRequestError(fmt.Sprintf("application name is not a valid alphanumeric sequence"), "")

	errorApplicationDescriptionSize = errors.NewBadRequestError(fmt.Sprintf("application description must be between %d and %d characters", applicationDescriptionMin, applicationDescriptionMax), "")
	errorApplicationDescriptionType = errors.NewBadRequestError(fmt.Sprintf("application description is not a valid alphanumeric sequence"), "")

	errorApplicationUserURLInvalid = errors.NewBadRequestError(fmt.Sprintf("application url is not a valid url"), "")

	errorApplicationIDIsAlreadySet = errors.NewBadRequestError(fmt.Sprintf("application id is already set"), "")

	errorApplicationAuthTokenUpdateNotAllowed = errors.NewBadRequestError(fmt.Sprintf("not allowed to update the application token"), "")
)

// CreateApplication validates an application on create
func CreateApplication(application *entity.Application) (errs []errors.Error) {
	if !StringLengthBetween(application.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, errorApplicationNameSize)
	}

	if !StringLengthBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errorApplicationDescriptionSize)
	}

	if !alphaNumExtraCharFirst.MatchString(application.Name) {
		errs = append(errs, errorApplicationNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(application.Description) {
		errs = append(errs, errorApplicationDescriptionType)
	}

	if application.ID != 0 {
		errs = append(errs, errorApplicationIDIsAlreadySet)
	}

	if application.AccountID == 0 {
		errs = append(errs, errorAccountIDZero)
	}

	if application.URL != "" && !IsValidURL(application.URL, true) {
		errs = append(errs, errorApplicationUserURLInvalid)
	}

	if len(application.Images) > 0 {
		if !checkImages(application.Images) {
			errs = append(errs, errorInvalidImageURL)
		}
	}

	return
}

// UpdateApplication validates an application on update
func UpdateApplication(existingApplication, updatedApplication *entity.Application) (errs []errors.Error) {
	if !StringLengthBetween(updatedApplication.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, errorApplicationNameSize)
	}

	if !StringLengthBetween(updatedApplication.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, errorApplicationDescriptionSize)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplication.Name) {
		errs = append(errs, errorApplicationNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplication.Description) {
		errs = append(errs, errorApplicationDescriptionType)
	}

	if updatedApplication.URL != "" && !IsValidURL(updatedApplication.URL, true) {
		errs = append(errs, errorApplicationUserURLInvalid)
	}

	if len(updatedApplication.Images) > 0 {
		if !checkImages(updatedApplication.Images) {
			errs = append(errs, errorInvalidImageURL)
		}
	}

	if existingApplication.AuthToken != updatedApplication.AuthToken {
		errs = append(errs, errorApplicationAuthTokenUpdateNotAllowed)
	}

	return
}
