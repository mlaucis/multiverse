/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Event interface
	Event interface {
		// Create adds an event to the database and returns the created event or an error
		Create(event *entity.Event, retrieve bool) (evn *entity.Event, err errors.Error)

		// Read returns the event matching the ID or an error
		Read(accountID, applicationID, userID, eventID int64) (event *entity.Event, err errors.Error)

		// Update updates an event in the database and returns the updated event or an error
		Update(existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err errors.Error)

		// Delete deletes the event matching the IDs or an error
		Delete(*entity.Event) errors.Error

		// List returns all events from a certain user
		List(accountID, applicationID, userID int64) (events []*entity.Event, err errors.Error)

		// ConnectionList returns all events from connections
		ConnectionList(accountID, applicationID, userID int64) (events []*entity.Event, err errors.Error)

		// WriteEventToConnectionsLists takes an event and writes it to the user connections list
		WriteToConnectionsLists(event *entity.Event, key string) errors.Error

		// DeleteEventFromConnectionsLists takes a user id and key and deletes it to the user connections list
		DeleteFromConnectionsLists(accountID, applicationID, userID int64, key string) (err errors.Error)

		// GeoSearch retrieves all the events from an application within a radius of the provided coordinates
		GeoSearch(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err errors.Error)

		// ObjectSearch returns all the events for a specific object
		ObjectSearch(accountID, applicationID int64, objectKey string) ([]*entity.Event, errors.Error)

		// LocationSearch returns all the events for a specific object
		LocationSearch(accountID, applicationID int64, locationKey string) ([]*entity.Event, errors.Error)
	}
)
