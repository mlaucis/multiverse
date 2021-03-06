package validator

import (
	"strings"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/utils"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
	storageHelper "github.com/tapglue/multiverse/v04/storage/helper"
)

const (
	userNameMin = 2
	userNameMax = 40
)

// CreateUser validates a user on create
func CreateUser(datastore core.ApplicationUser, accountID, applicationID int64, user *entity.ApplicationUser) (errs []errors.Error) {
	if user.FirstName != "" {
		if !StringLengthBetween(user.FirstName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.ErrApplicationUserFirstNameSize.SetCurrentLocation())
		}
	}

	if user.LastName != "" {
		if !StringLengthBetween(user.LastName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.ErrApplicationUserLastNameSize.SetCurrentLocation())
		}
	}

	if user.Username != "" {
		if !StringLengthBetween(user.Username, userNameMin, userNameMax) {
			errs = append(errs, errmsg.ErrApplicationUserUsernameSize.SetCurrentLocation())
		}
	}

	if user.Username == "" && user.Email == "" {
		errs = append(errs, errmsg.ErrApplicationUsernameAndEmailAreEmpty.SetCurrentLocation())
	}

	if user.Password == "" {
		errs = append(errs, errmsg.ErrAuthPasswordEmpty.SetCurrentLocation())
	}

	if user.Email != "" && !IsValidEmail(user.Email) {
		errs = append(errs, errmsg.ErrApplicationUserEmailInvalid.SetCurrentLocation())
	}

	if user.Email != "" {
		if isDuplicate, err := DuplicateApplicationUserEmail(datastore, accountID, applicationID, user.Email); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists.SetCurrentLocation())
			} else {
				errs = append(errs, err...)
			}
		}
	}

	if user.Username != "" {
		if isDuplicate, err := DuplicateApplicationUserUsername(datastore, accountID, applicationID, user.Username); isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserUsernameInUse.SetCurrentLocation())
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
			errs = append(errs, errmsg.ErrApplicationUserFirstNameSize.SetCurrentLocation())
		}
	}

	if updatedApplicationUser.LastName != "" {
		if !StringLengthBetween(updatedApplicationUser.LastName, userNameMin, userNameMax) {
			errs = append(errs, errmsg.ErrApplicationUserLastNameSize.SetCurrentLocation())
		}
	}

	if updatedApplicationUser.Username != "" {
		if !StringLengthBetween(updatedApplicationUser.Username, userNameMin, userNameMax) {
			errs = append(errs, errmsg.ErrApplicationUserUsernameSize.SetCurrentLocation())
		}
	}

	if updatedApplicationUser.Username == "" && updatedApplicationUser.Email == "" {
		errs = append(errs, errmsg.ErrApplicationUsernameAndEmailAreEmpty.SetCurrentLocation())
	}

	if updatedApplicationUser.Email != "" && !IsValidEmail(updatedApplicationUser.Email) {
		errs = append(errs, errmsg.ErrApplicationUserEmailInvalid.SetCurrentLocation())
	}

	if updatedApplicationUser.Email != "" && existingApplicationUser.Email != updatedApplicationUser.Email {
		isDuplicate, err := DuplicateApplicationUserEmail(datastore, accountID, applicationID, updatedApplicationUser.Email)
		if isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserEmailAlreadyExists.SetCurrentLocation())
			} else if err != nil {
				errs = append(errs, err...)
			}
		}
	}

	if updatedApplicationUser.Username != "" && existingApplicationUser.Username != updatedApplicationUser.Username {
		isDuplicate, err := DuplicateApplicationUserUsername(datastore, accountID, applicationID, updatedApplicationUser.Username)
		if isDuplicate || err != nil {
			if isDuplicate {
				errs = append(errs, errmsg.ErrApplicationUserUsernameInUse.SetCurrentLocation())
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
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}
	passwordParts := strings.SplitN(string(pass), ":", 3)
	if len(passwordParts) != 3 {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage("invalid password parts").SetCurrentLocation()}
	}

	salt, err := utils.Base64Decode(passwordParts[0])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	timestamp, err := utils.Base64Decode(passwordParts[1])
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	encryptedPassword, err := storageHelper.GenerateStrongEncryptedPassword(password, string(salt), string(timestamp))
	if err != nil {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage(err.Error()).SetCurrentLocation()}
	}

	if encryptedPassword != passwordParts[2] {
		return []errors.Error{errmsg.ErrAuthGeneric.UpdateInternalMessage("password mismatch").SetCurrentLocation()}
	}

	return nil
}

// DuplicateApplicationUserEmail checks if the user email is duplicate within the application or not
func DuplicateApplicationUserEmail(datastore core.ApplicationUser, accountID, applicationID int64, email string) (isDuplicate bool, errs []errors.Error) {
	if userExists, err := datastore.ExistsByEmail(accountID, applicationID, email); userExists || err != nil {
		if err != nil {
			return false, err
		} else if userExists {
			return true, []errors.Error{errmsg.ErrApplicationUserEmailAlreadyExists.SetCurrentLocation()}
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
			return true, []errors.Error{errmsg.ErrApplicationUserUsernameInUse.SetCurrentLocation()}
		}
	}

	return false, nil
}
