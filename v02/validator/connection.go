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

var (
	errorUserFromIDZero = errors.NewBadRequestError(fmt.Sprintf("user from id can't be 0"), "")
	errorUserFromIDType = errors.NewBadRequestError(fmt.Sprintf("user from id is not a valid integer"), "")

	errorUserToIDZero = errors.NewBadRequestError(fmt.Sprintf("user to id can't be 0"), "")
	errorUserToIDType = errors.NewBadRequestError(fmt.Sprintf("user to id is not a valid integer"), "")
)

// CreateConnection validates a connection on create
func CreateConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s does not exists", connection.UserFromID), ""))
		}
	}
	userFrom, err := datastore.Read(accountID, applicationID, connection.UserFromID)
	if err != nil {
		errs = append(errs, err...)
	}
	if !userFrom.Activated {
		errs = append(errs, errors.NewInternalError(fmt.Sprintf("user %s is not activated", userFrom.Username), ""))
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s does not exists", connection.UserFromID), ""))
		}
	}
	userTo, err := datastore.Read(accountID, applicationID, connection.UserToID)
	if err != nil {
		errs = append(errs, err...)
	}
	if userTo == nil {
		errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s not found", connection.UserToID), ""))
	} else if !userTo.Activated {
		errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s is not activated", connection.UserToID), ""))
	}

	return
}

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s does not exists", connection.UserFromID), ""))
		}
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s does not exists", connection.UserFromID), ""))
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
			errs = append(errs, errors.NewNotFoundError(fmt.Sprintf("user %s does not exists", updatedConnection.UserFromID), ""))
		}
	}

	return
}
