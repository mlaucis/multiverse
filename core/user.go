/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"
	"time"

	"fmt"

	"github.com/tapglue/backend/core/entity"
)

func generateUserID(applicationID int64) (int64, error) {
	return storageEngine.Incr(storageClient.GenerateApplicationUserID(applicationID)).Result()
}

// ReadUser returns the user matching the ID or an error
func ReadUser(accountID, applicationID, userID int64) (user *entity.User, err error) {
	key := storageClient.User(accountID, applicationID, userID)

	result, err := storageEngine.Get(key).Result()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &user); err != nil {
		return nil, err
	}

	return
}

// UpdateUser updates a user in the database and returns the updates user or an error
func UpdateUser(user *entity.User, retrieve bool) (usr *entity.User, err error) {
	user.UpdatedAt = time.Now()

	val, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	key := storageClient.User(user.AccountID, user.ApplicationID, user.ID)
	exist, err := storageEngine.Exists(key).Result()
	if !exist {
		return nil, fmt.Errorf("user does not exist")
	}
	if err != nil {
		return nil, err
	}

	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	if !user.Enabled {
		listKey := storageClient.Users(user.AccountID, user.ApplicationID)
		if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return user, nil
	}

	return ReadUser(user.AccountID, user.ApplicationID, user.ID)
}

// DeleteUser deletes the user matching the IDs or an error
func DeleteUser(accountID, applicationID, userID int64) (err error) {
	key := storageClient.User(accountID, applicationID, userID)
	result, err := storageEngine.Del(key).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		return fmt.Errorf("The resource for the provided id doesn't exist")
	}

	listKey := storageClient.Users(accountID, applicationID)
	if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
		return err
	}

	// TODO: Remove Users Connections?
	// TODO: Remove Users Connection Lists?
	// TODO: Remove User in other Users Connection Lists?
	// TODO: Remove Users Events?
	// TODO: Remove Users Events from Lists?

	return nil
}

// ReadUserList returns all users from a certain account
func ReadUserList(accountID, applicationID int64) (users []*entity.User, err error) {
	key := storageClient.Users(accountID, applicationID)

	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		err := errors.New("There are no users for this app")
		return nil, err
	}

	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

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
	user.Enabled = true
	user.CreatedAt = time.Now()
	user.UpdatedAt, user.ReceivedAt = user.CreatedAt, user.CreatedAt

	if user.ID, err = generateUserID(user.ApplicationID); err != nil {
		return nil, err
	}

	val, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	key := storageClient.User(user.AccountID, user.ApplicationID, user.ID)

	exist, err := storageEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, fmt.Errorf("user already exists")
	}
	if err != nil {
		return nil, err
	}

	listKey := storageClient.Users(user.AccountID, user.ApplicationID)

	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return user, nil
	}

	return ReadUser(user.AccountID, user.ApplicationID, user.ID)
}

// CreateUserSession handles the creation of a user session and returns the session token
func CreateUserSession(user *entity.User) (string, error) {
	// TODO properly generate this
	// It should look like this: base64(dateGenerated:randomToken)

	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	return "fake-token", nil
}

// RefreshUserSession generates a new session token for the user session
func RefreshUserSession(sessionToken string, user *entity.User) (string, error) {
	// TODO rate limit this to one per 6 hours as we generally want to have sessions that last at least a day/month/year/decade?
	return "new-fake-token", nil
}

// DestroyUserSession removes the user session
func DestroyUserSession(sessionToken string, user *entity.User) error {
	// TODO properly implement this
	return nil
}
