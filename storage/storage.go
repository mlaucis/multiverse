/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package storage holds common functions regardless of the storage engine used
package storage

import (
	"fmt"

	red "gopkg.in/redis.v2"
)

type (
	// Client structure holds the storage engine and functions needed to operate the backend
	Client struct {
		engine *red.Client
	}
)

// Defining keys
const (
	idAccountKey         = "ids:accounts"
	idAccountUserKey     = "ids:account:%d:user"
	idAccountAppKey      = "ids:account:%d:application"
	idApplicationUserKey = "ids:application:%d:user"
	idApplicationEvent   = "ids:application:%d:events"

	accountKey = "account:%d"

	accountUserKey  = "account:%d:user:%d"
	accountUsersKey = "account:%d:users"

	accountAppKey  = "account:%d:app:%d"
	accountAppsKey = "account:%d:apps"

	connectionKey      = "app:%d:user:%d:follows:%d"
	connectionsKey     = "app:%d:user:%d:connections"
	connectionUsersKey = "app:%d:user:%d:follows_users"
	followedByUsersKey = "app:%d:user:%d:followed_by_users"

	userKey  = "app:%d:user:%d"
	usersKey = "app:%d:users"

	eventKey  = "app:%d:user:%d:event_%d"
	eventsKey = "app:%d:user:%d:events"

	connectionEventsKey     = "app:%d:user:%d:connection_events"
	connectionEventsKeyLoop = "%s:connection_events"
)

var (
	instance *Client
)

// GenerateAccountID generates a new account ID
func (client *Client) GenerateAccountID() (int64, error) {
	return client.engine.Incr(idAccountKey).Result()
}

// GenerateAccountUserID generates a new account user id for a specified account
func (client *Client) GenerateAccountUserID(accountID int64) (int64, error) {
	return client.engine.Incr(fmt.Sprintf(idAccountUserKey, accountID)).Result()
}

// GenerateApplicationID generates a new application ID
func (client *Client) GenerateApplicationID(accountID int64) (int64, error) {
	return client.engine.Incr(fmt.Sprintf(idAccountAppKey, accountID)).Result()
}

// GenerateApplicationUserID generates the user id in the specified app
func (client *Client) GenerateApplicationUserID(applicationID int64) string {
	return fmt.Sprintf(idApplicationUserKey, applicationID)
}

// GenerateApplicationEventID generates the event id in the specified app
func (client *Client) GenerateApplicationEventID(applicationID int64) string {
	return fmt.Sprintf(idApplicationEvent, applicationID)
}

// AccountKey returns the key for a specified account
func (client *Client) AccountKey(accountID int64) string {
	return fmt.Sprintf(accountKey, accountID)
}

// AccountUserKey returns the key for a specific user of an account
func (client *Client) AccountUserKey(accountID, accountUserID int64) string {
	return fmt.Sprintf(accountUserKey, accountID, accountUserID)
}

// AccountUsersKey returns the key for account users
func (client *Client) AccountUsersKey(accountID int64) string {
	return fmt.Sprintf(accountUsersKey, accountID)
}

// AccountAppKey returns the key for one account app
func (client *Client) AccountAppKey(accountID, applicationID int64) string {
	return fmt.Sprintf(accountAppKey, accountID, applicationID)
}

// AccountAppsKey returns the key for one account app
func (client *Client) AccountAppsKey(accountID int64) string {
	return fmt.Sprintf(accountAppsKey, accountID)
}

// ConnectionKey gets the key for the connection
func (client *Client) ConnectionKey(applicationID, userFromID, userToID int64) string {
	return fmt.Sprintf(connectionKey, applicationID, userFromID, userToID)
}

// ConnectionsKey replace this
func (client *Client) ConnectionsKey(applicationID, userFromID int64) string {
	return fmt.Sprintf(connectionsKey, applicationID, userFromID)
}

// ConnectionUsersKey replace this
func (client *Client) ConnectionUsersKey(applicationID, userFromID int64) string {
	return fmt.Sprintf(connectionUsersKey, applicationID, userFromID)
}

// FollowedByUsersKey replace this
func (client *Client) FollowedByUsersKey(applicationID, userToID int64) string {
	return fmt.Sprintf(followedByUsersKey, applicationID, userToID)
}

// UserKey gets the key for the user
func (client *Client) UserKey(applicationID, userID int64) string {
	return fmt.Sprintf(userKey, applicationID, userID)
}

// UsersKey replace this
func (client *Client) UsersKey(applicationID int64) string {
	return fmt.Sprintf(usersKey, applicationID)
}

// EventKey replace this
func (client *Client) EventKey(applicationID, userID, eventID int64) string {
	return fmt.Sprintf(eventKey, applicationID, userID, eventID)
}

// EventsKey replace this
func (client *Client) EventsKey(applicationID, userID int64) string {
	return fmt.Sprintf(eventsKey, applicationID, userID)
}

// ConnectionEventsKey replace this
func (client *Client) ConnectionEventsKey(applicationID, userID int64) string {
	return fmt.Sprintf(connectionEventsKey, applicationID, userID)
}

// ConnectionEventsKeyLoop replace this
func (client *Client) ConnectionEventsKeyLoop(userID string) string {
	return fmt.Sprintf(connectionEventsKeyLoop, userID)
}

// Engine returns the storage engine used
func (client *Client) Engine() *red.Client {
	return client.engine
}

// Init initializes the storage package with the required storage engine
func Init(engine *red.Client) *Client {
	if instance == nil {
		instance = &Client{
			engine: engine,
		}
	}

	return instance
}
