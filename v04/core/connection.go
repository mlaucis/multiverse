package core

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/entity"
)

// Connection interface
type Connection interface {
	// Create creates a user connection and returns an error if it occured
	Create(
		accountID, applicationID int64,
		connection *entity.Connection,
	) (err []errors.Error)

	// Read returns the connection between two users of the requested type, if any, between two users
	Read(
		accountID, applicationID int64,
		userFromID, userToID uint64,
		connectionType entity.ConnectionTypeType,
	) (connection *entity.Connection, err []errors.Error)

	// Update updates a connection in the database and returns the updated connection user or an error
	Update(
		accountID, applicationID int64,
		existingConnection, updatedConnection *entity.Connection,
		retrieve bool,
	) (connection *entity.Connection, err []errors.Error)

	// Delete deletes the connection matching the IDs or an error
	Delete(
		accountID, applicationID int64,
		userFromID, userToID uint64,
		connectionType entity.ConnectionTypeType,
	) (err []errors.Error)

	// Following returns all connections from a certain user
	Following(
		accountID, applicationID int64,
		userID uint64,
	) (followingIDs []uint64, err []errors.Error)

	// FollowedBy returns all connections from a certain user
	FollowedBy(
		accountID, applicationID int64,
		userID uint64,
	) (followedByIDs []uint64, err []errors.Error)

	// FollowedBy returns all friends from a certain user
	Friends(
		accountID, applicationID int64,
		userID uint64,
	) (friendsIDs []uint64, err []errors.Error)

	// FriendsAndFollowingIDs returns all friends and people followed by the specified user
	FriendsAndFollowingIDs(
		accountID, applicationID int64,
		userID uint64,
	) (friendsAndFollowingIDs []uint64, err []errors.Error)

	// SocialConnect creates the connections between a user and his other social peers
	SocialConnect(
		accountID, applicationID int64,
		user *entity.ApplicationUser,
		platform string,
		socialFriendsIDs []string,
		connectionType entity.ConnectionTypeType,
		connectionState entity.ConnectionStateType,
	) (connectedUserIDs []uint64, err []errors.Error)

	// CreateMultiple will connect a user with its friends from another social network
	CreateMultiple(
		accountID, applicationID int64,
		user *entity.ApplicationUser,
		connectionType entity.ConnectionTypeType,
		connectionState entity.ConnectionStateType,
		ourStoredUsersIDs []uint64,
	) (connectedUsersIDs []uint64, err []errors.Error)

	// Relation returns the relation between two users
	Relation(
		accountID, applicationID int64,
		userFromID, userToID uint64,
	) (relation *entity.Relation, err []errors.Error)

	// Exists checks if a connection exists between two users
	Exists(
		accountID, applicationID int64,
		userFromID, userToID uint64,
		connectionType entity.ConnectionTypeType,
	) (exists bool, err []errors.Error)

	// ConnectionsByState returns the connections that correspond to a specific state
	ConnectionsByState(
		accountID, applicationID int64,
		userID uint64,
		state entity.ConnectionStateType,
	) (connections []*entity.Connection, err []errors.Error)
}
