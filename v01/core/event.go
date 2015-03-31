/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"time"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v01/entity"

	"github.com/tapglue/georedis"
	red "gopkg.in/redis.v2"
)

// ReadEvent returns the event matching the ID or an error
func ReadEvent(accountID, applicationID, userID, eventID int64) (event *entity.Event, err *tgerrors.TGError) {
	key := storageClient.Event(accountID, applicationID, userID, eventID)

	result, er := storageEngine.Get(key).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the event (1)", er.Error())
	}

	if er = json.Unmarshal([]byte(result), &event); er != nil {
		return nil, tgerrors.NewInternalError("failed to read the event (2)", er.Error())
	}

	return
}

// UpdateEvent updates an event in the database and returns the updated event or an error
func UpdateEvent(existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err *tgerrors.TGError) {
	updatedEvent.UpdatedAt = time.Now()

	val, er := json.Marshal(updatedEvent)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to update the event (1)", er.Error())
	}

	key := storageClient.Event(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.UserID, updatedEvent.ID)
	if er = storageEngine.Set(key, string(val)).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to update the event (1)", er.Error())
	}

	if existingEvent.Latitude != 0 &&
		existingEvent.Longitude != 0 {
		removeEventGeo(key, &updatedEvent)
	}

	if updatedEvent.Enabled && (existingEvent.Latitude != updatedEvent.Latitude ||
		existingEvent.Longitude != updatedEvent.Longitude) {
		addEventGeo(key, &updatedEvent)
	}

	if updatedEvent.Object != nil && updatedEvent.Enabled {
		if existingEvent.Object != nil {
			if er := removeEventObject(key, &existingEvent); er != nil {
				return nil, er
			}
		}

		if er := addEventObject(key, &updatedEvent); er != nil {
			return nil, er
		}
	}

	if existingEvent.Location != updatedEvent.Location {
		if existingEvent.Location != "" {
			if er := removeEventLocation(key, &updatedEvent); er != nil {
				return nil, er
			}
		}

		if updatedEvent.Location != "" && updatedEvent.Enabled {
			if er := addEventLocation(key, &updatedEvent); er != nil {
				return nil, er
			}
		}
	}

	if !updatedEvent.Enabled {
		listKey := storageClient.Events(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.UserID)
		if er = storageEngine.ZRem(listKey, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to read the event (1)", er.Error())
		}
		if err = DeleteEventFromConnectionsLists(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.UserID, key); err != nil {
			return nil, tgerrors.NewInternalError("failed to read the event (1)", err.Error())
		}

		if existingEvent.Latitude != 0 && existingEvent.Longitude != 0 {
			removeEventGeo(key, &updatedEvent)
		}

		if existingEvent.Object != nil {
			removeEventObject(key, &existingEvent)
		}

		if existingEvent.Location != "" {
			removeEventLocation(key, &updatedEvent)
		}
	}

	if !retrieve {
		return &updatedEvent, nil
	}

	return ReadEvent(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.UserID, updatedEvent.ID)
}

// DeleteEvent deletes the event matching the IDs or an error
func DeleteEvent(accountID, applicationID, userID, eventID int64) (err *tgerrors.TGError) {
	key := storageClient.Event(accountID, applicationID, userID, eventID)

	event, err := ReadEvent(accountID, applicationID, userID, eventID)
	if err != nil {
		return err
	}

	result, er := storageEngine.Del(key).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to delete the event (1)", er.Error())
	}

	if result != 1 {
		return tgerrors.NewInternalError("failed to delete the event (2)", "event already deleted")
	}

	listKey := storageClient.Events(accountID, applicationID, userID)
	if er = storageEngine.ZRem(listKey, key).Err(); er != nil {
		return tgerrors.NewInternalError("failed to read the event (1)", er.Error())
	}

	if err = DeleteEventFromConnectionsLists(accountID, applicationID, userID, key); err != nil {
		return
	}

	removeEventGeo(key, event)
	removeEventObject(key, event)
	removeEventLocation(key, event)

	return nil
}

// ReadEventList returns all events from a certain user
func ReadEventList(accountID, applicationID, userID int64) (events []*entity.Event, err *tgerrors.TGError) {
	key := storageClient.Events(accountID, applicationID, userID)

	result, er := storageEngine.ZRevRange(key, "0", "-1").Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the event list (1)", er.Error())
	}

	if len(result) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, er := storageEngine.MGet(result...).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the event list (2)", er.Error())
	}

	return toEvents(resultList)
}

