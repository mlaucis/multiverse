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
	typeMin = 1
	typeMax = 30
)

// CreateEvent validates an event on create
func CreateEvent(datastore core.ApplicationUser, accountID, applicationID int64, event *entity.Event) (errs []errors.Error) {
	if !StringLengthBetween(event.Type, typeMin, typeMax) {
		errs = append(errs, errmsg.ErrEventTypeSize)
	}

	if event.ID != "" {
		errs = append(errs, errmsg.ErrEventIDIsAlreadySet)
	}

	if event.Visibility == 0 {
		errs = append(errs, errmsg.ErrEventMissingVisiblity)
	} else if event.Visibility != 10 && event.Visibility != 20 && event.Visibility != 30 {
		errs = append(errs, errmsg.ErrEventInvalidVisiblity)
	}

	if len(errs) == 0 {
		// Run expensive check only if there are no existing errors
		if exists, err := datastore.ExistsByID(accountID, applicationID, event.UserID); !exists || err != nil {
			if err != nil {
				errs = append(errs, err...)
			} else {
				errs = append(errs, errmsg.ErrApplicationUserNotFound)
			}
		}
	}

	return
}

// UpdateEvent validates an event on update
func UpdateEvent(existingEvent, event *entity.Event) (errs []errors.Error) {
	if !StringLengthBetween(event.Type, typeMin, typeMax) {
		errs = append(errs, errmsg.ErrEventTypeSize)
	}

	if event.Visibility == 0 {
		errs = append(errs, errmsg.ErrEventMissingVisiblity)
	} else if event.Visibility != entity.EventPrivate && event.Visibility != entity.EventConnections && event.Visibility != entity.EventPublic {
		errs = append(errs, errmsg.ErrEventInvalidVisiblity)
	}

	// TODO define more rules for updating an event

	return
}
