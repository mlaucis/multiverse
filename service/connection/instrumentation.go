package connection

import (
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/tapglue/multiverse/errors"
)

const (
	fieldApp    = "app"
	fieldMethod = "method"
	fieldStore  = "store"
)

type instrumentStrangleService struct {
	StrangleService

	errCount  metrics.Counter
	opCount   metrics.Counter
	opLatency metrics.TimeHistogram
	store     string
}

// InstrumentMiddleware observes key aspects of Service operations and exposes
// Prometheus metrics.
func InstrumentMiddleware(ns string, store string) StrangleMiddleware {
	var (
		fieldKeys = []string{fieldApp, fieldMethod, fieldStore}
		namespace = strings.Replace(ns, "-", "_", -1)
		subsytem  = "service_connection"

		errCount = kitprometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsytem,
			Name:      "err_count",
			Help:      "Number of failed operations",
		}, fieldKeys)
		opCount = kitprometheus.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsytem,
			Name:      "op_count",
			Help:      "Number of operations performed",
		}, fieldKeys)
		opLatency = metrics.NewTimeHistogram(
			time.Second,
			kitprometheus.NewHistogram(
				prometheus.HistogramOpts{
					Namespace: namespace,
					Subsystem: subsytem,
					Name:      "op_latency_seconds",
					Help:      "Distribution of op duration in seconds",
				},
				fieldKeys,
			),
		)
	)

	return func(next StrangleService) StrangleService {
		return &instrumentStrangleService{
			errCount:        errCount,
			opCount:         opCount,
			opLatency:       opLatency,
			StrangleService: next,
			store:           store,
		}
	}
}

func (s *instrumentStrangleService) FriendsAndFollowingIDs(orgID, appID int64, id uint64) (ids []uint64, errs []errors.Error) {
	defer func(begin time.Time) {
		var (
			a = metrics.Field{
				Key:   fieldApp,
				Value: strconv.FormatInt(appID, 10),
			}
			m = metrics.Field{
				Key:   fieldMethod,
				Value: "FriendsAndFollowingIDs",
			}
			store = metrics.Field{
				Key:   fieldStore,
				Value: s.store,
			}
		)

		if errs != nil {
			s.errCount.With(a).With(m).With(store).Add(1)
		}

		s.opCount.With(a).With(m).With(store).Add(1)
		s.opLatency.With(a).With(m).With(store).Observe(time.Since(begin))
	}(time.Now())

	return s.StrangleService.FriendsAndFollowingIDs(orgID, appID, id)
}
