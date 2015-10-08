package redis

import (
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/tapglue/multiverse/limiter"
)

func TestLimiter(t *testing.T) {
	var (
		pool = redis.NewPool(func() (redis.Conn, error) {
			return redis.Dial("tcp", "127.0.0.1:6379")
		}, 10)
		limitee = &limiter.Limitee{
			Hash:       "token",
			Limit:      10,
			WindowSize: 1 * time.Second,
		}
		l = NewLimiter(pool, "limitertest")
	)

	limit, _, err := l.Request(limitee)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	for i := 0; i < 20; i++ {
		_, _, err := l.Request(limitee)
		if err != nil {
			t.Fatalf("request failed: %s", err)
		}
	}

	limit, _, err = l.Request(limitee)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if have, want := limit, int64(0); have != want {
		t.Errorf("have %v, want %v", have, want)
	}

	time.Sleep(1 * time.Second)

	limit, _, err = l.Request(limitee)
	if err != nil {
		t.Fatalf("request failed: %s", err)
	}

	if have, want := limit, int64(9); have != want {
		t.Errorf("have %v, want %v", have, want)
	}
}
