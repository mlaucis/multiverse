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

// CreateConnection validates a connection
func CreateConnection(connection *entity.Connection) error {
	errs := []*error{}

	// Validate ApplicationID
	if connection.ApplicationID == 0 {
		errs = append(errs, &errorApplicationIDZero)
	}

	if numInt.Match([]byte(fmt.Sprintf("%d", connection.ApplicationID))) {
		errs = append(errs, &errorApplicationIDType)
	}

	// Validate UserFromID
	if connection.UserFromID == 0 {
		errs = append(errs, &errorUserFromIDZero)
	}

	if numInt.Match([]byte(fmt.Sprintf("%d", connection.UserFromID))) {
		errs = append(errs, &errorUserFromIDType)
	}

	// Validate UserToID
	if connection.UserToID == 0 {
		errs = append(errs, &errorUserToIDZero)
	}

	if numInt.Match([]byte(fmt.Sprintf("%d", connection.UserToID))) {
		errs = append(errs, &errorUserToIDType)
	}

	// Validate Users
	if !userExists(connection.ApplicationID, connection.UserFromID) {
		errs = append(errs, &errorUserDoesNotExists)
	}

	if !userExists(connection.ApplicationID, connection.UserToID) {
		errs = append(errs, &errorUserDoesNotExists)
	}

	return packErrors(errs)
}
