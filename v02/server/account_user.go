/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

type (
	// AccountUser holds the account user routes
	AccountUser interface {
		// Read handles requests to a single account user
		Read(*context.Context) tgerrors.TGError

		// Update handles requests update an account user
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single account user
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests create an account user
		Create(*context.Context) tgerrors.TGError

		// List handles requests to list all account users
		List(*context.Context) tgerrors.TGError

		// Login handles the requests to login the user in the system
		Login(*context.Context) tgerrors.TGError

		// RefreshSession handles the requests to refresh the account user session token
		RefreshSession(*context.Context) tgerrors.TGError

		// Logout handles the requests to logout the account user from the system
		Logout(*context.Context) tgerrors.TGError

		// PopulateContext adds the accountUser to the context
		PopulateContext(*context.Context) tgerrors.TGError
	}
)
