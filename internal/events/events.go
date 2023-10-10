package events

import (
	"context"
	"sync"

	"github.com/soerenschneider/hermes/internal/notification"
)

type EventSource interface {
	Listen(ctx context.Context, cortex *notification.Dispatcher, wg *sync.WaitGroup) error
}
