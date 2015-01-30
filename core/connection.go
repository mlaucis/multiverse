/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"

	"github.com/tapglue/backend/core/entity"
)

// ReadConnectionList returns all connections from a certain user
func ReadConnectionList(applicationID, userID int64) (users []*entity.User, err error) {
	key := storageClient.ConnectionUsersKey(applicationID, userID)

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
	key := storageClient.ConnectionKey(connection.ApplicationID, connection.UserFromID, connection.UserToID)

	// Write resource
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := storageClient.ConnectionsKey(connection.ApplicationID, connection.UserFromID)

	// Write list
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	userListKey := storageClient.ConnectionUsersKey(connection.ApplicationID, connection.UserFromID)

	// Generate following key
	userKey := storageClient.UserKey(connection.ApplicationID, connection.UserToID)

	// Write list
	if err = storageEngine.LPush(userListKey, userKey).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	followerListKey := storageClient.FollowedByUsersKey(connection.ApplicationID, connection.UserToID)

	// Generate follower key
	followerKey := storageClient.UserKey(connection.ApplicationID, connection.UserFromID)

	// Write list
	if err = storageEngine.LPush(followerListKey, followerKey).Err(); err != nil {
		return nil, err
	}

	// TODO: Add events of user "user_to_id" to list of user "user_from_id" order by date

	if !retrieve {
		return connection, nil
	}

	// Return resource
	return connection, nil
}
