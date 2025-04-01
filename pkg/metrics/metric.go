package metrics

import "github.com/prometheus/client_golang/prometheus"

func init() {
	prometheus.MustRegister(HttpRequestsTotal)
	prometheus.MustRegister(HttpRequestDuration)
	prometheus.MustRegister(ApplicationErrorTotal)
	prometheus.MustRegister(ServiceRequestDuration)
	prometheus.MustRegister(DbQueryHistogram)
	prometheus.MustRegister(ApplicationPanicTotal)
}

var (
	// application monitoring metrics
	HttpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"path", "method", "status"},
	)
	HttpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Duration of HTTP requests in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"path", "method"},
	)
	// application error monitoring metrics
	ApplicationErrorTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "application_error_total",
			Help: "application error total",
		},
		[]string{"path", "method", "code"},
	)
	// service monitoring metrics
	ServiceRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "service_request_duration_seconds",
			Help:    "Histogram of HTTP request durations",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"service_name", "method", "endpoint", "status"},
	)
	// mysql monitoring metrics
	DbQueryHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "gorm_query_duration_seconds",
			Help:    "Duration of GORM database queries.",
			Buckets: prometheus.ExponentialBuckets(0.001, 2, 15), // 1ms ~ 32s
		},
		[]string{"operation", "table"},
	)
	// Application panic monitoring metrics
	ApplicationPanicTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "application_panic_total",
			Help: "application panic total",
		},
		[]string{"panic_location"},
	)
)
