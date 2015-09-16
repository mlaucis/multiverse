package validator

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v02/entity"
	"github.com/tapglue/multiverse/v02/errmsg"
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
		errs = append(errs, errmsg.ErrAccountNameSize)
	}

	if !StringLengthBetween(account.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, errmsg.ErrAccountDescriptionSize)
	}

	if account.ID != 0 {
		errs = append(errs, errmsg.ErrAccountIDIsAlreadySet)
	}

	if account.AuthToken != "" {
		errs = append(errs, errmsg.ErrAccountTokenAlreadySet)
	}

	if len(account.Images) > 0 {
		if !checkImages(account.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	return
}

// UpdateAccount validates an account on update
func UpdateAccount(existingAccount, updatedAccount *entity.Account) (errs []errors.Error) {
	if updatedAccount.ID == 0 {
		errs = append(errs, errmsg.ErrAccountIDZero)
	}

	if !StringLengthBetween(updatedAccount.Name, accountNameMin, accountNameMax) {
		errs = append(errs, errmsg.ErrAccountNameSize)
	}

	if !StringLengthBetween(updatedAccount.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, errmsg.ErrAccountDescriptionSize)
	}

	if len(updatedAccount.Images) > 0 {
		if !checkImages(updatedAccount.Images) {
			errs = append(errs, errmsg.ErrInvalidImageURL)
		}
	}

	return
}
