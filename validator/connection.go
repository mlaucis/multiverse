/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package validator

import (
	"fmt"

	"github.com/tapglue/backend/core/entity"
)

var (
	errorUserFromIDZero = fmt.Errorf("user from id can't be 0")
	errorUserFromIDType = fmt.Errorf("user from id is not a valid integer")

	errorUserToIDZero = fmt.Errorf("user to id can't be 0")
	errorUserToIDType = fmt.Errorf("user to id is not a valid integer")
)

// CreateConnection validates a connection on create
func CreateConnection(connection *entity.Connection) error {
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

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(connection *entity.Connection) error {
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
func UpdateConnection(connection *entity.Connection) error {
	errs := []*error{}

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
