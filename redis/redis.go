/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package redis provides the redis needed functions for redis
package redis

import (
	//"net"

	"gopkg.in/redis.v2"
)

var client *redis.Client

// Init initializes the redis client
func Init() {
	options := &redis.Options{
		Addr: "127.0.0.1:6379",
		/*Dialer: func() (net.Conn, error) {
			return net.Dial("tcp", "127.0.0.1:6379")
		},*/
		Password: "",
		DB:       0,
		PoolSize: 30,
	}
	client = redis.NewTCPClient(options)
}

func GetClient() *redis.Client {
	return client
}
