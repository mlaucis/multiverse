/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/tapglue/backend/errors"
	"github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/entity"

	red "gopkg.in/redis.v2"
)

// UpdateConnection updates a connection in the database and returns the updated connection user or an error
func UpdateConnection(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	updatedConnection.UpdatedAt = time.Now()
	var er error

	val, er := json.Marshal(updatedConnection)
	if err != nil {
		return nil, errors.NewInternalError("failed to update the connection (1)", er.Error())
	}

	key := storageClient.Connection(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID, updatedConnection.UserToID)
	exist, er := storageEngine.Exists(key).Result()
	if !exist {
		return nil, errors.NewNotFoundError("failed to update teh connection (2)", "connection not found")
	}
	if er != nil {
		return nil, errors.NewInternalError("failed to update the connection (3)", er.Error())
	}

	if er = storageEngine.Set(key, string(val)).Err(); er != nil {
		return nil, errors.NewInternalError("failed to update the connection (4)", er.Error())
	}

	if !updatedConnection.Enabled {
		listKey := storageClient.Connections(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID)
		if er = storageEngine.LRem(listKey, 0, key).Err(); er != nil {
			return nil, errors.NewInternalError("failed to update the connection (5)", er.Error())
		}
		userListKey := storageClient.ConnectionUsers(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID)
		userKey := storageClient.User(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserToID)
		if er = storageEngine.LRem(userListKey, 0, userKey).Err(); er != nil {
			return nil, errors.NewInternalError("failed to update the connection (6)", er.Error())
		}
		followerListKey := storageClient.FollowedByUsers(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserToID)
		followerKey := storageClient.User(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID)
		if er = storageEngine.LRem(followerListKey, 0, followerKey).Err(); er != nil {
			return nil, errors.NewInternalError("failed to update the connection (7)", er.Error())
		}
	}

	if !retrieve {
		return &updatedConnection, nil
	}

	return ReadConnection(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID, updatedConnection.UserToID)
}

// DeleteConnection deletes the connection matching the IDs or an error
func DeleteConnection(accountID, applicationID, userFromID, userToID int64) (err errors.Error) {
	key := storageClient.Connection(accountID, applicationID, userFromID, userToID)
	result, er := storageEngine.Del(key).Result()
	if er != nil {
		return errors.NewInternalError("failed to delete the connection (1)", er.Error())
	}

	if result != 1 {
		return errors.NewNotFoundError("failed to delete the connection (2)", "connection not found")
	}

	listKey := storageClient.Connections(accountID, applicationID, userFromID)
	if er = storageEngine.LRem(listKey, 0, key).Err(); er != nil {
		return errors.NewInternalError("failed to delete the connection (3)", er.Error())
	}
	userListKey := storageClient.ConnectionUsers(accountID, applicationID, userFromID)
	userKey := storageClient.User(accountID, applicationID, userToID)
	if er = storageEngine.LRem(userListKey, 0, userKey).Err(); er != nil {
		return errors.NewInternalError("failed to delete the connection (4)", er.Error())
	}
	followerListKey := storageClient.FollowedByUsers(accountID, applicationID, userToID)
	followerKey := storageClient.User(accountID, applicationID, userFromID)
	if er = storageEngine.LRem(followerListKey, 0, followerKey).Err(); er != nil {
		return errors.NewInternalError("failed to delete the connection (5)", er.Error())
	}

	if err := DeleteConnectionEventsFromLists(accountID, applicationID, userFromID, userToID); err != nil {
		return errors.NewInternalError("failed to delete the connection (6)", err.Error())
	}

	return nil
}

// ReadConnectionList returns all connections from a certain user
func ReadConnectionList(accountID, applicationID, userID int64) (users []*entity.User, err errors.Error) {
	key := storageClient.ConnectionUsers(accountID, applicationID, userID)
	result, er := storageEngine.LRange(key, 0, -1).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the connection list (1)", er.Error())
	}

	if len(result) == 0 {
		return []*entity.User{}, nil
	}

	return fetchAndDecodeMultipleUsers(result)
}

// ReadFollowedByList returns all connections from a certain user
func ReadFollowedByList(accountID, applicationID, userID int64) (users []*entity.User, err errors.Error) {
	key := storageClient.FollowedByUsers(accountID, applicationID, userID)
	result, er := storageEngine.LRange(key, 0, -1).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to read the followers list (1)", er.Error())
	}

	if len(result) == 0 {
		return []*entity.User{}, nil
	}

	return fetchAndDecodeMultipleUsers(result)
}

