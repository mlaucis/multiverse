/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/backend/core/entity"
)

// ReadAccountUser returns the account matching the ID or an error
func ReadAccountUser(accountID, accountUserID int64) (accountUser *entity.AccountUser, err error) {
	result, err := storageEngine.Get(storageClient.AccountUser(accountID, accountUserID)).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &accountUser); err != nil {
		return nil, err
	}

	return
}

// UpdateAccountUser update an account user in the database and returns the updated account user or an error
func UpdateAccountUser(accountUser *entity.AccountUser, retrieve bool) (accUser *entity.AccountUser, err error) {
	accountUser.UpdatedAt = time.Now()

	val, err := json.Marshal(accountUser)
	if err != nil {
		return nil, err
	}

	key := storageClient.AccountUser(accountUser.AccountID, accountUser.ID)
	exist, err := storageEngine.Exists(key).Result()
	if !exist {
		return nil, fmt.Errorf("account user does not exist")
	}
	if err != nil {
		return nil, err
	}

	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	if !accountUser.Enabled {
		listKey := storageClient.AccountUsers(accountUser.AccountID)
		if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return accountUser, nil
	}

	return ReadAccountUser(accountUser.AccountID, accountUser.ID)
}

// DeleteAccountUser deletes the account user matching the IDs or an error
func DeleteAccountUser(accountID, userID int64) (err error) {
	// TODO: Make not deletable if its the only account user of an account
	key := storageClient.AccountUser(accountID, userID)
	result, err := storageEngine.Del(key).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		return fmt.Errorf("The resource for the provided id doesn't exist")
	}

	listKey := storageClient.AccountUsers(accountID)
	if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
		return err
	}

	return nil
}

// ReadAccountUserList returns all the users from a certain account
func ReadAccountUserList(accountID int64) (accountUsers []*entity.AccountUser, err error) {
	result, err := storageEngine.LRange(storageClient.Account(accountID), 0, -1).Result()
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

	key := storageClient.AccountUser(accountUser.AccountID, accountUser.ID)
	exist, err := storageEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, fmt.Errorf("account user does not exists")
	}
	if err != nil {
		return nil, err
	}

	listKey := storageClient.AccountUsers(accountUser.AccountID)
	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return accountUser, nil
	}

	return ReadAccountUser(accountUser.AccountID, accountUser.ID)
}
