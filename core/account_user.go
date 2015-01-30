/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/redis"
)

// Defining redis keys
const (
	AccountUserKey  string = "account_%d_user_%d"
	AccountUsersKey string = "account_%d_users"
)

// generateAccountUserID generates a new account user ID
func generateAccountUserID(accountID int64) (int64, error) {
	incr := redis.Client().Incr(fmt.Sprintf("ids_account_%d_user", accountID))
	return incr.Result()
}

// ReadAccountUser returns the account matching the ID or an error
func ReadAccountUser(accountID int64, accountUserID int64) (accountUser *entity.AccountUser, err error) {
	// Generate resource key
	key := fmt.Sprintf(AccountUserKey, accountID, accountUserID)

	// Read from db
	result, err := redis.Client().Get(key).Result()
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
	// Generate resource key
	key := fmt.Sprintf(AccountUsersKey, accountID)

	// Read from db
	result, err := redis.Client().LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", result)

	// Read from db
	resultList, err := redis.Client().MGet(result...).Result()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%v\n", resultList)

	// Parse JSON
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
	// Generate account id
	if accountUser.ID, err = generateAccountUserID(accountUser.AccountID); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(accountUser)
	if err != nil {
		return nil, err
	}

	// Generate resource key
	key := fmt.Sprintf(AccountUserKey, accountUser.AccountID, accountUser.ID)

	// Write resource
	if err = redis.Client().Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	// Generate list key
	listKey := fmt.Sprintf(AccountUsersKey, accountUser.AccountID)

	// Write list
	if err = redis.Client().LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return accountUser, nil
	}

	// Return resource
	return ReadAccountUser(accountUser.AccountID, accountUser.ID)
}
