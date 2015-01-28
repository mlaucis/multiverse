/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/redis"
	"github.com/tapglue/backend/validator"
)

// getnewAccountID generates a new account ID
func getNewAccountID() (int64, error) {
	incr := redis.GetClient().Incr("ids_account")
	return incr.Result()
}

// GetAccountByID returns the account matching the ID or an error
func GetAccountByID(accountID int64) (account *entity.Account, err error) {
	result, err := redis.GetClient().Get(fmt.Sprintf("account_%d", accountID)).Result()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &account); err != nil {
		return nil, err
	}

	return
}

// AddAccount adds a new account to the database and returns the created account or an error
func AddAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err error) {
	if err = validator.ValidateAccount(account); err != nil {
		return nil, err
	}

	if account.ID, err = getNewAccountID(); err != nil {
		return nil, err
	}

	val, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	result, err := redis.GetClient().Set(fmt.Sprintf("account_%d", account.ID), string(val)).Result()
	if err != nil {
		return nil, err
	}

	fmt.Printf("%s\n", result)
	if !retrieve {
		return account, nil
	}

	// Return account
	return GetAccountByID(account.ID)
}
