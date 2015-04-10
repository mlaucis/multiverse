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
	"github.com/tapglue/backend/v01/entity"
)

// ReadAccountUser returns the account matching the ID or an error
func ReadAccountUser(accountID, accountUserID int64) (accountUser *entity.AccountUser, er tgerrors.TGError) {
	result, err := storageEngine.Get(storageClient.AccountUser(accountID, accountUserID)).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user (1)", err.Error())
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &accountUser); err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user (2)", err.Error())
	}

	return
}

// UpdateAccountUser update an account user in the database and returns the updated account user or an error
func UpdateAccountUser(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError) {
	updatedAccountUser.UpdatedAt = time.Now()

	if updatedAccountUser.Password == "" {
		updatedAccountUser.Password = existingAccountUser.Password
	} else if updatedAccountUser.Password != existingAccountUser.Password {
		updatedAccountUser.Password = storageClient.EncryptPassword(updatedAccountUser.Password)
	}

	val, err := json.Marshal(updatedAccountUser)
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to update the account user (1)", err.Error())
	}

	key := storageClient.AccountUser(updatedAccountUser.AccountID, updatedAccountUser.ID)
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, tgerrors.NewInternalError("failed to update the account user (2)", err.Error())
	}

	emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(existingAccountUser.Email))
	usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(existingAccountUser.Username))
	_, err = storageEngine.Del(emailListKey, usernameListKey).Result()
	// TODO handle this, maybe?

	if !updatedAccountUser.Enabled {
		listKey := storageClient.AccountUsers(updatedAccountUser.AccountID)
		if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
			return nil, tgerrors.NewInternalError("failed to update the account user (3)", err.Error())
		}
	} else {
		emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(updatedAccountUser.Email))
		err = storageEngine.HMSet(
			emailListKey,
			"acc", fmt.Sprintf("%d", updatedAccountUser.AccountID),
			"usr", fmt.Sprintf("%d", updatedAccountUser.ID),
		).Err()

		if err != nil {
			return nil, tgerrors.NewInternalError("failed to update the account user (4)", err.Error())
		}

		usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(updatedAccountUser.Username))
		err = storageEngine.HMSet(
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

	return ReadAccountUser(updatedAccountUser.AccountID, updatedAccountUser.ID)
}

// DeleteAccountUser deletes the account user matching the IDs or an error
func DeleteAccountUser(accountID, userID int64) tgerrors.TGError {
	// TODO: Make not deletable if its the only account user of an account
	accountUser, er := ReadAccountUser(accountID, userID)
	if er != nil {
		return er
	}

	key := storageClient.AccountUser(accountID, userID)
	result, err := storageEngine.Del(key).Result()
	if err != nil {
		return tgerrors.NewInternalError("failed to delete the account user (1)", err.Error())
	}

	if result != 1 {
		return tgerrors.NewNotFoundError("failed to delete the account user (2)", "account user not found")
	}

	listKey := storageClient.AccountUsers(accountID)
	if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
		return tgerrors.NewInternalError("failed to delete the account user (3)", err.Error())
	}

	emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(accountUser.Email))
	usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(accountUser.Username))
	_, err = storageEngine.Del(emailListKey, usernameListKey).Result()
	if err == nil {
		return nil
	}

	return tgerrors.NewInternalError("failed to delete the account user (4)", err.Error())
}

// ReadAccountUserList returns all the users from a certain account
func ReadAccountUserList(accountID int64) (accountUsers []*entity.AccountUser, er tgerrors.TGError) {
	result, err := storageEngine.LRange(storageClient.AccountUsers(accountID), 0, -1).Result()
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to read the account user list (1)", err.Error())
	}

	resultList, err := storageEngine.MGet(result...).Result()
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

