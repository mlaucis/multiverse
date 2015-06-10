/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
)

// IsValidLoginPayload checks if the login payload is valid
func IsValidLoginPayload(loginPayload *entity.LoginPayload) []errors.Error {
	if loginPayload.Email != "" && loginPayload.Username != "" {
		if loginPayload.EmailName == "" {
			return []errors.Error{errmsg.GotBothUsernameAndEmailError}
		}
	}

	if loginPayload.Email == "" && loginPayload.Username == "" && loginPayload.EmailName == "" {
		return []errors.Error{errmsg.GotNoUsernameOrEmailError}
	}

	if loginPayload.Email != "" {
		if !IsValidEmail(loginPayload.Email) {
			return []errors.Error{errmsg.InvalidEmailAddressError}
		}
	}

	if loginPayload.Username != "" {
		if !StringLengthBetween(loginPayload.Username, accountUserNameMin, accountUserNameMax) {
			return []errors.Error{errmsg.AccountUserUsernameSizeError}
		}

		if !alphaNumExtraCharFirst.Match([]byte(loginPayload.Username)) {
			return []errors.Error{errmsg.AccountUserUsernameTypeError}
		}
	}

	return nil
}
