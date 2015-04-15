/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package redis provides the redis needed functions for redis
package redis

import (
	"fmt"

	storageHelper "github.com/tapglue/backend/v02/storage/helper"

	red "gopkg.in/redis.v2"
)

type (
	// Client structure
	Client struct {
		datastore *red.Client
	}
)

// GenerateAccountID generates a new account ID
func (c *Client) GenerateAccountID() (int64, error) {
	return c.datastore.Incr(storageHelper.IDAccount).Result()
}

// GenerateAccountUserID generates a new account user id for a specified account
func (c *Client) GenerateAccountUserID(accountID int64) (int64, error) {
	return c.datastore.Incr(fmt.Sprintf(storageHelper.IDAccountUser, accountID)).Result()
}

// GenerateApplicationID generates a new application ID
func (c *Client) GenerateApplicationID(accountID int64) (int64, error) {
	return c.datastore.Incr(fmt.Sprintf(storageHelper.IDAccountApp, accountID)).Result()
}

// GenerateApplicationUserID generates the user id in the specified app
func (c *Client) GenerateApplicationUserID(applicationID int64) (int64, error) {
	return c.datastore.Incr(fmt.Sprintf(storageHelper.IDApplicationUser, applicationID)).Result()
}

// GenerateApplicationEventID generates the event id in the specified app
func (c *Client) GenerateApplicationEventID(applicationID int64) (int64, error) {
	return c.datastore.Incr(fmt.Sprintf(storageHelper.IDApplicationEvent, applicationID)).Result()
}

// Datastore returns the client datastore
func (c *Client) Datastore() *red.Client {
	return c.datastore
}

// New initializes the redis client
func New(address, password string, db int64, poolSize int) *Client {
	options := &red.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	return &Client{
		datastore: red.NewTCPClient(options),
	}
}
