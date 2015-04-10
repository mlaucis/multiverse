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
	accountUser struct {
		a       core.Account
		storage *storage.Client
		redis   *red.Client
	}
)

// Create adds a new account user to the database and returns the created account user or an error
func (au *accountUser) Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError) {
	var err error
	if accountUser.ID, err = au.storage.GenerateAccountUserID(accountUser.AccountID); err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (1)", err.Error())
	}

	accountUser.Enabled = true
	accountUser.CreatedAt = time.Now()
	accountUser.UpdatedAt = accountUser.CreatedAt
	accountUser.LastLogin, err = time.Parse(time.RFC3339, "0000-01-01T00:00:00Z")
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (2)", err.Error())
	}

	// Encrypt password
	accountUser.Password = au.storage.EncryptPassword(accountUser.Password)

	val, err := json.Marshal(accountUser)
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (3)", err.Error())
	}

	key := au.storage.AccountUser(accountUser.AccountID, accountUser.ID)
	exist, err := au.redis.SetNX(key, string(val)).Result()
	if !exist {
		return nil, tgerrors.NewInternalError("failed to create the account user (4)", "account user missing")
	}
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (5)", err.Error())
	}

	idListKey := au.storage.AccountUsers(accountUser.AccountID)
	if err = au.redis.LPush(idListKey, key).Err(); err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (6)", err.Error())
	}

	emailListKey := au.storage.AccountUserByEmail(utils.Base64Encode(accountUser.Email))
	err = au.redis.HMSet(
		emailListKey,
		"acc", fmt.Sprintf("%d", accountUser.AccountID),
		"usr", fmt.Sprintf("%d", accountUser.ID),
	).Err()

	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (7)", err.Error())
	}

	usernameListKey := au.storage.AccountUserByUsername(utils.Base64Encode(accountUser.Username))
	err = au.redis.HMSet(
		usernameListKey,
		"acc", fmt.Sprintf("%d", accountUser.AccountID),
		"usr", fmt.Sprintf("%d", accountUser.ID),
	).Err()

	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (8)", err.Error())
	}

	if !retrieve {
		return accountUser, nil
	}

	return au.Read(accountUser.AccountID, accountUser.ID)
}

// Read returns the account matching the ID or an error
func (au *accountUser) Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er tgerrors.TGError) {
	result, err := au.redis.Get(au.storage.AccountUser(accountID, accountUserID)).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user (1)", err.Error())
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &accountUser); err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user (2)", err.Error())
	}

	return
}

// Update update an account user in the database and returns the updated account user or an error
func (au *accountUser) Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError) {
	updatedAccountUser.UpdatedAt = time.Now()

	if updatedAccountUser.Password == "" {
		updatedAccountUser.Password = existingAccountUser.Password
	} else if updatedAccountUser.Password != existingAccountUser.Password {
		updatedAccountUser.Password = au.storage.EncryptPassword(updatedAccountUser.Password)
	}

	val, err := json.Marshal(updatedAccountUser)
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to update the account user (1)", err.Error())
	}

	key := au.storage.AccountUser(updatedAccountUser.AccountID, updatedAccountUser.ID)
	if err = au.redis.Set(key, string(val)).Err(); err != nil {
		return nil, tgerrors.NewInternalError("failed to update the account user (2)", err.Error())
	}

	emailListKey := au.storage.AccountUserByEmail(utils.Base64Encode(existingAccountUser.Email))
	usernameListKey := au.storage.AccountUserByUsername(utils.Base64Encode(existingAccountUser.Username))
	_, err = au.redis.Del(emailListKey, usernameListKey).Result()
	// TODO handle this, maybe?

	if !updatedAccountUser.Enabled {
		listKey := au.storage.AccountUsers(updatedAccountUser.AccountID)
		if err = au.redis.LRem(listKey, 0, key).Err(); err != nil {
			return nil, tgerrors.NewInternalError("failed to update the account user (3)", err.Error())
		}
	} else {
		emailListKey := au.storage.AccountUserByEmail(utils.Base64Encode(updatedAccountUser.Email))
		err = au.redis.HMSet(
			emailListKey,
			"acc", fmt.Sprintf("%d", updatedAccountUser.AccountID),
			"usr", fmt.Sprintf("%d", updatedAccountUser.ID),
		).Err()

		if err != nil {
			return nil, tgerrors.NewInternalError("failed to update the account user (4)", err.Error())
		}

		usernameListKey := au.storage.AccountUserByUsername(utils.Base64Encode(updatedAccountUser.Username))
		err = au.redis.HMSet(
			usernameListKey,
			"acc", fmt.Sprintf("%d", updatedAccountUser.AccountID),
			"usr", fmt.Sprintf("%d", updatedAccountUser.ID),
		).Err()

		if err != nil {
			return nil, tgerrors.NewInternalError("failed to update the account user (5)", err.Error())
		}
	}

	if !retrieve {
		return &updatedAccountUser, nil
	}

	return au.Read(updatedAccountUser.AccountID, updatedAccountUser.ID)
}

