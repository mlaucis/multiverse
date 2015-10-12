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

type applicationUser struct {
	c       core.Connection
	storage kinesis.Client
	ksis    *ksis.Kinesis
}

func (appu *applicationUser) Create(accountID, applicationID int64, user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) Read(accountID, applicationID int64, userID uint64, withStatistics bool) (user *entity.ApplicationUser, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) ReadMultiple(accountID, applicationID int64, userIDs []uint64) (users []*entity.ApplicationUser, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) Update(accountID, applicationID int64, existingUser, updatedUser entity.ApplicationUser, retrieve bool) (*entity.ApplicationUser, []errors.Error) {
	user := entity.ApplicationUserWithIDs{}
	user.OrgID = accountID
	user.AppID = applicationID
	user.ApplicationUser = updatedUser
	data, er := json.Marshal(user)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the user (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-user-update-%d-%d-%d", accountID, applicationID, updatedUser.ID)
	_, err := appu.storage.PackAndPutRecord(kinesis.StreamApplicationUserUpdate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return &updatedUser, nil
	}

	return nil, nil
}

func (appu *applicationUser) Delete(accountID, applicationID int64, userID uint64) []errors.Error {
	user := entity.ApplicationUserWithIDs{}
	user.OrgID = accountID
	user.AppID = applicationID
	user.ApplicationUser = entity.ApplicationUser{ID: userID}
	data, er := json.Marshal(user)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the user (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("application-user-delete-%d-%d-%d", accountID, applicationID, userID)
	_, err := appu.storage.PackAndPutRecord(kinesis.StreamApplicationUserDelete, partitionKey, data)
	if err != nil {
		return []errors.Error{err}
	}

	return nil
}

func (appu *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) CreateSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	return "", []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) RefreshSession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) (string, []errors.Error) {
	return "", []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) GetSession(accountID, applicationID int64, user *entity.ApplicationUser) (string, []errors.Error) {
	return "", []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) DestroySession(accountID, applicationID int64, sessionToken string, user *entity.ApplicationUser) []errors.Error {
	return []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) ExistsByID(accountID, applicationID int64, userID uint64) (bool, []errors.Error) {
	return false, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) FindBySession(accountID, applicationID int64, sessionKey string) (*entity.ApplicationUser, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) Search(accountID, applicationID int64, searchTerm string) ([]*entity.ApplicationUser, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
}

func (appu *applicationUser) FriendStatistics(accountID, applicationID int64, appUser *entity.ApplicationUser) []errors.Error {
	return []errors.Error{errmsg.ErrServerInvalidHandler.SetCurrentLocation()}
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
