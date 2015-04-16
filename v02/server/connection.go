/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

type (
	// Connection holds the routes for the connections
	Connection interface {
		// Update handles requests to update a user connection
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single connection
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests to create a user connection
		Create(*context.Context) tgerrors.TGError

		// List handles requests to list a users connections
		List(*context.Context) tgerrors.TGError

		// FollowedByList handles requests to list a users list of users who follow him
		FollowedByList(*context.Context) tgerrors.TGError

		// Confirm handles requests to confirm a user connection
		Confirm(*context.Context) tgerrors.TGError

		// CreateSocialConnections creates the social connections between users of the same social network
		CreateSocial(*context.Context) tgerrors.TGError
	}
)

// AcceptedPlatforms defines which social platforms we accept right now
var AcceptedPlatforms = map[string]bool{
	"facebook": true,
	"twitter":  true,
	"gplus":    true,
	"abook":    true,
}
