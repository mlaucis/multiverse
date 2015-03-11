/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tapglue/backend/core/entity"

	"github.com/tapglue/georedis"

	red "gopkg.in/redis.v2"
)

// ReadEvent returns the event matching the ID or an error
func ReadEvent(accountID, applicationID, userID, eventID int64) (event *entity.Event, err error) {
	key := storageClient.Event(accountID, applicationID, userID, eventID)

	result, err := storageEngine.Get(key).Result()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &event); err != nil {
		return nil, err
	}

	return
}

// UpdateEvent updates an event in the database and returns the updated event or an error
func UpdateEvent(event *entity.Event, retrieve bool) (evn *entity.Event, err error) {
	event.UpdatedAt = time.Now()

	val, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	key := storageClient.Event(event.AccountID, event.ApplicationID, event.UserID, event.ID)
	exist, err := storageEngine.Exists(key).Result()
	if !exist {
		return nil, fmt.Errorf("event does not exist")
	}
	if err != nil {
		return nil, err
	}

	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	coordinates := georedis.GeoKey{
		Lat:   event.Latitude,
		Lon:   event.Longitude,
		Label: key,
	}

	geoEventKey := storageClient.EventGeoKey(event.AccountID, event.ApplicationID)

	georedis.RemoveCoordinatesByKeys(storageEngine, geoEventKey, key)
	georedis.AddCoordinates(storageEngine, geoEventKey, 52, coordinates)

	objectEventKey := storageClient.EventObjectKey(event.AccountID, event.ApplicationID, event.Object.ID)
	if err = storageEngine.SRem(objectEventKey, key).Err(); err != nil {
		return nil, err
	}

	if err = storageEngine.SAdd(objectEventKey, key).Err(); err != nil {
		return nil, err
	}

	locationEventKey := storageClient.EventLocationKey(event.AccountID, event.ApplicationID, event.Location)
	if err = storageEngine.SRem(locationEventKey, key).Err(); err != nil {
		return nil, err
	}
	if err = storageEngine.SAdd(locationEventKey, key).Err(); err != nil {
		return nil, err
	}

	if !event.Enabled {
		listKey := storageClient.Events(event.AccountID, event.ApplicationID, event.UserID)
		if err = storageEngine.ZRem(listKey, key).Err(); err != nil {
			return nil, err
		}
		if err = DeleteEventFromConnectionsLists(event.AccountID, event.ApplicationID, event.UserID, key); err != nil {
			return nil, err
		}

	}

	if !retrieve {
		return event, nil
	}

	return ReadEvent(event.AccountID, event.ApplicationID, event.UserID, event.ID)
}

// DeleteEvent deletes the event matching the IDs or an error
func DeleteEvent(accountID, applicationID, userID, eventID int64) (err error) {
	key := storageClient.Event(accountID, applicationID, userID, eventID)

	event, err := ReadEvent(accountID, applicationID, userID, eventID)
	if err != nil {
		return err
	}

	result, err := storageEngine.Del(key).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		return fmt.Errorf("The resource for the provided id doesn't exist")
	}

	listKey := storageClient.Events(accountID, applicationID, userID)
	if err = storageEngine.ZRem(listKey, key).Err(); err != nil {
		return err
	}

	if err = DeleteEventFromConnectionsLists(accountID, applicationID, userID, key); err != nil {
		return err
	}

	geoEventKey := storageClient.EventGeoKey(accountID, applicationID)
	if _, err = georedis.RemoveCoordinatesByKeys(storageEngine, geoEventKey, key); err != nil {
		return err
	}

	objectEventKey := storageClient.EventObjectKey(accountID, applicationID, event.Object.ID)
	if err = storageEngine.SRem(objectEventKey, key).Err(); err != nil {
		return err
	}

	locationEventKey := storageClient.EventLocationKey(accountID, applicationID, event.Location)
	if err = storageEngine.SRem(locationEventKey, key).Err(); err != nil {
		return err
	}

	return nil
}

// ReadEventList returns all events from a certain user
func ReadEventList(accountID, applicationID, userID int64) (events []*entity.Event, err error) {
	key := storageClient.Events(accountID, applicationID, userID)

	result, err := storageEngine.ZRevRange(key, "0", "-1").Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		err := errors.New("There are no events for this user")
		return nil, err
	}

	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	return toEvents(resultList)
}

// ReadConnectionEventList returns all events from connections
func ReadConnectionEventList(accountID, applicationID, userID int64) (events []*entity.Event, err error) {
	key := storageClient.ConnectionEvents(accountID, applicationID, userID)

	result, err := storageEngine.ZRevRange(key, "0", "-1").Result()
	if err != nil {
		return nil, err
	}

	// TODO maybe this shouldn't be an error but rather return that there are no events from connections
	if len(result) == 0 {
		err := errors.New("There are no events from connections")
		return nil, err
	}

	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	return toEvents(resultList)
}

