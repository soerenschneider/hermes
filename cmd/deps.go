package main

import (
	"github.com/soerenschneider/hermes/internal/config"
	"github.com/soerenschneider/hermes/internal/notification"
	"github.com/soerenschneider/hermes/internal/queue"
	"github.com/soerenschneider/hermes/internal/queue/sqlite"
)

const defaultQueueSize = 500

func (d *Deps) buildCortex(conf *config.Config) (*notification.NotificationDispatcher, error) {
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
	if conf.Db == nil || conf.Db.Type == "memory" {
		return sqlite.New(":memory:")
	}

	return sqlite.New(conf.Db.Name)
}
