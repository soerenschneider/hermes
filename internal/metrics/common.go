package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	namespace             = "hermes"
	subsystemNotification = "notification"
)

var (
	NotificationValidationErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemNotification,
		Name:      "validation_errors_total",
		Help:      "Total errors validating notifications",
	}, []string{"service"})

	NotificationDispatchErrors = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemNotification,
		Name:      "dispatch_errors_total",
		Help:      "Total errors dispatching notifications",
	}, []string{"service"})

	NotificationDispatchRetries = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemNotification,
		Name:      "dispatch_backoff_retries_total",
		Help:      "Total amount of retries for sending notification",
	}, []string{"service"})

	NotificationDispatchTime = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: namespace,
		Subsystem: subsystemNotification,
		Name:      "dispatch_time_total",
		Buckets:   []float64{0.75, 0.9, 0.95, 0.99},
	}, []string{"service"})
)
