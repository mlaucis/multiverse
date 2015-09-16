package postgres

import (
	"encoding/json"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/entity"
)

func (p *pg) connectionCreate(msg string) []errors.Error {
	connection := &entity.ConnectionWithIDs{}
	err := json.Unmarshal([]byte(msg), connection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.Create(connection.OrgID, connection.AppID, &connection.Connection, false)
	return er
}

func (p *pg) connectionConfirm(msg string) []errors.Error {
	connection := &entity.ConnectionWithIDs{}
	err := json.Unmarshal([]byte(msg), connection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.Confirm(connection.OrgID, connection.AppID, &connection.Connection, false)
	return er
}

func (p *pg) connectionAutoConnect(msg string) []errors.Error {
	autoConnection := entity.AutoConnectSocialFriends{}
	err := json.Unmarshal([]byte(msg), autoConnection)
	if err != nil {
		return []errors.Error{errBadInputJSON.UpdateInternalMessage(err.Error())}
	}

	_, er := p.connection.AutoConnectSocialFriends(
		autoConnection.User.OrgID,
		autoConnection.User.AppID,
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
		socialConnection.User.OrgID,
		socialConnection.User.AppID,
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
		updatedConnection.OrgID,
		updatedConnection.AppID,
		updatedConnection.UserFromID,
		updatedConnection.UserToID)
	if er != nil {
		return er
	}

	_, er = p.connection.Update(
		updatedConnection.OrgID,
		updatedConnection.AppID,
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

	return p.connection.Delete(connection.OrgID, connection.AppID, &connection.Connection)
}
