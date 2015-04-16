/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/tgerrors"
)

type (
	// Event route handler
	Event interface {
		// Read handles requests to retrieve a single event
		Read(*context.Context) tgerrors.TGError

		// Update handles requests to update an event
		Update(*context.Context) tgerrors.TGError

		// Delete handles requests to delete a single event
		Delete(*context.Context) tgerrors.TGError

		// List handles requests to retrieve a users events
		List(*context.Context) tgerrors.TGError

		// ConnectionEventsList handles requests to retrieve a users connections events
		ConnectionEventsList(*context.Context) tgerrors.TGError

		// Create handles requests to create an event
		Create(*context.Context) tgerrors.TGError

		// SearchGeo handles requests to retrieve a users connections events
		SearchGeo(*context.Context) tgerrors.TGError

		// SearchObject handles requests to retrieve events in a certain location / radius
		SearchObject(*context.Context) tgerrors.TGError

		// SearchLocation handles requests to retrieve a users connections events
		SearchLocation(*context.Context) tgerrors.TGError
	}
)
