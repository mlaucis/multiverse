package handlers

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
)

// Organization holds the account routes
type Organization interface {
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
