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
	connection struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (c *connection) Create(accountID, applicationID int64, conn *entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	data, er := json.Marshal(conn)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionCreate, partitionKey, data)

	return nil, []errors.Error{err}
}

func (c *connection) Read(accountID, applicationID int64, userFromID, userToID string) (connection *entity.Connection, err []errors.Error) {
	return connection, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) Update(accountID, applicationID int64, existingConnection, updatedConnection entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	data, er := json.Marshal(updatedConnection)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionUpdate, partitionKey, data)

	return nil, []errors.Error{err}
}

func (c *connection) Delete(accountID, applicationID int64, connection *entity.Connection) []errors.Error {
	data, er := json.Marshal(connection)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionDelete, partitionKey, data)

	return []errors.Error{err}
}

func (c *connection) List(accountID, applicationID int64, userID string) (users []*entity.ApplicationUser, err []errors.Error) {
	return users, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) FollowedBy(accountID, applicationID int64, userID string) (users []*entity.ApplicationUser, err []errors.Error) {
	return users, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) Friends(accountID, applicationID int64, userID string) (users []*entity.ApplicationUser, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) FriendsAndFollowing(accountID, applicationID int64, userID string) ([]*entity.ApplicationUser, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) Confirm(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	data, er := json.Marshal(connection)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while confirming the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionConfirm, partitionKey, data)

	return nil, []errors.Error{err}
}

func (c *connection) WriteEventsToList(accountID, applicationID int64, connection *entity.Connection) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID int64, userFromID, userToID string) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, []errors.Error) {
	data, er := json.Marshal(struct {
		User             *entity.ApplicationUser `json:"user"`
		Platform         string                  `json:"platform"`
		SocialFriendsIDs []string                `json:"social_friends_ids"`
	}{
		User:             user,
		Platform:         platform,
		SocialFriendsIDs: socialFriendsIDs,
	})
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while confirming the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionSocialConnect, partitionKey, data)

	return nil, []errors.Error{err}
}

func (c *connection) AutoConnectSocialFriends(accountID, applicationID int64, user *entity.ApplicationUser, connectionType string, ourStoredUsersIDs []*entity.ApplicationUser) ([]*entity.ApplicationUser, []errors.Error) {
	data, er := json.Marshal(struct {
		User              *entity.ApplicationUser   `json:"user"`
		OurStoredUsersIDs []*entity.ApplicationUser `json:"our_stored_users_ids"`
	}{
		User:              user,
		OurStoredUsersIDs: ourStoredUsersIDs,
	})
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the connections via social platform (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionAutoConnect, partitionKey, data)

	return nil, []errors.Error{err}
}

// NewConnection creates a new Connection
func NewConnection(storageClient kinesis.Client) core.Connection {
	return &connection{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
