package handlers

import (
	"github.com/tapglue/backend/context"
	"github.com/tapglue/backend/errors"
)

// Member holds the account user routes
type Member interface {
	// Read handles requests to a single account user
	Read(*context.Context) []errors.Error

	// Update handles requests update an account user
	Update(*context.Context) []errors.Error

	// Delete handles requests to delete a single account user
	Delete(*context.Context) []errors.Error

	// Create handles requests create an account user
	Create(*context.Context) []errors.Error

	// List handles requests to list all account users
	List(*context.Context) []errors.Error

	// Login handles the requests to login the user in the system
	Login(*context.Context) []errors.Error

	// RefreshSession handles the requests to refresh the account user session token
	RefreshSession(*context.Context) []errors.Error

	// Logout handles the requests to logout the account user from the system
	Logout(*context.Context) []errors.Error

	// PopulateContext adds the accountUser to the context
	PopulateContext(*context.Context) []errors.Error
}
