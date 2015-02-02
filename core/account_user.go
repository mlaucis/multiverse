/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"time"

	"github.com/tapglue/backend/core/entity"
)

// ReadAccountUser returns the account matching the ID or an error
func ReadAccountUser(accountID, accountUserID int64) (accountUser *entity.AccountUser, err error) {
	result, err := storageEngine.Get(storageClient.AccountUserKey(accountID, accountUserID)).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &accountUser); err != nil {
		return nil, err
	}

	return
}

// ReadAccountUserList returns all the users from a certain account
func ReadAccountUserList(accountID int64) (accountUsers []*entity.AccountUser, err error) {
	result, err := storageEngine.LRange(storageClient.AccountKey(accountID), 0, -1).Result()
	if err != nil {
		return nil, err
	}

	resultList, err := storageEngine.MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	accountUser := &entity.AccountUser{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), accountUser); err != nil {
			return nil, err
		}
		accountUsers = append(accountUsers, accountUser)
		accountUser = &entity.AccountUser{}
	}

	return
}

// WriteAccountUser adds a new account user to the database and returns the created account user or an error
func WriteAccountUser(accountUser *entity.AccountUser, retrieve bool) (accUser *entity.AccountUser, err error) {
	if accountUser.ID, err = storageClient.GenerateAccountUserID(accountUser.AccountID); err != nil {
		return nil, err
	}

	accountUser.Enabled = true
	accountUser.CreatedAt = time.Now()
	accountUser.UpdatedAt, accountUser.ReceivedAt = accountUser.CreatedAt, accountUser.CreatedAt

	val, err := json.Marshal(accountUser)
	if err != nil {
		return nil, err
	}

	key := storageClient.AccountUserKey(accountUser.AccountID, accountUser.ID)
	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	listKey := storageClient.AccountUsersKey(accountUser.AccountID)
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return accountUser, nil
	}

	return ReadAccountUser(accountUser.AccountID, accountUser.ID)
}
