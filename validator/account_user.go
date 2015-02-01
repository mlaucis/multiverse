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
	errorAccountUserFirstNameNotString = fmt.Errorf("user first name is not a valid alphanumeric sequence")

	errorAccountUserLastNameSize      = fmt.Errorf("user last name must be between %d and %d characters", accountNameMin, accountNameMax)
	errorAccountUserLastNameNotString = fmt.Errorf("user last name is not a valid alphanumeric sequence")

	errorAccountUserUsernameSize      = fmt.Errorf("user username must be between %d and %d characters", accountNameMin, accountNameMax)
	errorAccountUserUsernameNotString = fmt.Errorf("user username is not a valid alphanumeric sequence")

	errorAccountUserURLInvalid   = fmt.Errorf("user url is not a valid url")
	errorAccountUserEmailInvalid = fmt.Errorf("user email is not valid")
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

	if !stringBetween(accountUser.Username, accountUserNameMin, accountUserNameMax) {
		errs = append(errs, &errorAccountUserUsernameSize)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.FirstName)) {
		errs = append(errs, &errorAccountUserFirstNameNotString)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.LastName)) {
		errs = append(errs, &errorAccountUserLastNameNotString)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.Username)) {
		errs = append(errs, &errorAccountUserUsernameNotString)
	}

	if accountUser.Email == "" || !email.Match(accountUser.Email) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	if accountUser.URL != "" && !url.Match([]byte(accountUser.URL)) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	if !accountExists(accountUser.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	return packErrors(errs)
}
