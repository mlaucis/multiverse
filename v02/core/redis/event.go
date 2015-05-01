package redis

import (
	"encoding/json"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/redis"

	"github.com/tapglue/georedis"
	red "gopkg.in/redis.v2"
)

type (
	event struct {
		storage redis.Client
		redis   *red.Client
	}
)

func (e *event) Create(accountID, applicationID int64, event *entity.Event, retrieve bool) (evn *entity.Event, err errors.Error) {
	event.Enabled = true
	event.CreatedAt = time.Now()
	event.UpdatedAt = event.CreatedAt
	var er error

	if event.ID, er = e.storage.GenerateApplicationEventID(applicationID); er != nil {
		return nil, errors.NewInternalError("failed to write the event (1)", er.Error())
	}

	val, er := json.Marshal(event)
	if er != nil {
		return nil, errors.NewInternalError("failed to write the event (2)", er.Error())
	}

	key := storageHelper.Event(accountID, applicationID, event.UserID, event.ID)
	if er = e.redis.Set(key, string(val)).Err(); err != nil {
		return nil, errors.NewInternalError("failed to write the event (3)", er.Error())
	}

	listKey := storageHelper.Events(accountID, applicationID, event.UserID)

	setVal := red.Z{Score: float64(event.CreatedAt.Unix()), Member: key}
	if er = e.redis.ZAdd(listKey, setVal).Err(); er != nil {
		return nil, errors.NewInternalError("failed to write the event (4)", er.Error())
	}

	if event.Latitude != 0 && event.Longitude != 0 {
		coordinates := georedis.GeoKey{
			Lat:   event.Latitude,
			Lon:   event.Longitude,
			Label: key,
		}

		geoEventKey := storageHelper.EventGeoKey(accountID, applicationID)
		georedis.AddCoordinates(e.redis, geoEventKey, 52, coordinates)
	}

	if event.Object != nil {
		objectEventKey := storageHelper.EventObjectKey(accountID, applicationID, event.Object.ID)
		if er = e.redis.SAdd(objectEventKey, key).Err(); er != nil {
			return nil, errors.NewInternalError("failed to write the event (5)", er.Error())
		}
	}

	if event.Location != "" {
		locationEventKey := storageHelper.EventLocationKey(accountID, applicationID, event.Location)
		if er = e.redis.SAdd(locationEventKey, key).Err(); er != nil {
			return nil, errors.NewInternalError("failed to write the event (6)", er.Error())
		}
	}

	if err = e.WriteToConnectionsLists(accountID, applicationID, event, key); err != nil {
		return nil, err
	}

	if !retrieve {
		return event, nil
	}

	return e.Read(accountID, applicationID, event.UserID, event.ID)
}

func (e *event) Read(accountID, applicationID int64, userID, eventID string) (event *entity.Event, err errors.Error) {
	key := storageHelper.Event(accountID, applicationID, userID, eventID)

	result, er := e.redis.Get(key).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the event (1)", er.Error())
	}

	if er = json.Unmarshal([]byte(result), &event); er != nil {
		return nil, errors.NewInternalError("failed to read the event (2)", er.Error())
	}

	return
}

func (e *event) Update(accountID, applicationID int64, existingEvent, updatedEvent entity.Event, retrieve bool) (evn *entity.Event, err errors.Error) {
	updatedEvent.UpdatedAt = time.Now()

	val, er := json.Marshal(updatedEvent)
	if er != nil {
		return nil, errors.NewInternalError("failed to update the event (1)", er.Error())
	}

	key := storageHelper.Event(accountID, applicationID, updatedEvent.UserID, updatedEvent.ID)
	if er = e.redis.Set(key, string(val)).Err(); er != nil {
		return nil, errors.NewInternalError("failed to update the event (1)", er.Error())
	}

	if existingEvent.Latitude != 0 &&
		existingEvent.Longitude != 0 {
		e.removeEventGeo(accountID, applicationID, key, &updatedEvent)
	}

	if updatedEvent.Enabled && (existingEvent.Latitude != updatedEvent.Latitude ||
		existingEvent.Longitude != updatedEvent.Longitude) {
		e.addEventGeo(accountID, applicationID, key, &updatedEvent)
	}

	if updatedEvent.Object != nil && updatedEvent.Enabled {
		if existingEvent.Object != nil {
			if er := e.removeEventObject(accountID, applicationID, key, &existingEvent); er != nil {
				return nil, er
			}
		}

		if er := e.addEventObject(accountID, applicationID, key, &updatedEvent); er != nil {
			return nil, er
		}
	}

	if existingEvent.Location != updatedEvent.Location {
		if existingEvent.Location != "" {
			if er := e.removeEventLocation(accountID, applicationID, key, &updatedEvent); er != nil {
				return nil, er
			}
		}

		if updatedEvent.Location != "" && updatedEvent.Enabled {
			if er := e.addEventLocation(accountID, applicationID, key, &updatedEvent); er != nil {
				return nil, er
			}
		}
	}

	if !updatedEvent.Enabled {
		listKey := storageHelper.Events(accountID, applicationID, updatedEvent.UserID)
		if er = e.redis.ZRem(listKey, key).Err(); er != nil {
			return nil, errors.NewInternalError("failed to read the event (1)", er.Error())
		}
		if err = e.DeleteFromConnectionsLists(accountID, applicationID, updatedEvent.UserID, key); err != nil {
			return nil, errors.NewInternalError("failed to read the event (1)", err.Error())
		}

		if existingEvent.Latitude != 0 && existingEvent.Longitude != 0 {
			e.removeEventGeo(accountID, applicationID, key, &updatedEvent)
		}

		if existingEvent.Object != nil {
			e.removeEventObject(accountID, applicationID, key, &existingEvent)
		}

		if existingEvent.Location != "" {
			e.removeEventLocation(accountID, applicationID, key, &updatedEvent)
		}
	}

	if !retrieve {
		return &updatedEvent, nil
	}

	return e.Read(accountID, applicationID, updatedEvent.UserID, updatedEvent.ID)
}

