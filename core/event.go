/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/redis"
)

// Defining keys
const (
	EventKey                string = "app_%d_user_%d_event_%d"
	EventsKey               string = "app_%d_user_%d_events"
	ConnectionEventsKey     string = "app_%d_user_%d_connection_events"
	ConnectionEventsKeyLoop string = "%s_connection_events"
)

// generateEventID generates a new event ID
func generateEventID(applicationID int64, userID int64) (int64, error) {
	incr := redis.Client().Incr(fmt.Sprintf("ids_application_%d_user_%d_event", applicationID, userID))
	return incr.Result()
}

// ReadEvent returns the event matching the ID or an error
func ReadEvent(applicationID int64, userID int64, eventID int64) (event *entity.Event, err error) {
	// Generate resource key
	key := fmt.Sprintf(EventKey, applicationID, userID, eventID)

	// Read from db
	result, err := redis.Client().Get(key).Result()
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
func ReadEventList(applicationID int64, userID int64) (events []*entity.Event, err error) {
	// Generate resource key
	key := fmt.Sprintf(EventsKey, applicationID, userID)

	// Read from db
	result, err := redis.Client().LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no events for this user")
		return nil, err
	}

	// Read from db
	resultList, err := redis.Client().MGet(result...).Result()
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
func ReadConnectionEventList(applicationID int64, userID int64) (events []*entity.Event, err error) {
	// Generate resource key
	key := fmt.Sprintf(ConnectionEventsKey, applicationID, userID)

	// Read from db
	result, err := redis.Client().LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no events from connections")
		return nil, err
	}

	// Read from db
	resultList, err := redis.Client().MGet(result...).Result()
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

// WriteEvent adds an event to the database and returns the created event or an error
func WriteEvent(event *entity.Event, retrieve bool) (evn *entity.Event, err error) {

	// Generate id
	if event.ID, err = generateEventID(event.ApplicationID, event.UserID); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(event)
	if err != nil {
		return nil, err
	}

	// Generate resource key
	key := fmt.Sprintf(EventKey, event.ApplicationID, event.UserID, event.ID)

	// Write resource
	if err = redis.Client().Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := fmt.Sprintf(EventsKey, event.ApplicationID, event.UserID)

	// Write list
	if err = redis.Client().LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	// Generate connections key
	connectionsKey := fmt.Sprintf(ConnectionUsersKey, event.ApplicationID, event.UserID)

	// Read connections
	connections, err := redis.Client().LRange(connectionsKey, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Write to connections lists
	for _, each := range connections {
		// Create Key
		feedKey := fmt.Sprintf(ConnectionEventsKeyLoop, each)
		// Write to lists
		if err = redis.Client().LPush(feedKey, key).Err(); err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return event, nil
	}

	// Return resource
	return ReadEvent(event.ApplicationID, event.UserID, event.ID)
}
