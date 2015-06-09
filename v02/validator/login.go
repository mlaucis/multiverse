/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

var (
	errGotBothUsernameAndEmail = errors.NewBadRequestError("both username and email are specified", "both username and email are specified")
	errGotNoUsernameOrEmail    = errors.NewBadRequestError("both username and email are empty", "both username and email are empty")
	errInvalidEmailAddress     = errors.NewBadRequestError("invalid email address", "invalid email address")
)

// IsValidLoginPayload checks if the login payload is valid
func IsValidLoginPayload(loginPayload *entity.LoginPayload) []errors.Error {
	if loginPayload.Email != "" && loginPayload.Username != "" {
		if loginPayload.EmailName == "" {
			return []errors.Error{errGotBothUsernameAndEmail}
		}
	}

	if loginPayload.Email == "" && loginPayload.Username == "" && loginPayload.EmailName == "" {
		return []errors.Error{errGotNoUsernameOrEmail}
	}

	if loginPayload.Email != "" {
		if !IsValidEmail(loginPayload.Email) {
			return []errors.Error{errInvalidEmailAddress}
		}
	}

	if loginPayload.Username != "" {
		if !StringLengthBetween(loginPayload.Username, accountUserNameMin, accountUserNameMax) {
			return []errors.Error{errors.NewFromError(errors.BadRequestError, errorAccountUserUsernameSize, false)}
		}

		if !alphaNumExtraCharFirst.Match([]byte(loginPayload.Username)) {
			return []errors.Error{errors.NewFromError(errors.BadRequestError, errorAccountUserUsernameType, false)}
		}
	}

	return nil
}
