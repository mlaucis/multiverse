/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/entity"
)

type (
	// Event interface
	Event interface {
		Create(event *entity.Event, retrieve bool) (evn *entity.Event, err tgerrors.TGError)
		Read(accountID, applicationID, userID, eventID int64) (event *entity.Event, err tgerrors.TGError)
		Update(existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err tgerrors.TGError)
		Delete(accountID, applicationID, userID, eventID int64) (err tgerrors.TGError)
		List(accountID, applicationID, userID int64) (events []*entity.Event, err tgerrors.TGError)
		ConnectionList(accountID, applicationID, userID int64) (events []*entity.Event, err tgerrors.TGError)
		WriteToConnectionsLists(event *entity.Event, key string) (err tgerrors.TGError)
		DeleteFromConnectionsLists(accountID, applicationID, userID int64, key string) (err tgerrors.TGError)
		GeoSearch(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err tgerrors.TGError)
		ObjectSearch(accountID, applicationID int64, objectKey string) ([]*entity.Event, tgerrors.TGError)
		LocationSearch(accountID, applicationID int64, locationKey string) ([]*entity.Event, tgerrors.TGError)
	}
)
