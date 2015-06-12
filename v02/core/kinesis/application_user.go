package kinesis

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
	applicationUser struct {
		c       core.Connection
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (appu *applicationUser) Create(accountID, applicationID int64, user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) Read(accountID, applicationID int64, userID string) (user *entity.ApplicationUser, err []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) ReadMultiple(accountID, applicationID int64, userIDs []string) (users []*entity.ApplicationUser, err []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) Update(accountID, applicationID int64, existingUser, updatedUser entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, []errors.Error) {
	data, er := json.Marshal(updatedUser)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the user (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-user-update-%d-%d-%d", accountID, applicationID, updatedUser.ID)
	_, err := appu.storage.PackAndPutRecord(kinesis.StreamApplicationUserUpdate, partitionKey, data)

	return nil, []errors.Error{err}
}

func (appu *applicationUser) Delete(accountID, applicationID int64, applicationUser *entity.ApplicationUser) []errors.Error {
	data, er := json.Marshal(applicationUser)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the user (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-user-delete-%d-%d-%d", accountID, applicationID, applicationUser.ID)
	_, err := appu.storage.PackAndPutRecord(kinesis.StreamApplicationUserDelete, partitionKey, data)

	return []errors.Error{err}
}

func (appu *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	return "", invalidHandlerError
}

func (appu *applicationUser) RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, []errors.Error) {
	return "", invalidHandlerError
}

func (appu *applicationUser) GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	return "", invalidHandlerError
}

func (appu *applicationUser) DestroySession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) []errors.Error {
	return invalidHandlerError
}

func (appu *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, []errors.Error) {
	return false, invalidHandlerError
}

func (appu *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, []errors.Error) {
	return false, invalidHandlerError
}

func (appu *applicationUser) ExistsByID(accountID, applicationID int64, userID string) (bool, []errors.Error) {
	return false, invalidHandlerError
}

func (appu *applicationUser) FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) Search(accountID, applicationID int64, searchTerm string) ([]*entity.ApplicationUser, []errors.Error) {
	return nil, invalidHandlerError
}

func (appu *applicationUser) FindByCustomID(accountID, applicationID int64, customID string) (*entity.ApplicationUser, []errors.Error) {
	return nil, invalidHandlerError
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
