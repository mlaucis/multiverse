/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

const (
	userNameMin = 2
	userNameMax = 40
)

var (
	errorUserFirstNameSize = fmt.Errorf("user first name must be between %d and %d characters", userNameMin, userNameMax)
	errorUserFirstNameType = fmt.Errorf("user first name is not a valid alphanumeric sequence")

	errorUserLastNameSize = fmt.Errorf("user last name must be between %d and %d characters", userNameMin, userNameMax)
	errorUserLastNameType = fmt.Errorf("user last name is not a valid alphanumeric sequence")

	errorUserUsernameSize = fmt.Errorf("user username must be between %d and %d characters", userNameMin, userNameMax)
	errorUserUsernameType = fmt.Errorf("user username is not a valid alphanumeric sequence")

	errorApplicationIDZero = fmt.Errorf("application id can't be 0")
	errorApplicationIDType = fmt.Errorf("application id is not a valid integer")

	errorAuthTokenInvalid = fmt.Errorf("auth token is invalid")
	errorUserURLInvalid   = fmt.Errorf("user url is not a valid url")
	errorUserEmailInvalid = fmt.Errorf("user email is not valid")

	errorUserIDIsAlreadySet = fmt.Errorf("user id is already set")
)

// CreateUser validates a user on create
func CreateUser(user *entity.User) error {
	errs := []*error{}

	if !stringLenghtBetween(user.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !stringLenghtBetween(user.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !stringLenghtBetween(user.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.FirstName)) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.LastName)) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.Username)) {
		errs = append(errs, &errorUserUsernameType)
	}

	if user.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if user.AuthToken == "" {
		errs = append(errs, &errorAuthTokenInvalid)
	}

	if user.Email == "" || !email.Match([]byte(user.Email)) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if user.URL != "" && !url.Match([]byte(user.URL)) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(user.Image) > 0 {
		for _, image := range user.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !applicationExists(user.AccountID, user.ApplicationID) {
		errs = append(errs, &errorApplicationDoesNotExists)
	}

	return packErrors(errs)
}

// UpdateUser validates a user on update
func UpdateUser(user *entity.User) error {
	errs := []*error{}

	if !stringLenghtBetween(user.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !stringLenghtBetween(user.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !stringLenghtBetween(user.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.FirstName)) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.LastName)) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(user.Username)) {
		errs = append(errs, &errorUserUsernameType)
	}

	if user.AuthToken == "" {
		errs = append(errs, &errorAuthTokenInvalid)
	}

	if user.Email == "" || !email.Match([]byte(user.Email)) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if user.URL != "" && !url.Match([]byte(user.URL)) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(user.Image) > 0 {
		for _, image := range user.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	if !applicationExists(user.AccountID, user.ApplicationID) {
		errs = append(errs, &errorApplicationDoesNotExists)
	}

	return packErrors(errs)
}
