package redis

import (
	"github.com/tapglue/backend/tgerrors"
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

func (a *account) Create(account *entity.Account, retrieve bool) (acc *entity.Account, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (a *account) Read(accountID int64) (account *entity.Account, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (a *account) Update(existingAccount, updatedAccount entity.Account, retrieve bool) (acc *entity.Account, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (a *account) Delete(accountID int64) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (a *account) Exists(accountID int64) bool {
	panic("not implemented yet")
	return false
}

// NewAccount creates a new Account
func NewAccount(storageClient kinesis.Client) core.Account {
	return &account{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
