package notification

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

func WithDeadLetterQueue(service string) DispatcherOpts {
	return func(c *NotificationDispatcher) error {
		if len(service) == 0 {
			return errors.New("empty deadletterqueue service provided")
		}

		service, ok := c.providers[service]
		if !ok {
			return fmt.Errorf("invalid service for deadletterqueue: no such service: %s", service)
		}
		log.Info().Msgf("Using service %s as dead letter retryQueue", service)
		c.deadLetterQueue = service
		return nil
	}
}
