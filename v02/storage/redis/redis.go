/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package redis provides the redis needed functions for redis
package redis

import (
	"fmt"

	storageHelper "github.com/tapglue/backend/v02/storage/helper"

	"errors"

	red "gopkg.in/redis.v2"
)

type (
	// Client structure
	Client interface {
		// GenerateAccountID generates a new account ID
		GenerateAccountID() (int64, error)

		// GenerateAccountUserID generates a new account user id for a specified account
		GenerateAccountUserID(accountID int64) (int64, error)

		// GenerateApplicationID generates a new application ID
		GenerateApplicationID(accountID int64) (int64, error)

		// GenerateApplicationUserID generates the user id in the specified app
		GenerateApplicationUserID(applicationID int64) (string, error)

		// GenerateApplicationEventID generates the event id in the specified app
		GenerateApplicationEventID(applicationID int64) (string, error)

		// Datastore returns the client datastore
		Datastore() *red.Client
	}

	cli struct {
		datastore *red.Client
	}
)

func (c *cli) GenerateAccountID() (int64, error) {
	return c.datastore.Incr(storageHelper.IDAccount).Result()
}

func (c *cli) GenerateAccountUserID(accountID int64) (int64, error) {
	return c.datastore.Incr(fmt.Sprintf(storageHelper.IDAccountUser, accountID)).Result()
}

func (c *cli) GenerateApplicationID(accountID int64) (int64, error) {
	return c.datastore.Incr(fmt.Sprintf(storageHelper.IDAccountApp, accountID)).Result()
}

func (c *cli) GenerateApplicationUserID(applicationID int64) (string, error) {
	return "", errors.New("needs a new implementation")
	//return c.datastore.Incr(fmt.Sprintf(storageHelper.IDApplicationUser, applicationID)).Result()
}

func (c *cli) GenerateApplicationEventID(applicationID int64) (string, error) {
	return "", errors.New("needs a new implementation")
	//return c.datastore.Incr(fmt.Sprintf(storageHelper.IDApplicationEvent, applicationID)).Result()
}

func (c *cli) Datastore() *red.Client {
	return c.datastore
}

// New initializes the redis client
func New(address, password string, db int64, poolSize int) Client {
	options := &red.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	return &cli{
		datastore: red.NewTCPClient(options),
	}
}
