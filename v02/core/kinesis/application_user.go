package redis

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/kinesis"

	"encoding/json"
	"fmt"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	applicationUser struct {
		c       core.Connection
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (appu *applicationUser) Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError) {
	data, er := json.Marshal(updatedUser)
	if er != nil {
		return tgerrors.NewInternalError("error while updating the user (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-user-update-%d-%d-%d", updatedUser.AccountID, updatedUser.ApplicationID, updatedUser.ID)
	_, err = appu.storage.PutRecord("application_user_update", partitionKey, data)

	return nil, err
}

func (appu *applicationUser) Delete(applicationUser *entity.ApplicationUser) (err tgerrors.TGError) {
	data, er := json.Marshal(applicationUser)
	if er != nil {
		return tgerrors.NewInternalError("error while deleting the user (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("application-user-delete-%d-%d-%d", applicationUser.AccountID, applicationUser.ApplicationID, applicationUser.ID)
	_, err = appu.storage.PutRecord("application_user_delete", partitionKey, data)

	return err
}

func (appu *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) CreateSession(user *entity.ApplicationUser) (string, tgerrors.TGError) {
	return "", tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, tgerrors.TGError) {
	return "", tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) GetSession(user *entity.ApplicationUser) (string, tgerrors.TGError) {
	return "", tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) DestroySession(sessionToken string, user *entity.ApplicationUser) tgerrors.TGError {
	return tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, tgerrors.TGError) {
	return false, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, tgerrors.TGError) {
	return nil, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, tgerrors.TGError) {
	return false, tgerrors.NewNotFoundError("not found", "invalid handler specified")
}

func (appu *applicationUser) ExistsByID(accountID, applicationID, userID int64) bool {
	panic("not implemented yet")
	return false
}

// NewApplicationUser creates a new Event
func NewApplicationUser(storageClient kinesis.Client) core.ApplicationUser {
	return &applicationUser{
		c:       NewConnection(storageClient),
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}

// NewApplicationUserWithConnection creates a new Event
func NewApplicationUserWithConnection(storageClient kinesis.Client, c core.Connection) core.ApplicationUser {
	return &applicationUser{
		c:       c,
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
