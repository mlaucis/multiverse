/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package redis

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
	connection struct {
		storage kinesis.Client
		ksis    *ksis.Kinesis
	}
)

func (c *connection) Create(conn *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	data, er := json.Marshal(conn)
	if er != nil {
		return nil, errors.NewInternalError("error while creating the connection (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", conn.AccountID, conn.ApplicationID)
	_, err = c.storage.PutRecord(kinesis.StreamConnectionCreate, partitionKey, data)

	return nil, err
}

func (c *connection) Read(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err errors.Error) {
	return connection, errors.NewInternalError("no suitable implementation found", "no suitable implementation found")
}

func (c *connection) Update(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	data, er := json.Marshal(updatedConnection)
	if er != nil {
		return nil, errors.NewInternalError("error while updating the connection (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", updatedConnection.AccountID, updatedConnection.ApplicationID)
	_, err = c.storage.PutRecord(kinesis.StreamConnectionUpdate, partitionKey, data)

	return nil, err
}

func (c *connection) Delete(connection *entity.Connection) (err errors.Error) {
	data, er := json.Marshal(connection)
	if er != nil {
		return errors.NewInternalError("error while deleting the connection (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", connection.AccountID, connection.ApplicationID)
	_, err = c.storage.PutRecord(kinesis.StreamConnectionDelete, partitionKey, data)

	return err
}

func (c *connection) List(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err errors.Error) {
	return users, errors.NewInternalError("no suitable implementation found", "no suitable implementation found")
}

func (c *connection) FollowedBy(accountID, applicationID, userID int64) (users []*entity.ApplicationUser, err errors.Error) {
	return users, errors.NewInternalError("no suitable implementation found", "no suitable implementation found")
}

func (c *connection) Confirm(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	data, er := json.Marshal(connection)
	if er != nil {
		return nil, errors.NewInternalError("error while confirming the connection (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", connection.AccountID, connection.ApplicationID)
	_, err = c.storage.PutRecord(kinesis.StreamConnectionConfirm, partitionKey, data)

	return nil, err
}

func (c *connection) WriteEventsToList(connection *entity.Connection) (err errors.Error) {
	return errors.NewInternalError("no suitable implementation found", "no suitable implementation found")
}

func (c *connection) DeleteEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err errors.Error) {
	return errors.NewInternalError("no suitable implementation found", "no suitable implementation found")
}

func (c *connection) SocialConnect(user *entity.ApplicationUser, platform string, socialFriendsIDs []string) (users []*entity.ApplicationUser, err errors.Error) {
	data, er := json.Marshal(struct {
		user             *entity.ApplicationUser
		platform         string
		socialFriendsIDs []string
	}{
		user:             user,
		platform:         platform,
		socialFriendsIDs: socialFriendsIDs,
	})
	if er != nil {
		return nil, errors.NewInternalError("error while confirming the connection (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", user.AccountID, user.ApplicationID)
	_, err = c.storage.PutRecord(kinesis.StreamConnectionSocialConnect, partitionKey, data)

	return nil, err
}

func (c *connection) AutoConnectSocialFriends(user *entity.ApplicationUser, ourStoredUsersIDs []interface{}) (users []*entity.ApplicationUser, err errors.Error) {
	data, er := json.Marshal(struct {
		user              *entity.ApplicationUser
		ourStoredUsersIDs []interface{}
	}{
		user:              user,
		ourStoredUsersIDs: ourStoredUsersIDs,
	})
	if er != nil {
		return nil, errors.NewInternalError("error while creating the connections via social platform (1)", er.Error())
	}

	partitionKey := fmt.Sprintf("partitionKey-%d-%d", user.AccountID, user.ApplicationID)
	_, err = c.storage.PutRecord(kinesis.StreamConnectionAutoConnect, partitionKey, data)

	return nil, err
}

// NewConnection creates a new Connection
func NewConnection(storageClient kinesis.Client) core.Connection {
	return &connection{
		storage: storageClient,
		ksis:    storageClient.Datastore(),
	}
}
