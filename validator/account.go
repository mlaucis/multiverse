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
	accountNameMax = 25

	accountDescriptionMin = 10
	accountDescriptionMax = 100
)

var (
	errorAccountNameSize      = fmt.Errorf("account name must be between %d and %d characters", accountNameMin, accountNameMax)
	errorAccountNameNotString = fmt.Errorf("account name is not a valid alphanumeric sequence")

	errorAccountDescriptionSize      = fmt.Errorf("account description must be between %d and %d characters", accountDescriptionMin, accountDescriptionMax)
	errorAccountDescriptionNotString = fmt.Errorf("account description is not a valid alphanumeric sequence")

	errorAccountIDIsAlreadySet = fmt.Errorf("account id is already set")
	errorAccountSetNotEnabled  = fmt.Errorf("account cannot be set as disabled")
)

// CreateAccount validates the account
func CreateAccount(account *entity.Account) error {
	errs := []*error{}

	if !stringBetween(account.Name, accountNameMin, accountNameMax) {
		errs = append(errs, &errorAccountNameSize)
	}

	if !stringBetween(account.Description, accountDescriptionMin, accountDescriptionMax) {
		errs = append(errs, &errorAccountDescriptionSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(account.Name)) {
		errs = append(errs, &errorAccountNameNotString)
	}

	if !alphaNumExtraCharFirst.Match([]byte(account.Description)) {
		errs = append(errs, &errorAccountDescriptionNotString)
	}

	if numFloat.Match([]byte(fmt.Sprintf("%d", account.ID))) {
		errs = append(errs, &errorAccountIDIsAlreadySet)
	}

	if !account.Enabled {
		errs = append(errs, &errorAccountSetNotEnabled)
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
