package redis_test

import (
	"testing"
	"time"

	rLimiter "github.com/tapglue/multiverse/limiter"
	"github.com/tapglue/multiverse/limiter/redis"

	redigo "github.com/garyburd/redigo/redis"
	. "gopkg.in/check.v1"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) {
	t.Parallel()
	TestingT(t)
}

type RedisSuite struct{}

var _ = Suite(&RedisSuite{})

func (s *RedisSuite) TestNegativeTTL(c *C) {
	connPool := redigo.NewPool(func() (redigo.Conn, error) {
		return redigo.Dial("tcp", "localhost:6379")
	}, 5)
	limiter := redis.NewLimiter(connPool, "demo")

	var (
		redisKey       = "demodemo2015-10-06 21:18:59.985307449 +0200 CEST"
		limit    int64 = 5
		i        int64 = 0
	)
	limitee := rLimiter.Limitee{
		Hash:       redisKey,
		Limit:      limit,
		WindowSize: 2,
	}
	x, _, er := limiter.Request(&limitee)
	c.Assert(er, IsNil)
	c.Assert(x, Equals, limit-1)
	for i = 0; i < limit; i++ {
		x, _, er := limiter.Request(&limitee)
		c.Assert(er, IsNil)
		c.Assert(x, Equals, limit-i-1)
	}

	// Wait for the lock to expire
	time.Sleep(2 * time.Second)

	limitee.Limit = limit
	x, _, er = limiter.Request(&limitee)
	c.Assert(er, IsNil)
	c.Assert(x, Equals, limit-1)
	for i = 0; i < limit; i++ {
		x, _, er := limiter.Request(&limitee)
		c.Assert(er, IsNil)
		c.Assert(x, Equals, limit-i-1)
	}
}
