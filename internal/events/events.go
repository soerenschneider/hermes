package events

import (
	"context"
	"sync"

	"github.com/soerenschneider/hermes/pkg"
)

type Dispatcher interface {
	Accept(notification pkg.NotificationRequest, eventSource string) error
}

type EventSource interface {
	Listen(ctx context.Context, dispatcher Dispatcher, wg *sync.WaitGroup) error
}
