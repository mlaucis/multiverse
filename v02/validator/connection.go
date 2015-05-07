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
	errorUserFromIDZero = fmt.Errorf("user from id can't be 0")
	errorUserFromIDType = fmt.Errorf("user from id is not a valid integer")

	errorUserToIDZero = fmt.Errorf("user to id can't be 0")
	errorUserToIDType = fmt.Errorf("user to id is not a valid integer")
)

// CreateConnection validates a connection on create
func CreateConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) errors.Error {
	errs := []*error{}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %s does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}
	userFrom, err := datastore.Read(accountID, applicationID, connection.UserFromID)
	if err != nil {
		er := err.Raw()
		errs = append(errs, &er)
	}
	if !userFrom.Activated {
		err := fmt.Errorf("user %s is not activated", userFrom.Username)
		errs = append(errs, &err)
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %s does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}
	userTo, err := datastore.Read(accountID, applicationID, connection.UserToID)
	if err != nil {
		er := err.Raw()
		errs = append(errs, &er)
	}
	if userTo == nil {
		err := fmt.Errorf("user %s not found", connection.UserToID)
		errs = append(errs, &err)
	} else if !userTo.Activated {
		err := fmt.Errorf("user %s is not activated", connection.UserToID)
		errs = append(errs, &err)
	}

	return packErrors(errs)
}

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) errors.Error {
	errs := []*error{}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %s does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %s does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}

// UpdateConnection validates a connection on update
func UpdateConnection(datastore core.ApplicationUser, accountID, applicationID int64, existingConnection, updatedConnection *entity.Connection) errors.Error {
	errs := []*error{}

	if exists, err := datastore.ExistsByID(accountID, applicationID, updatedConnection.UserToID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %s does not exists", updatedConnection.UserFromID)
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}
