package redis

import (
	"github.com/tapglue/backend/tgerrors"
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

func (appu *applicationUser) Create(user *entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (appu *applicationUser) Read(accountID, applicationID, userID int64) (user *entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (appu *applicationUser) Update(existingUser, updatedUser entity.ApplicationUser, retrieve bool) (usr *entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (appu *applicationUser) Delete(accountID, applicationID, userID int64) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (appu *applicationUser) List(accountID, applicationID int64) (users []*entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (appu *applicationUser) CreateSession(user *entity.ApplicationUser) (string, tgerrors.TGError) {
	panic("not implemented yet")
	return "", nil
}

func (appu *applicationUser) RefreshSession(sessionToken string, user *entity.ApplicationUser) (string, tgerrors.TGError) {
	panic("not implemented yet")
	return "", nil
}

func (appu *applicationUser) GetSession(user *entity.ApplicationUser) (string, tgerrors.TGError) {
	panic("not implemented yet")
	return "", nil
}

func (appu *applicationUser) DestroySession(sessionToken string, user *entity.ApplicationUser) tgerrors.TGError {
	panic("not implemented yet")
	return nil
}

func (appu *applicationUser) FindByEmail(accountID, applicationID int64, email string) (*entity.ApplicationUser, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (appu *applicationUser) ExistsByEmail(accountID, applicationID int64, email string) (bool, tgerrors.TGError) {
	panic("not implemented yet")
	return false, nil
}

func (appu *applicationUser) FindByUsername(accountID, applicationID int64, username string) (*entity.ApplicationUser, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (appu *applicationUser) ExistsByUsername(accountID, applicationID int64, username string) (bool, tgerrors.TGError) {
	panic("not implemented yet")
	return false, nil
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
