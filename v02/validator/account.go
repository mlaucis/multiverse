/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
)

const (
	accountNameMin = 3
	accountNameMax = 40

	accountDescriptionMin = 0
	accountDescriptionMax = 100
)

// CreateAccount validates an account on create
func CreateAccount(account *entity.Account) (errs []errors.Error) {
	if !StringLengthBetween(account.Name, accountNameMin, accountNameMax) {
		errs = append(errs, errmsg.AccountNameSizeError)
	}

	if !StringLengthBetween(account.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, errmsg.AccountDescriptionSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(account.Name) {
		errs = append(errs, errmsg.AccountNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(account.Description) {
		errs = append(errs, errmsg.AccountDescriptionTypeError)
	}

	if account.ID != 0 {
		errs = append(errs, errmsg.AccountIDIsAlreadySetError)
	}

	if account.AuthToken != "" {
		errs = append(errs, errmsg.AccountTokenAlreadySetError)
	}

	if len(account.Images) > 0 {
		if !checkImages(account.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	return
}

// UpdateAccount validates an account on update
func UpdateAccount(existingAccount, updatedAccount *entity.Account) (errs []errors.Error) {
	if updatedAccount.ID == 0 {
		errs = append(errs, errmsg.AccountIDZeroError)
	}

	if !StringLengthBetween(updatedAccount.Name, accountNameMin, accountNameMax) {
		errs = append(errs, errmsg.AccountNameSizeError)
	}

	if !StringLengthBetween(updatedAccount.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, errmsg.AccountDescriptionSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccount.Name) {
		errs = append(errs, errmsg.AccountNameTypeError)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccount.Description) {
		errs = append(errs, errmsg.AccountDescriptionTypeError)
	}

	if len(updatedAccount.Images) > 0 {
		if !checkImages(updatedAccount.Images) {
			errs = append(errs, errmsg.InvalidImageURLError)
		}
	}

	return
}
