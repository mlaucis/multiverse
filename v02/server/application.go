/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

type (
	// Application defines the routes for the application
	Application interface {
		// Read handles requests to a single application
		Read(*context.Context) tgerrors.TGError

		// Update handles requests updates an application
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single application
		Delete(*context.Context) tgerrors.TGError

		// Create handles requests create an application
		Create(*context.Context) tgerrors.TGError

		// List handles requests list all account applications
		List(*context.Context) tgerrors.TGError

		// PopulateContext adds the application to the context
		PopulateContext(ctx *context.Context) tgerrors.TGError
	}
)
