/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

const (
	accountUserNameMin = 2
	accountUserNameMax = 40
)

var (
	errorAccountUserFirstNameSize = fmt.Errorf("user first name must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserFirstNameType = fmt.Errorf("user first name is not a valid alphanumeric sequence")

	errorAccountUserLastNameSize = fmt.Errorf("user last name must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserLastNameType = fmt.Errorf("user last name is not a valid alphanumeric sequence")

	errorAccountUserUsernameSize = fmt.Errorf("user username must be between %d and %d characters", accountUserNameMin, accountUserNameMax)
	errorAccountUserUsernameType = fmt.Errorf("user username is not a valid alphanumeric sequence")

	errorAccountIDZero = fmt.Errorf("account id can't be 0")
	errorAccountIDType = fmt.Errorf("account id is not a valid integer")

	errorAccountUserURLInvalid   = fmt.Errorf("user url is not a valid url")
	errorAccountUserEmailInvalid = fmt.Errorf("user email is not valid")
)

// CreateAccountUser validates an account user
func CreateAccountUser(accountUser *entity.AccountUser) error {
	errs := []*error{}

	// Validate names
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
		errs = append(errs, &errorAccountUserFirstNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.LastName)) {
		errs = append(errs, &errorAccountUserLastNameType)
	}

	if !alphaNumExtraCharFirst.Match([]byte(accountUser.Username)) {
		errs = append(errs, &errorAccountUserUsernameType)
	}

	// Validate AccountID
	if accountUser.AccountID == 0 {
		errs = append(errs, &errorAccountIDZero)
	}

	if numInt.Match([]byte(fmt.Sprintf("%d", accountUser.AccountID))) {
		errs = append(errs, &errorAccountIDType)
	}

	// Validate Email
	if accountUser.Email == "" || !email.Match([]byte(accountUser.Email)) {
		errs = append(errs, &errorAccountUserEmailInvalid)
	}

	// Validate URL
	if accountUser.URL != "" && !url.Match([]byte(accountUser.URL)) {
		errs = append(errs, &errorAccountUserURLInvalid)
	}

	// Validate Image
	if len(accountUser.Image) > 0 {
		for _, image := range accountUser.Image {
			if !url.Match([]byte(image.URL)) {
				errs = append(errs, &errorInvalidImageURL)
			}
		}
	}

	// Validate Account
	if !accountExists(accountUser.AccountID) {
		errs = append(errs, &errorAccountDoesNotExists)
	}

	return packErrors(errs)
}