func (e *event) Delete(accountID, applicationID int64, event *entity.Event) (err errors.Error) {
	key := storageHelper.Event(accountID, applicationID, event.UserID, event.ID)
	result, er := e.redis.Del(key).Result()
	if er != nil {
		return errors.NewInternalError("failed to delete the event (1)", er.Error())
	}

	if result != 1 {
		return errors.NewInternalError("failed to delete the event (2)", "event already deleted")
	}

	listKey := storageHelper.Events(accountID, applicationID, event.UserID)
	if er = e.redis.ZRem(listKey, key).Err(); er != nil {
		return errors.NewInternalError("failed to read the event (1)", er.Error())
	}

	if err = e.DeleteFromConnectionsLists(accountID, applicationID, event.UserID, key); err != nil {
		return
	}

	e.removeEventGeo(accountID, applicationID, key, event)
	e.removeEventObject(accountID, applicationID, key, event)
	e.removeEventLocation(accountID, applicationID, key, event)

	return nil
}

func (e *event) List(accountID, applicationID int64, userID string) (events []*entity.Event, err errors.Error) {
	key := storageHelper.Events(accountID, applicationID, userID)

	result, er := e.redis.ZRevRange(key, "0", "-1").Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the event list (1)", er.Error())
	}

	if len(result) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, er := e.redis.MGet(result...).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the event list (2)", er.Error())
	}

	return e.toEvents(resultList)
}

func (e *event) ConnectionList(accountID, applicationID int64, userID string) (events []*entity.Event, err errors.Error) {
	key := storageHelper.ConnectionEvents(accountID, applicationID, userID)

	result, er := e.redis.ZRevRange(key, "0", "-1").Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the connections events (1)", er.Error())
	}

	// TODO maybe this shouldn't be an error but rather return that there are no events from connections
	if len(result) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, er := e.redis.MGet(result...).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the connections events (2)", er.Error())
	}

	return e.toEvents(resultList)
}

func (e *event) WriteToConnectionsLists(accountID, applicationID int64, event *entity.Event, key string) (err errors.Error) {
	connectionsKey := storageHelper.FollowedByUsers(accountID, applicationID, event.UserID)

	connections, er := e.redis.LRange(connectionsKey, 0, -1).Result()
	if er != nil {
		return errors.NewInternalError("failed to write the event to the lists (1)", er.Error())
	}

	for _, userKey := range connections {
		feedKey := storageHelper.ConnectionEventsLoop(userKey)

		val := red.Z{Score: float64(event.CreatedAt.Unix()), Member: key}
		if er = e.redis.ZAdd(feedKey, val).Err(); er != nil {
			return errors.NewInternalError("failed to write the event to the list (2)", er.Error())
		}
	}

	return nil
}

func (e *event) DeleteFromConnectionsLists(accountID, applicationID int64, userID, key string) (err errors.Error) {
	connectionsKey := storageHelper.FollowedByUsers(accountID, applicationID, userID)
	connections, er := e.redis.LRange(connectionsKey, 0, -1).Result()
	if er != nil {
		return errors.NewInternalError("failed to delete the event from list (1)", er.Error())
	}

	for _, userKey := range connections {
		feedKey := storageHelper.ConnectionEventsLoop(userKey)
		if er = e.redis.ZRem(feedKey, key).Err(); er != nil {
			return errors.NewInternalError("failed to delete the event from list (2)", er.Error())
		}
	}

	return nil
}

