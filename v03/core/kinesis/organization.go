package kinesis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v03/core"
	"github.com/tapglue/multiverse/v03/entity"
	"github.com/tapglue/multiverse/v03/errmsg"
	"github.com/tapglue/multiverse/v03/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type account struct {
	storage kinesis.Client
	ksis    *ksis.Kinesis
}

func (a *account) Create(account *entity.Organization, retrieve bool) (acc *entity.Organization, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (a *account) Read(accountID int64) (account *entity.Organization, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (a *account) Update(existingAccount, updatedAccount entity.Organization, retrieve bool) (*entity.Organization, []errors.Error) {
	data, er := json.Marshal(updatedAccount)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the account (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("account-%d-update", updatedAccount.ID)
	_, err := a.storage.PackAndPutRecord(kinesis.StreamAccountUpdate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return &updatedAccount, nil
	}

	return nil, nil
}

func (a *account) Delete(account *entity.Organization) []errors.Error {
	data, er := json.Marshal(account)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the account (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partition-%d-delete", account.ID)
	_, err := a.storage.PackAndPutRecord(kinesis.StreamAccountDelete, partitionKey, data)
	if err != nil {
		return []errors.Error{err}
	}

	return nil
}

func (a *account) Exists(accountID int64) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (a *account) FindByKey(authKey string) (*entity.Organization, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (a *account) ReadByPublicID(id string) (*entity.Organization, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

// NewAccount creates a new Account
func NewAccount(storageClient kinesis.Client) core.Organization {
	return &account{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