// Delete deletes the account user matching the IDs or an error
func (au *accountUser) Delete(accountID, userID int64) tgerrors.TGError {
	// TODO: Make not deletable if its the only account user of an account
	accountUser, er := au.Read(accountID, userID)
	if er != nil {
		return er
	}

	key := au.storage.AccountUser(accountID, userID)
	result, err := au.redis.Del(key).Result()
	if err != nil {
		return tgerrors.NewInternalError("failed to delete the account user (1)", err.Error())
	}

	if result != 1 {
		return tgerrors.NewNotFoundError("failed to delete the account user (2)", "account user not found")
	}

	listKey := au.storage.AccountUsers(accountID)
	if err = au.redis.LRem(listKey, 0, key).Err(); err != nil {
		return tgerrors.NewInternalError("failed to delete the account user (3)", err.Error())
	}

	emailListKey := au.storage.AccountUserByEmail(utils.Base64Encode(accountUser.Email))
	usernameListKey := au.storage.AccountUserByUsername(utils.Base64Encode(accountUser.Username))
	_, err = au.redis.Del(emailListKey, usernameListKey).Result()
	if err == nil {
		return nil
	}

	return tgerrors.NewInternalError("failed to delete the account user (4)", err.Error())
}

// List returns all the users from a certain account
func (au *accountUser) List(accountID int64) (accountUsers []*entity.AccountUser, er tgerrors.TGError) {
	result, err := au.redis.LRange(au.storage.AccountUsers(accountID), 0, -1).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user list (1)", err.Error())
	}

	resultList, err := au.redis.MGet(result...).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user list (2)", err.Error())
	}

	accountUser := &entity.AccountUser{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), accountUser); err != nil {
			return nil, tgerrors.NewInternalError("failed to read the account user list (3)", err.Error())
		}
		accountUsers = append(accountUsers, accountUser)
		accountUser = &entity.AccountUser{}
	}

	return
}

// CreateSession handles the creation of a user session and returns the session token
func (au *accountUser) CreateSession(user *entity.AccountUser) (string, tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := au.storage.AccountSessionKey(user.AccountID, user.ID)
	token := au.storage.GenerateAccountSessionID(user)

	_, err := au.redis.Set(sessionKey, token).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to create the account user session (1)", err.Error())
	}

	expired, err := au.redis.Expire(sessionKey, au.storage.SessionTimeoutDuration()).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to create the account user session (2)", err.Error())
	}
	if !expired {
		return "", tgerrors.NewInternalError("failed to create the account user session (3)", "failed to set expired stuff")
	}

	return token, nil
}

// RefreshSession generates a new session token for the user session
func (au *accountUser) RefreshSession(sessionToken string, user *entity.AccountUser) (string, tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := au.storage.AccountSessionKey(user.AccountID, user.ID)
	storedToken, err := au.redis.Get(sessionKey).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (1)", err.Error())
	}

	if storedToken != sessionToken {
		return "", tgerrors.NewInternalError("failed to refresh session token (2)\nsession token mismatch", err.Error())
	}

	token := au.storage.GenerateAccountSessionID(user)

	if err := au.redis.Set(sessionKey, token).Err(); err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (3)", err.Error())
	}

	expired, err := au.redis.Expire(sessionKey, au.storage.SessionTimeoutDuration()).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (4)", err.Error())
	}
	if !expired {
		tgerrors.NewInternalError("failed to refresh session token (5)", "could not set expire time")

	}

	return token, nil
}

