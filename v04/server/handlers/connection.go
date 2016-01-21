package handlers

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
)

// Connection holds the routes for the connections
type Connection interface {
	// Update handles requests to update a user connection
	Update(*context.Context) []errors.Error

	// Delete handles requests to delete a single connection
	Delete(*context.Context) []errors.Error

	// Create handles requests to create a user connection
	Create(*context.Context) []errors.Error

	// FollowingList handles requests to get a users connections
	FollowingList(*context.Context) []errors.Error

	// CurrentUserFollowingList handles requests to get the current users connections
	CurrentUserFollowingList(*context.Context) []errors.Error

	// FollowedByList handles requests to get a users list of users who follow him
	FollowedByList(*context.Context) []errors.Error

	// CurrentUserFollowedByList handles requests to get the current user list of users who follow him
	CurrentUserFollowedByList(*context.Context) []errors.Error

	// Friends handles requests to get a user list of friends
	Friends(*context.Context) []errors.Error

	// CurrentUserFriends handles requests to get the current user list of friends
	CurrentUserFriends(*context.Context) []errors.Error

	// CreateSocialConnections creates the social connections between users of the same social network
	CreateSocial(*context.Context) []errors.Error

	// CreateFriend is a an alias for creating a friend connection type
	CreateFriend(*context.Context) []errors.Error

	// CreateFollow is an alias for creating a follow connection type
	CreateFollow(*context.Context) []errors.Error

	// UserConnectionsByState retrieves the user connections by state
	UserConnectionsByState(*context.Context) []errors.Error

	// CurrentUserConnectionsByState retrieves the user connections by state
	CurrentUserConnectionsByState(*context.Context) []errors.Error
}
