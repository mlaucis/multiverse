/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package storage holds common functions regardless of the storage engine used
package storage

import (
	"encoding/base64"
	"fmt"

	"github.com/tapglue/backend/core/entity"

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
	idAccount          = "ids:acs"
	idAccountUser      = "ids:ac:%d:u"
	idAccountApp       = "ids:ac:%d:a"
	idApplicationUser  = "ids:a:%d:u"
	idApplicationEvent = "ids:a:%d:e"

	account = "acc:%d"

	accountUser  = "acc:%d:user:%d"
	accountUsers = "acc:%d:users"

	application  = "acc:%d:app:%d"
	applications = "acc:%d:apps"

	user  = "acc:%d:app:%d:user:%d"
	users = "acc:%d:app:%d:user"

	connection      = "acc:%d:app:%d:user:%d:connection:%d"
	connections     = "acc:%d:app:%d:user:%d:connections"
	followsUsers    = "acc:%d:app:%d:user:%d:followsUsers"
	followedByUsers = "acc:%d:app:%d:user:%d:follwedByUsers"

	event  = "acc:%d:app:%d:user:%d:event:%d"
	events = "acc:%d:app:%d:user:%d:events"

	connectionEvents     = "acc:%d:app:%d:user:%d:connectionEvents"
	connectionEventsLoop = "%s:connectionEvents"
)

var (
	instance *Client
)

// GenerateAccountID generates a new account ID
func (client *Client) GenerateAccountID() (int64, error) {
	return client.engine.Incr(idAccount).Result()
}

// GenerateAccountToken returns a token for the specified account
func (client *Client) GenerateAccountToken(account *entity.Account) (string, error) {
	return fmt.Sprintf(
		"token_%d_%s",
		account.ID,
		base64Encode(account.Name),
	), nil
}

// GenerateAccountUserID generates a new account user id for a specified account
func (client *Client) GenerateAccountUserID(accountID int64) (int64, error) {
	return client.engine.Incr(fmt.Sprintf(idAccountUser, accountID)).Result()
}

// GenerateApplicationID generates a new application ID
func (client *Client) GenerateApplicationID(accountID int64) (int64, error) {
	return client.engine.Incr(fmt.Sprintf(idAccountApp, accountID)).Result()
}

// GenerateApplicationToken returns a token for the specified application of an account
func (client *Client) GenerateApplicationToken(application *entity.Application) (string, error) {
	return fmt.Sprintf(
		"token_%d_%d_%s",
		application.AccountID,
		application.ID,
		base64Encode(application.Name),
	), nil
}

// GenerateApplicationUserID generates the user id in the specified app
func (client *Client) GenerateApplicationUserID(applicationID int64) string {
	return fmt.Sprintf(idApplicationUser, applicationID)
}

// GenerateApplicationEventID generates the event id in the specified app
func (client *Client) GenerateApplicationEventID(applicationID int64) string {
	return fmt.Sprintf(idApplicationEvent, applicationID)
}

// Account returns the key for a specified account
func (client *Client) Account(accountID int64) string {
	return fmt.Sprintf(account, accountID)
}

// AccountUser returns the key for a specific user of an account
func (client *Client) AccountUser(accountID, accountUserID int64) string {
	return fmt.Sprintf(accountUser, accountID, accountUserID)
}

// AccountUsers returns the key for account users
func (client *Client) AccountUsers(accountID int64) string {
	return fmt.Sprintf(accountUsers, accountID)
}

// Application returns the key for one account app
func (client *Client) Application(accountID, applicationID int64) string {
	return fmt.Sprintf(application, accountID, applicationID)
}

// Applications returns the key for one account app
func (client *Client) Applications(accountID int64) string {
	return fmt.Sprintf(applications, accountID)
}

// Connection gets the key for the connection
func (client *Client) Connection(accountID, applicationID, userFromID, userToID int64) string {
	return fmt.Sprintf(connection, accountID, applicationID, userFromID, userToID)
}

// Connections gets the key for the connections list
func (client *Client) Connections(accountID, applicationID, userFromID int64) string {
	return fmt.Sprintf(connections, accountID, applicationID, userFromID)
}

// ConnectionUsers gets the key for the connectioned users list
func (client *Client) ConnectionUsers(accountID, applicationID, userFromID int64) string {
	return fmt.Sprintf(followsUsers, accountID, applicationID, userFromID)
}

// FollowedByUsers gets the key for the list of followers
func (client *Client) FollowedByUsers(accountID, applicationID, userToID int64) string {
	return fmt.Sprintf(followedByUsers, accountID, applicationID, userToID)
}

// User gets the key for the user
func (client *Client) User(accountID, applicationID, userID int64) string {
	return fmt.Sprintf(user, accountID, applicationID, userID)
}

// Users gets the key the app users list
func (client *Client) Users(accountID, applicationID int64) string {
	return fmt.Sprintf(users, accountID, applicationID)
}

// Event gets the key for an event
func (client *Client) Event(accountID, applicationID, userID, eventID int64) string {
	return fmt.Sprintf(event, accountID, applicationID, userID, eventID)
}

// Events get the key for the events list
func (client *Client) Events(accountID, applicationID, userID int64) string {
	return fmt.Sprintf(events, accountID, applicationID, userID)
}

// ConnectionEvents get the key for the connections events list
func (client *Client) ConnectionEvents(accountID, applicationID, userID int64) string {
	return fmt.Sprintf(connectionEvents, accountID, applicationID, userID)
}

// ConnectionEventsLoop gets the key for looping through connections
func (client *Client) ConnectionEventsLoop(userID string) string {
	return fmt.Sprintf(connectionEventsLoop, userID)
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

func base64Encode(value string) string {
	return base64.StdEncoding.EncodeToString([]byte(value))
}
