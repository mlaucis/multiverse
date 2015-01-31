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
	UserToKey          string = "app_%d_user_%d"
	ConnectionKey      string = "app_%d_user_%d_connection_%d"
	ConnectionsKey     string = "app_%d_user_%d_connections"
	ConnectionUsersKey string = "app_%d_user_%d_connection_users"
)

// ReadConnectionList returns all connections from a certain user
func ReadConnectionList(applicationID int64, userID int64) (users []*entity.User, err error) {
	// Generate resource key
	key := fmt.Sprintf(ConnectionUsersKey, applicationID, userID)

	// Read from db
	result, err := redis.Client().LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no connections for this user")
		return nil, err
	}

	// Read from db
	resultList, err := redis.Client().MGet(result...).Result()
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
	key := fmt.Sprintf(ConnectionKey, connection.ApplicationID, connection.UserFromID, connection.UserToID)

	// Write resource
	if err = redis.Client().Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := fmt.Sprintf(ConnectionsKey, connection.ApplicationID, connection.UserFromID)

	// Write list
	if err = redis.Client().LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	userListKey := fmt.Sprintf(ConnectionUsersKey, connection.ApplicationID, connection.UserFromID)

	// Generate user key
	userToKey := fmt.Sprintf(UserToKey, connection.ApplicationID, connection.UserToID)

	// Write list
	if err = redis.Client().LPush(userListKey, userToKey).Err(); err != nil {
		return nil, err
	}

	// TODO: Add events of user "user_to_id" to list of user "user_from_id" order by date

	if !retrieve {
		return connection, nil
	}

	// Return resource
	return connection, nil
}
