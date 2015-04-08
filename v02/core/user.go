/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/entity"
)

// ReadApplicationUser returns the user matching the ID or an error
func ReadApplicationUser(accountID, applicationID, userID int64) (user *entity.User, err *tgerrors.TGError) {
	key := storageClient.User(accountID, applicationID, userID)

	result, er := redisEngine.Get(key).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read application user (1)", er.Error())
	}

	if er = json.Unmarshal([]byte(result), &user); er != nil {
		return nil, tgerrors.NewInternalError("failed to read application user (2)", er.Error())
	}

	return
}

// UpdateUser updates a user in the database and returns the updates user or an error
func UpdateUser(existingUser, updatedUser entity.User, retrieve bool) (usr *entity.User, err *tgerrors.TGError) {

	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	} else if updatedUser.Password != existingUser.Password {
		// Encrypt password - we should do this only if the password changes
		updatedUser.Password = storageClient.EncryptPassword(updatedUser.Password)
	}

	val, er := json.Marshal(updatedUser)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application user (1)", er.Error())
	}

	key := storageClient.User(updatedUser.AccountID, updatedUser.ApplicationID, updatedUser.ID)
	if er = redisEngine.Set(key, string(val)).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application user (2)", er.Error())
	}

	if existingUser.Email != updatedUser.Email {
		emailListKey := storageClient.ApplicationUserByEmail(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(existingUser.Email))
		_, er = redisEngine.Del(emailListKey).Result()

		emailListKey = storageClient.ApplicationUserByEmail(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(updatedUser.Email))
		er = redisEngine.Set(emailListKey, fmt.Sprintf("%d", updatedUser.ID)).Err()
		if er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (3)", er.Error())
		}
	}

	if existingUser.Username != updatedUser.Username {
		usernameListKey := storageClient.ApplicationUserByUsername(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(existingUser.Username))
		_, er = redisEngine.Del(usernameListKey).Result()

		usernameListKey = storageClient.ApplicationUserByUsername(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(updatedUser.Username))
		er = redisEngine.Set(usernameListKey, fmt.Sprintf("%d", updatedUser.ID)).Err()

		if er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (4)", er.Error())
		}
	}

	if !updatedUser.Enabled {
		listKey := storageClient.Users(updatedUser.AccountID, updatedUser.ApplicationID)
		if er = redisEngine.LRem(listKey, 0, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (5)", er.Error())
		}
	} else {
		listKey := storageClient.Users(updatedUser.AccountID, updatedUser.ApplicationID)
		if er = redisEngine.LPush(listKey, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (6)", er.Error())
		}
	}

	if !retrieve {
		return &updatedUser, nil
	}

	return ReadApplicationUser(updatedUser.AccountID, updatedUser.ApplicationID, updatedUser.ID)
}

// DeleteUser deletes the user matching the IDs or an error
func DeleteUser(accountID, applicationID, userID int64) (err *tgerrors.TGError) {
	user, err := ReadApplicationUser(accountID, applicationID, userID)
	if err != nil {
		return err
	}

	disabledUser := *user
	disabledUser.Enabled = false
	disabledUser.Password = ""
	_, err = UpdateUser(*user, disabledUser, false)

	// TODO: Remove Users Connections?
	// TODO: Remove Users Connection Lists?
	// TODO: Remove User in other Users Connection Lists?
	// TODO: Remove Users Events?
	// TODO: Remove Users Events from Lists?

	return

	// TODO Figure out if we should just simply remove the user or not

	/*key := storageClient.User(accountID, applicationID, userID)
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

	emailListKey := storageClient.AccountUserByEmail(Base64Encode(user.Email))
	usernameListKey := storageClient.AccountUserByUsername(Base64Encode(user.Username))
	_, err = storageEngine.Del(emailListKey, usernameListKey).Result()

	return nil
	*/
}

