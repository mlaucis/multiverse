/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package validator defines all the functions needed to validate
package validator

import (
	"fmt"
	"net/mail"
	"net/url"
	"regexp"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
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

	errorInvalidImageURL           = fmt.Errorf("image url is not valid")
	errorAccountDoesNotExists      = fmt.Errorf("account does not exists")
	errorApplicationDoesNotExists  = fmt.Errorf("application does not exists")
	errorUserEmailAlreadyExists    = errors.NewBadRequestError("user already exists (1)", "user already exists (1)")
	errorUserUsernameAlreadyExists = errors.NewBadRequestError("user already exists (2)", "user already exists (2)")
	errorEmailAddressInUse         = fmt.Errorf("email address already in use")
	errorUsernameInUse             = fmt.Errorf("username already in use")
)

// packErrors prints errors happened during validation
func packErrors(errs []*error) errors.Error {
	if len(errs) == 0 {
		return nil
	}

	er := ""
	for _, e := range errs {
		er += (*e).Error() + "\n"
	}

	return errors.NewBadRequestError(er[:len(er)-1], er[:len(er)-1])
}

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

func checkImages(images []*entity.Image) bool {
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
