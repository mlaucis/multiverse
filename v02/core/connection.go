/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Connection interface
	Connection interface {
		// Create adds a user connection and returns the created connection or an error
		Create(connection *entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError)

		// Read returns the connection, if any, between two users
		Read(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err tgerrors.TGError)

		// Update updates a connection in the database and returns the updated connection user or an error
		Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError)

		// Delete deletes the connection matching the IDs or an error
		Delete(*entity.Connection) (err tgerrors.TGError)

		// List returns all connections from a certain user
		List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err tgerrors.TGError)

		// FollowedBy returns all connections from a certain user
		FollowedBy(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err tgerrors.TGError)

		// Confirm confirms a user connection and returns the connection or an error
		Confirm(connection *entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError)

		// WriteEventsToList takes a connection and writes the events to the lists
		WriteEventsToList(connection *entity.Connection) (err tgerrors.TGError)

		// DeleteEventsFromLists takes a connection and deletes the events from the lists
		DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err tgerrors.TGError)

		// SocialConnect creates the connections between a user and his other social peers
		SocialConnect(user *entity.ApplicationUser, platform string, socialFriendsIDs []string) ([]*entity.ApplicationUser, tgerrors.TGError)

		// AutoConnectSocialFriends will connect a user with their its friends on from another social network
		AutoConnectSocialFriends(user *entity.ApplicationUser, ourStoredUsersIDs []interface{}) (users []*entity.ApplicationUser, err tgerrors.TGError)
	}
)
