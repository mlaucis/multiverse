/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"

	"github.com/tapglue/backend/core/entity"
)

// Defining keys
const (
	_UserKey  string = "app_%d_user_%d"
	_UsersKey string = "app_%d_users"
)

// generateUserID generates a new user ID
func generateUserID(applicationID int64) (int64, error) {
	return storageEngine.Incr(storageClient.GenerateApplicationUserID(applicationID)).Result()
}

// ReadUser returns the user matching the ID or an error
func ReadUser(applicationID int64, userID int64) (user *entity.User, err error) {
	// Generate resource key
	key := storageClient.UserKey(applicationID, userID)

	// Read from db
	result, err := storageEngine.Get(key).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &user); err != nil {
		return nil, err
	}

	return
}

// ReadUserList returns all users from a certain account
func ReadUserList(applicationID int64) (users []*entity.User, err error) {
	// Generate resource key
	key := storageClient.UsersKey(applicationID)

	// Read from db
	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	// Return no elements
	if len(result) == 0 {
		err := errors.New("There are no users for this app")
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

// WriteUser adds a user to the database and returns the created user or an error
func WriteUser(user *entity.User, retrieve bool) (usr *entity.User, err error) {
	// Generate id
	if user.ID, err = generateUserID(user.ApplicationID); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	// Generate resource key
	key := storageClient.UserKey(user.ApplicationID, user.ID)

	// Write resource
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := storageClient.UsersKey(user.ApplicationID)

	// Write list
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return user, nil
	}

	// Return resource
	return ReadUser(user.ApplicationID, user.ID)
}
