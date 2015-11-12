package validator

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
)

// CreateConnection validates a connection on create
func CreateConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
	if !connection.IsValidState() {
		return []errors.Error{errmsg.ErrConnectionStateInvalid.UpdateInternalMessage("connection state is invalid. got:" + string(connection.State)).SetCurrentLocation()}
	}

	if connection.UserFromID == connection.UserToID {
		return []errors.Error{errmsg.ErrConnectionSelfConnectingUser.SetCurrentLocation()}
	}

	if connection.Type != entity.ConnectionTypeFriend && connection.Type != entity.ConnectionTypeFollow {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong.UpdateMessage("unexpected connection type " + string(connection.Type)).SetCurrentLocation()}
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ErrApplicationUserNotFound.SetCurrentLocation())
		}
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ErrApplicationUserNotFound.SetCurrentLocation())
		}
	}

	return
}

// UpdateConnection validates a connection on update
func UpdateConnection(datastore core.ApplicationUser, accountID, applicationID int64, existingConnection, updatedConnection *entity.Connection) (errs []errors.Error) {
	if !updatedConnection.IsValidType() {
		errs = append(errs, errmsg.ErrConnectionTypeIsWrong.UpdateMessage("unexpected connection type "+string(updatedConnection.Type)).SetCurrentLocation())
	}

	if !updatedConnection.IsValidState() {
		errs = append(errs, errmsg.ErrConnectionStateInvalid.UpdateMessage("unexpected connection state "+string(updatedConnection.State)).SetCurrentLocation())
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, updatedConnection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ErrApplicationUserNotFound.SetCurrentLocation())
		}
	}

	return
}
