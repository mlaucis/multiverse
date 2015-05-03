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
	account struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (a *account) Create(account *entity.Account, retrieve bool) (acc *entity.Account, err errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (a *account) Read(accountID int64) (account *entity.Account, err errors.Error) {
	return nil, errors.NewNotFoundError("not found", "invalid handler specified")
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (acc *entity.Account, err errors.Error) {
	data, er := json.Marshal(updatedAccount)
	if er != nil {
		return nil, errors.NewInternalError("error while updating the account (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("account-%d-update", updatedAccount.ID)
	_, err = a.storage.PackAndPutRecord(kinesis.StreamAccountUpdate, partitionKey, data)

	return nil, err
}

func (a *account) Delete(account *entity.Account) (err errors.Error) {
	data, er := json.Marshal(account)
	if er != nil {
		return errors.NewInternalError("error while deleting the account (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partition-%d-delete", account.ID)
	_, err = a.storage.PackAndPutRecord(kinesis.StreamAccountDelete, partitionKey, data)

	return err
}

func (a *account) Exists(accountID int64) (bool, errors.Error) {
	return false, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (a *account) FindByKey(authKey string) (*entity.Account, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

func (a *account) ReadByPublicID(id string) (*entity.Account, errors.Error) {
	return nil, errors.NewInternalError("not implemented yet", "not implemented yet")
}

// NewAccount creates a new Account
func NewAccount(storageClient kinesis.Client) core.Account {
	return &account{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
