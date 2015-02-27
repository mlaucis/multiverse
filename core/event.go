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
func ReadConnectionEventList(accountID, applicationID, userID int64) (events []*entity.Event, err error) {
	key := storageClient.ConnectionEvents(accountID, applicationID, userID)

	result, err := storageEngine.ZRevRange(key, "0", "-1").Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		err := errors.New("There are no events from connections")
		return nil, err
	}

	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

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
