/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package handlers

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
)

type (
	// Connection holds the routes for the connections
	Connection interface {
		// Update handles requests to update a user connection
		Update(*context.Context) []errors.Error

		// Delete handles requests to delete a single connection
		Delete(*context.Context) []errors.Error

		// Create handles requests to create a user connection
		Create(*context.Context) []errors.Error

		// List handles requests to get a users connections
		List(*context.Context) []errors.Error

		// List handles requests to get the current users connections
		CurrentUserList(*context.Context) []errors.Error

		// FollowedByList handles requests to get a users list of users who follow him
		FollowedByList(*context.Context) []errors.Error

		// CurrentUserFollowedByList handles requests to get the current user list of users who follow him
		CurrentUserFollowedByList(*context.Context) []errors.Error

		// Friends handles requests to get a user list of friends
		Friends(*context.Context) []errors.Error

		// CurrentUserFriends handles requests to get the current user list of friends
		CurrentUserFriends(*context.Context) []errors.Error

		// Confirm handles requests to confirm a user connection
		Confirm(*context.Context) []errors.Error

		// CreateSocialConnections creates the social connections between users of the same social network
		CreateSocial(*context.Context) []errors.Error
	}
)
