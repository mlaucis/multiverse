package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage"

	red "gopkg.in/redis.v2"
)

type (
	applicationUser struct {
		c       core.Connection
		storage *storage.Client
		redis   *red.Client
	}
)

// WriteUser adds a user to the database and returns the created user or an error
func (appu *applicationUser) Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError) {
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

	if user.ID, er = appu.storage.GenerateApplicationUserID(user.ApplicationID); er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (2)", er.Error())
	}

	// Encrypt password
	user.Password = appu.storage.EncryptPassword(user.Password)

	val, er := json.Marshal(user)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (3)", er.Error())
	}

	key := appu.storage.User(user.AccountID, user.ApplicationID, user.ID)

	exist, er := appu.redis.SetNX(key, string(val)).Result()
	if !exist {
		return nil, tgerrors.NewInternalError("failed to write the application user (4)", "duplicate user")
	}
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (5)", er.Error())
	}

	stringUserID := fmt.Sprintf("%d", user.ID)

	emailListKey := appu.storage.ApplicationUserByEmail(user.AccountID, user.ApplicationID, utils.Base64Encode(user.Email))
	result, er := appu.redis.SetNX(emailListKey, stringUserID).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (6)", er.Error())
	}
	if !result {
		return nil, tgerrors.NewInternalError("failed to write the application user (7)", "duplicate user by e-mail")
	}

	usernameListKey := appu.storage.ApplicationUserByUsername(user.AccountID, user.ApplicationID, utils.Base64Encode(user.Username))
	result, er = appu.redis.SetNX(usernameListKey, stringUserID).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (8)", er.Error())
	}
	if !result {
		return nil, tgerrors.NewInternalError("failed to write the application user (9)", "duplicate user by username")
	}

	socialValues := []string{}
	applicationSocialKey := ""
	for idx := range user.SocialIDs {
		applicationSocialKey = appu.storage.SocialConnection(
			user.AccountID,
			user.ApplicationID,
			idx,
			utils.Base64Encode(user.SocialIDs[idx]),
		)
		socialValues = append(socialValues, applicationSocialKey, stringUserID)
	}

	if applicationSocialKey != "" {
		er := appu.redis.MSet(socialValues...).Err()
		if er != nil {
			return nil, tgerrors.NewInternalError("failed to write the application user (10)", er.Error())
		}
	}

	if len(user.SocialConnectionsIDs) > 0 {
		existingSocialIDsKeys := []string{}
		applicationSocialKey := ""
		for socialPlatform := range user.SocialConnectionsIDs {
			for idx := range user.SocialConnectionsIDs[socialPlatform] {
				applicationSocialKey = appu.storage.SocialConnection(
					user.AccountID,
					user.ApplicationID,
					socialPlatform,
					utils.Base64Encode(user.SocialConnectionsIDs[socialPlatform][idx]),
				)
				existingSocialIDsKeys = append(existingSocialIDsKeys, applicationSocialKey)
			}
		}

		if applicationSocialKey != "" {
			existingSocialIDs, er := appu.redis.MGet(existingSocialIDsKeys...).Result()
			if er != nil {
				return nil, tgerrors.NewInternalError("failed to write the application user (11)", er.Error())
			}
			if len(existingSocialIDs) > 0 {
				user.Connections, err = appu.c.AutoConnectSocialFriends(user, existingSocialIDs)
				if err != nil {
					return
				}
			}
		}
	}

	listKey := appu.storage.Users(user.AccountID, user.ApplicationID)
	if er = appu.redis.LPush(listKey, key).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to write the application user (12)", er.Error())
	}

	if !retrieve {
		return user, err
	}

	return appu.Read(user.AccountID, user.ApplicationID, user.ID)
}

// ReadApplicationUser returns the user matching the ID or an error
func (appu *applicationUser) Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err tgerrors.TGError) {
	key := appu.storage.User(accountID, applicationID, userID)

	result, er := appu.redis.Get(key).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read application user (1)", er.Error())
	}

	if er = json.Unmarshal([]byte(result), &user); er != nil {
		return nil, tgerrors.NewInternalError("failed to read application user (2)", er.Error())
	}

	return
}

