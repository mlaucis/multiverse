package core

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v03/entity"
)

type (
	// Connection interface
	Connection interface {
		// Create adds a user connection and returns the created connection or an error
		Create(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (con *entity.Connection, err []errors.Error)

		// Read returns the connection, if any, between two users
		Read(accountID, applicationID int64, userFromID, userToID uint64) (connection *entity.Connection, err []errors.Error)

		// Update updates a connection in the database and returns the updated connection user or an error
		Update(accountID, applicationID int64, existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err []errors.Error)

		// Delete deletes the connection matching the IDs or an error
		Delete(accountID, applicationID int64, conn *entity.Connection) []errors.Error

		// List returns all connections from a certain user
		List(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, err []errors.Error)

		// FollowedBy returns all connections from a certain user
		FollowedBy(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, err []errors.Error)

		// FollowedBy returns all friends from a certain user
		Friends(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, err []errors.Error)

		// FriendsAndFollowing returns all friends and people followed by the specified user
		FriendsAndFollowing(accountID, applicationID int64, userID uint64) ([]*entity.ApplicationUser, []errors.Error)

		// Confirm confirms a user connection and returns the connection or an error
		Confirm(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (con *entity.Connection, err []errors.Error)

		// WriteEventsToList takes a connection and writes the events to the lists
		WriteEventsToList(accountID, applicationID int64, conn *entity.Connection) []errors.Error

		// DeleteEventsFromLists takes a connection and deletes the events from the lists
		DeleteEventsFromLists(accountID, applicationID int64, userFromID, userToID uint64) []errors.Error

		// SocialConnect creates the connections between a user and his other social peers
		SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, []errors.Error)

		// AutoConnectSocialFriends will connect a user with their its friends on from another social network
		AutoConnectSocialFriends(accountID, applicationID int64, user *entity.ApplicationUser, connectionType string, ourStoredUsersIDs []*entity.ApplicationUser) ([]*entity.ApplicationUser, []errors.Error)

		// Relation returns the relation between two users
		Relation(accountID, applicationID int64, userFromID, userToID uint64) (*entity.Relation, []errors.Error)
	}
)
