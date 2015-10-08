// Package redis provides a Redis implementation of the rate limiting interfaces
package redis

import (
	"fmt"
	"time"

	"github.com/tapglue/multiverse/limiter"

	"github.com/garyburd/redigo/redis"
)

type rateLimiter struct {
	prefix string
	pool   *redis.Pool
}

// NewLimiter returns a Redis Limiter implementation.
func NewLimiter(pool *redis.Pool, prefix string) limiter.Limiter {
	return &rateLimiter{
		prefix: prefix,
		pool:   pool,
	}
}

func (rateLimiter *rateLimiter) Request(limitee *limiter.Limitee) (int64, time.Time, error) {
	var (
		conn    = rateLimiter.pool.Get()
		expires = time.Now().Add(limitee.WindowSize)
		key     = fmt.Sprintf("%s:%s", rateLimiter.prefix, limitee.Hash)
	)
	defer conn.Close()

	quota, err := getQuota(conn, key)
	if err != nil {
		return 0, time.Now(), err
	}

	ttl, err := getTTL(conn, key)
	if err != nil {
		return 0, time.Now(), err
	}

	if ttl < 0 {
		quota = limitee.Limit - 1

		_, err := conn.Do("SET", key, quota, "EX", uint64(limitee.WindowSize/time.Second))
		if err != nil {
			return 0, time.Now(), err
		}

		return quota, expires, nil
	}

	return quota, time.Now().Add(ttl), nil
}

func getQuota(conn redis.Conn, key string) (int64, error) {
	// DECR on non-existent keys will set them to `-1` we can make use of that to
	// determine if we have to reset the quota.
	res, err := conn.Do("DECR", key)
	if err != nil {
		return 0, err
	}

	return res.(int64), nil
}

func getTTL(conn redis.Conn, key string) (time.Duration, error) {
	// TTL returns -2 for a key that doesn't exist and -1 if none is set.
	res, err := conn.Do("TTL", key)
	if err != nil {
		return 0, err
	}

	return time.Duration(res.(int64)) * time.Second, nil
}