// ReadConnectionEventList returns all events from connections
func ReadConnectionEventList(accountID, applicationID, userID int64) (events []*entity.Event, err *tgerrors.TGError) {
	key := storageClient.ConnectionEvents(accountID, applicationID, userID)

	result, er := storageEngine.ZRevRange(key, "0", "-1").Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the connections events (1)", er.Error())
	}

	// TODO maybe this shouldn't be an error but rather return that there are no events from connections
	if len(result) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, er := storageEngine.MGet(result...).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the connections events (2)", er.Error())
	}

	return toEvents(resultList)
}

// WriteEvent adds an event to the database and returns the created event or an error
func WriteEvent(event *entity.Event, retrieve bool) (evn *entity.Event, err *tgerrors.TGError) {
	event.Enabled = true
	event.CreatedAt = time.Now()
	event.UpdatedAt = event.CreatedAt
	var er error

	if event.ID, er = storageClient.GenerateApplicationEventID(event.ApplicationID); er != nil {
		return nil, tgerrors.NewInternalError("failed to write the event (1)", er.Error())
	}

	val, er := json.Marshal(event)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the event (2)", er.Error())
	}

	key := storageClient.Event(event.AccountID, event.ApplicationID, event.UserID, event.ID)
	if er = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, tgerrors.NewInternalError("failed to write the event (3)", er.Error())
	}

	listKey := storageClient.Events(event.AccountID, event.ApplicationID, event.UserID)

	setVal := red.Z{Score: float64(event.CreatedAt.Unix()), Member: key}
	if er = storageEngine.ZAdd(listKey, setVal).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to write the event (4)", er.Error())
	}

	if event.Latitude != 0 && event.Longitude != 0 {
		coordinates := georedis.GeoKey{
			Lat:   event.Latitude,
			Lon:   event.Longitude,
			Label: key,
		}

		geoEventKey := storageClient.EventGeoKey(event.AccountID, event.ApplicationID)
		georedis.AddCoordinates(storageEngine, geoEventKey, 52, coordinates)
	}

	if event.Object != nil {
		objectEventKey := storageClient.EventObjectKey(event.AccountID, event.ApplicationID, event.Object.ID)
		if er = storageEngine.SAdd(objectEventKey, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to write the event (5)", er.Error())
		}
	}

	if event.Location != "" {
		locationEventKey := storageClient.EventLocationKey(event.AccountID, event.ApplicationID, event.Location)
		if er = storageEngine.SAdd(locationEventKey, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to write the event (6)", er.Error())
		}
	}

	if err = WriteEventToConnectionsLists(event, key); err != nil {
		return nil, err
	}

	if !retrieve {
		return event, nil
	}

	return ReadEvent(event.AccountID, event.ApplicationID, event.UserID, event.ID)
}

// WriteEventToConnectionsLists takes an event and writes it to the user connections list
func WriteEventToConnectionsLists(event *entity.Event, key string) (err *tgerrors.TGError) {
	connectionsKey := storageClient.FollowedByUsers(event.AccountID, event.ApplicationID, event.UserID)

	connections, er := storageEngine.LRange(connectionsKey, 0, -1).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to write the event to the lists (1)", er.Error())
	}

	for _, userKey := range connections {
		feedKey := storageClient.ConnectionEventsLoop(userKey)

		val := red.Z{Score: float64(event.CreatedAt.Unix()), Member: key}
		if er = storageEngine.ZAdd(feedKey, val).Err(); er != nil {
			return tgerrors.NewInternalError("failed to write the event to the list (2)", er.Error())
		}
	}

	return nil
}

// DeleteEventFromConnectionsLists takes a user id and key and deletes it to the user connections list
func DeleteEventFromConnectionsLists(accountID, applicationID, userID int64, key string) (err *tgerrors.TGError) {
	connectionsKey := storageClient.FollowedByUsers(accountID, applicationID, userID)
	connections, er := storageEngine.LRange(connectionsKey, 0, -1).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to delete the event from list (1)", er.Error())
	}

	for _, userKey := range connections {
		feedKey := storageClient.ConnectionEventsLoop(userKey)
		if er = storageEngine.ZRem(feedKey, key).Err(); er != nil {
			return tgerrors.NewInternalError("failed to delete the event from list (2)", er.Error())
		}
	}

	return nil
}

// SearchGeoEvents retrieves all the events from an application within a radius of the provided coordinates
func SearchGeoEvents(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err *tgerrors.TGError) {
	geoEventKey := storageClient.EventGeoKey(accountID, applicationID)

	eventKeys, er := georedis.SearchByRadius(storageEngine, geoEventKey, latitude, longitude, radius, 52)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to search for events by geo (1)", er.Error())
	}

	resultList, er := storageEngine.MGet(eventKeys...).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to search for events by geo (2)", er.Error())
	}

	return toEvents(resultList)
}

