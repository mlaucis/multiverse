package validator

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/core"
	"github.com/tapglue/multiverse/v04/entity"
	"github.com/tapglue/multiverse/v04/errmsg"
)

// CreateConnection validates a connection on create
func CreateConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
	if connection.UserFromID == connection.UserToID {
		return []errors.Error{errmsg.ErrConnectionSelfConnectingUser.SetCurrentLocation()}
	}

	if connection.Type != entity.ConnectionTypeFriend && connection.Type != entity.ConnectionTypeFollow {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong.UpdateMessage("unexpected connection type " + connection.Type).SetCurrentLocation()}
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserFromID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ErrApplicationUserNotFound.SetCurrentLocation())
		}
	}
	userFrom, err := datastore.Read(accountID, applicationID, connection.UserFromID, false)
	if err != nil {
		errs = append(errs, err...)
	}
	if !userFrom.Activated {
		errs = append(errs, errmsg.ErrApplicationUserNotActivated.SetCurrentLocation())
	}

	if exists, err := datastore.ExistsByID(accountID, applicationID, connection.UserToID); !exists || err != nil {
		if err != nil {
			errs = append(errs, err...)
		} else {
			errs = append(errs, errmsg.ErrApplicationUserNotFound.SetCurrentLocation())
		}
	}
	userTo, err := datastore.Read(accountID, applicationID, connection.UserToID, false)
	if err != nil {
		errs = append(errs, err...)
	}
	if userTo == nil {
		errs = append(errs, errmsg.ErrApplicationUserNotActivated.SetCurrentLocation())
	} else if !userTo.Activated {
		errs = append(errs, errmsg.ErrApplicationUserNotActivated.SetCurrentLocation())
	}

	return
}

// ConfirmConnection validates a connection on confirmation
func ConfirmConnection(datastore core.ApplicationUser, accountID, applicationID int64, connection *entity.Connection) (errs []errors.Error) {
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
	if updatedConnection.Type != entity.ConnectionTypeFriend && updatedConnection.Type != entity.ConnectionTypeFollow {
		return []errors.Error{errmsg.ErrConnectionTypeIsWrong.UpdateMessage("unexpected connection type " + updatedConnection.Type).SetCurrentLocation()}
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
