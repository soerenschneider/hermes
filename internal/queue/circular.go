package queue

import (
	"errors"

	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/pkg"

	"github.com/adrianbrad/queue"
	"github.com/rs/zerolog/log"
)

type CircularQueue struct {
	queue *queue.Circular[pkg.Notification]
}

func NewQueue(size int) (*CircularQueue, error) {
	if size < 100 {
		return nil, errors.New("size must not be less than 100")
	}

	metrics.QueueCapacity.Set(float64(size))
	return &CircularQueue{
		queue: queue.NewCircular([]pkg.Notification{}, size),
	}, nil
}

func (q *CircularQueue) Offer(item pkg.Notification) error {
	log.Debug().Msg("Inserting item to queue")
	err := q.queue.Offer(item)
	metrics.QueueSize.Set(float64(q.queue.Size()))
	return err
}

func (q *CircularQueue) Get() (pkg.Notification, error) {
	log.Debug().Msg("Removing head from from queue")
	item, err := q.queue.Get()
	metrics.QueueSize.Set(float64(q.queue.Size()))
	return item, err
}

func (q *CircularQueue) IsEmpty() bool {
	return q.queue.IsEmpty()
}
