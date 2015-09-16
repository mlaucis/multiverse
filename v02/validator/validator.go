// Package validator defines all the functions needed to validate
package validator

import (
	"net/mail"
	"net/url"

	"github.com/tapglue/multiverse/v02/entity"

	"github.com/satori/go.uuid"
)

// IsValidURL checks is an url is valid
func IsValidURL(checkURL string, absolute bool) bool {
	url, err := url.Parse(checkURL)
	if err != nil {
		return false
	}

	if absolute {
		return url.IsAbs()
	}

	return true
}

func checkImages(images map[string]*entity.Image) bool {
	for idx := range images {
		if u, err := url.Parse(images[idx].URL); err != nil || !u.IsAbs() {
			return false
		}
	}

	return true
}

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
