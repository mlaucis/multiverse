package server

import (
	"net/http"
	"strconv"

	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	namespace    = "api"
	subsystem    = "intaker"
	fieldKeys    = []string{"route", "api_version", "status"}
	requestCount = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "request_count",
		Help:      "Number of requests received",
	}, fieldKeys)
	responseBytes = kitprometheus.NewCounterFrom(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystem,
		Name:      "response_bytes",
		Help:      "Bytes returned as response bodies",
	}, fieldKeys)
)

func metricHandler(route, apiVersion string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rc := &responseRecorder{ResponseWriter: w}

		next.ServeHTTP(rc, r)

		status := strconv.Itoa(rc.statusCode)

		requestCount.With(
			"route", route,
			"api_version", apiVersion,
			"status", status,
		).Add(1)
		responseBytes.With(
			"route", route,
			"api_version", apiVersion,
			"status", status,
		).Add(float64(rc.bytesWritten))
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
