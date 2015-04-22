/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/entity"
)

// ReadAccount returns the account matching the ID or an error
func ReadAccount(accountID int64) (account *entity.Account, err errors.Error) {
	result, er := storageEngine.Get(storageClient.Account(accountID)).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to retrieve the account (1)", er.Error())
	}

	if er := json.Unmarshal([]byte(result), &account); er != nil {
		return nil, errors.NewInternalError("failed to retrieve the account (2)", er.Error())
	}

	return
}

// UpdateAccount updates the account matching the ID or an error
func UpdateAccount(existingAccount, updatedAccount entity.Account, retrieve bool) (acc *entity.Account, err errors.Error) {
	updatedAccount.UpdatedAt = time.Now()

	val, er := json.Marshal(updatedAccount)
	if er != nil {
		return nil, errors.NewInternalError("failed to update the account (1)", er.Error())
	}

	key := storageClient.Account(updatedAccount.ID)
	exist, er := storageEngine.Exists(key).Result()
	if !exist {
		return nil, errors.NewInternalError("failed to update the account (2)\naccount does not exist", "account does not exist")
	}
	if er != nil {
		return nil, errors.NewInternalError("failed to update the account (3)", er.Error())
	}

	if er = storageEngine.Set(key, string(val)).Err(); er != nil {
		return nil, errors.NewInternalError("failed to update the account (4)", err.Error())
	}

	if !retrieve {
		return &updatedAccount, nil
	}

	return ReadAccount(updatedAccount.ID)
}

// DeleteAccount deletes the account matching the ID or an error
func DeleteAccount(accountID int64) (err errors.Error) {
	result, er := storageEngine.Del(storageClient.Account(accountID)).Result()
	if er != nil {
		return errors.NewInternalError("failed to delete the account (1)", er.Error())
	}

	// TODO: Disable Account users
	// TODO: Disable Applications
	// TODO: Disable Applications Users
	// TODO: Disable Applications Events

	if result != 1 {
		return errors.NewNotFoundError("The resource for the provided id doesn't exist", fmt.Sprintf("unexisting account for id %d", accountID))
	}

	return nil
}

// WriteAccount adds a new account to the database and returns the created account or an error
func WriteAccount(account *entity.Account, retrieve bool) (acc *entity.Account, err errors.Error) {
	var er error
	if account.ID, er = storageClient.GenerateAccountID(); er != nil {
		return nil, errors.NewInternalError("failed to write the account (1)", er.Error())
	}

	account.AuthToken = storageClient.GenerateAccountSecretKey(account)
	account.Enabled = true
	account.CreatedAt = time.Now()
	account.UpdatedAt = account.CreatedAt

	val, er := json.Marshal(account)
	if er != nil {
		return nil, errors.NewInternalError("failed to write the account (2)", er.Error())
	}

	// TODO this should never happen, maybe we should panic instead just to catch it better?
	exist, er := storageEngine.SetNX(storageClient.Account(account.ID), string(val)).Result()
	if !exist {
		return nil, errors.NewInternalError("failed to write the account (3)", "account id already present")
	}
	if er != nil {
		return nil, errors.NewInternalError("failed to write the account (4)", er.Error())
	}

	// Store the token details in redis
	_, er = storageEngine.HMSet(
		"tokens:"+utils.Base64Encode(account.AuthToken),
		"acc", strconv.FormatInt(account.ID, 10),
	).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to write the account (5)", er.Error())
	}

	if !retrieve {
		return account, nil
	}

	return ReadAccount(account.ID)
}
