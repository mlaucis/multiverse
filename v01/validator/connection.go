/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/core"
	"github.com/tapglue/backend/v01/entity"
)

var (
	errorUserFromIDZero = fmt.Errorf("user from id can't be 0")
	errorUserFromIDType = fmt.Errorf("user from id is not a valid integer")

	errorUserToIDZero = fmt.Errorf("user to id can't be 0")
	errorUserToIDType = fmt.Errorf("user to id is not a valid integer")
)

// CreateConnection validates a connection on create
func CreateConnection(connection *entity.Connection) tgerrors.TGError {
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

	if !UserExists(connection.AccountID, connection.ApplicationID, connection.UserFromID) {
		errs = append(errs, &errorUserDoesNotExists)
	}
	userFrom, err := core.ReadApplicationUser(connection.AccountID, connection.ApplicationID, connection.UserFromID)
	if err != nil {
		errs = append(errs, &errorUserDoesNotExists)
	}
	if !userFrom.Activated {
		err := fmt.Errorf("user %s is not activated", userFrom.Username)
		errs = append(errs, &err)
	}

	if !UserExists(connection.AccountID, connection.ApplicationID, connection.UserToID) {
		errs = append(errs, &errorUserDoesNotExists)
	}
	userTo, err := core.ReadApplicationUser(connection.AccountID, connection.ApplicationID, connection.UserToID)
	if err != nil {
		errs = append(errs, &errorUserDoesNotExists)
	}
	if !userTo.Activated {
		err := fmt.Errorf("user %s is not activated", userTo.Username)
		errs = append(errs, &err)
	}

	return packErrors(errs)
}

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(connection *entity.Connection) tgerrors.TGError {
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

	if !UserExists(connection.AccountID, connection.ApplicationID, connection.UserFromID) {
		errs = append(errs, &errorUserDoesNotExists)
	}

	if !UserExists(connection.AccountID, connection.ApplicationID, connection.UserToID) {
		errs = append(errs, &errorUserDoesNotExists)
	}

	return packErrors(errs)
}

// UpdateConnection validates a connection on update
func UpdateConnection(existingConnection, updatedConnection *entity.Connection) tgerrors.TGError {
	errs := []*error{}

	if updatedConnection.UserFromID == 0 {
		errs = append(errs, &errorUserFromIDZero)
	}

	if updatedConnection.UserToID == 0 {
		errs = append(errs, &errorUserToIDZero)
	}

	if !UserExists(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserToID) {
		errs = append(errs, &errorUserDoesNotExists)
	}

	return packErrors(errs)
}
