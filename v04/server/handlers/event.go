package handlers

import (
	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v04/context"
)

// Event route handler
type Event interface {
	// CurrentUserRead handles requests to retrieve a single event for the current user
	CurrentUserRead(*context.Context) []errors.Error

	// Read handles requests to retrieve a single event
	Read(*context.Context) []errors.Error

	// Update handles requests to update an event
	Update(*context.Context) []errors.Error

	// CurrentUserUpdate handles requests to update an event for the current user
	CurrentUserUpdate(*context.Context) []errors.Error

	// Delete handles requests to delete a single event
	Delete(*context.Context) []errors.Error

	// CurrentUserDelete handles requests to delete a single event from a user
	CurrentUserDelete(*context.Context) []errors.Error

	// List handles requests to retrieve a users events
	List(*context.Context) []errors.Error

	// CurrentUserList handles requests to retrieve the current user events
	CurrentUserList(*context.Context) []errors.Error

	// Feed handles requests to retrieve a users connections events
	Feed(*context.Context) []errors.Error

	// CurrentUserCreate handles requests to create an event for the current user
	CurrentUserCreate(*context.Context) []errors.Error

	// Create handles requests to create an event for a user
	Create(*context.Context) []errors.Error

	// Search handles requests to retrieve events that match a certain query
	Search(*context.Context) []errors.Error

	// UnreadFeed will return only the events in the feed that not read yet by the user and flag them as read
	UnreadFeed(*context.Context) []errors.Error

	// UnreadFeedCount will return the count of the events in the feed that not read yet by the user
	UnreadFeedCount(*context.Context) []errors.Error
}
