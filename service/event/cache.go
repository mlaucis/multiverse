package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/tapglue/multiverse/platform/metrics"
)

const (
	cacheKeySeparator = "."
	cachePrefixCount  = "cache.events.count"
	cacheTTLDefault   = 300

	commandEX  = "EX"
	commandGET = "GET"
	commandSET = "SET"
)

type cacheService struct {
	next Service
	pool *redis.Pool
}

// CacheServiceMiddleware adds caching capabilities to the Service by using
// read-through and write-through methods to store results of heavy computation
// with sensible TTLs.
func CacheServiceMiddleware(pool *redis.Pool) ServiceMiddleware {
	return func(next Service) Service {
		return &cacheService{
			next: next,
			pool: pool,
		}
	}
}

func (s *cacheService) ActiveUserIDs(ns string, p Period) (ids []uint64, err error) {
	return s.next.ActiveUserIDs(ns, p)
}

func (s *cacheService) Count(ns string, opts QueryOptions) (int, error) {
	var (
		count = 0
		con   = s.pool.Get()
		key   = cacheKey(opts)
	)
	defer con.Close()

	res, err := con.Do(commandGET, key)
	if err != nil {
		return 0, fmt.Errorf("cache get failed: %s", err)
	}

	if res == nil {
		count, err = s.next.Count(ns, opts)
		if err != nil {
			return 0, err
		}

		_, err = con.Do(commandSET, key, uint64(count), commandEX, cacheTTLDefault)
		if err != nil {
			return 0, fmt.Errorf("cache set failed: %s", err)
		}
	} else {
		var c uint64 = 0

		_, err = redis.Scan([]interface{}{res}, &c)
		if err != nil {
			return 0, fmt.Errorf("cache scan failed: %s", err)
		}

		count = int(c)
	}

	return count, nil
}

func (s *cacheService) CreatedByDay(
	ns string,
	start, end time.Time,
) (ts metrics.Timeseries, err error) {
	return s.next.CreatedByDay(ns, start, end)
}

func (s *cacheService) Put(ns string, input *Event) (output *Event, err error) {
	return s.next.Put(ns, input)
}

func (s *cacheService) Query(ns string, opts QueryOptions) (list List, err error) {
	return s.next.Query(ns, opts)
}

func (s *cacheService) Setup(ns string) (err error) {
	return s.next.Setup(ns)
}

func (s *cacheService) Teardown(ns string) (err error) {
	return s.next.Teardown(ns)
}

func cacheKey(opts QueryOptions) string {
	ps := []string{
		cachePrefixCount,
	}

	if len(opts.Types) == 1 {
		ps = append(ps, opts.Types[0])
	}

	if len(opts.ObjectIDs) == 1 {
		ps = append(ps, fmt.Sprintf("%d", opts.ObjectIDs[0]))
	}

	return strings.Join(ps, cacheKeySeparator)
}
