/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/utils"
)

// ReadUser returns the user matching the ID or an error
func ReadApplicationUser(accountID, applicationID, userID int64) (user *entity.User, err error) {
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
	storedUser, err := ReadApplicationUser(user.AccountID, user.ApplicationID, user.ID)
	if err != nil {
		return nil, err
	}

	storedUser = user

	// Encrypt password - we should do this only if the password changes
	storedUser.Password = storageClient.EncryptPassword(user.Password)

	val, err := json.Marshal(storedUser)
	if err != nil {
		return nil, err
	}

	existingUser, err := ReadApplicationUser(user.AccountID, user.ApplicationID, user.ID)
	if err != nil {
		return nil, err
	}

	key := storageClient.User(user.AccountID, user.ApplicationID, user.ID)
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(existingUser.Email))
	usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(existingUser.Username))
	_, err = storageEngine.Del(emailListKey, usernameListKey).Result()

	if !storedUser.Enabled {
		listKey := storageClient.Users(user.AccountID, user.ApplicationID)
		if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
			return nil, err
		}
	} else {
		emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(user.Email))
		err = storageEngine.Set(emailListKey, fmt.Sprintf("%d", user.ID)).Err()
		if err != nil {
			return nil, err
		}

		usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(user.Username))
		err = storageEngine.Set(usernameListKey, fmt.Sprintf("%d", user.ID)).Err()

		if err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return storedUser, nil
	}

	return ReadApplicationUser(user.AccountID, user.ApplicationID, user.ID)
}

// DeleteUser deletes the user matching the IDs or an error
func DeleteUser(accountID, applicationID, userID int64) (err error) {
	user, err := ReadApplicationUser(accountID, applicationID, userID)
	if err != nil {
		return err
	}

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

	emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(user.Email))
	usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(user.Username))
	_, err = storageEngine.Del(emailListKey, usernameListKey).Result()

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
	user.UpdatedAt = user.CreatedAt
	user.LastLogin, err = time.Parse(time.RFC3339, "0000-01-01T00:00:00Z")
	if err != nil {
		return nil, err
	}

	if user.ID, err = storageClient.GenerateApplicationUserID(user.ApplicationID); err != nil {
		return nil, err
	}

	// Encrypt password
	user.Password = storageClient.EncryptPassword(user.Password)

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

	emailListKey := storageClient.ApplicationUserByEmail(user.AccountID, user.ApplicationID, utils.Base64Encode(user.Email))
	result, err := storageEngine.SetNX(emailListKey, fmt.Sprintf("%d", user.ID)).Result()
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, fmt.Errorf("user with email address already exists")
	}

	usernameListKey := storageClient.ApplicationUserByUsername(user.AccountID, user.ApplicationID, utils.Base64Encode(user.Username))
	result, err = storageEngine.SetNX(usernameListKey, fmt.Sprintf("%d", user.ID)).Result()
	if err != nil {
		return nil, err
	}
	if !result {
		return nil, fmt.Errorf("user with email address already exists")
	}

	listKey := storageClient.Users(user.AccountID, user.ApplicationID)
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return user, nil
	}

	return ReadApplicationUser(user.AccountID, user.ApplicationID, user.ID)
}

// CreateUserSession handles the creation of a user session and returns the session token
func CreateApplicationUserSession(user *entity.User) (string, error) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)
	token := storageClient.GenerateApplicationSessionID(user)

	if err := storageEngine.Set(sessionKey, token).Err(); err != nil {
		return "", err
	}

	expired, err := storageEngine.Expire(sessionKey, storageClient.SessionTimeoutDuration()).Result()
	if err != nil {
		return "", err
	}
	if !expired {
		return "", fmt.Errorf("could not set expire time")
	}

	return token, nil
}

// RefreshUserSession generates a new session token for the user session
func RefreshApplicationUserSession(sessionToken string, user *entity.User) (string, error) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)

	storedToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", err
	}

	if storedToken != sessionToken {
		return "", fmt.Errorf("session token mismatch")
	}

	token := storageClient.GenerateApplicationSessionID(user)

	if err := storageEngine.Set(sessionKey, token).Err(); err != nil {
		return "", err
	}

	expired, err := storageEngine.Expire(sessionKey, storageClient.SessionTimeoutDuration()).Result()
	if err != nil {
		return "", err
	}
	if !expired {
		return "", fmt.Errorf("could not set expire time")
	}

	return token, nil
}

// GetApplicationUserSession returns the application user session
func GetApplicationUserSession(user *entity.User) (string, error) {
	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)
	storedSessionToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", fmt.Errorf("could not fetch session from storage")
	}

	if storedSessionToken == "" {
		return "", fmt.Errorf("session not found")
	}

	return storedSessionToken, nil
}

// DestroyUserSession removes the user session
func DestroyApplicationUserSession(sessionToken string, user *entity.User) error {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)

	storedToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return err
	}

	if storedToken != sessionToken {
		return fmt.Errorf("session token mismatch")
	}

	result, err := storageEngine.Del(sessionKey).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		return fmt.Errorf("invalid session")
	}

	return nil
}

// ApplicationuserByEmailExists checks if an application user exists by searching it via the email
func ApplicationUserByEmailExists(accountID, applicationID int64, email string) (bool, error) {
	emailListKey := storageClient.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))

	return storageEngine.Exists(emailListKey).Result()
}

// FindApplicationUserByEmail returns an application user by its email
func FindApplicationUserByEmail(accountID, applicationID int64, email string) (*entity.User, error) {
	emailListKey := storageClient.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))

	return findApplicationUserByKey(accountID, applicationID, emailListKey)
}

// FindApplicationUserByUsername returns an application user by its username
func FindApplicationUserByUsername(accountID, applicationID int64, username string) (*entity.User, error) {
	usernameListKey := storageClient.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(username))

	return findApplicationUserByKey(accountID, applicationID, usernameListKey)
}

// findApplicationUserByKey returns an application user regardless of the key used to search for him
func findApplicationUserByKey(accountID, applicationID int64, bucketName string) (*entity.User, error) {
	storedValue, err := storageEngine.Get(bucketName).Result()
	if err != nil {
		return nil, err
	}

	userID, err := strconv.ParseInt(storedValue, 10, 64)
	if err != nil {
		return nil, err
	}

	applicationUser, err := ReadApplicationUser(accountID, applicationID, userID)
	if err != nil {
		return nil, err
	}

	return applicationUser, nil
}
