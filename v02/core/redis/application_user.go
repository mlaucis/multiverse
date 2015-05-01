package redis

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/redis"

	red "gopkg.in/redis.v2"
)

type (
	applicationUser struct {
		c       core.Connection
		storage redis.Client
		redis   *red.Client
	}
)

func (appu *applicationUser) Create(accountID, applicationID int64, user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err errors.Error) {
	// TODO We should introduce an option for the application to either allow for activated/deactivated behavior
	// and if they chose it, then we need to provide an endpoint to activate a user or not
	//user.Activated = true

	var er error
	user.Enabled = true
	user.CreatedAt = time.Now()
	user.UpdatedAt = user.CreatedAt
	user.LastLogin, er = time.Parse(time.RFC3339, "0000-01-01T00:00:00Z")
	if er != nil {
		return nil, errors.NewInternalError("failed to write the application user (1)", er.Error())
	}

	if user.ID, er = appu.storage.GenerateApplicationUserID(applicationID); er != nil {
		return nil, errors.NewInternalError("failed to write the application user (2)", er.Error())
	}

	// Encrypt password
	user.Password = storageHelper.EncryptPassword(user.Password)

	val, er := json.Marshal(user)
	if er != nil {
		return nil, errors.NewInternalError("failed to write the application user (3)", er.Error())
	}

	key := storageHelper.User(accountID, applicationID, user.ID)

	exist, er := appu.redis.SetNX(key, string(val)).Result()
	if !exist {
		return nil, errors.NewInternalError("failed to write the application user (4)", "duplicate user")
	}
	if err != nil {
		return nil, errors.NewInternalError("failed to write the application user (5)", er.Error())
	}

	stringUserID := fmt.Sprintf("%d", user.ID)

	emailListKey := storageHelper.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(user.Email))
	result, er := appu.redis.SetNX(emailListKey, stringUserID).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to write the application user (6)", er.Error())
	}
	if !result {
		return nil, errors.NewInternalError("failed to write the application user (7)", "duplicate user by e-mail")
	}

	usernameListKey := storageHelper.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(user.Username))
	result, er = appu.redis.SetNX(usernameListKey, stringUserID).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to write the application user (8)", er.Error())
	}
	if !result {
		return nil, errors.NewInternalError("failed to write the application user (9)", "duplicate user by username")
	}

	socialValues := []string{}
	applicationSocialKey := ""
	for idx := range user.SocialIDs {
		applicationSocialKey = storageHelper.SocialConnection(
			accountID,
			applicationID,
			idx,
			utils.Base64Encode(user.SocialIDs[idx]),
		)
		socialValues = append(socialValues, applicationSocialKey, stringUserID)
	}

	if applicationSocialKey != "" {
		er := appu.redis.MSet(socialValues...).Err()
		if er != nil {
			return nil, errors.NewInternalError("failed to write the application user (10)", er.Error())
		}
	}

	if len(user.SocialConnectionsIDs) > 0 {
		existingSocialIDsKeys := []string{}
		applicationSocialKey := ""
		for socialPlatform := range user.SocialConnectionsIDs {
			for idx := range user.SocialConnectionsIDs[socialPlatform] {
				applicationSocialKey = storageHelper.SocialConnection(
					accountID,
					applicationID,
					socialPlatform,
					utils.Base64Encode(user.SocialConnectionsIDs[socialPlatform][idx]),
				)
				existingSocialIDsKeys = append(existingSocialIDsKeys, applicationSocialKey)
			}
		}

		if applicationSocialKey != "" {
			existingSocialIDs, er := appu.redis.MGet(existingSocialIDsKeys...).Result()
			if er != nil {
				return nil, errors.NewInternalError("failed to write the application user (11)", er.Error())
			}
			if len(existingSocialIDs) > 0 {
				// TODO FIX THIS
				//user.Connections, err = appu.c.AutoConnectSocialFriends(user, existingSocialIDs)
				if err != nil {
					return
				}
			}
		}
	}

	listKey := storageHelper.Users(accountID, applicationID)
	if er = appu.redis.LPush(listKey, key).Err(); er != nil {
		return nil, errors.NewInternalError("failed to write the application user (12)", er.Error())
	}

	if !retrieve {
		return user, err
	}

	return appu.Read(accountID, applicationID, user.ID)
}

func (appu *applicationUser) Read(accountID, applicationID int64, userID string) (user *entity.ApplicationUser, err errors.Error) {
	key := storageHelper.User(accountID, applicationID, userID)

	result, er := appu.redis.Get(key).Result()
	if err != nil {
		return nil, errors.NewInternalError("failed to read application user (1)", er.Error())
	}

	if er = json.Unmarshal([]byte(result), &user); er != nil {
		return nil, errors.NewInternalError("failed to read application user (2)", er.Error())
	}

	return
}

