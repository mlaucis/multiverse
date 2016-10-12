package event

import (
	"fmt"
	"strings"
	"time"

	"github.com/tapglue/multiverse/platform/cache"
	"github.com/tapglue/multiverse/platform/metrics"
)

const (
	cachePrefixCount = "events.count"
)

type cacheService struct {
	countCache cache.CountService
	next       Service
}

// CacheServiceMiddleware adds caching capabilities to the Service by using
// read-through and write-through methods to store results of heavy computation
// with sensible TTLs.
func CacheServiceMiddleware(countCache cache.CountService) ServiceMiddleware {
	return func(next Service) Service {
		return &cacheService{
			countCache: countCache,
			next:       next,
		}
	}
}

func (s *cacheService) ActiveUserIDs(ns string, p Period) (ids []uint64, err error) {
	return s.next.ActiveUserIDs(ns, p)
}

func (s *cacheService) Count(ns string, opts QueryOptions) (int, error) {
	var (
		key = cacheKey(opts)
	)

	count, err := s.countCache.Get(ns, key)
	if err == nil {
		return count, nil
	}

	if !cache.IsKeyNotFound(err) {
		return -1, err
	}

	count, err = s.next.Count(ns, opts)
	if err != nil {
		return -1, err
	}

	err = s.countCache.Set(ns, key, count)

	return count, err
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

	return strings.Join(ps, cache.KeySeparator)
}
