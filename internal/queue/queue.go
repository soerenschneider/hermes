package queue

import (
	"github.com/soerenschneider/hermes/pkg"
)

type Queue interface {
	Offer(item pkg.Notification) error
	Get() (pkg.Notification, error)
	IsEmpty() bool
}
