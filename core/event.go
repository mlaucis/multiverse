/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"

	"github.com/tapglue/backend/core/entity"
	red "gopkg.in/redis.v2"
)

// generateEventID generates a new event ID
func generateEventID(applicationID int64) (int64, error) {
	return storageEngine.Incr(storageClient.GenerateApplicationEventID(applicationID)).Result()
}

// ReadEvent returns the event matching the ID or an error
func ReadEvent(applicationID, userID, eventID int64) (event *entity.Event, err error) {
	// Generate resource key
	key := storageClient.EventKey(applicationID, userID, eventID)

	// Read from db
	result, err := storageEngine.Get(key).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &event); err != nil {
		return nil, err
	}

	return
}

// ReadEventList returns all events from a certain user
func ReadEventList(applicationID, userID int64) (events []*entity.Event, err error) {
	// Generate resource key
	key := storageClient.EventsKey(applicationID, userID)

	// Read from db
	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no events for this user")
		return nil, err
	}

	// Read from db
	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	event := &entity.Event{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), event); err != nil {
			return nil, err
		}
		events = append(events, event)
		event = &entity.Event{}
	}

	return
}

// ReadConnectionEventList returns all events from connections
func ReadConnectionEventList(applicationID, userID int64) (events []*entity.Event, err error) {
	// Generate resource key
	key := storageClient.ConnectionEventsKey(applicationID, userID)

	// Read from db
	result, err := storageEngine.ZRevRange(key, "0", "-1").Result()
	if err != nil {
		panic(err)
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no events from connections")
		panic(err)
		return nil, err
	}

	// Read from db
	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		panic(err)
		return nil, err
	}

	// Parse JSON
	event := &entity.Event{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), event); err != nil {
			panic(err)
			return nil, err
		}
		events = append(events, event)
		event = &entity.Event{}
	}

	return
}

// WriteEvent adds an event to the database and returns the created event or an error
func WriteEvent(event *entity.Event, retrieve bool) (evn *entity.Event, err error) {

	// Generate id
	if event.ID, err = generateEventID(event.ApplicationID); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	// Generate resource key
	key := storageClient.EventKey(event.ApplicationID, event.UserID, event.ID)

	// Write resource
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := storageClient.EventsKey(event.ApplicationID, event.UserID)

	// Write list
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	// Generate connections key
	connectionsKey := storageClient.FollowedByUsersKey(event.ApplicationID, event.UserID)

	// Read connections
	connections, err := storageEngine.LRange(connectionsKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Write to connections lists
	for _, userID := range connections {
		// Create Key
		feedKey := storageClient.ConnectionEventsKeyLoop(userID)

		// Write to lists
		val := red.Z{Score: float64(event.ReceivedAt.Unix()), Member: key}
		if err = storageEngine.ZAdd(feedKey, val).Err(); err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return event, nil
	}

	// Return resource
	return ReadEvent(event.ApplicationID, event.UserID, event.ID)
}
