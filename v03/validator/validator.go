// Package validator defines all the functions needed to validate
package validator

import (
	"net/mail"

	"github.com/satori/go.uuid"
)

// IsValidEmail checks if a string is a valid email address
func IsValidEmail(eMail string) bool {
	_, err := mail.ParseAddress(eMail)
	return err == nil
}

// StringLengthBetween validates the a strings length
func StringLengthBetween(value string, minLength, maxLength int) bool {
	valueLen := len(value)

	if valueLen < minLength {
		return false
	}

	if valueLen > maxLength {
		return false
	}

	return true
}

// IsValidUUID5 checks if the str is a valid UUID or not
func IsValidUUID5(str string) bool {
	uid, err := uuid.FromString(str)

	return err == nil && uid.Version() == 5
}