func (e *event) GeoSearch(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err errors.Error) {
	geoEventKey := storageHelper.EventGeoKey(accountID, applicationID)

	eventKeys, er := georedis.SearchByRadius(e.redis, geoEventKey, latitude, longitude, radius, 52)
	if er != nil {
		return nil, errors.NewInternalError("failed to search for events by geo (1)", er.Error())
	}

	resultList, er := e.redis.MGet(eventKeys...).Result()
	if err != nil {
		return nil, errors.NewInternalError("failed to search for events by geo (2)", er.Error())
	}

	return e.toEvents(resultList)
}

func (e *event) ObjectSearch(accountID, applicationID int64, objectKey string) ([]*entity.Event, errors.Error) {
	objectEventKey := storageHelper.EventObjectKey(accountID, applicationID, objectKey)

	return e.fetchEventsFromKeys(accountID, applicationID, objectEventKey)
}

func (e *event) LocationSearch(accountID, applicationID int64, locationKey string) ([]*entity.Event, errors.Error) {
	locationEventKey := storageHelper.EventLocationKey(accountID, applicationID, locationKey)

	return e.fetchEventsFromKeys(accountID, applicationID, locationEventKey)
}

// fetchEventsFromKeys returns all the events matching a certain search key from the specified bucket
func (e *event) fetchEventsFromKeys(accountID, applicationID int64, bucketName string) ([]*entity.Event, errors.Error) {
	_, keys, er := e.redis.SScan(bucketName, 0, "", 300).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the events (1)", er.Error())
	}

	if len(keys) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, err := e.redis.MGet(keys...).Result()
	if err != nil {
		return nil, errors.NewInternalError("failed to read the events (2)", er.Error())
	}

	return e.toEvents(resultList)
}

// toEvents converts the events from json format to go structs
func (e *event) toEvents(resultList []interface{}) ([]*entity.Event, errors.Error) {
	events := []*entity.Event{}
	for _, result := range resultList {
		if result == nil {
			continue
		}
		event := &entity.Event{}
		if er := json.Unmarshal([]byte(result.(string)), event); er != nil {
			return []*entity.Event{}, errors.NewInternalError("failed to read the event from list (1)", er.Error())
		}
		events = append(events, event)
		event = &entity.Event{}
	}

	return events, nil
}

func (e *event) addEventGeo(accountID, applicationID int64, key string, updatedEvent *entity.Event) errors.Error {
	coordinates := georedis.GeoKey{
		Lat:   updatedEvent.Latitude,
		Lon:   updatedEvent.Longitude,
		Label: key,
	}

	geoEventKey := storageHelper.EventGeoKey(accountID, applicationID)
	_, er := georedis.AddCoordinates(e.redis, geoEventKey, 52, coordinates)
	if er == nil {
		return nil
	}
	return errors.NewInternalError("failed to add the event by geo (1)", er.Error())
}

func (e *event) removeEventGeo(accountID, applicationID int64, key string, updatedEvent *entity.Event) errors.Error {
	geoEventKey := storageHelper.EventGeoKey(accountID, applicationID)
	_, er := georedis.RemoveCoordinatesByKeys(e.redis, geoEventKey, key)
	if er == nil {
		return nil
	}
	return errors.NewInternalError("failed to remove the event by geo (1)", er.Error())
}

func (e *event) addEventObject(accountID, applicationID int64, key string, updatedEvent *entity.Event) errors.Error {
	objectEventKey := storageHelper.EventObjectKey(accountID, applicationID, updatedEvent.Object.ID)
	er := e.redis.SAdd(objectEventKey, key).Err()
	if er == nil {
		return nil
	}
	return errors.NewInternalError("failed to add the event by object (1)", er.Error())
}

func (e *event) removeEventObject(accountID, applicationID int64, key string, updatedEvent *entity.Event) errors.Error {
	objectEventKey := storageHelper.EventObjectKey(accountID, applicationID, updatedEvent.Object.ID)
	er := e.redis.SRem(objectEventKey, key).Err()
	if er == nil {
		return nil
	}
	return errors.NewInternalError("failed to remove the event by geo (1)", er.Error())
}

func (e *event) addEventLocation(accountID, applicationID int64, key string, updatedEvent *entity.Event) errors.Error {
	locationEventKey := storageHelper.EventLocationKey(accountID, applicationID, updatedEvent.Location)
	er := e.redis.SAdd(locationEventKey, key).Err()
	if er == nil {
		return nil
	}
	return errors.NewInternalError("failed to add the event by location (1)", er.Error())
}

func (e *event) removeEventLocation(accountID, applicationID int64, key string, updatedEvent *entity.Event) errors.Error {
	locationEventKey := storageHelper.EventLocationKey(accountID, applicationID, updatedEvent.Location)
	er := e.redis.SRem(locationEventKey, key).Err()
	if er == nil {
		return nil
	}
	return errors.NewInternalError("failed to remove the event by location (1)", er.Error())
}

// NewEvent creates a new Event
func NewEvent(storageClient redis.Client) core.Event {
	return &event{
		storage: storageClient,
		redis:   storageClient.Datastore(),
	}
}
