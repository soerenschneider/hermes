package notification

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/internal/queue"
	"github.com/soerenschneider/hermes/internal/validation"
	"github.com/soerenschneider/hermes/pkg"

	"github.com/cenkalti/backoff/v4"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

const (
	maxQueueRetries           = 15
	defaultQueueRetryInterval = 1 * time.Minute
)

type DispatcherOpts func(*Dispatcher) error

type Dispatcher struct {
	providers map[string]NotificationProvider

	retryQueue               queue.Queue
	retryQueueInterval       time.Duration
	retryQueueReconciliation sync.Once

	backoffImpl backoff.BackOff

	acceptBuffer    chan pkg.Notification
	deadLetterQueue NotificationProvider
}

func NewDispatcher(providers map[string]NotificationProvider, queueImpl queue.Queue, opts ...DispatcherOpts) (*Dispatcher, error) {
	if len(providers) == 0 {
		return nil, errors.New("no notification services provided")
	}

	dispatcher := &Dispatcher{
		backoffImpl:        backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3),
		providers:          providers,
		retryQueue:         queueImpl,
		retryQueueInterval: defaultQueueRetryInterval,
		acceptBuffer:       make(chan pkg.Notification, 30),
	}

	var errs error
	for _, opt := range opts {
		if err := opt(dispatcher); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return dispatcher, errs
}

func (d *Dispatcher) Accept(notification pkg.NotificationRequest, eventSource string) error {
	if err := validation.Validate(notification); err != nil {
		metrics.NotificationValidationErrors.WithLabelValues(eventSource).Inc()
		return err
	}

	if !d.hasServiceDefined(notification.ServiceId) {
		return fmt.Errorf("no dead letter retryQueue defined and service %q not known: %w", notification.ServiceId, ErrServiceNotFound)
	}

	d.acceptBuffer <- pkg.FromNotification(notification)
	return nil
}

func (d *Dispatcher) StartQueueReconciliation(ctx context.Context) {
	d.retryQueueReconciliation.Do(func() {
		ticker := time.NewTicker(d.retryQueueInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				for !d.retryQueue.IsEmpty() {
					msg, err := d.retryQueue.Get()
					if err != nil {
						metrics.QueueErrors.WithLabelValues("dequeue").Inc()
						log.Debug().Err(err).Msg("could not dequeue message")
						continue
					}

					d.acceptBuffer <- msg
				}
			}
		}
	})
}

func (d *Dispatcher) Listen(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case item := <-d.acceptBuffer:
			func(item pkg.Notification) {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				defer cancel()
				svc := d.getNotificationService(item.ServiceId)
				d.send(ctx, svc, item)
			}(item)
		}
	}
}

func (d *Dispatcher) getNotificationService(serviceId string) NotificationProvider {
	svc, ok := d.providers[serviceId]
	if ok {
		return svc
	}

	if d.deadLetterQueue == nil {
		return nil
	}

	return d.deadLetterQueue
}

func (d *Dispatcher) hasServiceDefined(serviceId string) bool {
	_, ok := d.providers[serviceId]
	if ok {
		return true
	}

	return d.deadLetterQueue != nil
}

func (d *Dispatcher) send(ctx context.Context, svc NotificationProvider, item pkg.Notification) {
	dispatch := func() error {
		metrics.NotificationDispatchRetries.WithLabelValues(item.ServiceId).Inc()
		return svc.Send(ctx, item.Subject, item.Message)
	}

	if err := backoff.Retry(dispatch, d.backoffImpl); err != nil {
		item.UnsuccessfulAttempts += 1
		metrics.NotificationDispatchErrors.WithLabelValues(item.ServiceId).Inc()
		log.Warn().Err(err).Int("retries", item.UnsuccessfulAttempts).Msg("could not dispatch message, adding to retryQueue")

		if item.UnsuccessfulAttempts >= maxQueueRetries {
			log.Warn().Msg("Dropping message")
			return
		}

		if err := d.retryQueue.Offer(item); err != nil {
			metrics.QueueErrors.WithLabelValues("offer").Inc()
			log.Error().Err(err).Msg("could not enqueue message, this should not happen")
		}
	}
}
