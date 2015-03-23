/**
 * @author Onur Akpolat <onurakpolat@gmail.com>
 */

package core

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	. "github.com/tapglue/backend/utils"
	"github.com/tapglue/backend/v01/entity"

	red "gopkg.in/redis.v2"
)

// UpdateConnection updates a connection in the database and returns the updated connection user or an error
func UpdateConnection(existingConnection, updatedConnection entity.Connection, retrieve bool) (con *entity.Connection, err error) {
	updatedConnection.UpdatedAt = time.Now()

	val, err := json.Marshal(updatedConnection)
	if err != nil {
		return nil, err
	}

	key := storageClient.Connection(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID, updatedConnection.UserToID)
	exist, err := storageEngine.Exists(key).Result()
	if !exist {
		return nil, fmt.Errorf("connection does not exist")
	}
	if err != nil {
		return nil, err
	}

	if err = storageEngine.Set(key, string(val)).Err(); err != nil {
		return nil, err
	}

	if !updatedConnection.Enabled {
		listKey := storageClient.Connections(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID)
		if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
			return nil, err
		}
		userListKey := storageClient.ConnectionUsers(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID)
		userKey := storageClient.User(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserToID)
		if err = storageEngine.LRem(userListKey, 0, userKey).Err(); err != nil {
			return nil, err
		}
		followerListKey := storageClient.FollowedByUsers(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserToID)
		followerKey := storageClient.User(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID)
		if err = storageEngine.LRem(followerListKey, 0, followerKey).Err(); err != nil {
			return nil, err
		}
	}

	if !retrieve {
		return &updatedConnection, nil
	}

	return ReadConnection(updatedConnection.AccountID, updatedConnection.ApplicationID, updatedConnection.UserFromID, updatedConnection.UserToID)
}

// DeleteConnection deletes the connection matching the IDs or an error
func DeleteConnection(accountID, applicationID, userFromID, userToID int64) (err error) {
	key := storageClient.Connection(accountID, applicationID, userFromID, userToID)
	result, err := storageEngine.Del(key).Result()
	if err != nil {
		return err
	}

	if result != 1 {
		return fmt.Errorf("The resource for the provided id doesn't exist")
	}

	listKey := storageClient.Connections(accountID, applicationID, userFromID)
	if err = storageEngine.LRem(listKey, 0, key).Err(); err != nil {
		return err
	}
	userListKey := storageClient.ConnectionUsers(accountID, applicationID, userFromID)
	userKey := storageClient.User(accountID, applicationID, userToID)
	if err = storageEngine.LRem(userListKey, 0, userKey).Err(); err != nil {
		return err
	}
	followerListKey := storageClient.FollowedByUsers(accountID, applicationID, userToID)
	followerKey := storageClient.User(accountID, applicationID, userFromID)
	if err = storageEngine.LRem(followerListKey, 0, followerKey).Err(); err != nil {
		return err
	}

	if err = DeleteConnectionEventsFromLists(accountID, applicationID, userFromID, userToID); err != nil {
		return err
	}

	return nil
}

// ReadConnectionList returns all connections from a certain user
func ReadConnectionList(accountID, applicationID, userID int64) (users []*entity.User, err error) {
	key := storageClient.ConnectionUsers(accountID, applicationID, userID)
	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []*entity.User{}, nil
	}

	return fetchAndDecodeMultipleUsers(result)
}

