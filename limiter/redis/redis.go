/**
 * @author Florin Patan <florinpatan@gmail.com>
 */

// Package redis provides a Redis implementation of the rate limiting interfaces
// NOTE: As of 13.05.2015 this is not a very strict implementation, meaning that with
// sufficient concurrency levels / high latency the limits will be broken.
// Also, limit creation point might not be as accurate as possible.
// TODO Have proper locks in places so that this doesn't happen
package redis

import (
	"errors"
	"strconv"
	"time"

	"github.com/tapglue/backend/limiter"

	red "github.com/garyburd/redigo/redis"
)

type (
	rateLimiter struct {
		bucket string
		conn   red.Conn
	}
)

var (
	errTime = time.Date(2015, 5, 1, 1, 2, 3, 4, time.UTC)
)

func (rateLimiter *rateLimiter) Request(limitee *limiter.Limitee) (int64, time.Time, error) {
	hash := rateLimiter.bucket + limitee.Hash

	remaining, err := rateLimiter.conn.Do("GET", hash)
	if err != nil {
		return 0, errTime, err
	}

	if remaining == nil {
		return rateLimiter.create(limitee)
	}

	left, err := strconv.ParseInt(string(remaining.([]uint8)), 10, 64)
	if err != nil {
		return 0, errTime, errors.New("something went wrong")
	}

	if left > 0 {
		return rateLimiter.decrement(limitee, left)
	}

	if left <= 0 {
		return rateLimiter.expiresIn(limitee)
	}

	return 0, errTime, errors.New("something went wrong")
}

func (rateLimiter *rateLimiter) decrement(limitee *limiter.Limitee, value int64) (int64, time.Time, error) {
	hash := rateLimiter.bucket + limitee.Hash
	expiry, err := rateLimiter.conn.Do("TTL", hash)
	if err != nil {
		return 0, errTime, err
	}

	_, err = rateLimiter.conn.Do("DECR", hash)
	if err != nil {
		return 0, errTime, err
	}

	return value - 1, time.Now().Add(time.Duration(expiry.(int64)) * time.Second), nil
}

func (rateLimiter *rateLimiter) create(limitee *limiter.Limitee) (int64, time.Time, error) {
	hash := rateLimiter.bucket + limitee.Hash
	limit := limitee.Limit
	response, err := rateLimiter.conn.Do("SET", hash, limit, "EX", limitee.WindowSize, "NX")
	if err != nil {
		return 0, errTime, err
	}

	// Check if this was set by someone else meanwhile
	if response == nil {
		return rateLimiter.decrement(limitee, limitee.Limit)
	}

	return limit - 1, time.Now().Add(time.Duration(limitee.WindowSize) * time.Second), nil
}

func (rateLimiter *rateLimiter) expiresIn(limitee *limiter.Limitee) (int64, time.Time, error) {
	hash := rateLimiter.bucket + limitee.Hash
	expiry, err := rateLimiter.conn.Do("TTL", hash)
	if err != nil {
		return 0, errTime, err
	}

	return 0, time.Now().Add(time.Duration(expiry.(int64)) * time.Second), nil
}

// NewLimiter creates a new Limiter implementation using Redis
func NewLimiter(conn red.Conn, bucketName string) limiter.Limiter {
	return &rateLimiter{
		bucket: bucketName,
		conn:   conn,
	}
}