package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	RequestCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "endpoint"},
	)

	ShortenedLinksTotal = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "shortened_links_total",
			Help: "Total number of shortened links",
		},
	)

	RequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response durations for HTTP requests",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint"},
	)

	ErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_errors_total",
			Help: "Total number of HTTP errors",
		},
		[]string{"code"},
	)

	UniqueUsersTotal = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "unique_users_total",
			Help: "Total number of unique users",
		},
	)
)

func init() {
	prometheus.MustRegister(RequestCount)
	prometheus.MustRegister(ShortenedLinksTotal)
	prometheus.MustRegister(RequestDuration)
	prometheus.MustRegister(ErrorCount)
	prometheus.MustRegister(UniqueUsersTotal)
}
