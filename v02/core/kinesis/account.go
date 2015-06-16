package kinesis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/errmsg"
	"github.com/tapglue/backend/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	account struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (a *account) Create(account *entity.Account, retrieve bool) (acc *entity.Account, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (a *account) Read(accountID int64) (account *entity.Account, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (*entity.Account, []errors.Error) {
	data, er := json.Marshal(updatedAccount)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the account (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("account-%d-update", updatedAccount.ID)
	_, err := a.storage.PackAndPutRecord(kinesis.StreamAccountUpdate, partitionKey, data)

	return nil, []errors.Error{err}
}

func (a *account) Delete(account *entity.Account) []errors.Error {
	data, er := json.Marshal(account)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the account (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partition-%d-delete", account.ID)
	_, err := a.storage.PackAndPutRecord(kinesis.StreamAccountDelete, partitionKey, data)

	return []errors.Error{err}
}

func (a *account) Exists(accountID int64) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (a *account) FindByKey(authKey string) (*entity.Account, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (a *account) ReadByPublicID(id string) (*entity.Account, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

// NewAccount creates a new Account
func NewAccount(storageClient kinesis.Client) core.Account {
	return &account{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
