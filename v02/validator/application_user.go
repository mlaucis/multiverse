/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"
	"strings"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
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
func CreateUser(datastore core.ApplicationUser, user *entity.ApplicationUser) errors.Error {
	errs := []*error{}

	if !StringLengthBetween(user.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !StringLengthBetween(user.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !StringLengthBetween(user.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.MatchString(user.FirstName) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(user.LastName) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(user.Username) {
		errs = append(errs, &errorUserUsernameType)
	}

	if user.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if user.Email == "" || !IsValidEmail(user.Email) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if user.URL != "" && !IsValidURL(user.URL, true) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(user.Image) > 0 {
		if !checkImages(user.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if isDuplicate, err := DuplicateApplicationUserEmail(datastore, user.AccountID, user.ApplicationID, user.Email); isDuplicate || err != nil {
		if isDuplicate {
			rawErr := errorUserEmailAlreadyExists.Raw()
			errs = append(errs, &rawErr)
		} else {
			rawErr := err.Raw()
			errs = append(errs, &rawErr)
		}
	}

	if isDuplicate, err := DuplicateApplicationUserUsername(datastore, user.AccountID, user.ApplicationID, user.Username); isDuplicate || err != nil {
		if isDuplicate {
			rawErr := errorUserUsernameAlreadyExists.Raw()
			errs = append(errs, &rawErr)
		} else {
			rawErr := err.Raw()
			errs = append(errs, &rawErr)
		}
	}

	return packErrors(errs)
}

// UpdateUser validates a user on update
func UpdateUser(datastore core.ApplicationUser, existingApplicationUser, updatedApplicationUser *entity.ApplicationUser) errors.Error {
	errs := []*error{}

	if !StringLengthBetween(updatedApplicationUser.FirstName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserFirstNameSize)
	}

	if !StringLengthBetween(updatedApplicationUser.LastName, userNameMin, userNameMax) {
		errs = append(errs, &errorUserLastNameSize)
	}

	if !StringLengthBetween(updatedApplicationUser.Username, userNameMin, userNameMax) {
		errs = append(errs, &errorUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.FirstName) {
		errs = append(errs, &errorUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.LastName) {
		errs = append(errs, &errorUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.Username) {
		errs = append(errs, &errorUserUsernameType)
	}

	if updatedApplicationUser.Email == "" || !IsValidEmail(updatedApplicationUser.Email) {
		errs = append(errs, &errorUserEmailInvalid)
	}

	if updatedApplicationUser.URL != "" && !IsValidURL(updatedApplicationUser.URL, true) {
		errs = append(errs, &errorUserURLInvalid)
	}

	if len(updatedApplicationUser.Image) > 0 {
		if !checkImages(updatedApplicationUser.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if existingApplicationUser.Email != updatedApplicationUser.Email {
		isDuplicate, err := DuplicateApplicationUserEmail(datastore, updatedApplicationUser.AccountID, updatedApplicationUser.ApplicationID, updatedApplicationUser.Email)
		if isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorEmailAddressInUse)
			} else if err != nil {
				rawErr := err.Raw()
				errs = append(errs, &rawErr)
			}
		}
	}

	if existingApplicationUser.Username != updatedApplicationUser.Username {
		isDuplicate, err := DuplicateApplicationUserUsername(datastore, updatedApplicationUser.AccountID, updatedApplicationUser.ApplicationID, updatedApplicationUser.Username)
		if isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorUsernameInUse)
			} else if err != nil {
				rawErr := err.Raw()
				errs = append(errs, &rawErr)
			}
		}
	}

	return packErrors(errs)
}

// ApplicationUserCredentialsValid checks is a certain user has the right credentials
func ApplicationUserCredentialsValid(password string, user *entity.ApplicationUser) errors.Error {
	pass, err := utils.Base64Decode(user.Password)
	if err != nil {
		return errors.NewInternalError("failed to check the account user credentials (1)", err.Error())
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return errors.NewInternalError("failed to check the account user credentials (2)", "invalid password parts")
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return errors.NewInternalError("failed to check the account user credentials (3)", err.Error())
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return errors.NewInternalError("failed to check the account user credentials (4)", err.Error())
	}

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return errors.NewInternalError("failed to check the account user credentials (5)\ninvalid user credentials", "password mismatch")
	}

	return nil
}

// DuplicateApplicationUserEmail checks if the user email is duplicate within the application or not
func DuplicateApplicationUserEmail(datastore core.ApplicationUser, accountID, applicationID int64, email string) (bool, errors.Error) {
	if userExists, err := datastore.ExistsByEmail(accountID, applicationID, email); userExists || err != nil {
		if err != nil {
			return false, errors.NewInternalError("failed to perform email validation (1)", err.Error())
		} else if userExists {
			return true, errorUserEmailAlreadyExists
		}
	}

	return false, nil
}

// DuplicateApplicationUserUsername checks if the username is duplicate within the application or not
func DuplicateApplicationUserUsername(datastore core.ApplicationUser, accountID, applicationID int64, username string) (bool, errors.Error) {
	if userExists, err := datastore.ExistsByUsername(accountID, applicationID, username); userExists || err != nil {
		if err != nil {
			return false, errors.NewInternalError("failed to perform username validation (1)", err.Error())
		} else if userExists {
			return true, errorUserUsernameAlreadyExists
		}
	}

	return false, nil
}
