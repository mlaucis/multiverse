/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package redis

import (
	"github.com/tapglue/backend/tgerrors"
	"github.com/tapglue/backend/v02/core"
	"github.com/tapglue/backend/v02/entity"
	"github.com/tapglue/backend/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	connection struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (c *connection) Create(connection *entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) Read(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) Delete(accountID, applicationID, userFromID, userToID int64) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (c *connection) List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) FollowedBy(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) Confirm(connection *entity.Connection, retrieve bool) (con *entity.Connection, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) WriteEventsToList(connection *entity.Connection) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err tgerrors.TGError) {
	panic("not implemented yet")
	return nil
}

func (c *connection) SocialConnect(user *entity.ApplicationUser, platform string, socialFriendsIDs []string) ([]*entity.ApplicationUser, tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

func (c *connection) AutoConnectSocialFriends(user *entity.ApplicationUser, ourStoredUsersIDs []interface{}) (users []*entity.ApplicationUser, err tgerrors.TGError) {
	panic("not implemented yet")
	return nil, nil
}

// NewConnection creates a new Connection
func NewConnection(storageClient kinesis.Client) core.Connection {
	return &connection{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