// ReadUserList returns all users from a certain account
func ReadUserList(accountID, applicationID int64) (users []*entity.User, err *tgerrors.TGError) {
	key := storageClient.Users(accountID, applicationID)

	result, er := redisEngine.LRange(key, 0, -1).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the application user list (1)", er.Error())
	}

	if len(result) == 0 {
		return users, nil
	}

	resultList, er := redisEngine.MGet(result...).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the application user list (2)", er.Error())
	}

	user := &entity.User{}
	for _, result := range resultList {
		if er = json.Unmarshal([]byte(result.(string)), user); er != nil {
			return nil, tgerrors.NewInternalError("failed to read the application user list (3)", er.Error())
		}
		users = append(users, user)
		user = &entity.User{}
	}

	return
}

// WriteUser adds a user to the database and returns the created user or an error
func WriteUser(user *entity.User, retrieve bool) (usr *entity.User, err *tgerrors.TGError) {
	// TODO We should introduce an option for the application to either allow for activated/deactivated behavior
	// and if they chose it, then we need to provide an endpoint to activate a user or not
	//user.Activated = true

	var er error
	user.Enabled = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.LastLogin, er = time.Parse(time.RFC3339, "0000-01-01T00:00:00Z")
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (1)", er.Error())
	}

	if user.ID, er = storageClient.GenerateApplicationUserID(user.ApplicationID); er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (2)", er.Error())
	}

	// Encrypt password
	user.Password = storageClient.EncryptPassword(user.Password)

	val, er := json.Marshal(user)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (3)", er.Error())
	}

	key := storageClient.User(user.AccountID, user.ApplicationID, user.ID)

	exist, er := redisEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, tgerrors.NewInternalError("failed to write the application user (4)", "duplicate user")
	}
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (5)", er.Error())
	}

	stringUserID := fmt.Sprintf("%d", user.ID)

	emailListKey := storageClient.ApplicationUserByEmail(user.AccountID, user.ApplicationID, utils.Base64Encode(user.Email))
	result, er := redisEngine.SetNX(emailListKey, stringUserID).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (6)", er.Error())
	}
	if !result {
		return nil, tgerrors.NewInternalError("failed to write the application user (7)", "duplicate user by e-mail")
	}

	usernameListKey := storageClient.ApplicationUserByUsername(user.AccountID, user.ApplicationID, utils.Base64Encode(user.Username))
	result, er = redisEngine.SetNX(usernameListKey, stringUserID).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (8)", er.Error())
	}
	if !result {
		return nil, tgerrors.NewInternalError("failed to write the application user (9)", "duplicate user by username")
	}

	socialValues := []string{}
	applicationSocialKey := ""
	for idx := range user.SocialIDs {
		applicationSocialKey = storageClient.SocialConnection(
			user.AccountID,
			user.ApplicationID,
			idx,
			utils.Base64Encode(user.SocialIDs[idx]),
		)
		socialValues = append(socialValues, applicationSocialKey, stringUserID)
	}

	if applicationSocialKey != "" {
		er := redisEngine.MSet(socialValues...).Err()
		if er != nil {
			return nil, tgerrors.NewInternalError("failed to write the application user (10)", er.Error())
		}
	}

	if len(user.SocialConnectionsIDs) > 0 {
		existingSocialIDsKeys := []string{}
		applicationSocialKey := ""
		for socialPlatform := range user.SocialConnectionsIDs {
			for idx := range user.SocialConnectionsIDs[socialPlatform] {
				applicationSocialKey = storageClient.SocialConnection(
					user.AccountID,
					user.ApplicationID,
					socialPlatform,
					utils.Base64Encode(user.SocialConnectionsIDs[socialPlatform][idx]),
				)
				existingSocialIDsKeys = append(existingSocialIDsKeys, applicationSocialKey)
			}
		}

		if applicationSocialKey != "" {
			existingSocialIDs, er := redisEngine.MGet(existingSocialIDsKeys...).Result()
			if er != nil {
				return nil, tgerrors.NewInternalError("failed to write the application user (11)", er.Error())
			}
			if len(existingSocialIDs) > 0 {
				user.Connections, err = autoConnectSocialFriends(user, existingSocialIDs)
				if err != nil {
					return
				}
			}
		}
	}

	listKey := storageClient.Users(user.AccountID, user.ApplicationID)
	if er = redisEngine.LPush(listKey, key).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (12)", er.Error())
	}

	if !retrieve {
		return user, err
	}

	return ReadApplicationUser(user.AccountID, user.ApplicationID, user.ID)
}

