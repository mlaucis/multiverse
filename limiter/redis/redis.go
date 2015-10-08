// Package redis provides a Redis implementation of the rate limiting interfaces
// NOTE: As of 13.05.2015 this is not a very strict implementation, meaning that with
// sufficient concurrency levels / high latency the limits will be broken.
// Also, limit creation point might not be as accurate as possible.
// TODO Have proper locks in places so that this doesn't happen
package redis

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tapglue/multiverse/limiter"

	"github.com/garyburd/redigo/redis"
)

type rateLimiter struct {
	bucket   string
	connPool *redis.Pool
}

// NewLimiter creates a new Limiter implementation using Redis
func NewLimiter(connPool *redis.Pool, bucketName string) limiter.Limiter {
	return &rateLimiter{
		bucket:   bucketName,
		connPool: connPool,
	}
}

func (rateLimiter *rateLimiter) Request(limitee *limiter.Limitee) (int64, time.Time, error) {
	var (
		key     = fmt.Sprintf("%s:%s", rateLimiter.bucket, limitee.Hash)
		conn    = rateLimiter.connPool.Get()
		expires = time.Now().Add(time.Duration(limitee.WindowSize))
		left    = int64(-1)
	)
	defer conn.Close()

	res, err := conn.Do("GET", key)
	if err != nil {
		return 0, time.Now(), err
	}

	if res != nil {
		left, err = strconv.ParseInt(string(res.([]uint8)), 10, 64)
		if err != nil {
			return 0, expires, fmt.Errorf("parsing counter failed: %s", err)
		}

		expiry, err := conn.Do("TTL", key)
		if err != nil {
			return 0, time.Now(), err
		}
		expires = time.Now().Add(time.Duration(expiry.(int64)) * time.Second)
	}

	if left == -1 {
		_, err := conn.Do("SET", key, limitee.Limit-1, "EX", limitee.WindowSize, "NX")
		if err != nil {
			return 0, expires, err
		}

		return limitee.Limit - 1, expires, nil
	}

	if left > 0 {
		_, err = conn.Do("DECR", key)
		if err != nil {
			return 0, expires, err
		}

		return left - 1, expires, nil
	}

	return left, expires, nil
}
