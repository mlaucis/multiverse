/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

var (
	errGotBothUsernameAndEmail = fmt.Errorf("both username and email are specified. please use only one of them")
	errGotNoUsernameOrEmail    = fmt.Errorf("both username and email are empty. please use one of them")
	errInvalidEmailAddress     = fmt.Errorf("invalid email address")
)

func IsValidLoginPayload(loginPayload *entity.LoginPayload) error {
	if loginPayload.Email != "" && loginPayload.Username != "" {
		return errGotBothUsernameAndEmail
	}

	if loginPayload.Email == "" && loginPayload.Username == "" {
		return errGotNoUsernameOrEmail
	}

	if loginPayload.Email != "" {
		if !IsValidEmail(loginPayload.Email) {
			return errInvalidEmailAddress
		}
	}

	if loginPayload.Username != "" {
		if !StringLengthBetween(loginPayload.Username, accountUserNameMin, accountUserNameMax) {
			return errorAccountUserUsernameSize
		}

		if !alphaNumExtraCharFirst.Match([]byte(loginPayload.Username)) {
			return errorAccountUserUsernameType
		}
	}

	return nil
}