// CreateApplicationUserSession handles the creation of a user session and returns the session token
func CreateApplicationUserSession(user *entity.User) (string, *tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)
	token := storageClient.GenerateApplicationSessionID(user)

	if er := redisEngine.Set(sessionKey, token).Err(); er != nil {
		return "", tgerrors.NewInternalError("failed to create the application user session (1)", er.Error())
	}

	expired, er := redisEngine.Expire(sessionKey, storageClient.SessionTimeoutDuration()).Result()
	if er != nil {
		return "", tgerrors.NewInternalError("failed to create the application user session (2)", er.Error())
	}
	if !expired {
		return "", tgerrors.NewInternalError("failed to create the application user session (3)", "failed to set the expired")
	}

	return token, nil
}

// RefreshApplicationUserSession generates a new session token for the user session
func RefreshApplicationUserSession(sessionToken string, user *entity.User) (string, *tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)

	storedToken, er := redisEngine.Get(sessionKey).Result()
	if er != nil {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (1)", er.Error())
	}

	if storedToken != sessionToken {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (2)", "session token mismatch")
	}

	token := storageClient.GenerateApplicationSessionID(user)

	if er := redisEngine.Set(sessionKey, token).Err(); er != nil {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (3)", er.Error())
	}

	expired, er := redisEngine.Expire(sessionKey, storageClient.SessionTimeoutDuration()).Result()
	if er != nil {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (4)", er.Error())
	}
	if !expired {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (5)", "failed to set expired")
	}

	return token, nil
}

// GetApplicationUserSession returns the application user session
func GetApplicationUserSession(user *entity.User) (string, error) {
	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)
	storedSessionToken, err := redisEngine.Get(sessionKey).Result()
	if err != nil {
		return "", fmt.Errorf("could not fetch session from storage")
	}

	if storedSessionToken == "" {
		return "", fmt.Errorf("session not found")
	}

	return storedSessionToken, nil
}

// DestroyApplicationUserSession removes the user session
func DestroyApplicationUserSession(sessionToken string, user *entity.User) *tgerrors.TGError {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	sessionKey := storageClient.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)

	storedToken, er := redisEngine.Get(sessionKey).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to destroy the application user session (1)", er.Error())
	}

	if storedToken != sessionToken {
		return tgerrors.NewInternalError("failed to destroy the application user session (2)", "session token mismatch")
	}

	result, er := redisEngine.Del(sessionKey).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to destroy the application user session (3)", er.Error())
	}

	if result != 1 {
		return tgerrors.NewInternalError("failed to destroy the application user session (4)", er.Error())
	}

	return nil
}

// ApplicationUserByEmailExists checks if an application user exists by searching it via the email
func ApplicationUserByEmailExists(accountID, applicationID int64, email string) (bool, error) {
	emailListKey := storageClient.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))

	return redisEngine.Exists(emailListKey).Result()
}

// FindApplicationUserByEmail returns an application user by its email
func FindApplicationUserByEmail(accountID, applicationID int64, email string) (*entity.User, *tgerrors.TGError) {
	emailListKey := storageClient.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))

	return findApplicationUserByKey(accountID, applicationID, emailListKey)
}

// FindApplicationUserByUsername returns an application user by its username
func FindApplicationUserByUsername(accountID, applicationID int64, username string) (*entity.User, *tgerrors.TGError) {
	usernameListKey := storageClient.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(username))

	return findApplicationUserByKey(accountID, applicationID, usernameListKey)
}

// findApplicationUserByKey returns an application user regardless of the key used to search for him
func findApplicationUserByKey(accountID, applicationID int64, bucketName string) (*entity.User, *tgerrors.TGError) {
	storedValue, er := redisEngine.Get(bucketName).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to retrieve the application user (1)", er.Error())
	}

	userID, er := strconv.ParseInt(storedValue, 10, 64)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to retrieve the application user (2)", er.Error())
	}

	applicationUser, err := ReadApplicationUser(accountID, applicationID, userID)
	if err != nil {
		return nil, err
	}

	return applicationUser, nil
}