func (appu *applicationUser) Update(accountID, applicationID int64, existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err errors.Error) {

	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	} else if updatedUser.Password != existingUser.Password {
		// Encrypt password - we should do this only if the password changes
		updatedUser.Password = storageHelper.EncryptPassword(updatedUser.Password)
	}

	val, er := json.Marshal(updatedUser)
	if er != nil {
		return nil, errors.NewInternalError("failed to update the application user (1)", er.Error())
	}

	key := storageHelper.User(accountID, applicationID, updatedUser.ID)
	if er = appu.redis.Set(key, string(val)).Err(); er != nil {
		return nil, errors.NewInternalError("failed to update the application user (2)", er.Error())
	}

	if existingUser.Email != updatedUser.Email {
		emailListKey := storageHelper.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(existingUser.Email))
		_, er = appu.redis.Del(emailListKey).Result()

		emailListKey = storageHelper.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(updatedUser.Email))
		er = appu.redis.Set(emailListKey, fmt.Sprintf("%d", updatedUser.ID)).Err()
		if er != nil {
			return nil, errors.NewInternalError("failed to update the application user (3)", er.Error())
		}
	}

	if existingUser.Username != updatedUser.Username {
		usernameListKey := storageHelper.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(existingUser.Username))
		_, er = appu.redis.Del(usernameListKey).Result()

		usernameListKey = storageHelper.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(updatedUser.Username))
		er = appu.redis.Set(usernameListKey, fmt.Sprintf("%d", updatedUser.ID)).Err()

		if er != nil {
			return nil, errors.NewInternalError("failed to update the application user (4)", er.Error())
		}
	}

	if !updatedUser.Enabled {
		listKey := storageHelper.Users(accountID, applicationID)
		if er = appu.redis.LRem(listKey, 0, key).Err(); er != nil {
			return nil, errors.NewInternalError("failed to update the application user (5)", er.Error())
		}
	} else {
		listKey := storageHelper.Users(accountID, applicationID)
		if er = appu.redis.LPush(listKey, key).Err(); er != nil {
			return nil, errors.NewInternalError("failed to update the application user (6)", er.Error())
		}
	}

	if !retrieve {
		return &updatedUser, nil
	}

	return appu.Read(accountID, applicationID, updatedUser.ID)
}

func (appu *applicationUser) Delete(accountID, applicationID int64, applicationUser *entity.ApplicationUser) (err errors.Error) {
	user, err := appu.Read(accountID, applicationID, applicationUser.ID)
	if err != nil {
		return err
	}

	disabledUser := *user
	disabledUser.Enabled = false
	disabledUser.Password = ""
	_, err = appu.Update(accountID, applicationID, *user, disabledUser, false)

	// TODO: Remove Users Connections?
	// TODO: Remove Users Connection Lists?
	// TODO: Remove User in other Users Connection Lists?
	// TODO: Remove Users Events?
	// TODO: Remove Users Events from Lists?

	return

	// TODO Figure out if we should just simply remove the user or not

	/*key := appu.storage.User(accountID, applicationID, userID)
	result, err := storageEngine.Del(key).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		return fmt.Errorf("The resource for the provided id doesn't exist")
	}

	listKey := appu.storage.Users(accountID, applicationID)
	if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
		return err
	}

	emailListKey := appu.storage.AccountUserByEmail(Base64Encode(user.Email))
	usernameListKey := appu.storage.AccountUserByUsername(Base64Encode(user.Username))
	_, err = storageEngine.Del(emailListKey, usernameListKey).Result()

	return nil
	*/
}

func (appu *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err errors.Error) {
	key := storageHelper.Users(accountID, applicationID)

	result, er := appu.redis.LRange(key, 0, -1).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the application user list (1)", er.Error())
	}

	if len(result) == 0 {
		return users, nil
	}

	resultList, er := appu.redis.MGet(result...).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the application user list (2)", er.Error())
	}

	user := &entity.ApplicationUser{}
	for _, result := range resultList {
		if er = json.Unmarshal([]byte(result.(string)), user); er != nil {
			return nil, errors.NewInternalError("failed to read the application user list (3)", er.Error())
		}
		users = append(users, user)
		user = &entity.ApplicationUser{}
	}

	return
}

func (appu *applicationUser) CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, errors.Error) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageHelper.ApplicationSessionKey(accountID, applicationID, user.ID)
	token := storageHelper.GenerateApplicationSessionID(user)

	if er := appu.redis.Set(sessionKey, token).Err(); er != nil {
		return "", errors.NewInternalError("failed to create the application user session (1)", er.Error())
	}

	expired, er := appu.redis.Expire(sessionKey, storageHelper.SessionTimeoutDuration()).Result()
	if er != nil {
		return "", errors.NewInternalError("failed to create the application user session (2)", er.Error())
	}
	if !expired {
		return "", errors.NewInternalError("failed to create the application user session (3)", "failed to set the expired")
	}

	return token, nil
}

