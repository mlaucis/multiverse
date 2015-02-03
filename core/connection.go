/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"

	"fmt"

	"github.com/tapglue/backend/core/entity"
	red "gopkg.in/redis.v2"
)

// ReadConnectionList returns all connections from a certain user
func ReadConnectionList(applicationID, userID int64) (users []*entity.User, err error) {
	key := storageClient.ConnectionUsers(applicationID, userID)

	// Read from db
	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no connections for this user")
		return nil, err
	}

	// Read from db
	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	user := &entity.User{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), user); err != nil {
			return nil, err
		}
		users = append(users, user)
		user = &entity.User{}
	}

	return
}

// WriteConnection adds a user connection to the database and returns the created user connection or an error
func WriteConnection(connection *entity.Connection, retrieve bool) (con *entity.Connection, err error) {
	// Encode JSON
	val, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	// Generate resource key
	key := storageClient.Connection(connection.ApplicationID, connection.UserFromID, connection.UserToID)

	// Write resource
	exist, err := storageEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, fmt.Errorf("user connection already exists")
	}
	if err != nil {
		return nil, err
	}

	// Generate list key
	listKey := storageClient.Connections(connection.ApplicationID, connection.UserFromID)

	// Write list
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	userListKey := storageClient.ConnectionUsers(connection.ApplicationID, connection.UserFromID)

	// Generate following key
	userKey := storageClient.User(connection.ApplicationID, connection.UserToID)

	// Write list
	if err = storageEngine.LPush(userListKey, userKey).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	followerListKey := storageClient.FollowedByUsers(connection.ApplicationID, connection.UserToID)

	// Generate follower key
	followerKey := storageClient.User(connection.ApplicationID, connection.UserFromID)

	// Write list
	if err = storageEngine.LPush(followerListKey, followerKey).Err(); err != nil {
		return nil, err
	}

	// Write connection events to list
	if err = WriteConnectionEventsToList(connection); err != nil {
		return nil, err
	}

	if !retrieve {
		return connection, nil
	}

	// Return resource
	return connection, nil
}

// WriteConnectionEventsToList takes a connection and writes the events to the lists
func WriteConnectionEventsToList(connection *entity.Connection) (err error) {

	// Generate list key (UserFromID connection events)
	connectionEventsKey := storageClient.ConnectionEvents(connection.ApplicationID, connection.UserFromID)

	// Generate list key (UserToID events)
	eventsKey := storageClient.Events(connection.ApplicationID, connection.UserToID)

	// Read events
	events, err := storageEngine.ZRevRangeWithScores(eventsKey, "0", "-1").Result()
	if err != nil {
		return err
	}

	// Sync if events exist
	if len(events) >= 1 {
		var vals []red.Z

		for _, eventKey := range events {
			val := red.Z{Score: float64(eventKey.Score), Member: eventKey.Member}
			vals = append(vals, val)
		}

		// Write list
		if err = storageEngine.ZAdd(connectionEventsKey, vals...).Err(); err != nil {
			return err
		}
	}

	return nil
}
