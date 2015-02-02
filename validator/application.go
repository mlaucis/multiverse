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
)

// CreateApplication validates an application on create
func CreateApplication(application *entity.Application) error {
	errs := []*error{}

	if !stringBetween(application.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, &errorApplicationNameSize)
	}

	if !stringBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, &errorApplicationDescriptionSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(application.Name)) {
		errs = append(errs, &errorApplicationNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(application.Description)) {
		errs = append(errs, &errorApplicationDescriptionType)
	}

	if application.ID != 0 {
		errs = append(errs, &errorApplicationIDIsAlreadySet)
	}

	if application.AccountID == 0 {
		errs = append(errs, &errorAccountIDZero)
	}

	if application.URL != "" && !url.Match([]byte(application.URL)) {
		errs = append(errs, &errorApplicationUserURLInvalid)
	}

	if len(application.Image) > 0 {
		for _, image := range application.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !accountExists(application.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	return packErrors(errs)
}

// UpdateApplication validates an application on update
func UpdateApplication(application *entity.Application) error {
	errs := []*error{}

	if !stringBetween(application.Name, applicationNameMin, applicationNameMax) {
		errs = append(errs, &errorApplicationNameSize)
	}

	if !stringBetween(application.Description, applicationDescriptionMin, applicationDescriptionMax) {
		errs = append(errs, &errorApplicationDescriptionSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(application.Name)) {
		errs = append(errs, &errorApplicationNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(application.Description)) {
		errs = append(errs, &errorApplicationDescriptionType)
	}

	if application.URL != "" && !url.Match([]byte(application.URL)) {
		errs = append(errs, &errorApplicationUserURLInvalid)
	}

	if len(application.Image) > 0 {
		for _, image := range application.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !accountExists(application.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	return packErrors(errs)
}
