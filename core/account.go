/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/core/entity"
	"github.com/tapglue/backend/redis"
)

// Defining redis keys
const AccountKey string = "account_%d"

// generateAccountID generates a new account ID
func generateAccountID() (int64, error) {
	incr := redis.Client().Incr("ids_account")
	return incr.Result()
}

// ReadAccount returns the account matching the ID or an error
func ReadAccount(accountID int64) (account *entity.Account, err error) {
	// Read from db
	result, err := redis.Client().Get(fmt.Sprintf(AccountKey, accountID)).Result()
	if err != nil {
		return nil, err
	}

	// Parse JSON
	if err = json.Unmarshal([]byte(result), &account); err != nil {
		return nil, err
	}

	return
}

// WriteAccount adds a new account to the database and returns the created account or an error
func WriteAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err error) {
	// Generate account id
	if account.ID, err = generateAccountID(); err != nil {
		return nil, err
	}

	// Encode JSON
	val, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Write
	if err = redis.Client().Set(fmt.Sprintf(AccountKey, account.ID), string(val)).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return account, nil
	}

	// Return resource
	return ReadAccount(account.ID)
}