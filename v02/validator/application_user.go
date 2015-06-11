/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"strings"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
)

const (
	userNameMin = 2
	userNameMax = 40
)

// CreateUser validates a user on create
func CreateUser(datastore core.ApplicationUser, accountID, applicationID int64, user *entity.ApplicationUser) (errs []errors.Error) {
	if user.FirstName != "" {
		if !StringLengthBetween(user.FirstName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.UserFirstNameSizeError)
		}

		if !alphaNumExtraCharFirst.MatchString(user.FirstName) {
			errs = append(errs, errmsg.UserFirstNameTypeError)
		}
	}

	if user.LastName != "" {
		if !StringLengthBetween(user.LastName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.UserLastNameSizeError)
		}

		if !alphaNumExtraCharFirst.MatchString(user.LastName) {
			errs = append(errs, errmsg.UserLastNameTypeError)
		}
	}

	if user.Username != "" {
		if !StringLengthBetween(user.Username, userNameMin, userNameMax) {
			errs = append(errs, errmsg.UserUsernameSizeError)
		}

		if !alphaNumExtraCharFirst.MatchString(user.Username) {
			errs = append(errs, errmsg.UserUsernameTypeError)
		}
	}

	if user.Username == "" && user.Email == "" {
		errs = append(errs, errmsg.UsernameAndEmailAreEmptyError)
	}

	if user.Email != "" && !IsValidEmail(user.Email) {
		errs = append(errs, errmsg.UserEmailInvalidError)
	}

	if user.URL != "" && !IsValidURL(user.URL, true) {
		errs = append(errs, errmsg.UserURLInvalidError)
	}

	if len(user.Images) > 0 {
		if !checkImages(user.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	if user.Email != "" {
		if isDuplicate, err := DuplicateApplicationUserEmail(datastore, accountID, applicationID, user.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.UserEmailAlreadyExistsError)
			} else {
				errs = append(errs, err...)
			}
		}
	}

	if user.Username != "" {
		if isDuplicate, err := DuplicateApplicationUserUsername(datastore, accountID, applicationID, user.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.UserUsernameAlreadyExistsError)
			} else {
				errs = append(errs, err...)
			}
		}
	}

	return
}

// UpdateUser validates a user on update
func UpdateUser(datastore core.ApplicationUser, accountID, applicationID int64, existingApplicationUser, updatedApplicationUser *entity.ApplicationUser) (errs []errors.Error) {
	if updatedApplicationUser.FirstName != "" {
		if !StringLengthBetween(updatedApplicationUser.FirstName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.UserFirstNameSizeError)
		}

		if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.FirstName) {
			errs = append(errs, errmsg.UserFirstNameTypeError)
		}
	}

	if updatedApplicationUser.LastName != "" {
		if !StringLengthBetween(updatedApplicationUser.LastName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.UserLastNameSizeError)
		}

		if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.LastName) {
			errs = append(errs, errmsg.UserLastNameTypeError)
		}
	}

	if updatedApplicationUser.Username != "" {
		if !StringLengthBetween(updatedApplicationUser.Username, userNameMin, userNameMax) {
			errs = append(errs, errmsg.UserUsernameSizeError)
		}

		if !alphaNumExtraCharFirst.MatchString(updatedApplicationUser.Username) {
			errs = append(errs, errmsg.UserUsernameTypeError)
		}
	}

	if updatedApplicationUser.Username == "" && updatedApplicationUser.Email == "" {
		errs = append(errs, errmsg.UsernameAndEmailAreEmptyError)
	}

	if updatedApplicationUser.Email != "" && !IsValidEmail(updatedApplicationUser.Email) {
		errs = append(errs, errmsg.UserEmailInvalidError)
	}

	if updatedApplicationUser.URL != "" && !IsValidURL(updatedApplicationUser.URL, true) {
		errs = append(errs, errmsg.UserURLInvalidError)
	}

	if len(updatedApplicationUser.Images) > 0 {
		if !checkImages(updatedApplicationUser.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	if updatedApplicationUser.Email != "" && existingApplicationUser.Email != updatedApplicationUser.Email {
		isDuplicate, err := DuplicateApplicationUserEmail(datastore, accountID, applicationID, updatedApplicationUser.Email)
		if isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.EmailAddressInUseError)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if updatedApplicationUser.Username != "" && existingApplicationUser.Username != updatedApplicationUser.Username {
		isDuplicate, err := DuplicateApplicationUserUsername(datastore, accountID, applicationID, updatedApplicationUser.Username)
		if isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.UsernameInUseError)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	return
}

// ApplicationUserCredentialsValid checks is a certain user has the right credentials
func ApplicationUserCredentialsValid(password string, user *entity.ApplicationUser) (errs []errors.Error) {
	pass, err := utils.Base64Decode(user.Password)
	if err != nil {
		return []errors.Error{errmsg.GenericAuthenticationError.UpdateInternalMessage(err.Error())}
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return []errors.Error{errmsg.GenericAuthenticationError.UpdateInternalMessage("invalid password parts")}
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return []errors.Error{errmsg.GenericAuthenticationError.UpdateInternalMessage(err.Error())}
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return []errors.Error{errmsg.GenericAuthenticationError.UpdateInternalMessage(err.Error())}
	}

	encryptedPassword, err := storageHelper.GenerateStrongEncryptedPassword(password, string(salt), string(timestamp))
	if err != nil {
		return []errors.Error{errmsg.GenericAuthenticationError.UpdateInternalMessage(err.Error())}
	}

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errmsg.GenericAuthenticationError.UpdateInternalMessage("password mismatch")}
	}

	return nil
}

// DuplicateApplicationUserEmail checks if the user email is duplicate within the application or not
func DuplicateApplicationUserEmail(datastore core.ApplicationUser, accountID, applicationID int64, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(accountID, applicationID, email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.UserEmailAlreadyExistsError}
		}
	}

	return false, nil
}

// DuplicateApplicationUserUsername checks if the username is duplicate within the application or not
func DuplicateApplicationUserUsername(datastore core.ApplicationUser, accountID, applicationID int64, username string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByUsername(accountID, applicationID, username); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.UserUsernameAlreadyExistsError}
		}
	}

	return false, nil
}
