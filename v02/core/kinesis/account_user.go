package redis

import (
	"github.com/tapglue/backend/tgerrors"
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

func (au *accountUser) Create(accountUser *entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) Read(accountID, accountUserID int64) (accountUser *entity.AccountUser, er tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) Update(existingAccountUser, updatedAccountUser entity.AccountUser, retrieve bool) (*entity.AccountUser, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) Delete(accountID, userID int64) tgerrors.TGError {
	return tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) List(accountID int64) (accountUsers []*entity.AccountUser, er tgerrors.TGError) {
	return accountUsers, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) CreateSession(user *entity.AccountUser) (string, tgerrors.TGError) {
	return "", tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) RefreshSession(sessionToken string, user *entity.AccountUser) (string, tgerrors.TGError) {
	return "", tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) DestroySession(sessionToken string, user *entity.AccountUser) tgerrors.TGError {
	return tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) GetSession(user *entity.AccountUser) (string, tgerrors.TGError) {
	return "", tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) FindByEmail(email string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	return nil, nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) ExistsByEmail(email string) (bool, tgerrors.TGError) {
	return false, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) FindByUsername(username string) (*entity.Account, *entity.AccountUser, tgerrors.TGError) {
	return nil, nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) ExistsByUsername(username string) (bool, tgerrors.TGError) {
	return false, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (au *accountUser) ExistsByID(accountID, userID int64) bool {
	panic("not implemented yet")
	return false
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