// ReadFollowedByList returns all connections from a certain user
func ReadFollowedByList(accountID, applicationID, userID int64) (users []*entity.User, err error) {
	key := storageClient.FollowedByUsers(accountID, applicationID, userID)
	result, err := storageEngine.LRange(key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	if len(result) == 0 {
		return []*entity.User{}, nil
	}

	return fetchAndDecodeMultipleUsers(result)
}

// WriteConnection adds a user connection and returns the created connection or an error
func WriteConnection(connection *entity.Connection, retrieve bool) (con *entity.Connection, err error) {
	// We confirm the connection in the past forcefully so that we can update it at the confirmation time
	connection.ConfirmedAt = time.Date(0, time.January, 1, 0, 0, 0, 0, time.UTC)
	connection.Enabled = false
	connection.CreatedAt = time.Now()
	connection.UpdatedAt = connection.CreatedAt

	val, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	key := storageClient.Connection(connection.AccountID, connection.ApplicationID, connection.UserFromID, connection.UserToID)

	exist, err := storageEngine.SetNX(key, string(val)).Result()
	if !exist {
		return nil, fmt.Errorf("user connection already exists")
	}
	if err != nil {
		return nil, err
	}

	return connection, nil
}

// ConfirmConnection confirms a user connection and returns the connection or an error
func ConfirmConnection(connection *entity.Connection, retrieve bool) (con *entity.Connection, err error) {
	// We confirm the connection in the past forcefully so that we can update it at the confirmation time
	connection.Enabled = true
	connection.ConfirmedAt = time.Now()
	connection.UpdatedAt = connection.ConfirmedAt

	val, err := json.Marshal(connection)
	if err != nil {
		return nil, err
	}

	key := storageClient.Connection(connection.AccountID, connection.ApplicationID, connection.UserFromID, connection.UserToID)

	cmd := red.NewStringCmd("SET", key, string(val), "XX")
	storageEngine.Process(cmd)
	err = cmd.Err()
	if err != nil {
		return nil, err
	}

	listKey := storageClient.Connections(connection.AccountID, connection.ApplicationID, connection.UserFromID)

	if err = storageEngine.LPush(listKey, key).Err(); err != nil {
		return nil, err
	}

	userListKey := storageClient.ConnectionUsers(connection.AccountID, connection.ApplicationID, connection.UserFromID)

	userKey := storageClient.User(connection.AccountID, connection.ApplicationID, connection.UserToID)

	if err = storageEngine.LPush(userListKey, userKey).Err(); err != nil {
		return nil, err
	}

	followerListKey := storageClient.FollowedByUsers(connection.AccountID, connection.ApplicationID, connection.UserToID)

	followerKey := storageClient.User(connection.AccountID, connection.ApplicationID, connection.UserFromID)

	if err = storageEngine.LPush(followerListKey, followerKey).Err(); err != nil {
		return nil, err
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
func WriteConnectionEventsToList(connection *entity.Connection) (err error) {
	connectionEventsKey := storageClient.ConnectionEvents(connection.AccountID, connection.ApplicationID, connection.UserFromID)

	eventsKey := storageClient.Events(connection.AccountID, connection.ApplicationID, connection.UserToID)

	events, err := storageEngine.ZRevRangeWithScores(eventsKey, "0", "-1").Result()
	if err != nil {
		return err
	}

	if len(events) >= 1 {
		var vals []red.Z

		for _, eventKey := range events {
			val := red.Z{Score: float64(eventKey.Score), Member: eventKey.Member}
			vals = append(vals, val)
		}

		if err = storageEngine.ZAdd(connectionEventsKey, vals...).Err(); err != nil {
			return err
		}
	}

	return nil
}

// DeleteConnectionEventsFromLists takes a connection and deletes the events from the lists
func DeleteConnectionEventsFromLists(accountID, applicationID, userFromID, userToID int64) (err error) {
	connectionEventsKey := storageClient.ConnectionEvents(accountID, applicationID, userFromID)

	eventsKey := storageClient.Events(accountID, applicationID, userToID)

	events, err := storageEngine.ZRevRangeWithScores(eventsKey, "0", "-1").Result()
	if err != nil {
		return err
	}

	if len(events) >= 1 {
		var members []string

		for _, eventKey := range events {
			member := eventKey.Member
			members = append(members, member)
		}

		if err = storageEngine.ZRem(connectionEventsKey, members...).Err(); err != nil {
			return err
		}
	}

	return nil
}

// ReadConnection returns the connection, if any, between two users
func ReadConnection(accountID, applicationID, userFromID, userToID int64) (connection *entity.Connection, err error) {
	key := storageClient.Connection(accountID, applicationID, userFromID, userToID)

	result, err := storageEngine.Get(key).Result()
	if err != nil {
		return
	}
	if result == "" {
		return nil, nil
	}

	connection = &entity.Connection{}
	err = json.Unmarshal([]byte(result), connection)
	return
}

// SocialConnect creates the connections between a user and his other social peers
func SocialConnect(user *entity.User, platform string, socialFriendsIDs []string) ([]*entity.User, error) {
	result := []*entity.User{}

	encodedSocialFriendsIDs := []string{}
	for idx := range socialFriendsIDs {
		encodedSocialFriendsIDs = append(encodedSocialFriendsIDs, storageClient.SocialConnection(
			user.AccountID,
			user.ApplicationID,
			platform,
			Base64Encode(socialFriendsIDs[idx])))
	}

	ourStoredUsersIDs, err := storageEngine.MGet(encodedSocialFriendsIDs...).Result()
	if err != nil {
		return result, err
	}

	if len(ourStoredUsersIDs) == 0 {
		return result, nil
	}

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

		_, err = WriteConnection(connection, false)
		if err != nil {
			continue
		}

		_, err = ConfirmConnection(connection, false)
		if err != nil {
			continue
		}

		ourUserKeys = append(
			ourUserKeys,
			storageClient.User(user.AccountID, user.ApplicationID, userID),
		)
	}

	return fetchAndDecodeMultipleUsers(ourUserKeys)
}

func fetchAndDecodeMultipleUsers(keys []string) (users []*entity.User, err error) {
	fmt.Scan()
	if len(keys) == 0 {
		return []*entity.User{}, nil
	}

	resultList, err := storageEngine.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	user := &entity.User{}
	for _, result := range resultList {
		if err = json.Unmarshal([]byte(result.(string)), user); err != nil {
			return nil, err
		}
		users = append(users, user)
		user = &entity.User{}
	}

	return
}
