/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"
	"strings"

	"github.com/tapglue/backend/tgerrors"
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
	errorAccountUserFirstNameSize = fmt.Errorf("user first name must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserFirstNameType = fmt.Errorf("user first name is not a valid alphanumeric sequence")

	errorAccountUserLastNameSize = fmt.Errorf("user last name must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserLastNameType = fmt.Errorf("user last name is not a valid alphanumeric sequence")

	errorAccountUserUsernameSize = fmt.Errorf("user username must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserUsernameType = fmt.Errorf("user username is not a valid alphanumeric sequence")

	errorAccountUserPasswordSize = fmt.Errorf("user password must be between %d and %d characters", accountUserPasswordMin, accountUserPasswordMax)

	errorAccountIDZero = fmt.Errorf("account id can't be 0")
	errorAccountIDType = fmt.Errorf("account id is not a valid integer")

	errorAccountUserURLInvalid   = fmt.Errorf("user url is not a valid url")
	errorAccountUserEmailInvalid = fmt.Errorf("user email is not valid")
)

// CreateAccountUser validates an account user on create
func CreateAccountUser(datastore core.AccountUser, accountUser *entity.AccountUser) tgerrors.TGError {
	errs := []*error{}

	if !StringLengthBetween(accountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserFirstNameSize)
	}

	if !StringLengthBetween(accountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserLastNameSize)
	}

	if !StringLengthBetween(accountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserUsernameSize)
	}

	if !StringLengthBetween(accountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
		errs = append(errs, &errorAccountUserPasswordSize)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.FirstName) {
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.LastName) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(accountUser.Username) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars

	if accountUser.AccountID == 0 {
		errs = append(errs, &errorAccountIDZero)
	}

	if accountUser.Email == "" || !IsValidEmail(accountUser.Email) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !IsValidURL(accountUser.URL, false) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if len(accountUser.Image) > 0 {
		if !checkImages(accountUser.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if isDuplicate, err := DuplicateAccountUserEmail(datastore, accountUser.Email); isDuplicate || err != nil {
		if isDuplicate {
			rawErr := errorUserEmailAlreadyExists.RawError()
			errs = append(errs, &rawErr)
		} else {
			rawErr := err.RawError()
			errs = append(errs, &rawErr)
		}
	}

	if isDuplicate, err := DuplicateAccountUserUsername(datastore, accountUser.Username); isDuplicate || err != nil {
		if isDuplicate {
			rawErr := errorUserEmailAlreadyExists.RawError()
			errs = append(errs, &rawErr)
		} else {
			rawErr := err.RawError()
			errs = append(errs, &rawErr)
		}
	}

	return packErrors(errs)
}

// UpdateAccountUser validates an account user on update
func UpdateAccountUser(datastore core.AccountUser, existingAccountUser, updatedAccountUser *entity.AccountUser) tgerrors.TGError {
	errs := []*error{}

	if !StringLengthBetween(updatedAccountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserFirstNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserLastNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserUsernameSize)
	}

	if updatedAccountUser.Password != "" {
		if !StringLengthBetween(updatedAccountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
			errs = append(errs, &errorAccountUserPasswordSize)
		}
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.FirstName) {
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.LastName) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccountUser.Username) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// TODO add validation for password rules such as use all type of chars
	if updatedAccountUser.Email == "" || !IsValidEmail(updatedAccountUser.Email) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if updatedAccountUser.URL != "" && !IsValidURL(updatedAccountUser.URL, true) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if len(updatedAccountUser.Image) > 0 {
		if !checkImages(updatedAccountUser.Image) {
			errs = append(errs, &errorInvalidImageURL)
		}
	}

	if existingAccountUser.Email != updatedAccountUser.Email {
		if isDuplicate, err := DuplicateAccountUserEmail(datastore, updatedAccountUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorEmailAddressInUse)
			} else if err != nil {
				rawErr := err.RawError()
				errs = append(errs, &rawErr)
			}
		}
	}

	if existingAccountUser.Username != updatedAccountUser.Username {
		if isDuplicate, err := DuplicateAccountUserUsername(datastore, updatedAccountUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, &errorUsernameInUse)
			} else if err != nil {
				rawErr := err.RawError()
				errs = append(errs, &rawErr)
			}
		}
	}

	return packErrors(errs)
}

// AccountUserCredentialsValid checks is a certain user has the right credentials
func AccountUserCredentialsValid(password string, user *entity.AccountUser) tgerrors.TGError {
	pass, err := utils.Base64Decode(user.Password)
	if err != nil {
		return tgerrors.NewInternalError("failed to validate account user credentials (1)", err.Error())
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return tgerrors.NewInternalError("failed to validate account user credentials (2)", "invalid password parts")
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return tgerrors.NewInternalError("failed to validate account user credentials (3)", err.Error())
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return tgerrors.NewInternalError("failed to validate account user credentials (4)", err.Error())
	}

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return tgerrors.NewInternalError("failed to validate account user credentials (5)", "different passwords")
	}

	return nil
}

// DuplicateAccountUserEmail checks if the user e-mail is duplicate within the provided account
func DuplicateAccountUserEmail(datastore core.AccountUser, email string) (bool, tgerrors.TGError) {
	if userExists, err := datastore.ExistsByEmail(email); userExists || err != nil {
		if err != nil {
			return false, tgerrors.NewInternalError("failed while retrieving the e-mail address", err.Error())
		} else if userExists {
			return true, errorUserEmailAlreadyExists
		}
	}

	return false, nil
}

// DuplicateAccountUserUsername checks if the username is duplicate within the provided account
func DuplicateAccountUserUsername(datastore core.AccountUser, username string) (bool, tgerrors.TGError) {
	if userExists, err := datastore.ExistsByUsername(username); userExists || err != nil {
		if err != nil {
			return false, tgerrors.NewInternalError("failed while retrieving the username", err.Error())
		} else if userExists {
			return true, errorUserUsernameAlreadyExists
		}
	}

	return false, nil
}
