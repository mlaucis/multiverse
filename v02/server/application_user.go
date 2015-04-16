/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

type (
	// ApplicationUser defines the application user routes
	ApplicationUser interface {
		// Read handles requests to retrieve a single user
		Read(*context.Context) tgerrors.TGError

		// Update handles requests to update a user
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single user
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests to create a user
		Create(*context.Context) tgerrors.TGError

		// Login handles the requests to login the user in the system
		Login(*context.Context) tgerrors.TGError

		// RefreshSession handles the requests to refresh the user session token
		RefreshSession(*context.Context) tgerrors.TGError

		// Logout handles the requests to logout the user from the system
		Logout(*context.Context) tgerrors.TGError

		// PopulateContext adds the applicationUser to the context
		PopulateContext(*context.Context) tgerrors.TGError
	}
)