// SearchObjectEvents returns all the events for a specific object
func SearchObjectEvents(accountID, applicationID int64, objectKey string) ([]*entity.Event, *tgerrors.TGError) {
	objectEventKey := storageClient.EventObjectKey(accountID, applicationID, objectKey)

	return fetchEventsFromKeys(accountID, applicationID, objectEventKey)
}

// SearchLocationEvents returns all the events for a specific object
func SearchLocationEvents(accountID, applicationID int64, locationKey string) ([]*entity.Event, *tgerrors.TGError) {
	locationEventKey := storageClient.EventLocationKey(accountID, applicationID, locationKey)

	return fetchEventsFromKeys(accountID, applicationID, locationEventKey)
}

// fetchEventsFromKeys returns all the events matching a certain search key from the specified bucket
func fetchEventsFromKeys(accountID, applicationID int64, bucketName string) ([]*entity.Event, *tgerrors.TGError) {
	_, keys, er := storageEngine.SScan(bucketName, 0, "", 300).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the events (1)", er.Error())
	}

	if len(keys) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, err := storageEngine.MGet(keys...).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read the events (2)", er.Error())
	}

	return toEvents(resultList)
}

// toEvents converts the events from json format to go structs
func toEvents(resultList []interface{}) ([]*entity.Event, *tgerrors.TGError) {
	events := []*entity.Event{}
	for _, result := range resultList {
		if result == nil {
			continue
		}
		event := &entity.Event{}
		if er := json.Unmarshal([]byte(result.(string)), event); er != nil {
			return []*entity.Event{}, tgerrors.NewInternalError("failed to read the event from list (1)", er.Error())
		}
		events = append(events, event)
		event = &entity.Event{}
	}

	return events, nil
}

func addEventGeo(key string, updatedEvent *entity.Event) *tgerrors.TGError {
	coordinates := georedis.GeoKey{
		Lat:   updatedEvent.Latitude,
		Lon:   updatedEvent.Longitude,
		Label: key,
	}

	geoEventKey := storageClient.EventGeoKey(updatedEvent.AccountID, updatedEvent.ApplicationID)
	_, er := georedis.AddCoordinates(storageEngine, geoEventKey, 52, coordinates)
	if er == nil {
		return nil
	}
	return tgerrors.NewInternalError("failed to add the event by geo (1)", er.Error())
}

func removeEventGeo(key string, updatedEvent *entity.Event) *tgerrors.TGError {
	geoEventKey := storageClient.EventGeoKey(updatedEvent.AccountID, updatedEvent.ApplicationID)
	_, er := georedis.RemoveCoordinatesByKeys(storageEngine, geoEventKey, key)
	if er == nil {
		return nil
	}
	return tgerrors.NewInternalError("failed to remove the event by geo (1)", er.Error())
}

func addEventObject(key string, updatedEvent *entity.Event) *tgerrors.TGError {
	objectEventKey := storageClient.EventObjectKey(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.Object.ID)
	er := storageEngine.SAdd(objectEventKey, key).Err()
	if er == nil {
		return nil
	}
	return tgerrors.NewInternalError("failed to add the event by object (1)", er.Error())
}

func removeEventObject(key string, updatedEvent *entity.Event) *tgerrors.TGError {
	objectEventKey := storageClient.EventObjectKey(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.Object.ID)
	er := storageEngine.SRem(objectEventKey, key).Err()
	if er == nil {
		return nil
	}
	return tgerrors.NewInternalError("failed to remove the event by geo (1)", er.Error())
}

func addEventLocation(key string, updatedEvent *entity.Event) *tgerrors.TGError {
	locationEventKey := storageClient.EventLocationKey(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.Location)
	er := storageEngine.SAdd(locationEventKey, key).Err()
	if er == nil {
		return nil
	}
	return tgerrors.NewInternalError("failed to add the event by location (1)", er.Error())
}

func removeEventLocation(key string, updatedEvent *entity.Event) *tgerrors.TGError {
	locationEventKey := storageClient.EventLocationKey(updatedEvent.AccountID, updatedEvent.ApplicationID, updatedEvent.Location)
	er := storageEngine.SRem(locationEventKey, key).Err()
	if er == nil {
		return nil
	}
	return tgerrors.NewInternalError("failed to remove the event by location (1)", er.Error())
}
