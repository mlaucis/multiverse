/**
 * @author Florin Patan <florinpatan@gmail.com>
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
	accountUserNameMin = 2
	accountUserNameMax = 40

	accountUserPasswordMin = 4
	accountUserPasswordMax = 60
)

// CreateAccountUser validates an account user on create
func CreateAccountUser(datastore core.AccountUser, accountUser *entity.AccountUser) (errs []errors.Error) {
	if !StringLengthBetween(accountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.AccountUserFirstNameSizeError)
	}

	if !StringLengthBetween(accountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.AccountUserLastNameSizeError)
	}

	if !StringLengthBetween(accountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.AccountUserUsernameSizeError)
	}

	if !StringLengthBetween(accountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
		errs = append(errs, errmsg.AccountUserPasswordSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.FirstName) {
		errs = append(errs, errmsg.AccountUserFirstNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.LastName) {
		errs = append(errs, errmsg.AccountUserLastNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.Username) {
		errs = append(errs, errmsg.AccountUserUsernameTypeError)
	}

	// TODO add validation for password rules such as use all type of chars

	if accountUser.AccountID == 0 {
		errs = append(errs, errmsg.AccountIDZeroError)
	}

	if accountUser.Email == "" || !IsValidEmail(accountUser.Email) {
		errs = append(errs, errmsg.AccountUserEmailInvalidError)
	}

	if accountUser.URL != "" && !IsValidURL(accountUser.URL, false) {
		errs = append(errs, errmsg.AccountUserURLInvalidError)
	}

	if len(accountUser.Images) > 0 {
		if !checkImages(accountUser.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	if isDuplicate, err := DuplicateAccountUserEmail(datastore, accountUser.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.UserEmailAlreadyExistsError)
		} else {
			errs = append(errs, err...)
		}
	}

	if isDuplicate, err := DuplicateAccountUserUsername(datastore, accountUser.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.UserEmailAlreadyExistsError)
		} else {
			errs = append(errs, err...)
		}
	}

	return
}

// UpdateAccountUser validates an account user on update
func UpdateAccountUser(datastore core.AccountUser, existingAccountUser, updatedAccountUser *entity.AccountUser) (errs []errors.Error) {
	if !StringLengthBetween(updatedAccountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.AccountUserFirstNameSizeError)
	}

	if !StringLengthBetween(updatedAccountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.AccountUserLastNameSizeError)
	}

	if !StringLengthBetween(updatedAccountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.AccountUserUsernameSizeError)
	}

	if updatedAccountUser.Password != "" {
		if !StringLengthBetween(updatedAccountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
			errs = append(errs, errmsg.AccountUserPasswordSizeError)
		}
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.FirstName) {
		errs = append(errs, errmsg.AccountUserFirstNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.LastName) {
		errs = append(errs, errmsg.AccountUserLastNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.Username) {
		errs = append(errs, errmsg.AccountUserUsernameTypeError)
	}

	// TODO add validation for password rules such as use all type of chars
	if updatedAccountUser.Email == "" || !IsValidEmail(updatedAccountUser.Email) {
		errs = append(errs, errmsg.AccountUserEmailInvalidError)
	}

	if updatedAccountUser.URL != "" && !IsValidURL(updatedAccountUser.URL, true) {
		errs = append(errs, errmsg.AccountUserURLInvalidError)
	}

	if len(updatedAccountUser.Images) > 0 {
		if !checkImages(updatedAccountUser.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	if existingAccountUser.Email != updatedAccountUser.Email {
		if isDuplicate, err := DuplicateAccountUserEmail(datastore, updatedAccountUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.EmailAddressInUseError)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if existingAccountUser.Username != updatedAccountUser.Username {
		if isDuplicate, err := DuplicateAccountUserUsername(datastore, updatedAccountUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.UsernameInUseError)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	return
}

// AccountUserCredentialsValid checks is a certain user has the right credentials
func AccountUserCredentialsValid(password string, user *entity.AccountUser) (errs []errors.Error) {
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

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errmsg.PasswordMismatchError}
	}

	return
}

// DuplicateAccountUserEmail checks if the user e-mail is duplicate within the provided account
func DuplicateAccountUserEmail(datastore core.AccountUser, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.UserEmailAlreadyExistsError}
		}
	}

	return false, nil
}

// DuplicateAccountUserUsername checks if the username is duplicate within the provided account
func DuplicateAccountUserUsername(datastore core.AccountUser, username string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByUsername(username); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.UserUsernameAlreadyExistsError}
		}
	}

	return false, nil
}
