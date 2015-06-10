/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
)

const (
	verbMin = 1
	verbMax = 30
)

// CreateEvent validates an event on create
func CreateEvent(datastore core.ApplicationUser, accountID, applicationID int64, event *entity.Event) (errs []errors.Error) {
	if !StringLengthBetween(event.Type, verbMin, verbMax) {
		errs = append(errs, errmsg.VerbSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(event.Type) {
		errs = append(errs, errmsg.VerbTypeError)
	}

	if event.ID != "" {
		errs = append(errs, errmsg.EventIDIsAlreadySetError)
	}

	if event.Visibility == 0 {
		errs = append(errs, errmsg.EventMissingVisiblityError)
	} else if event.Visibility != 10 && event.Visibility != 20 && event.Visibility != 30 {
		errs = append(errs, errmsg.EventInvalidVisiblityError)
	}

	if len(errs) == 0 {
		// Run expensive check only if there are no existing errors
		if exists, err := datastore.ExistsByID(accountID, applicationID, event.UserID); !exists || err != nil {
			if err != nil {
				errs = append(errs, err...)
			} else {
				errs = append(errs, errmsg.ApplicationUserNotFoundError)
			}
		}
	}

	return
}

// UpdateEvent validates an event on update
func UpdateEvent(existingEvent, event *entity.Event) (errs []errors.Error) {
	if !StringLengthBetween(event.Type, verbMin, verbMax) {
		errs = append(errs, errmsg.VerbSizeError)
	}

	if !alphaNumExtraCharFirst.MatchString(event.Type) {
		errs = append(errs, errmsg.VerbTypeError)
	}

	if event.Visibility == 0 {
		errs = append(errs, errmsg.EventMissingVisiblityError)
	} else if event.Visibility != entity.EventPrivate && event.Visibility != entity.EventConnections && event.Visibility != entity.EventPublic {
		errs = append(errs, errmsg.EventInvalidVisiblityError)
	}

	// TODO define more rules for updating an event

	return
}
