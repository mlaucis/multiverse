/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

type (
	// Account holds the account routes
	Account interface {
		// Read handles requests to a single account
		Read(*context.Context) tgerrors.TGError

		// Update handles requests to update a single account
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single account
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests create an account
		Create(*context.Context) tgerrors.TGError

		// PopulateContext adds the account to the context
		PopulateContext(*context.Context) tgerrors.TGError
	}
)
