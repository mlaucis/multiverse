/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

const (
	accountNameMin = 3
	accountNameMax = 40

	accountDescriptionMin = 0
	accountDescriptionMax = 100
)

var (
	errorAccountNameSize = fmt.Errorf("account name must be between %d and %d characters", accountNameMin, accountNameMax)
	errorAccountNameType = fmt.Errorf("account name is not a valid alphanumeric sequence")

	errorAccountDescriptionSize = fmt.Errorf("account description must be between %d and %d characters", accountDescriptionMin, accountDescriptionMax)
	errorAccountDescriptionType = fmt.Errorf("account description is not a valid alphanumeric sequence")

	errorAccountIDIsAlreadySet  = fmt.Errorf("account id is already set")
	errorAccountSetNotEnabled   = fmt.Errorf("account cannot be set as disabled")
	errorAccountTokenAlreadySet = fmt.Errorf("account token is already set")
)

// CreateAccount validates an account on create
func CreateAccount(account *entity.Account) error {
	errs := []*error{}

	if !stringLenghtBetween(account.Name, accountNameMin, accountNameMax) {
		errs = append(errs, &errorAccountNameSize)
	}

	if !stringLenghtBetween(account.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, &errorAccountDescriptionSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(account.Name)) {
		errs = append(errs, &errorAccountNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(account.Description)) {
		errs = append(errs, &errorAccountDescriptionType)
	}

	if account.ID != 0 {
		errs = append(errs, &errorAccountIDIsAlreadySet)
	}

	if account.AuthToken != "" {
		errs = append(errs, &errorAccountTokenAlreadySet)
	}

	if len(account.Image) > 0 {
		for _, image := range account.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	return packErrors(errs)
}

// UpdateAccount validates an account on update
func UpdateAccount(account *entity.Account) error {
	errs := []*error{}

	if !stringLenghtBetween(account.Name, accountNameMin, accountNameMax) {
		errs = append(errs, &errorAccountNameSize)
	}

	if !stringLenghtBetween(account.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, &errorAccountDescriptionSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(account.Name)) {
		errs = append(errs, &errorAccountNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(account.Description)) {
		errs = append(errs, &errorAccountDescriptionType)
	}

	if len(account.Image) > 0 {
		for _, image := range account.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	return packErrors(errs)
}
