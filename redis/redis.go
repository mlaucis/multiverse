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
func Init(address string) {
	options := &redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
		PoolSize: 30,
	}

	redisClient = &cli{
		client: redis.NewTCPClient(options),
	}
}

// Client returns the redis client
func Client() *redis.Client {
	return redisClient.client
}
