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
		errs = append(errs, errmsg.ErrAccountUserFirstNameSize)
	}

	if !StringLengthBetween(accountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.ErrAccountUserLastNameSize)
	}

	if !StringLengthBetween(accountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.ErrAccountUserUsernameSize)
	}

	if !StringLengthBetween(accountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
		errs = append(errs, errmsg.ErrAccountUserPasswordSize)
	}

	if accountUser.AccountID == 0 {
		errs = append(errs, errmsg.ErrAccountIDZero)
	}

	if accountUser.Email == "" || !IsValidEmail(accountUser.Email) {
		errs = append(errs, errmsg.ErrAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !IsValidURL(accountUser.URL, false) {
		errs = append(errs, errmsg.ErrAccountUserURLInvalid)
	}

	if len(accountUser.Images) > 0 {
		if !checkImages(accountUser.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	if isDuplicate, err := DuplicateAccountUserEmail(datastore, accountUser.Email); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists)
		} else {
			errs = append(errs, err...)
		}
	}

	if isDuplicate, err := DuplicateAccountUserUsername(datastore, accountUser.Username); isDuplicate || err != nil {
		if isDuplicate {
			errs = append(errs, errmsg.ErrApplicationUserUsernameInUse)
		} else {
			errs = append(errs, err...)
		}
	}

	return
}

// UpdateAccountUser validates an account user on update
func UpdateAccountUser(datastore core.AccountUser, existingAccountUser, updatedAccountUser *entity.AccountUser) (errs []errors.Error) {
	if !StringLengthBetween(updatedAccountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.ErrAccountUserFirstNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.ErrAccountUserLastNameSize)
	}

	if !StringLengthBetween(updatedAccountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, errmsg.ErrAccountUserUsernameSize)
	}

	if updatedAccountUser.Password != "" {
		if !StringLengthBetween(updatedAccountUser.Password, accountUserPasswordMin, accountUserPasswordMax) {
			errs = append(errs, errmsg.ErrAccountUserPasswordSize)
		}
	}

	if updatedAccountUser.Email == "" || !IsValidEmail(updatedAccountUser.Email) {
		errs = append(errs, errmsg.ErrAccountUserEmailInvalid)
	}

	if updatedAccountUser.URL != "" && !IsValidURL(updatedAccountUser.URL, true) {
		errs = append(errs, errmsg.ErrAccountUserURLInvalid)
	}

	if len(updatedAccountUser.Images) > 0 {
		if !checkImages(updatedAccountUser.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	if existingAccountUser.Email != updatedAccountUser.Email {
		if isDuplicate, err := DuplicateAccountUserEmail(datastore, updatedAccountUser.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists)
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if existingAccountUser.Username != updatedAccountUser.Username {
		if isDuplicate, err := DuplicateAccountUserUsername(datastore, updatedAccountUser.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserUsernameInUse)
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
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error())}
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage("invalid password parts")}
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error())}
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error())}
	}

	encryptedPassword := storageHelper.GenerateEncryptedPassword(password, string(salt), string(timestamp))

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errmsg.ErrAuthPasswordMismatch}
	}

	return
}

// DuplicateAccountUserEmail checks if the user e-mail is duplicate within the provided account
func DuplicateAccountUserEmail(datastore core.AccountUser, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.ErrApplicationUserEmailAlreadyExists}
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
			return true, []errors.Error{errmsg.ErrApplicationUserUsernameInUse}
		}
	}

	return false, nil
}
