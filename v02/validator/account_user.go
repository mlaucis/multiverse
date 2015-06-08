/**
 * @author Florin Patan <florinpatan@gmail.com>
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
	accountUserNameMin = 2
	accountUserNameMax = 40

	accountUserPasswordMin = 4
	accountUserPasswordMax = 60
)

var (
	errorAccountUserFirstNameSize = errors.NewBadRequestError(fmt.Sprintf("user first name must be between %d and %d characters", accountUserNameMin, accountUserNameMax), "")
	errorAccountUserFirstNameType = errors.NewBadRequestError(fmt.Sprintf("user first name is not a valid alphanumeric sequence"), "")

	errorAccountUserLastNameSize = errors.NewBadRequestError(fmt.Sprintf("user last name must be between %d and %d characters", accountUserNameMin, accountUserNameMax), "")
	errorAccountUserLastNameType = errors.NewBadRequestError(fmt.Sprintf("user last name is not a valid alphanumeric sequence"), "")

	errorAccountUserUsernameSize = errors.NewBadRequestError(fmt.Sprintf("user username must be between %d and %d characters", accountUserNameMin, accountUserNameMax), "")
	errorAccountUserUsernameType = errors.NewBadRequestError(fmt.Sprintf("user username is not a valid alphanumeric sequence"), "")

	errorAccountUserPasswordSize = errors.NewBadRequestError(fmt.Sprintf("user password must be between %d and %d characters", accountUserPasswordMin, accountUserPasswordMax), "")

	errorAccountIDZero = errors.NewBadRequestError(fmt.Sprintf("account id can't be 0"), "")
	errorAccountIDType = errors.NewBadRequestError(fmt.Sprintf("account id is not a valid integer"), "")

	errorAccountUserURLInvalid   = errors.NewBadRequestError(fmt.Sprintf("user url is not a valid url"), "")
	errorAccountUserEmailInvalid = errors.NewBadRequestError(fmt.Sprintf("user email is not valid"), "")
)

// CreateAccountUser validates an account user on create
func CreateAccountUser(datastore core.AccountUser, accountUser *entity.AccountUser) (errs []errors.Error) {
	if !StringLengthBetween(accountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errorAccountUserFirstNameSize)
	}

	if !StringLengthBetween(accountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errorAccountUserLastNameSize)
	}

	if !StringLengthBetween(accountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errorAccountUserUsernameSize)
	}

	if !StringLengthBetween(accountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
		errs = append(errs, errorAccountUserPasswordSize)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.FirstName) {
		errs = append(errs, errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.LastName) {
		errs = append(errs, errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.Username) {
		errs = append(errs, errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars

	if accountUser.AccountID == 0 {
		errs = append(errs, errorAccountIDZero)
	}

	if accountUser.Email == "" || !IsValidEmail(accountUser.Email) {
		errs = append(errs, errorAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !IsValidURL(accountUser.URL, false) {
		errs = append(errs, errorAccountUserURLInvalid)
	}

	if len(accountUser.Images) > 0 {
		if !checkImages(accountUser.Images) {
			errs = append(errs, errorInvalidImageURL)
		}
	}

	if isDuplicate, err := DuplicateAccountUserEmail(datastore, accountUser.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errorUserEmailAlreadyExists)
		} else {
			errs = append(errs, err...)
		}
	}

	if isDuplicate, err := DuplicateAccountUserUsername(datastore, accountUser.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errorUserEmailAlreadyExists)
		} else {
			errs = append(errs, err...)
		}
	}

	return
}

// UpdateAccountUser validates an account user on update
func UpdateAccountUser(datastore core.AccountUser, existingAccountUser, updatedAccountUser *entity.AccountUser) (errs []errors.Error) {
	if !StringLengthBetween(updatedAccountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errorAccountUserFirstNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errorAccountUserLastNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errorAccountUserUsernameSize)
	}

	if updatedAccountUser.Password != "" {
		if !StringLengthBetween(updatedAccountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
			errs = append(errs, errorAccountUserPasswordSize)
		}
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.FirstName) {
		errs = append(errs, errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.LastName) {
		errs = append(errs, errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.Username) {
		errs = append(errs, errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars
	if updatedAccountUser.Email == "" || !IsValidEmail(updatedAccountUser.Email) {
		errs = append(errs, errorAccountUserEmailInvalid)
	}

	if updatedAccountUser.URL != "" && !IsValidURL(updatedAccountUser.URL, true) {
		errs = append(errs, errorAccountUserURLInvalid)
	}

	if len(updatedAccountUser.Images) > 0 {
		if !checkImages(updatedAccountUser.Images) {
			errs = append(errs, errorInvalidImageURL)
		}
	}

	if existingAccountUser.Email != updatedAccountUser.Email {
		if isDuplicate, err := DuplicateAccountUserEmail(datastore, updatedAccountUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errorEmailAddressInUse)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if existingAccountUser.Username != updatedAccountUser.Username {
		if isDuplicate, err := DuplicateAccountUserUsername(datastore, updatedAccountUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errorUsernameInUse)
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
		return []errors.Error{errors.NewBadRequestError("failed to validate account user credentials (1)", err.Error())}
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return []errors.Error{errors.NewBadRequestError("failed to validate account user credentials (2)", "invalid password parts")}
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return []errors.Error{errors.NewBadRequestError("failed to validate account user credentials (3)", err.Error())}
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return []errors.Error{errors.NewBadRequestError("failed to validate account user credentials (4)", err.Error())}
	}

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errors.NewBadRequestError("failed to validate account user credentials (5)", "different passwords")}
	}

	return
}

// DuplicateAccountUserEmail checks if the user e-mail is duplicate within the provided account
func DuplicateAccountUserEmail(datastore core.AccountUser, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errorUserEmailAlreadyExists}
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
			return true, []errors.Error{errorUserUsernameAlreadyExists}
		}
	}

	return false, nil
}
