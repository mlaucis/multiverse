package core

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Event interface
	Event interface {
		// Create adds an event to the database and returns the created event or an error
		Create(accountID, applicationID int64, currentUserID uint64, event *entity.Event, retrieve bool) (evn *entity.Event, err []errors.Error)

		// Read returns the event matching the ID or an error
		Read(accountID, applicationID int64, userID, currentUserID, eventID uint64) (event *entity.Event, err []errors.Error)

		// Update updates an event in the database and returns the updated event or an error
		Update(accountID, applicationID int64, currentUserID uint64, existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err []errors.Error)

		// Delete deletes the event matching the IDs or an error
		Delete(accountID, applicationID int64, currentUserID uint64, evt *entity.Event) []errors.Error

		// List returns all events from a certain user
		List(accountID, applicationID int64, userID, currentUserID uint64) (events []*entity.Event, err []errors.Error)

		// ConnectionList returns all events from connections
		UserFeed(accountID, applicationID int64, user *entity.ApplicationUser) (count int, events []*entity.Event, err []errors.Error)

		// UnreadFeed returns only the events that would form a feed but have not been retrieved yet
		UnreadFeed(accountID, applicationID int64, user *entity.ApplicationUser) (count int, events []*entity.Event, err []errors.Error)

		// UnreadFeedCount returns the number of events since the last time either UserFeed() or UnreadFeed() was executed
		UnreadFeedCount(accountID, applicationID int64, user *entity.ApplicationUser) (count int, err []errors.Error)

		// WriteEventToConnectionsLists takes an event and writes it to the user connections list
		WriteToConnectionsLists(accountID, applicationID int64, event *entity.Event, key string) []errors.Error

		// DeleteEventFromConnectionsLists takes a user id and key and deletes it to the user connections list
		DeleteFromConnectionsLists(accountID, applicationID int64, userID uint64, key string) (err []errors.Error)

		// GeoSearch retrieves all the events from an application within a radius of the provided coordinates
		GeoSearch(accountID, applicationID int64, currentUserID uint64, latitude, longitude, radius float64, nearest int64) (events []*entity.Event, err []errors.Error)

		// ObjectSearch returns all the events for a specific object
		ObjectSearch(accountID, applicationID int64, currentUserID uint64, objectKey string) ([]*entity.Event, []errors.Error)

		// LocationSearch returns all the events for a specific object
		LocationSearch(accountID, applicationID int64, currentUserID uint64, locationKey string) ([]*entity.Event, []errors.Error)
	}
)
