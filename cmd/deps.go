package main

import (
	"github.com/soerenschneider/hermes/internal/config"
	"github.com/soerenschneider/hermes/internal/notification"
	"github.com/soerenschneider/hermes/internal/queue"
)

const defaultQueueSize = 500

func (d *Deps) buildCortex(conf *config.Config) (*notification.Dispatcher, error) {
	providers, err := buildProviders(conf)
	if err != nil {
		return nil, err
	}

	var opts []notification.DispatcherOpts
	if len(conf.DeadLetterQueue) > 0 {
		opts = append(opts, notification.WithDeadLetterQueue(conf.DeadLetterQueue))
	}

	return notification.NewDispatcher(providers, d.queue, opts...)
}

func (d *Deps) buildQueue(conf *config.Config) (queue.Queue, error) {
	return queue.NewQueue(defaultQueueSize)
}