func (appu *applicationUser) RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, errors.Error) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageHelper.ApplicationSessionKey(accountID, applicationID, user.ID)

	storedToken, er := appu.redis.Get(sessionKey).Result()
	if er != nil {
		return "", errors.NewInternalError("failed to refresh the application user session (1)", er.Error())
	}

	if storedToken != sessionToken {
		return "", errors.NewInternalError("failed to refresh the application user session (2)", "session token mismatch")
	}

	token := storageHelper.GenerateApplicationSessionID(user)

	if er := appu.redis.Set(sessionKey, token).Err(); er != nil {
		return "", errors.NewInternalError("failed to refresh the application user session (3)", er.Error())
	}

	expired, er := appu.redis.Expire(sessionKey, storageHelper.SessionTimeoutDuration()).Result()
	if er != nil {
		return "", errors.NewInternalError("failed to refresh the application user session (4)", er.Error())
	}
	if !expired {
		return "", errors.NewInternalError("failed to refresh the application user session (5)", "failed to set expired")
	}

	return token, nil
}

func (appu *applicationUser) GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, errors.Error) {
	sessionKey := storageHelper.ApplicationSessionKey(accountID, applicationID, user.ID)
	storedSessionToken, err := appu.redis.Get(sessionKey).Result()
	if err != nil {
		return "", errors.NewInternalError("error while fetching session", "could not fetch session from storage")
	}

	if storedSessionToken == "" {
		return "", errors.NewInternalError("session not found", "session not found")
	}

	return storedSessionToken, nil
}

func (appu *applicationUser) DestroySession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) errors.Error {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	sessionKey := storageHelper.ApplicationSessionKey(accountID, applicationID, user.ID)

	storedToken, er := appu.redis.Get(sessionKey).Result()
	if er != nil {
		return errors.NewInternalError("failed to destroy the application user session (1)", er.Error())
	}

	if storedToken != sessionToken {
		return errors.NewInternalError("failed to destroy the application user session (2)", "session token mismatch")
	}

	result, er := appu.redis.Del(sessionKey).Result()
	if er != nil {
		return errors.NewInternalError("failed to destroy the application user session (3)", er.Error())
	}

	if result != 1 {
		return errors.NewInternalError("failed to destroy the application user session (4)", er.Error())
	}

	return nil
}

func (appu *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, errors.Error) {
	emailListKey := storageHelper.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))

	return appu.findApplicationUserByKey(accountID, applicationID, emailListKey)
}

func (appu *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, errors.Error) {
	emailListKey := storageHelper.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))
	return appu.existsByKey(emailListKey)
}

func (appu *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, errors.Error) {
	usernameListKey := storageHelper.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(username))

	return appu.findApplicationUserByKey(accountID, applicationID, usernameListKey)
}

func (appu *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, errors.Error) {
	usernameListKey := storageHelper.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(username))
	return appu.existsByKey(usernameListKey)
}

func (appu *applicationUser) ExistsByID(accountID, applicationID int64, userID string) (bool, errors.Error) {
	user, err := appu.Read(accountID, applicationID, userID)
	if err != nil {
		return false, err
	}

	return user.Enabled, nil
}

// findApplicationUserByKey returns an application user regardless of the key used to search for him
func (appu *applicationUser) findApplicationUserByKey(accountID, applicationID int64, bucketName string) (*entity.ApplicationUser, errors.Error) {
	storedValue, er := appu.redis.Get(bucketName).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to retrieve the application user (1)", er.Error())
	}

	applicationUser, err := appu.Read(accountID, applicationID, storedValue)
	if err != nil {
		return nil, err
	}

	return applicationUser, nil
}

func (appu *applicationUser) existsByKey(bucketName string) (bool, errors.Error) {
	exists, err := appu.redis.Exists(bucketName).Result()
	if err != nil {
		return false, errors.NewInternalError("unexpected errror", err.Error())
	}

	return exists, nil
}

func (appu *applicationUser) FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewApplicationUser creates a new Event
func NewApplicationUser(storageClient redis.Client) core.ApplicationUser {
	return &applicationUser{
		c:       NewConnection(storageClient),
		storage: storageClient,
		redis:   storageClient.Datastore(),
	}
}

// NewApplicationUserWithConnection creates a new Event
func NewApplicationUserWithConnection(storageClient redis.Client, c core.Connection) core.ApplicationUser {
	return &applicationUser{
		c:       c,
		storage: storageClient,
		redis:   storageClient.Datastore(),
	}
}
