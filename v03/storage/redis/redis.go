// Package redis provides a quick abstraction layer for common used redis functions
package redis

import (
	"time"

	"github.com/tapglue/backend/config"

	redigo "github.com/garyburd/redigo/redis"
)

// NewRedigoPool creates a new redis connection pool using redigo driver
func NewRedigoPool(conf *config.Redis) *redigo.Pool {
	return &redigo.Pool{
		MaxIdle:     conf.PoolSize,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redigo.Conn, error) {
			c, err := redigo.Dial("tcp", conf.Hosts[0])
			if err != nil {
				return nil, err
			}
			if conf.Password != "" {
				if _, err := c.Do("AUTH", conf.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
