/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package redis provides the redis needed functions for redis
package redis

import (
	"gopkg.in/redis.v2"
)

type (
	cli struct {
		client *redis.Client
	}
)

var redisClient *cli

// Init initializes the redis client
func Init(address, password string, db int64, poolSize int) {
	options := &redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		PoolSize: poolSize,
	}

	redisClient = &cli{
		client: redis.NewTCPClient(options),
	}
}

// Client returns the redis client
func Client() *redis.Client {
	return redisClient.client
}
