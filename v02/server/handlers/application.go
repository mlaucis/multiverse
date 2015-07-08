/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package handlers

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
)

type (
	// Application defines the routes for the application
	Application interface {
		// Read handles requests to a single application
		Read(*context.Context) []errors.Error

		// Update handles requests updates an application
		Update(*context.Context) []errors.Error

		// Delete handles requests to delete a single application
		Delete(*context.Context) []errors.Error

		// Create handles requests create an application
		Create(*context.Context) []errors.Error

		// List handles requests list all account applications
		List(*context.Context) []errors.Error

		// PopulateContext adds the application to the context
		PopulateContext(ctx *context.Context) []errors.Error

		// PopulateContextFromID adds the application to the context
		PopulateContextFromID(ctx *context.Context) []errors.Error
	}
)
