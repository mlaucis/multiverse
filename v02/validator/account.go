/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

const (
	accountNameMin = 3
	accountNameMax = 40

	accountDescriptionMin = 0
	accountDescriptionMax = 100
)

var (
	errorAccountNameSize = errors.NewBadRequestError(fmt.Sprintf("account name must be between %d and %d characters", accountNameMin, accountNameMax), "")
	errorAccountNameType = errors.NewBadRequestError(fmt.Sprintf("account name is not a valid alphanumeric sequence"), "")

	errorAccountDescriptionSize = errors.NewBadRequestError(fmt.Sprintf("account description must be between %d and %d characters", accountDescriptionMin, accountDescriptionMax), "")
	errorAccountDescriptionType = errors.NewBadRequestError(fmt.Sprintf("account description is not a valid alphanumeric sequence"), "")

	errorAccountIDIsAlreadySet  = errors.NewBadRequestError(fmt.Sprintf("account id is already set"), "")
	errorAccountSetNotEnabled   = errors.NewBadRequestError(fmt.Sprintf("account cannot be set as disabled"), "")
	errorAccountTokenAlreadySet = errors.NewBadRequestError(fmt.Sprintf("account token is already set"), "")
)

// CreateAccount validates an account on create
func CreateAccount(account *entity.Account) (errs []errors.Error) {
	if !StringLengthBetween(account.Name, accountNameMin, accountNameMax) {
		errs = append(errs, errorAccountNameSize)
	}

	if !StringLengthBetween(account.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, errorAccountDescriptionSize)
	}

	if !alphaNumExtraCharFirst.MatchString(account.Name) {
		errs = append(errs, errorAccountNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(account.Description) {
		errs = append(errs, errorAccountDescriptionType)
	}

	if account.ID != 0 {
		errs = append(errs, errorAccountIDIsAlreadySet)
	}

	if account.AuthToken != "" {
		errs = append(errs, errorAccountTokenAlreadySet)
	}

	if len(account.Images) > 0 {
		if !checkImages(account.Images) {
			errs = append(errs, errorInvalidImageURL)
		}
	}

	return
}

// UpdateAccount validates an account on update
func UpdateAccount(existingAccount, updatedAccount *entity.Account) (errs []errors.Error) {
	if updatedAccount.ID == 0 {
		errs = append(errs, errorAccountIDZero)
	}

	if !StringLengthBetween(updatedAccount.Name, accountNameMin, accountNameMax) {
		errs = append(errs, errorAccountNameSize)
	}

	if !StringLengthBetween(updatedAccount.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, errorAccountDescriptionSize)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccount.Name) {
		errs = append(errs, errorAccountNameType)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedAccount.Description) {
		errs = append(errs, errorAccountDescriptionType)
	}

	if len(updatedAccount.Images) > 0 {
		if !checkImages(updatedAccount.Images) {
			errs = append(errs, errorInvalidImageURL)
		}
	}

	return
}