// WriteConnection adds a user connection and returns the created connection or an error
func WriteConnection(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	// We confirm the connection in the past forcefully so that we can update it at the confirmation time
	connection.ConfirmedAt = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
	connection.Enabled = false
	connection.CreatedAt = time.Now()
	connection.UpdatedAt = connection.CreatedAt

	val, er := json.Marshal(connection)
	if er != nil {
		return nil, errors.NewInternalError("failed to write the user connection (1)", er.Error())
	}

	key := storageClient.Connection(connection.AccountID, connection.ApplicationID, connection.UserFromID, connection.UserToID)
	exist, er := storageEngine.SetNX(key, string(val)).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to write the user connection (2)", er.Error())
	}
	if !exist {
		return nil, errors.NewInternalError("failed to write the user connection (3)", "connection does not exist")
	}

	return connection, nil
}

// ConfirmConnection confirms a user connection and returns the connection or an error
func ConfirmConnection(connection *entity.Connection, retrieve bool) (con *entity.Connection, err errors.Error) {
	// We confirm the connection in the past forcefully so that we can update it at the confirmation time
	connection.Enabled = true
	connection.ConfirmedAt = time.Now()
	connection.UpdatedAt = connection.ConfirmedAt

	val, er := json.Marshal(connection)
	if er != nil {
		return nil, errors.NewInternalError("failed to confirm the connection (1)", er.Error())
	}

	key := storageClient.Connection(connection.AccountID, connection.ApplicationID, connection.UserFromID, connection.UserToID)

	cmd := red.NewStringCmd("SET", key, string(val), "XX")
	storageEngine.Process(cmd)
	er = cmd.Err()
	if er != nil {
		return nil, errors.NewInternalError("failed to confirm the connection (2)", er.Error())
	}

	listKey := storageClient.Connections(connection.AccountID, connection.ApplicationID, connection.UserFromID)
	if er = storageEngine.LPush(listKey, key).Err(); er != nil {
		return nil, errors.NewInternalError("failed to confirm the connection (3)", er.Error())
	}

	userListKey := storageClient.ConnectionUsers(connection.AccountID, connection.ApplicationID, connection.UserFromID)
	userKey := storageClient.User(connection.AccountID, connection.ApplicationID, connection.UserToID)
	if er = storageEngine.LPush(userListKey, userKey).Err(); er != nil {
		return nil, errors.NewInternalError("failed to confirm the connection (4)", er.Error())
	}

	followerListKey := storageClient.FollowedByUsers(connection.AccountID, connection.ApplicationID, connection.UserToID)
	followerKey := storageClient.User(connection.AccountID, connection.ApplicationID, connection.UserFromID)
	if er = storageEngine.LPush(followerListKey, followerKey).Err(); er != nil {
		return nil, errors.NewInternalError("failed to confirm the connection (5)", er.Error())
	}

	if err = WriteConnectionEventsToList(connection); err != nil {
		return nil, err
	}

	if !retrieve {
		return connection, nil
	}

	return connection, nil
}

// WriteConnectionEventsToList takes a connection and writes the events to the lists
func WriteConnectionEventsToList(connection *entity.Connection) (err errors.Error) {
	connectionEventsKey := storageClient.ConnectionEvents(connection.AccountID, connection.ApplicationID, connection.UserFromID)

	eventsKey := storageClient.Events(connection.AccountID, connection.ApplicationID, connection.UserToID)

	events, er := storageEngine.ZRevRangeWithScores(eventsKey, "0", "-1").Result()
	if er != nil {
		return errors.NewInternalError("failed to write the event to the list", er.Error())
	}

	if len(events) >= 1 {
		var vals []red.Z

		for _, eventKey := range events {
			val := red.Z{Score: float64(eventKey.Score), Member: eventKey.Member}
			vals = append(vals, val)
		}

		if er = storageEngine.ZAdd(connectionEventsKey, vals...).Err(); er != nil {
			return errors.NewInternalError("failed to write the event to the list", er.Error())
		}
	}

	return
}

