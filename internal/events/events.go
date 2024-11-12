package events

import (
	"context"
	"sync"

	"github.com/soerenschneider/hermes/internal/domain"
)

type Dispatcher interface {
	Accept(notification domain.NotificationRequest, eventSource string) error
}

type EventSource interface {
	Listen(ctx context.Context, dispatcher Dispatcher, wg *sync.WaitGroup) error
}
