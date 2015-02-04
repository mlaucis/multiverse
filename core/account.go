/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/tapglue/backend/core/entity"
)

// ReadAccount returns the account matching the ID or an error
func ReadAccount(accountID int64) (account *entity.Account, err error) {
	result, err := storageEngine.Get(storageClient.Account(accountID)).Result()
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
	account.UpdatedAt = time.Now()

	val, err := json.Marshal(account)
	if err != nil {
		return nil, err
	}

	key := storageClient.Account(account.ID)
	exist, err := storageEngine.Exists(key).Result()
	if !exist {
		return nil, fmt.Errorf("account does not exist")
	}
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
func DeleteAccount(accountID int64) (err error) {
	result, err := storageEngine.Del(storageClient.Account(accountID)).Result()
	if err != nil {
		return err
	}

	// TODO: Disable Account users
	// TODO: Disable Applications

	if result != 1 {
		return fmt.Errorf("The resource for the provided id doesn't exist")
	}

	return nil
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

	// TODO this should never happen, maybe we should panic instead just to catch it better?
	exist, err := storageEngine.SetNX(storageClient.Account(account.ID), string(val)).Result()
	if !exist {
		return nil, fmt.Errorf("account already exists")
	}
	if err != nil {
		return nil, err
	}

	if !retrieve {
		return account, nil
	}

	return ReadAccount(account.ID)
}