// WriteEvent adds an event to the database and returns the created event or an error
func WriteEvent(event *entity.Event, retrieve bool) (evn *entity.Event, err error) {
	event.Enabled = true
	event.CreatedAt = time.Now()
	event.UpdatedAt, event.ReceivedAt = event.CreatedAt, event.CreatedAt

	if event.ID, err = storageClient.GenerateApplicationEventID(event.ApplicationID); err != nil {
		return nil, err
	}

	val, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	key := storageClient.Event(event.AccountID, event.ApplicationID, event.UserID, event.ID)
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	listKey := storageClient.Events(event.AccountID, event.ApplicationID, event.UserID)

	setVal := red.Z{Score: float64(event.ReceivedAt.Unix()), Member: key}
	if err = storageEngine.ZAdd(listKey, setVal).Err(); err != nil {
		return nil, err
	}

	coordinates := georedis.GeoKey{
		Lat:   event.Latitude,
		Lon:   event.Longitude,
		Label: key,
	}

	geoEventKey := storageClient.EventGeoKey(event.AccountID, event.ApplicationID)
	georedis.AddCoordinates(storageEngine, geoEventKey, 52, coordinates)

	objectEventKey := storageClient.EventObjectKey(event.AccountID, event.ApplicationID, event.Object.ID)
	if err = storageEngine.SAdd(objectEventKey, key).Err(); err != nil {
		return nil, err
	}

	locationEventKey := storageClient.EventLocationKey(event.AccountID, event.ApplicationID, event.Location)
	if err = storageEngine.SAdd(locationEventKey, key).Err(); err != nil {
		return nil, err
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
func WriteEventToConnectionsLists(event *entity.Event, key string) (err error) {
	connectionsKey := storageClient.FollowedByUsers(event.AccountID, event.ApplicationID, event.UserID)

	connections, err := storageEngine.LRange(connectionsKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, userKey := range connections {
		feedKey := storageClient.ConnectionEventsLoop(userKey)

		val := red.Z{Score: float64(event.ReceivedAt.Unix()), Member: key}
		if err = storageEngine.ZAdd(feedKey, val).Err(); err != nil {
			return err
		}
	}

	return nil
}

// DeleteEventFromConnectionsLists takes a user id and key and deletes it to the user connections list
func DeleteEventFromConnectionsLists(accountID, applicationID, userID int64, key string) (err error) {
	connectionsKey := storageClient.FollowedByUsers(accountID, applicationID, userID)

	connections, err := storageEngine.LRange(connectionsKey, 0, -1).Result()
	if err != nil {
		return err
	}

	for _, userKey := range connections {
		feedKey := storageClient.ConnectionEventsLoop(userKey)
		if err = storageEngine.ZRem(feedKey, key).Err(); err != nil {
			return err
		}
	}

	return nil
}

// SearchGeoEvents retrieves all the events from an application within a radius of the provided coordinates
func SearchGeoEvents(accountID, applicationID int64, latitude, longitude, radius float64) (events []*entity.Event, err error) {
	geoEventKey := storageClient.EventGeoKey(accountID, applicationID)

	eventKeys, err := georedis.SearchByRadius(storageEngine, geoEventKey, latitude, longitude, radius, 52)
	if err != nil {
		return events, err
	}

	resultList, err := storageEngine.MGet(eventKeys...).Result()
	if err != nil {
		return nil, err
	}

	return toEvents(resultList)
}

// SearchObjectEvents returns all the events for a specific object
func SearchObjectEvents(accountID, applicationID int64, objectKey string) ([]*entity.Event, error) {
	objectEventKey := storageClient.EventObjectKey(accountID, applicationID, objectKey)

	return fetchEventsFromKeys(accountID, applicationID, objectEventKey, objectKey)
}

// SearchLocationEvents returns all the events for a specific object
func SearchLocationEvents(accountID, applicationID int64, locationKey string) ([]*entity.Event, error) {
	locationEventKey := storageClient.EventLocationKey(accountID, applicationID, locationKey)

	return fetchEventsFromKeys(accountID, applicationID, locationEventKey, locationKey)
}

// fetchEventsFromKeys returns all the events matching a certain search key from the specified bucket
func fetchEventsFromKeys(accountID, applicationID int64, bucketName, searchKey string) ([]*entity.Event, error) {
	_, keys, err := storageEngine.SScan(bucketName, 0, searchKey, 300).Result()
	if err != nil {
		return nil, err
	}

	if len(keys) == 0 {
		return []*entity.Event{}, nil
	}

	resultList, err := storageEngine.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	return toEvents(resultList)
}

// toEvents converts the events from json format to go structs
func toEvents(resultList []interface{}) ([]*entity.Event, error) {
	events := []*entity.Event{}
	for _, result := range resultList {
		event := &entity.Event{}
		if err := json.Unmarshal([]byte(result.(string)), event); err != nil {
			return []*entity.Event{}, err
		}
		events = append(events, event)
		event = &entity.Event{}
	}

	return events, nil
}
