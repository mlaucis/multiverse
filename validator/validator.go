/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package validator defines all the functions needed to validate
package validator

import (
	"fmt"
	"regexp"

	"github.com/tapglue/backend/core"
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

	email = regexp.MustCompile("^(((([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|((\\x22)((((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(([\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(\\([\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(((\\x20|\\x09)*(\\x0d\\x0a))?(\\x20|\\x09)+)?(\\x22)))@((([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])([a-zA-Z]|\\d|-|\\.|_|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*([a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")
	url   = regexp.MustCompile(`^((ftp|http|https):\/\/)?(\S+(:\S*)?@)?((([1-9]\d?|1\d\d|2[01]\d|22[0-3])(\.(1?\d{1,2}|2[0-4]\d|25[0-5])){2}(?:\.([0-9]\d?|1\d\d|2[0-4]\d|25[0-4]))|((www\.)?)?(([a-z\x{00a1}-\x{ffff}0-9]+-?-?_?)*[a-z\x{00a1}-\x{ffff}0-9]+)(?:\.([a-z\x{00a1}-\x{ffff}]{2,}))?)|localhost)(:(\d{1,5}))?((\/|\?|#)[^\s]*)?$`)

	errorInvalidImageURL          = fmt.Errorf("image url is not valid")
	errorAccountDoesNotExists     = fmt.Errorf("account does not exists")
	errorApplicationDoesNotExists = fmt.Errorf("application does not exists")
	errorUserDoesNotExists        = fmt.Errorf("user does not exists")
)

// packErrors prints errors happened during validation
func packErrors(errs []*error) error {
	if len(errs) == 0 {
		return nil
	}

	er := ""
	for _, e := range errs {
		er += (*e).Error() + "\n"
	}

	return fmt.Errorf(er[:len(er)-1])
}

// stringBetween validates the a strings length
func stringBetween(value string, minLength, maxLength int) bool {
	valueLen := len(value)

	if valueLen < minLength {
		return false
	}

	if valueLen > maxLength {
		return false
	}

	return true
}

// accountExists validates if an account exists and returns the account or an error
func accountExists(accountID int64) bool {
	account, err := core.ReadAccount(accountID)
	if err != nil {
		return false
	}

	return account.Enabled
}

// userExists validates if a user exists and returns it or an error
func userExists(applicationID, userID int64) bool {
	user, err := core.ReadUser(applicationID, userID)
	if err != nil {
		return false
	}

	return user.Enabled
}
