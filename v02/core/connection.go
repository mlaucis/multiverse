/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Connection interface
	Connection interface {
		// Create adds a user connection and returns the created connection or an error
		Create(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error)

		// Read returns the connection, if any, between two users
		Read(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err errors.Error)

		// Update updates a connection in the database and returns the updated connection user or an error
		Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error)

		// Delete deletes the connection matching the IDs or an error
		Delete(*entity.Connection) errors.Error

		// List returns all connections from a certain user
		List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err errors.Error)

		// FollowedBy returns all connections from a certain user
		FollowedBy(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err errors.Error)

		// Confirm confirms a user connection and returns the connection or an error
		Confirm(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error)

		// WriteEventsToList takes a connection and writes the events to the lists
		WriteEventsToList(*entity.Connection) errors.Error

		// DeleteEventsFromLists takes a connection and deletes the events from the lists
		DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) errors.Error

		// SocialConnect creates the connections between a user and his other social peers
		SocialConnect(user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, errors.Error)

		// AutoConnectSocialFriends will connect a user with their its friends on from another social network
		AutoConnectSocialFriends(user *entity.ApplicationUser, connectionType string, ourStoredUsersIDs []*entity.ApplicationUser) ([]*entity.ApplicationUser, errors.Error)
	}
)
