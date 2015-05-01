/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
)

const (
	verbMin = 1
	verbMax = 30
)

var (
	errorVerbSize = fmt.Errorf("verb must be between %d and %d characters", verbMin, verbMax)
	errorVerbType = fmt.Errorf("verb is not a valid alphanumeric sequence")

	errorUserIDZero = fmt.Errorf("user id can't be 0")
	errorUserIDType = fmt.Errorf("user id is not a valid integer")

	errorEventIDIsAlreadySet = fmt.Errorf("event id is already set")
)

// CreateEvent validates an event on create
func CreateEvent(datastore core.ApplicationUser, accountID, applicationID int64, event *entity.Event) errors.Error {
	errs := []*error{}

	if event.UserID == 0 {
		errs = append(errs, &errorUserIDZero)
	}

	if !StringLengthBetween(event.Verb, verbMin, verbMax) {
		errs = append(errs, &errorVerbSize)
	}

	if !alphaNumExtraCharFirst.MatchString(event.Verb) {
		errs = append(errs, &errorVerbType)
	}

	if event.ID != 0 {
		errs = append(errs, &errorEventIDIsAlreadySet)
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, event.UserID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %d does not exists", event.UserID)
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}

// UpdateEvent validates an event on update
func UpdateEvent(existingEvent, updatedEvent *entity.Event) errors.Error {
	errs := []*error{}

	if !StringLengthBetween(updatedEvent.Verb, verbMin, verbMax) {
		errs = append(errs, &errorVerbSize)
	}

	if !alphaNumExtraCharFirst.MatchString(updatedEvent.Verb) {
		errs = append(errs, &errorVerbType)
	}

	// TODO define more rules for updating an event

	return packErrors(errs)
}