// UpdateUser updates a user in the database and returns the updates user or an error
func (appu *applicationUser) Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError) {

	if updatedUser.Password == "" {
		updatedUser.Password = existingUser.Password
	} else if updatedUser.Password != existingUser.Password {
		// Encrypt password - we should do this only if the password changes
		updatedUser.Password = appu.storage.EncryptPassword(updatedUser.Password)
	}

	val, er := json.Marshal(updatedUser)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application user (1)", er.Error())
	}

	key := appu.storage.User(updatedUser.AccountID, updatedUser.ApplicationID, updatedUser.ID)
	if er = appu.redis.Set(key, string(val)).Err(); er != nil {
		return nil, tgerrors.NewInternalError("failed to update the application user (2)", er.Error())
	}

	if existingUser.Email != updatedUser.Email {
		emailListKey := appu.storage.ApplicationUserByEmail(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(existingUser.Email))
		_, er = appu.redis.Del(emailListKey).Result()

		emailListKey = appu.storage.ApplicationUserByEmail(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(updatedUser.Email))
		er = appu.redis.Set(emailListKey, fmt.Sprintf("%d", updatedUser.ID)).Err()
		if er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (3)", er.Error())
		}
	}

	if existingUser.Username != updatedUser.Username {
		usernameListKey := appu.storage.ApplicationUserByUsername(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(existingUser.Username))
		_, er = appu.redis.Del(usernameListKey).Result()

		usernameListKey = appu.storage.ApplicationUserByUsername(existingUser.AccountID, existingUser.ApplicationID, utils.Base64Encode(updatedUser.Username))
		er = appu.redis.Set(usernameListKey, fmt.Sprintf("%d", updatedUser.ID)).Err()

		if er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (4)", er.Error())
		}
	}

	if !updatedUser.Enabled {
		listKey := appu.storage.Users(updatedUser.AccountID, updatedUser.ApplicationID)
		if er = appu.redis.LRem(listKey, 0, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (5)", er.Error())
		}
	} else {
		listKey := appu.storage.Users(updatedUser.AccountID, updatedUser.ApplicationID)
		if er = appu.redis.LPush(listKey, key).Err(); er != nil {
			return nil, tgerrors.NewInternalError("failed to update the application user (6)", er.Error())
		}
	}

	if !retrieve {
		return &updatedUser, nil
	}

	return appu.Read(updatedUser.AccountID, updatedUser.ApplicationID, updatedUser.ID)
}

// DeleteUser deletes the user matching the IDs or an error
func (appu *applicationUser) Delete(accountID, applicationID, userID int64) (err tgerrors.TGError) {
	user, err := appu.Read(accountID, applicationID, userID)
	if err != nil {
		return err
	}

	disabledUser := *user
	disabledUser.Enabled = false
	disabledUser.Password = ""
	_, err = appu.Update(*user, disabledUser, false)

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

// ReadUserList returns all users from a certain account
func (appu *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err tgerrors.TGError) {
	key := appu.storage.Users(accountID, applicationID)

	result, er := appu.redis.LRange(key, 0, -1).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the application user list (1)", er.Error())
	}

	if len(result) == 0 {
		return users, nil
	}

	resultList, er := appu.redis.MGet(result...).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to read the application user list (2)", er.Error())
	}

	user := &entity.ApplicationUser{}
	for _, result := range resultList {
		if er = json.Unmarshal([]byte(result.(string)), user); er != nil {
			return nil, tgerrors.NewInternalError("failed to read the application user list (3)", er.Error())
		}
		users = append(users, user)
		user = &entity.ApplicationUser{}
	}

	return
}

// CreateApplicationUserSession handles the creation of a user session and returns the session token
func (appu *applicationUser) CreateSession(user *entity.ApplicationUser) (string, tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := appu.storage.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)
	token := appu.storage.GenerateApplicationSessionID(user)

	if er := appu.redis.Set(sessionKey, token).Err(); er != nil {
		return "", tgerrors.NewInternalError("failed to create the application user session (1)", er.Error())
	}

	expired, er := appu.redis.Expire(sessionKey, appu.storage.SessionTimeoutDuration()).Result()
	if er != nil {
		return "", tgerrors.NewInternalError("failed to create the application user session (2)", er.Error())
	}
	if !expired {
		return "", tgerrors.NewInternalError("failed to create the application user session (3)", "failed to set the expired")
	}

	return token, nil
}

