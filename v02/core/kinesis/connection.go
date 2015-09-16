package kinesis

import (
	"encoding/json"
	"fmt"

	"github.com/tapglue/multiverse/errors"
	"github.com/tapglue/multiverse/v02/core"
	"github.com/tapglue/multiverse/v02/entity"
	"github.com/tapglue/multiverse/v02/errmsg"
	"github.com/tapglue/multiverse/v02/storage/kinesis"

	ksis "github.com/sendgridlabs/go-kinesis"
)

type (
	connection struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (c *connection) Create(accountID, applicationID int64, conn *entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	con := entity.ConnectionWithIDs{}
	con.AccountID = accountID
	con.ApplicationID = applicationID
	con.Connection = *conn
	data, er := json.Marshal(con)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionCreate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return conn, nil
	}

	return nil, nil
}

func (c *connection) Read(accountID, applicationID int64, userFromID, userToID uint64) (connection *entity.Connection, err []errors.Error) {
	return connection, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) Update(accountID, applicationID int64, existingConnection, updatedConnection entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	con := entity.ConnectionWithIDs{}
	con.AccountID = accountID
	con.ApplicationID = applicationID
	con.Connection = updatedConnection
	data, er := json.Marshal(con)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while updating the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionUpdate, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return &updatedConnection, nil
	}

	return nil, nil
}

func (c *connection) Delete(accountID, applicationID int64, connection *entity.Connection) []errors.Error {
	con := entity.ConnectionWithIDs{}
	con.AccountID = accountID
	con.ApplicationID = applicationID
	con.Connection = *connection
	data, er := json.Marshal(con)
	if er != nil {
		return []errors.Error{errors.NewInternalError(0, "error while deleting the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionDelete, partitionKey, data)
	if err != nil {
		return []errors.Error{err}
	}

	return nil
}

func (c *connection) List(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, err []errors.Error) {
	return users, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) FollowedBy(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, err []errors.Error) {
	return users, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) Friends(accountID, applicationID int64, userID uint64) (users []*entity.ApplicationUser, err []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) FriendsAndFollowing(accountID, applicationID int64, userID uint64) ([]*entity.ApplicationUser, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) Confirm(accountID, applicationID int64, connection *entity.Connection, retrieve bool) (*entity.Connection, []errors.Error) {
	con := entity.ConnectionWithIDs{}
	con.AccountID = accountID
	con.ApplicationID = applicationID
	con.Connection = *connection
	data, er := json.Marshal(con)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while confirming the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionConfirm, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	if retrieve {
		return connection, nil
	}

	return nil, nil
}

func (c *connection) WriteEventsToList(accountID, applicationID int64, connection *entity.Connection) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID int64, userFromID, userToID uint64) (err []errors.Error) {
	return []errors.Error{errmsg.ErrServerInvalidHandler}
}

func (c *connection) SocialConnect(accountID, applicationID int64, user *entity.ApplicationUser, platform string, socialFriendsIDs []string, connectionType string) ([]*entity.ApplicationUser, []errors.Error) {
	msg := entity.SocialConnection{
		Platform:         platform,
		Type:             connectionType,
		SocialFriendsIDs: socialFriendsIDs,
	}
	msg.User.AccountID = accountID
	msg.User.ApplicationID = applicationID
	msg.User.ApplicationUser = *user
	data, er := json.Marshal(msg)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while confirming the connection (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionSocialConnect, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	return nil, nil
}

func (c *connection) AutoConnectSocialFriends(accountID, applicationID int64, user *entity.ApplicationUser, connectionType string, ourStoredUsersIDs []*entity.ApplicationUser) ([]*entity.ApplicationUser, []errors.Error) {
	msg := entity.AutoConnectSocialFriends{
		Type:              connectionType,
		OurStoredUsersIDs: ourStoredUsersIDs,
	}
	msg.User.AccountID = accountID
	msg.User.ApplicationID = applicationID
	msg.User.ApplicationUser = *user
	data, er := json.Marshal(msg)
	if er != nil {
		return nil, []errors.Error{errors.NewInternalError(0, "error while creating the connections via social platform (1)", er.Error())}
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", accountID, applicationID)
	_, err := c.storage.PackAndPutRecord(kinesis.StreamConnectionAutoConnect, partitionKey, data)
	if err != nil {
		return nil, []errors.Error{err}
	}

	return nil, nil
}

func (c *connection) Relation(accountID, applicationID int64, userFromID, userToID uint64) (*entity.Relation, []errors.Error) {
	return nil, []errors.Error{errmsg.ErrServerNotImplementedYet}
}

// NewConnection creates a new Connection
func NewConnection(storageClient kinesis.Client) core.Connection {
	return &connection{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
