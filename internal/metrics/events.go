package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const subsystemEvents = "events"

var (
	AcceptedNotifications = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemEvents,
		Name:      "accepted_notifications_total",
		Help:      "Total number of accepted notifications",
	}, []string{"subsystem"})

	NotificationGarbageData = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemEvents,
		Name:      "invalid_notifications_total",
		Help:      "Total http invalid messages",
	}, []string{"subsystem"})

	SmtpSessions = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemEvents,
		Name:      "smtp_sessions_total",
		Help:      "Total sessions opened",
	})

	SmtpAuthFailure = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: namespace,
		Subsystem: subsystemEvents,
		Name:      "smtp_auth_failures_total",
		Help:      "Total smtp auth failures",
	})
)
