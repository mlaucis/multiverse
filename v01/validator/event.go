/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/entity"
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
func CreateEvent(event *entity.Event) *tgerrors.TGError {
	errs := []*error{}

	if event.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

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

	if !UserExists(event.AccountID, event.ApplicationID, event.UserID) {
		errs = append(errs, &errorUserDoesNotExists)
	}

	return packErrors(errs)
}

// UpdateEvent validates an event on update
func UpdateEvent(existingEvent, updatedEvent *entity.Event) *tgerrors.TGError {
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
