/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

const (
	applicationNameMin = 2
	applicationNameMax = 40

	applicationDescriptionMin = 0
	applicationDescriptionMax = 100
)

var (
	errorApplicationNameSize = fmt.Errorf("application name must be between %d and %d characters", applicationNameMin, applicationNameMax)
	errorApplicationNameType = fmt.Errorf("application name is not a valid alphanumeric sequence")

	errorApplicationDescriptionSize = fmt.Errorf("application description must be between %d and %d characters", applicationDescriptionMin, applicationDescriptionMax)
	errorApplicationDescriptionType = fmt.Errorf("application description is not a valid alphanumeric sequence")

	errorApplicationUserURLInvalid = fmt.Errorf("application url is not a valid url")

	errorApplicationIDIsAlreadySet = fmt.Errorf("application id is already set")

	errorApplicationAuthTokenUpdateNotAllowed = fmt.Errorf("not allowed to update the application token")
)

// CreateApplication validates an application on create
func CreateApplication(application *entity.Application) error {
	errs := []*error{}

	if !StringLengthBetween(application.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, &errorApplicationNameSize)
	}

	if !StringLengthBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, &errorApplicationDescriptionSize)
	}

	if !alphaNumExtraCharFirst.MatchString(application.Name) {
		errs = append(errs, &errorApplicationNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(application.Description) {
		errs = append(errs, &errorApplicationDescriptionType)
	}

	if application.ID != 0 {
		errs = append(errs, &errorApplicationIDIsAlreadySet)
	}

	if application.AccountID == 0 {
		errs = append(errs, &errorAccountIDZero)
	}

	if application.URL != "" && !IsValidURL(application.URL, true) {
		errs = append(errs, &errorApplicationUserURLInvalid)
	}

	if len(application.Image) > 0 {
		if !checkImages(application.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	return packErrors(errs)
}

// UpdateApplication validates an application on update
func UpdateApplication(existingApplication, updatedApplication *entity.Application) error {
	errs := []*error{}

	if !StringLengthBetween(updatedApplication.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, &errorApplicationNameSize)
	}

	if !StringLengthBetween(updatedApplication.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, &errorApplicationDescriptionSize)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplication.Name) {
		errs = append(errs, &errorApplicationNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplication.Description) {
		errs = append(errs, &errorApplicationDescriptionType)
	}

	if updatedApplication.URL != "" && !IsValidURL(updatedApplication.URL, true) {
		errs = append(errs, &errorApplicationUserURLInvalid)
	}

	if len(updatedApplication.Image) > 0 {
		if !checkImages(updatedApplication.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if existingApplication.AuthToken != updatedApplication.AuthToken {
		errs = append(errs, &errorApplicationAuthTokenUpdateNotAllowed)
	}

	return packErrors(errs)
}
