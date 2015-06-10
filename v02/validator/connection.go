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

// CreateConnection validates a connection on create
func CreateConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ApplicationUserNotFoundError)
		}
	}
	userFrom, err := datastore.Read(accountID, applicationID, connection.UserFromID)
	if err != nil {
		errs = append(errs, err...)
	}
	if !userFrom.Activated {
		errs = append(errs, errmsg.ApplicationUserNotActivated)
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ApplicationUserNotFoundError)
		}
	}
	userTo, err := datastore.Read(accountID, applicationID, connection.UserToID)
	if err != nil {
		errs = append(errs, err...)
	}
	if userTo == nil {
		errs = append(errs, errmsg.ApplicationUserNotActivated)
	} else if !userTo.Activated {
		errs = append(errs, errmsg.ApplicationUserNotActivated)
	}

	return
}

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ApplicationUserNotFoundError)
		}
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ApplicationUserNotFoundError)
		}
	}

	return
}

// UpdateConnection validates a connection on update
func UpdateConnection(datastore core.ApplicationUser, accountID, applicationID int64, existingConnection, updatedConnection *entity.Connection) (errs []errors.Error) {
	if exists, err := datastore.ExistsByID(accountID, applicationID, updatedConnection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ApplicationUserNotFoundError)
		}
	}

	return
}
