/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package server

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
)

type (
	// Account holds the account routes
	Account interface {
		// Read handles requests to a single account
		Read(*context.Context) []errors.Error

		// Update handles requests to update a single account
		Update(*context.Context) []errors.Error

		// Delete handles requests to delete a single account
		Delete(*context.Context) []errors.Error

		// Create handles requests create an account
		Create(*context.Context) []errors.Error

		// PopulateContext adds the account to the context
		PopulateContext(*context.Context) []errors.Error
	}
)
