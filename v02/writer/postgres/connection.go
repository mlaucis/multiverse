package postgres

import (
	"encoding/json"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

func (p *pg) connectionCreate(msg string) []errors.Error {
	connection := &entity.ConnectionWithIDs{}
	err := json.Unmarshal([]byte(msg), connection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.Create(connection.AccountID, connection.ApplicationID, &connection.Connection, false)
	return er
}

func (p *pg) connectionConfirm(msg string) []errors.Error {
	connection := &entity.ConnectionWithIDs{}
	err := json.Unmarshal([]byte(msg), connection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.Confirm(connection.AccountID, connection.ApplicationID, &connection.Connection, false)
	return er
}

func (p *pg) connectionAutoConnect(msg string) []errors.Error {
	autoConnection := entity.AutoConnectSocialFriends{}
	err := json.Unmarshal([]byte(msg), autoConnection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.AutoConnectSocialFriends(
		autoConnection.User.AccountID,
		autoConnection.User.ApplicationID,
		&autoConnection.User.ApplicationUser,
		autoConnection.Type,
		autoConnection.OurStoredUsersIDs)
	return er
}

func (p *pg) connectionSocialConnect(msg string) []errors.Error {
	socialConnection := &entity.SocialConnection{}
	err := json.Unmarshal([]byte(msg), socialConnection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.SocialConnect(
		socialConnection.User.AccountID,
		socialConnection.User.ApplicationID,
		&socialConnection.User.ApplicationUser,
		socialConnection.Platform,
		socialConnection.SocialFriendsIDs,
		socialConnection.Type,
	)
	return er
}

func (p *pg) connectionUpdate(msg string) []errors.Error {
	updatedConnection := entity.ConnectionWithIDs{}
	err := json.Unmarshal([]byte(msg), &updatedConnection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	existingConnection, er := p.connection.Read(
		updatedConnection.AccountID,
		updatedConnection.ApplicationID,
		updatedConnection.UserFromID,
		updatedConnection.UserToID)
	if er != nil {
		return er
	}

	_, er = p.connection.Update(
		updatedConnection.AccountID,
		updatedConnection.ApplicationID,
		*existingConnection,
		updatedConnection.Connection,
		false)
	return er
}

func (p *pg) connectionDelete(msg string) []errors.Error {
	connection := &entity.ConnectionWithIDs{}
	err := json.Unmarshal([]byte(msg), connection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	return p.connection.Delete(connection.AccountID, connection.ApplicationID, &connection.Connection)
}
