/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/entity"
)

var (
	errGotBothUsernameAndEmail = tgerrors.NewBadRequestError("both username and email are specified", "both username and email are specified")
	errGotNoUsernameOrEmail    = tgerrors.NewBadRequestError("both username and email are empty", "both username and email are empty")
	errInvalidEmailAddress     = tgerrors.NewBadRequestError("invalid email address", "invalid email address")
)

// IsValidLoginPayload checks if the login payload is valid
func IsValidLoginPayload(loginPayload *entity.LoginPayload) *tgerrors.TGError {
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
			return tgerrors.NewFromError(tgerrors.TGBadRequestError, errorAccountUserUsernameSize, false)
		}

		if !alphaNumExtraCharFirst.Match([]byte(loginPayload.Username)) {
			return tgerrors.NewFromError(tgerrors.TGBadRequestError, errorAccountUserUsernameType, false)
		}
	}

	return nil
}