// DestroySession removes the user session
func (au *accountUser) DestroySession(sessionToken string, user *entity.AccountUser) tgerrors.TGError {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	sessionKey := au.storage.AccountSessionKey(user.AccountID, user.ID)

	storedToken, err := au.redis.Get(sessionKey).Result()
	if err != nil {
		return tgerrors.NewInternalError("failed to destroy the session token (1)", err.Error())
	}

	if storedToken != sessionToken {
		return tgerrors.NewInternalError("failed to destroy the session token (2)", "session token mismatch")
	}

	result, err := au.redis.Del(sessionKey).Result()
	if err != nil {
		return tgerrors.NewInternalError("failed to destroy the session token (3)", err.Error())
	}

	if result != 1 {
		return tgerrors.NewInternalError("failed to destroy the session token (4)", "invalid session")
	}

	return nil
}

// GetSession retrieves the account user session token
func (au *accountUser) GetSession(user *entity.AccountUser) (string, tgerrors.TGError) {
	sessionKey := au.storage.AccountSessionKey(user.AccountID, user.ID)
	storedToken, err := au.redis.Get(sessionKey).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (1)", err.Error())
	}

	return storedToken, nil
}

// FindByEmail returns the account and account user for a certain e-mail address
func (au *accountUser) FindByEmail(email string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	emailListKey := au.storage.AccountUserByEmail(utils.Base64Encode(email))
	return au.findByKey(emailListKey)
}

func (au *accountUser) ExistsByEmail(email string) (bool, tgerrors.TGError) {
	emailListKey := au.storage.AccountUserByEmail(utils.Base64Encode(email))
	return au.existsByKey(emailListKey)
}

// FindByUsername returns the account and account user for a certain username
func (au *accountUser) FindByUsername(username string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	usernameListKey := au.storage.AccountUserByUsername(utils.Base64Encode(username))
	return au.findByKey(usernameListKey)
}

func (au *accountUser) ExistsByUsername(username string) (bool, tgerrors.TGError) {
	usernameListKey := au.storage.AccountUserByUsername(utils.Base64Encode(username))
	return au.existsByKey(usernameListKey)
}

// findByKey retrieves an account and accountUser that are stored by their key, regardless of the specified key
func (au *accountUser) findByKey(bucketName string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	details, err := au.redis.HMGet(bucketName, "acc", "usr").Result()
	if err != nil {
		return nil, nil, tgerrors.NewInternalError("failed to find the account user (1)", err.Error())
	}

	if len(details) != 2 || details[0] == nil || details[1] == nil {
		return nil, nil, tgerrors.NewInternalError("failed to find the account user (2)", "mismatching or nil parts")
	}

	accountID, err := strconv.ParseInt(details[0].(string), 10, 64)
	if err != nil {
		return nil, nil, tgerrors.NewInternalError("failed to find the account user (3)", err.Error())
	}

	userID, err := strconv.ParseInt(details[1].(string), 10, 64)
	if err != nil {
		return nil, nil, tgerrors.NewInternalError("failed to find the account user (4)", err.Error())
	}

	account, er := au.a.Read(accountID)
	if err != nil {
		return nil, nil, er
	}

	accountUser, er := au.Read(accountID, userID)
	if err != nil {
		return nil, nil, er
	}

	return account, accountUser, nil
}

// existsByKey checks if an AccountUser exists by a certain key
func (au *accountUser) existsByKey(bucketName string) (bool, tgerrors.TGError) {
	exists, err := au.redis.Exists(bucketName).Result()
	if err != nil {
		return false, tgerrors.NewInternalError("failed to find the account user (1)", err.Error())
	}

	return exists, nil
}

// NewAccountUser creates a new AccountUser
func NewAccountUser(storageClient *storage.Client, storageEngine *red.Client) core.AccountUser {
	return &accountUser{
		a:       NewAccount(storageClient, storageEngine),
		storage: storageClient,
		redis:   storageEngine,
	}
}

// NewAccountUserWithAccount creates a new AccountUser
func NewAccountUserWithAccount(storageClient *storage.Client, storageEngine *red.Client, a core.Account) core.AccountUser {
	return &accountUser{
		a:       a,
		storage: storageClient,
		redis:   storageEngine,
	}
}
