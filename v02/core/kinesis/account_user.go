package redis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	accountUser struct {
		a       core.Account
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (au *accountUser) Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, errors.Error) {
	data, er := json.Marshal(accountUser)
	if er != nil {
		return nil, errors.NewInternalError("error while creating the account user (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("account-user-%d-%d", accountUser.AccountID, accountUser.ID)
	_, err := au.storage.PackAndPutRecord(kinesis.StreamAccountUserCreate, partitionKey, data)

	return nil, err
}

func (au *accountUser) Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, errors.Error) {
	data, er := json.Marshal(updatedAccountUser)
	if er != nil {
		return nil, errors.NewInternalError("error while updating the account user (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("account-user-%d-%d", updatedAccountUser.AccountID, updatedAccountUser.ID)
	_, err := au.storage.PackAndPutRecord(kinesis.StreamAccountUserUpdate, partitionKey, data)

	return nil, err
}

func (au *accountUser) Delete(accountUser *entity.AccountUser) errors.Error {
	data, er := json.Marshal(accountUser)
	if er != nil {
		return errors.NewInternalError("error while creating the event (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("account-user-%d-%d", accountUser.AccountID, accountUser.ID)
	_, err := au.storage.PackAndPutRecord(kinesis.StreamAccountUserDelete, partitionKey, data)

	return err
}

func (au *accountUser) List(accountID int64) (accountUsers []*entity.AccountUser, er errors.Error) {
	return accountUsers, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) CreateSession(user *entity.AccountUser) (string, errors.Error) {
	return "", errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) RefreshSession(sessionToken string, user *entity.AccountUser) (string, errors.Error) {
	return "", errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) DestroySession(sessionToken string, user *entity.AccountUser) errors.Error {
	return errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) GetSession(user *entity.AccountUser) (string, errors.Error) {
	return "", errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) FindByEmail(email string) (*entity.Account, *entity.AccountUser, errors.Error) {
	return nil, nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) ExistsByEmail(email string) (bool, errors.Error) {
	return false, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) FindByUsername(username string) (*entity.Account, *entity.AccountUser, errors.Error) {
	return nil, nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) ExistsByUsername(username string) (bool, errors.Error) {
	return false, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) ExistsByID(accountID, userID int64) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (au *accountUser) FindBySession(sessionKey string) (*entity.AccountUser, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewAccountUser creates a new AccountUser
func NewAccountUser(storageClient kinesis.Client) core.AccountUser {
	return &accountUser{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
		a:       NewAccount(storageClient),
	}
}

// NewAccountUserWithAccount creates a new AccountUser
func NewAccountUserWithAccount(storageClient kinesis.Client, a core.Account) core.AccountUser {
	return &accountUser{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
		a:       a,
	}
}
