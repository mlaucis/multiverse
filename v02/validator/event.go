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

	errorEventIDIsAlreadySet   = fmt.Errorf("event id is already set")
	errorEventMissingVisiblity = fmt.Errorf("event visibility is missing")
	errorEventInvalidVisiblity = fmt.Errorf("event visibility is invalid")
)

// CreateEvent validates an event on create
func CreateEvent(datastore core.ApplicationUser, accountID, applicationID int64, event *entity.Event) errors.Error {
	errs := []*error{}

	if !StringLengthBetween(event.Verb, verbMin, verbMax) {
		errs = append(errs, &errorVerbSize)
	}

	if !alphaNumExtraCharFirst.MatchString(event.Verb) {
		errs = append(errs, &errorVerbType)
	}

	if event.ID != "" {
		errs = append(errs, &errorEventIDIsAlreadySet)
	}

	if event.Visibility == 0 {
		errs = append(errs, &errorEventMissingVisiblity)
	} else if event.Visibility != 10 && event.Visibility != 20 && event.Visibility != 30 {
		errs = append(errs, &errorEventInvalidVisiblity)
	}

	if len(errs) == 0 {
		// Run expensive check only if there are no existing errors
		if exists, err := datastore.ExistsByID(accountID, applicationID, event.UserID); !exists || err != nil {
			if err != nil {
				er := err.Raw()
				errs = append(errs, &er)
			} else {
				err := fmt.Errorf("user %d does not exists", event.UserID)
				errs = append(errs, &err)
			}
		}
	}

	return packErrors(errs)
}

// UpdateEvent validates an event on update
func UpdateEvent(existingEvent, event *entity.Event) errors.Error {
	errs := []*error{}

	if !StringLengthBetween(event.Verb, verbMin, verbMax) {
		errs = append(errs, &errorVerbSize)
	}

	if !alphaNumExtraCharFirst.MatchString(event.Verb) {
		errs = append(errs, &errorVerbType)
	}

	if event.Visibility == 0 {
		errs = append(errs, &errorEventMissingVisiblity)
	} else if event.Visibility != entity.EventPrivate && event.Visibility != entity.EventConnections && event.Visibility != entity.EventPublic {
		errs = append(errs, &errorEventInvalidVisiblity)
	}

	// TODO define more rules for updating an event

	return packErrors(errs)
}
