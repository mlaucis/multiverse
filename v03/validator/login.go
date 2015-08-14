package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/entity"
	"github.com/tapglue/backend/v03/errmsg"
)

// IsValidLoginPayload checks if the login payload is valid
func IsValidLoginPayload(loginPayload *entity.LoginPayload) []errors.Error {
	if loginPayload.Email != "" && loginPayload.Username != "" {
		if loginPayload.EmailName == "" {
			return []errors.Error{errmsg.ErrAuthGotBothUsernameAndEmail}
		}
	}

	if loginPayload.Email == "" && loginPayload.Username == "" && loginPayload.EmailName == "" {
		return []errors.Error{errmsg.ErrAuthGotNoUsernameOrEmail}
	}

	if loginPayload.Email != "" {
		if !IsValidEmail(loginPayload.Email) {
			return []errors.Error{errmsg.ErrAuthInvalidEmailAddress}
		}
	}

	if loginPayload.Username != "" {
		if !StringLengthBetween(loginPayload.Username, memberNameMin, memberNameMax) {
			return []errors.Error{errmsg.ErrMemberUsernameSize}
		}
	}

	return nil
}
