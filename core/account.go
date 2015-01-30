/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"encoding/json"

	"github.com/tapglue/backend/core/entity"
)

// ReadAccount returns the account matching the ID or an error
func ReadAccount(accountID int64) (account *entity.Account, err error) {
	result, err := storageEngine.Get(storageClient.AccountKey(accountID)).Result()
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(result), &account); err != nil {
		return nil, err
	}

	return
}

// WriteAccount adds a new account to the database and returns the created account or an error
func WriteAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err error) {
	// Generate id
	if account.ID, err = storageClient.GenerateAccountID(); err != nil {
		return nil, err
	}

	val, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	// Write
	if err = storageEngine.Set(storageClient.AccountKey(account.ID), string(val)).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return account, nil
	}

	return ReadAccount(account.ID)
}
