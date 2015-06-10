/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package validator defines all the functions needed to validate
package validator

import (
	"net/mail"
	"net/url"
	"regexp"

	"github.com/tapglue/backend/v02/entity"

	"github.com/satori/go.uuid"
)

var (
	alpha                  = regexp.MustCompile("^[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEEa-zA-Z]+$")
	alphaNum               = regexp.MustCompile("^[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]+$")
	alphaNumExtra          = regexp.MustCompile("^[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z ]+$")
	alphaNumCharFirst      = regexp.MustCompile("^[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEEa-zA-Z ][\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z]+$")
	alphaNumExtraCharFirst = regexp.MustCompile("^[\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEEa-zA-Z ][\u0020-\u007E\uFF61-\uFF9F\uFFA0-\uFFDC\uFFE8-\uFFEE0-9a-zA-Z ]+$")

	num      = regexp.MustCompile("^[-]?[0-9]+$")
	numInt   = regexp.MustCompile("^(?:[-+]?(?:0|[1-9][0-9]*))$")
	numFloat = regexp.MustCompile("^(?:[-+]?(?:[0-9]+))?(?:\\.[0-9]*)?(?:[eE][\\+\\-]?(?:[0-9]+))?$")
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