// RefreshApplicationUserSession generates a new session token for the user session
func (appu *applicationUser) RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := appu.storage.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)

	storedToken, er := appu.redis.Get(sessionKey).Result()
	if er != nil {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (1)", er.Error())
	}

	if storedToken != sessionToken {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (2)", "session token mismatch")
	}

	token := appu.storage.GenerateApplicationSessionID(user)

	if er := appu.redis.Set(sessionKey, token).Err(); er != nil {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (3)", er.Error())
	}

	expired, er := appu.redis.Expire(sessionKey, appu.storage.SessionTimeoutDuration()).Result()
	if er != nil {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (4)", er.Error())
	}
	if !expired {
		return "", tgerrors.NewInternalError("failed to refresh the application user session (5)", "failed to set expired")
	}

	return token, nil
}

// GetApplicationUserSession returns the application user session
func (appu *applicationUser) GetSession(user *entity.ApplicationUser) (string, tgerrors.TGError) {
	sessionKey := appu.storage.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)
	storedSessionToken, err := appu.redis.Get(sessionKey).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("error while fetching session", "could not fetch session from storage")
	}

	if storedSessionToken == "" {
		return "", tgerrors.NewInternalError("session not found", "session not found")
	}

	return storedSessionToken, nil
}

// DestroyApplicationUserSession removes the user session
func (appu *applicationUser) DestroySession(sessionToken string, user *entity.ApplicationUser) tgerrors.TGError {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	sessionKey := appu.storage.ApplicationSessionKey(user.AccountID, user.ApplicationID, user.ID)

	storedToken, er := appu.redis.Get(sessionKey).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to destroy the application user session (1)", er.Error())
	}

	if storedToken != sessionToken {
		return tgerrors.NewInternalError("failed to destroy the application user session (2)", "session token mismatch")
	}

	result, er := appu.redis.Del(sessionKey).Result()
	if er != nil {
		return tgerrors.NewInternalError("failed to destroy the application user session (3)", er.Error())
	}

	if result != 1 {
		return tgerrors.NewInternalError("failed to destroy the application user session (4)", er.Error())
	}

	return nil
}

// FindApplicationUserByEmail returns an application user by its email
func (appu *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, tgerrors.TGError) {
	emailListKey := appu.storage.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))

	return appu.findApplicationUserByKey(accountID, applicationID, emailListKey)
}

// ExistsByEmail checks if an application user exists by searching it via the email
func (appu *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, tgerrors.TGError) {
	emailListKey := appu.storage.ApplicationUserByEmail(accountID, applicationID, utils.Base64Encode(email))
	return appu.existsByKey(emailListKey)
}

// FindApplicationUserByUsername returns an application user by its username
func (appu *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, tgerrors.TGError) {
	usernameListKey := appu.storage.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(username))

	return appu.findApplicationUserByKey(accountID, applicationID, usernameListKey)
}

// ExistsByEmail checks if an application user exists by searching it via the email
func (appu *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, tgerrors.TGError) {
	usernameListKey := appu.storage.ApplicationUserByUsername(accountID, applicationID, utils.Base64Encode(username))
	return appu.existsByKey(usernameListKey)
}

// findApplicationUserByKey returns an application user regardless of the key used to search for him
func (appu *applicationUser) findApplicationUserByKey(accountID, applicationID int64, bucketName string) (*entity.ApplicationUser, tgerrors.TGError) {
	storedValue, er := appu.redis.Get(bucketName).Result()
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to retrieve the application user (1)", er.Error())
	}

	userID, er := strconv.ParseInt(storedValue, 10, 64)
	if er != nil {
		return nil, tgerrors.NewInternalError("failed to retrieve the application user (2)", er.Error())
	}

	applicationUser, err := appu.Read(accountID, applicationID, userID)
	if err != nil {
		return nil, err
	}

	return applicationUser, nil
}

func (appu *applicationUser) existsByKey(bucketName string) (bool, tgerrors.TGError) {
	exists, err := appu.redis.Exists(bucketName).Result()
	if err != nil {
		return false, tgerrors.NewInternalError("unexpected errror", err.Error())
	}

	return exists, nil
}

// NewApplicationUser creates a new Event
func NewApplicationUser(storageClient *storage.Client, storageEngine *red.Client) core.ApplicationUser {
	return &applicationUser{
		c:       NewConnection(storageClient, storageEngine),
		storage: storageClient,
		redis:   storageEngine,
	}
}

// NewApplicationUserWithConnection creates a new Event
func NewApplicationUserWithConnection(storageClient *storage.Client, storageEngine *red.Client, c core.Connection) core.ApplicationUser {
	return &applicationUser{
		c:       c,
		storage: storageClient,
		redis:   storageEngine,
	}
}
