package redis

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	storageHelper "github.com/tapglue/backend/v02/storage/helper"
	"github.com/tapglue/backend/v02/storage/redis"

	red "gopkg.in/redis.v2"
)

type (
	account struct {
		storage redis.Client
		redis   *red.Client
	}
)

func (a *account) Create(account *entity.Account, retrieve bool) (acc *entity.Account, err errors.Error) {
	var er error
	if account.ID, er = a.storage.GenerateAccountID(); er != nil {
		return nil, errors.NewInternalError("failed to write the account (1)", er.Error())
	}

	account.AuthToken = storageHelper.GenerateAccountSecretKey(account)
	account.Enabled = true
	account.CreatedAt = time.Now()
	account.UpdatedAt = account.CreatedAt

	val, er := json.Marshal(account)
	if er != nil {
		return nil, errors.NewInternalError("failed to write the account (2)", er.Error())
	}

	// TODO this should never happen, maybe we should panic instead just to catch it better?
	exist, er := a.redis.SetNX(storageHelper.Account(account.ID), string(val)).Result()
	if !exist {
		return nil, errors.NewInternalError("failed to write the account (3)", "account id already present")
	}
	if er != nil {
		return nil, errors.NewInternalError("failed to write the account (4)", er.Error())
	}

	// Store the token details in redis
	_, er = a.redis.HMSet(
		"tokens:"+utils.Base64Encode(account.AuthToken),
		"acc", strconv.FormatInt(account.ID, 10),
	).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to write the account (5)", er.Error())
	}

	if !retrieve {
		return account, nil
	}

	return a.Read(account.ID)
}

func (a *account) Read(accountID int64) (account *entity.Account, err errors.Error) {
	result, er := a.redis.Get(storageHelper.Account(accountID)).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to retrieve the account (1)", er.Error())
	}

	if er := json.Unmarshal([]byte(result), &account); er != nil {
		return nil, errors.NewInternalError("failed to retrieve the account (2)", er.Error())
	}

	return
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (acc *entity.Account, err errors.Error) {
	updatedAccount.UpdatedAt = time.Now()

	val, er := json.Marshal(updatedAccount)
	if er != nil {
		return nil, errors.NewInternalError("failed to update the account (1)", er.Error())
	}

	key := storageHelper.Account(updatedAccount.ID)
	exist, er := a.redis.Exists(key).Result()
	if !exist {
		return nil, errors.NewInternalError("failed to update the account (2)\naccount does not exist", "account does not exist")
	}
	if er != nil {
		return nil, errors.NewInternalError("failed to update the account (3)", er.Error())
	}

	if er = a.redis.Set(key, string(val)).Err(); er != nil {
		return nil, errors.NewInternalError("failed to update the account (4)", err.Error())
	}

	if !retrieve {
		return &updatedAccount, nil
	}

	return a.Read(updatedAccount.ID)
}

func (a *account) Delete(account *entity.Account) (err errors.Error) {
	result, er := a.redis.Del(storageHelper.Account(account.ID)).Result()
	if er != nil {
		return errors.NewInternalError("failed to delete the account (1)", er.Error())
	}

	// TODO: Disable Account users
	// TODO: Disable Applications
	// TODO: Disable Applications Users
	// TODO: Disable Applications Events

	if result != 1 {
		return errors.NewNotFoundError("The resource for the provided id doesn't exist", fmt.Sprintf("unexisting account for id %d", account.ID))
	}

	return nil
}

func (a *account) Exists(accountID int64) bool {
	account, err := a.Read(accountID)
	if err != nil {
		return false
	}

	return account.Enabled
}

// NewAccount creates a new Account
func NewAccount(storageClient redis.Client) core.Account {
	return &account{
		storage: storageClient,
		redis:   storageClient.Datastore(),
	}
}
