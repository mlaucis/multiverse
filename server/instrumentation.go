package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-kit/kit/metrics"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace    = "api"
	subsystem    = "intaker"
	fieldKeys    = []string{"route", "api_version", "status", "user_agent"}
	requestCount = kitprometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_count",
		Help:      "Number of requests received",
	}, fieldKeys)
	requestLatency = metrics.NewTimeHistogram(
		time.Microsecond,
		kitprometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Subsystem: subsystem,
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds",
			},
			fieldKeys,
		),
	)
	responseBytes = kitprometheus.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "response_bytes",
		Help:      "Bytes returned as response bodies",
	}, fieldKeys)
)

func metricHandler(route, apiVersion string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			begin = time.Now()
			rc    = &responseRecorder{ResponseWriter: w}
		)

		next.ServeHTTP(rc, r)

		var (
			routeNameField = metrics.Field{
				Key:   "route",
				Value: route,
			}
			routeVersionField = metrics.Field{
				Key:   "api_version",
				Value: apiVersion,
			}
			statusField = metrics.Field{
				Key:   "status",
				Value: strconv.Itoa(rc.statusCode),
			}
			uaField = metrics.Field{
				Key:   "user_agent",
				Value: r.Header.Get("User-Agent"),
			}
		)

		requestCount.With(routeNameField).With(routeVersionField).With(statusField).With(uaField).Add(1)
		requestLatency.With(routeNameField).With(routeVersionField).With(statusField).With(uaField).Observe(time.Since(begin))
		responseBytes.With(routeNameField).With(routeVersionField).With(statusField).With(uaField).Add(uint64(rc.bytesWritten))
	}
}

type responseRecorder struct {
	http.ResponseWriter

	bytesWritten int
	statusCode   int
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	n, err := r.ResponseWriter.Write(b)
	r.bytesWritten += n
	return n, err
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}