// DeleteConnectionEventsFromLists takes a connection and deletes the events from the lists
func DeleteConnectionEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err errors.Error) {
	connectionEventsKey := storageClient.ConnectionEvents(accountID, applicationID, userFromID)

	eventsKey := storageClient.Events(accountID, applicationID, userToID)

	events, er := storageEngine.ZRevRangeWithScores(eventsKey, "0", "-1").Result()
	if er != nil {
		return errors.NewInternalError("failed to delete the event from connections (1)", er.Error())
	}

	if len(events) >= 1 {
		var members []string

		for _, eventKey := range events {
			member := eventKey.Member
			members = append(members, member)
		}

		if er = storageEngine.ZRem(connectionEventsKey, members...).Err(); er != nil {
			return errors.NewInternalError("failed to delete the event from the connections (2)", er.Error())
		}
	}

	return nil
}

// ReadConnection returns the connection, if any, between two users
func ReadConnection(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err errors.Error) {
	key := storageClient.Connection(accountID, applicationID, userFromID, userToID)

	result, er := storageEngine.Get(key).Result()
	if er != nil {
		if er.Error() == "redis: nil" {
			return nil, nil
		}
		return nil, errors.NewInternalError("failed to read the connection (1)", er.Error())
	}
	if result == "" {
		return nil, nil
	}

	connection = &entity.Connection{}
	er = json.Unmarshal([]byte(result), connection)
	if er == nil {
		return
	}
	return nil, errors.NewInternalError("failed to read the connection (3)", er.Error())
}

// SocialConnect creates the connections between a user and his other social peers
func SocialConnect(user *entity.User, platform string, socialFriendsIDs []string) ([]*entity.User, errors.Error) {
	result := []*entity.User{}

	encodedSocialFriendsIDs := []string{}
	for idx := range socialFriendsIDs {
		encodedSocialFriendsIDs = append(encodedSocialFriendsIDs, storageClient.SocialConnection(
			user.AccountID,
			user.ApplicationID,
			platform,
			utils.Base64Encode(socialFriendsIDs[idx])))
	}

	ourStoredUsersIDs, er := storageEngine.MGet(encodedSocialFriendsIDs...).Result()
	if er != nil {
		return result, errors.NewInternalError("social connection failed (1)", er.Error())
	}

	if len(ourStoredUsersIDs) == 0 {
		return result, nil
	}

	return autoConnectSocialFriends(user, ourStoredUsersIDs)
}

func autoConnectSocialFriends(user *entity.User, ourStoredUsersIDs []interface{}) (users []*entity.User, err errors.Error) {
	ourUserKeys := []string{}
	for idx := range ourStoredUsersIDs {
		userID, err := strconv.ParseInt(ourStoredUsersIDs[idx].(string), 10, 64)
		if err != nil {
			continue
		}

		key := storageClient.Connection(user.AccountID, user.ApplicationID, user.ID, userID)
		if exists, err := storageEngine.Exists(key).Result(); exists || err != nil {
			// We don't want to update existing connections as we don't know if the user disabled them willingly or not
			// TODO Figure out if this is the right thing to do
			continue
		}

		connection := &entity.Connection{
			AccountID:     user.AccountID,
			ApplicationID: user.ApplicationID,
			UserFromID:    user.ID,
			UserToID:      userID,
		}

		_, er := WriteConnection(connection, false)
		if er != nil {
			continue
		}

		_, er = ConfirmConnection(connection, false)
		if er != nil {
			continue
		}

		connection = &entity.Connection{
			AccountID:     user.AccountID,
			ApplicationID: user.ApplicationID,
			UserFromID:    userID,
			UserToID:      user.ID,
		}

		_, er = WriteConnection(connection, false)
		if er != nil {
			continue
		}

		_, er = ConfirmConnection(connection, false)
		if er != nil {
			continue
		}

		ourUserKeys = append(
			ourUserKeys,
			storageClient.User(user.AccountID, user.ApplicationID, userID),
		)
	}

	return fetchAndDecodeMultipleUsers(ourUserKeys)
}

func fetchAndDecodeMultipleUsers(keys []string) (users []*entity.User, err errors.Error) {
	if len(keys) == 0 {
		return []*entity.User{}, nil
	}

	resultList, er := storageEngine.MGet(keys...).Result()
	if er != nil {
		return nil, errors.NewInternalError("failed to perform operation on user list (1)", er.Error())
	}

	user := &entity.User{}
	for _, result := range resultList {
		if er = json.Unmarshal([]byte(result.(string)), user); er != nil {
			return nil, errors.NewInternalError("failed to perform operation on user list (2)", er.Error())
		}
		users = append(users, user)
		user = &entity.User{}
	}

	return
}
