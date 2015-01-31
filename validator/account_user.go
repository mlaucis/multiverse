/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

const (
	accountUserNameMin = 3
	accountUserNameMax = 25
)

var (
	errorAccountUserFirstNameSize      = fmt.Errorf("user first name must be between %d and %d characters", accountNameMin, accountNameMax)
	errorAccountUserLastNameSize       = fmt.Errorf("user last name must be between %d and %d characters", accountNameMin, accountNameMax)
	errorAccountUserFirstNameNotString = fmt.Errorf("account first name is not a valid alphanumeric sequence")
	errorAccountUserLastNameNotString  = fmt.Errorf("account last name is not a valid alphanumeric sequence")
)

// CreateAccountUser validates an account user
func CreateAccountUser(accountUser *entity.AccountUser) error {
	errs := []*error{}

	if !stringBetween(accountUser.FirstName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserFirstNameSize)
	}

	if !stringBetween(accountUser.LastName, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserLastNameSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.FirstName)) {
		errs = append(errs, &errorAccountUserFirstNameNotString)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.LastName)) {
		errs = append(errs, &errorAccountUserLastNameNotString)
	}

	return packErrors(errs)
}