// WriteAccountUser adds a new account user to the database and returns the created account user or an error
func WriteAccountUser(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError) {
	var err error
	if accountUser.ID, err = storageClient.GenerateAccountUserID(accountUser.AccountID); err != nil {
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
	accountUser.Password = storageClient.EncryptPassword(accountUser.Password)

	val, err := json.Marshal(accountUser)
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (3)", err.Error())
	}

	key := storageClient.AccountUser(accountUser.AccountID, accountUser.ID)
	exist, err := storageEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, tgerrors.NewInternalError("failed to create the account user (4)", "account user missing")
	}
	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (5)", err.Error())
	}

	idListKey := storageClient.AccountUsers(accountUser.AccountID)
	if err = storageEngine.LPush(idListKey, key).Err(); err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (6)", err.Error())
	}

	emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(accountUser.Email))
	err = storageEngine.HMSet(
		emailListKey,
		"acc", fmt.Sprintf("%d", accountUser.AccountID),
		"usr", fmt.Sprintf("%d", accountUser.ID),
	).Err()

	if err != nil {
		return nil, tgerrors.NewInternalError("failed to create the account user (7)", err.Error())
	}

	usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(accountUser.Username))
	err = storageEngine.HMSet(
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

	return ReadAccountUser(accountUser.AccountID, accountUser.ID)
}

// CreateAccountUserSession handles the creation of a user session and returns the session token
func CreateAccountUserSession(user *entity.AccountUser) (string, tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageClient.AccountSessionKey(user.AccountID, user.ID)
	token := storageClient.GenerateAccountSessionID(user)

	_, err := storageEngine.Set(sessionKey, token).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to create the account user session (1)", err.Error())
	}

	expired, err := storageEngine.Expire(sessionKey, storageClient.SessionTimeoutDuration()).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to create the account user session (2)", err.Error())
	}
	if !expired {
		return "", tgerrors.NewInternalError("failed to create the account user session (3)", "failed to set expired stuff")
	}

	return token, nil
}

// RefreshAccountUserSession generates a new session token for the user session
func RefreshAccountUserSession(sessionToken string, user *entity.AccountUser) (string, tgerrors.TGError) {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	// TODO rate limit this to be at least x minutes after the logout
	// TODO do we customize the key session timeout per app

	sessionKey := storageClient.AccountSessionKey(user.AccountID, user.ID)

	storedToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (1)", err.Error())
	}

	if storedToken != sessionToken {
		return "", tgerrors.NewInternalError("failed to refresh session token (2)\nsession token mismatch", err.Error())
	}

	token := storageClient.GenerateAccountSessionID(user)

	if err := storageEngine.Set(sessionKey, token).Err(); err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (3)", err.Error())
	}

	expired, err := storageEngine.Expire(sessionKey, storageClient.SessionTimeoutDuration()).Result()
	if err != nil {
		return "", tgerrors.NewInternalError("failed to refresh session token (4)", err.Error())
	}
	if !expired {
		tgerrors.NewInternalError("failed to refresh session token (5)", "could not set expire time")

	}

	return token, nil
}

// DestroyAccountUserSession removes the user session
func DestroyAccountUserSession(sessionToken string, user *entity.AccountUser) tgerrors.TGError {
	// TODO support multiple sessions?
	// TODO rate limit this to x / per day?
	sessionKey := storageClient.AccountSessionKey(user.AccountID, user.ID)

	storedToken, err := storageEngine.Get(sessionKey).Result()
	if err != nil {
		return tgerrors.NewInternalError("failed to destroy the session token (1)", err.Error())
	}

	if storedToken != sessionToken {
		return tgerrors.NewInternalError("failed to destroy the session token (2)", "session token mismatch")
	}

	result, err := storageEngine.Del(sessionKey).Result()
	if err != nil {
		return tgerrors.NewInternalError("failed to destroy the session token (3)", err.Error())
	}

	if result != 1 {
		return tgerrors.NewInternalError("failed to destroy the session token (4)", "invalid session")
	}

	return nil
}

// FindAccountAndUserByEmail returns the account and account user for a certain e-mail address
func FindAccountAndUserByEmail(email string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	emailListKey := storageClient.AccountUserByEmail(utils.Base64Encode(email))

	return findAccountByKey(emailListKey)
}

// FindAccountAndUserByUsername returns the account and account user for a certain username
func FindAccountAndUserByUsername(username string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	usernameListKey := storageClient.AccountUserByUsername(utils.Base64Encode(username))

	return findAccountByKey(usernameListKey)
}

// findAccountByKey retrieves an account and accountUser that are stored by their key, regardless of the specified key
func findAccountByKey(bucketName string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {

	details, err := storageEngine.HMGet(bucketName, "acc", "usr").Result()
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

	account, er := ReadAccount(accountID)
	if err != nil {
		return nil, nil, er
	}

	accountUser, er := ReadAccountUser(accountID, userID)
	if err != nil {
		return nil, nil, er
	}

	return account, accountUser, nil
}
