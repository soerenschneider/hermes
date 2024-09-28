package notification

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

func WithDeadLetterQueue(serviceId string) DispatcherOpts {
	return func(c *NotificationDispatcher) error {
		if len(serviceId) == 0 {
			return errors.New("empty deadletterqueue service provided")
		}

		service, ok := c.providers[serviceId]
		if !ok {
			return fmt.Errorf("invalid service for deadletterqueue: no such service: %q", serviceId)
		}
		log.Info().Msgf("Using service %q as dead letter retryQueue", serviceId)
		c.deadLetterQueue = service
		return nil
	}
}
