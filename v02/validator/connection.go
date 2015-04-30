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
func CreateConnection(datastore core.ApplicationUser, connection *entity.Connection) errors.Error {
	errs := []*error{}

	if connection.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if connection.UserFromID == 0 {
		errs = append(errs, &errorUserFromIDZero)
	}

	if connection.UserToID == 0 {
		errs = append(errs, &errorUserToIDZero)
	}

	if exists, err := datastore.ExistsByID(connection.AccountID, connection.ApplicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %d does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}
	userFrom, err := datastore.Read(connection.AccountID, connection.ApplicationID, connection.UserFromID)
	if err != nil {
		er := err.Raw()
		errs = append(errs, &er)
	}
	if !userFrom.Activated {
		err := fmt.Errorf("user %s is not activated", userFrom.Username)
		errs = append(errs, &err)
	}

	if exists, err := datastore.ExistsByID(connection.AccountID, connection.ApplicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %d does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}
	userTo, err := datastore.Read(connection.AccountID, connection.ApplicationID, connection.UserToID)
	if err != nil {
		er := err.Raw()
		errs = append(errs, &er)
	}
	if !userTo.Activated {
		err := fmt.Errorf("user %s is not activated", userTo.Username)
		errs = append(errs, &err)
	}

	return packErrors(errs)
}

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(datastore core.ApplicationUser, connection *entity.Connection) errors.Error {
	errs := []*error{}

	if connection.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if connection.UserFromID == 0 {
		errs = append(errs, &errorUserFromIDZero)
	}

	if connection.UserToID == 0 {
		errs = append(errs, &errorUserToIDZero)
	}

	if exists, err := datastore.ExistsByID(connection.AccountID, connection.ApplicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %d does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}

	if exists, err := datastore.ExistsByID(connection.AccountID, connection.ApplicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %d does not exists", connection.UserFromID)
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}

// UpdateConnection validates a connection on update
func UpdateConnection(datastore core.ApplicationUser, existingConnection, updatedConnection *entity.Connection) errors.Error {
	errs := []*error{}

	if updatedConnection.UserFromID == 0 {
		errs = append(errs, &errorUserFromIDZero)
	}

	if updatedConnection.UserToID == 0 {
		errs = append(errs, &errorUserToIDZero)
	}

	if exists, err := datastore.ExistsByID(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserToID); !exists || err != nil {
		if err != nil {
			er := err.Raw()
			errs = append(errs, &er)
		} else {
			err := fmt.Errorf("user %d does not exists", updatedConnection.UserFromID)
			errs = append(errs, &err)
		}
	}

	return packErrors(errs)
}
