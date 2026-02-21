package telemetry

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics holds all Prometheus metrics exposed by NetSentry.
type Metrics struct {
	ValidationTotal    prometheus.Counter
	ValidationDuration prometheus.Histogram
	PolicyViolations   *prometheus.CounterVec
	ActiveValidations  prometheus.Gauge
}

// NewMetrics registers and returns the application Prometheus metrics.
func NewMetrics(namespace string) *Metrics {
	if namespace == "" {
		namespace = "netsentry"
	}
	return &Metrics{
		ValidationTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "validations_total",
			Help:      "Total number of validation runs completed.",
		}),
		ValidationDuration: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Name:      "validation_duration_seconds",
			Help:      "Duration of validation runs in seconds.",
			Buckets:   prometheus.DefBuckets,
		}),
		PolicyViolations: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "policy_violations_total",
			Help:      "Total policy violations by severity.",
		}, []string{"severity"}),
		ActiveValidations: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "active_validations",
			Help:      "Number of currently running validation operations.",
		}),
	}
}

// Handler returns the Prometheus HTTP handler for the /metrics endpoint.
func MetricsHandler() http.Handler {
	return promhttp.Handler()
}
