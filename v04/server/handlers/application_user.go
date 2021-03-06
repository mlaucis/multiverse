package handlers

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
)

// ApplicationUser defines the application user routes
type ApplicationUser interface {
	// Read handles requests to retrieve a single user
	Read(*context.Context) []errors.Error

	// ReadCurrent handles erquests to retrieve the current usr
	ReadCurrent(*context.Context) []errors.Error

	// UpdateCurrent handles requests to update the current user
	UpdateCurrent(*context.Context) []errors.Error

	// Delete handles requests to delete a user
	Delete(*context.Context) []errors.Error

	// DeleteCurrent handles requests to delete the current user
	DeleteCurrent(*context.Context) []errors.Error

	// Create handles requests to create a user
	Create(*context.Context) []errors.Error

	// Login handles the requests to login the user in the system
	Login(*context.Context) []errors.Error

	// RefreshSession handles the requests to refresh the user session token
	RefreshSession(*context.Context) []errors.Error

	// Logout handles the requests to logout the user from the system
	Logout(*context.Context) []errors.Error

	// Search handles the requests to search for an application user
	Search(*context.Context) []errors.Error

	// PopulateContext adds the applicationUser to the context
	PopulateContext(*context.Context) []errors.Error
}
