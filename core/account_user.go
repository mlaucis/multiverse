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
const AccountUserKey string = "account_%d_user%d"

// getnewAccountUserID generates a new account user ID
func getNewAccountUserID(accountID int64) (int64, error) {
	incr := redis.Client().Incr(fmt.Sprintf("ids_account_%d_user", accountID))
	return incr.Result()
}

// GetAccountUserByID returns the account matching the ID or an error
func GetAccountUserByID(accountID int64, accountUserID int64) (accountUser *entity.AccountUser, err error) {
	// Read from db
	result, err := redis.Client().Get(fmt.Sprintf(AccountUserKey, accountID, accountUserID)).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &accountUser); err != nil {
		return nil, err
	}

	return
}

// AddAccount adds a new account user to the database and returns the created account or an error
func AddAccountUser(accountID int64, accountUser *entity.AccountUser, retrieve bool) (accUser *entity.AccountUser, err error) {
	// Validate account
	// if err = validator.ValidateAccountUser(accountUser); err != nil {
	// 	return nil, err
	// }

	// Generate account id
	if accountUser.ID, err = getNewAccountUserID(accountID); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(accountUser)
	if err != nil {
		return nil, err
	}

	// Write to db
	if err = redis.Client().Set(fmt.Sprintf(AccountUserKey, accountID, accountUser.ID), string(val)).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return accountUser, nil
	}

	// Return resource
	return GetAccountUserByID(accountID, accountUser.ID)
}
