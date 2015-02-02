/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"encoding/json"

	"time"

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

// UpdateAccount updates the account matching the ID or an error
func UpdateAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err error) {
	if account.Token, err = storageClient.GenerateAccountToken(account); err != nil {
		return nil, err
	}

	key := storageClient.AccountKey(account.ID)

	if exist, err := storageEngine.Exists(key).Result(); !exist {
		// TODO send proper HTTP code
		return nil, err
	}

	account.UpdatedAt = time.Now()

	val, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	if !retrieve {
		return account, nil
	}

	return ReadAccount(account.ID)
}

// DeleteAccount deletes the account matching the ID or an error
func DeleteAccount(accountID int64) (rs string, err error) {
	result, err := storageEngine.Del(storageClient.AccountKey(accountID)).Result()

	switch result {
	case 0:
		rs = "The resource for the provided id doesn't exist"
	case 1:
		rs = "Resource was deleted successfully"
	}

	return rs, err
}

// WriteAccount adds a new account to the database and returns the created account or an error
func WriteAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err error) {
	if account.ID, err = storageClient.GenerateAccountID(); err != nil {
		return nil, err
	}

	if account.Token, err = storageClient.GenerateAccountToken(account); err != nil {
		return nil, err
	}

	account.Enabled = true
	account.CreatedAt = time.Now()
	account.UpdatedAt, account.ReceivedAt = account.CreatedAt, account.CreatedAt

	val, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	if exist, err := storageEngine.SetNX(storageClient.AccountKey(account.ID), string(val)).Result(); !exist {
		// TODO Wrong HTTP CODE is SENT (200 instead of 4XX)
		return nil, err
	}

	if !retrieve {
		return account, nil
	}

	return ReadAccount(account.ID)
}
