/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
)

type (
	// Event route handler
	Event interface {
		// Read handles requests to retrieve a single event
		Read(*context.Context) errors.Error

		// Update handles requests to update an event
		Update(*context.Context) errors.Error

		// Delete handles requests to delete a single event
		Delete(*context.Context) errors.Error

		// List handles requests to retrieve a users events
		List(*context.Context) errors.Error

		// ConnectionEventsList handles requests to retrieve a users connections events
		ConnectionEventsList(*context.Context) errors.Error

		// Create handles requests to create an event
		Create(*context.Context) errors.Error

		// SearchGeo handles requests to retrieve a users connections events
		SearchGeo(*context.Context) errors.Error

		// SearchObject handles requests to retrieve events in a certain location / radius
		SearchObject(*context.Context) errors.Error

		// SearchLocation handles requests to retrieve a users connections events
		SearchLocation(*context.Context) errors.Error
	}
)
