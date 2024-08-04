package rabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cenkalti/backoff/v4"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/hermes/internal/events"
	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/pkg"
	"go.uber.org/multierr"
)

type RabbitMqConnection struct {
	BrokerHost string
	Port       int
	Username   string
	Password   string
	Vhost      string
	UseSsl     bool

	CertFile string
	KeyFile  string
}

type RabbitMqEventListener struct {
	connection   RabbitMqConnection
	queueName    string
	consumerName string
}

type RabbitMqOpts func(listener *RabbitMqEventListener) error

func New(conn RabbitMqConnection, queueName string, opts ...RabbitMqOpts) (*RabbitMqEventListener, error) {
	ret := &RabbitMqEventListener{
		connection:   conn,
		queueName:    queueName,
		consumerName: "hermes",
	}

	var errs error
	for _, opt := range opts {
		if err := opt(ret); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return ret, errs
}

func (e *RabbitMqEventListener) buildConnectionString() string {
	protocol := "amqp"
	if e.connection.UseSsl {
		protocol = "amqps"
	}

	if len(e.connection.CertFile) > 0 && len(e.connection.KeyFile) > 0 {
		return fmt.Sprintf("amqps://%s:%d%s", e.connection.BrokerHost, e.connection.Port, e.connection.Vhost)
	}

	return fmt.Sprintf("%s://%s:%s@%s:%d%s", protocol, e.connection.Username, e.connection.Password, e.connection.BrokerHost, e.connection.Port, e.connection.Vhost)
}

func (e *RabbitMqEventListener) Listen(ctx context.Context, dispatcher events.Dispatcher, wg *sync.WaitGroup) error {
	wg.Add(1)
	defer wg.Done()

	impl := backoff.NewExponentialBackOff()
	operation := func() error {
		err := e.listen(ctx, dispatcher)
		if err != nil {
			log.Error().Str("component", "rabbitmq").Err(err).Msg("error while listening on rabbitmq events")
			var amqpErr *amqp.Error
			if errors.As(err, &amqpErr) {
				metrics.RabbitMqErrors.WithLabelValues(strconv.Itoa(amqpErr.Code)).Inc()
				if amqpErr.Code == 401 || amqpErr.Code == 403 {
					return backoff.Permanent(amqpErr)
				}
			}
		}

		if ctx.Err() != nil {
			return backoff.Permanent(err)
		}

		return err
	}
	notify := func(err error, d time.Duration) {
		log.Error().Err(err).Str("component", "rabbitmq").Msgf("Error after %v", d)
	}

	cont := true
	for cont {
		select {
		case <-ctx.Done():
			log.Debug().Str("component", "rabbitmq").Msg("Packing up")
			cont = false
		default:
			if err := backoff.RetryNotify(operation, impl, notify); err != nil {
				return fmt.Errorf("too many errors trying to listen on rabbitmq: %w", err)
			}
		}
	}

	return nil
}

func (e *RabbitMqEventListener) listen(ctx context.Context, dispatcher events.Dispatcher) error {
	conn, err := amqp.Dial(e.buildConnectionString())
	if err != nil {
		return err
	}
	defer conn.Close()
	conNotify := conn.NotifyClose(make(chan *amqp.Error, 1))

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	chNotify := ch.NotifyClose(make(chan *amqp.Error, 1))

	msgs, err := ch.Consume(
		e.queueName,
		e.consumerName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return backoff.Permanent(err)
	}

	log.Info().Str("component", "rabbitmq").Msg("Listening for messages...")
	for {
		select {
		case err := <-conNotify:
			log.Warn().Err(err).Str("component", "rabbitmq").Msg("connection closed")
			metrics.RabbitMqDisconnects.WithLabelValues("connection").Inc()
			return err
		case err := <-chNotify:
			log.Warn().Err(err).Str("component", "rabbitmq").Msg("channel closed")
			metrics.RabbitMqDisconnects.WithLabelValues("channel").Inc()
			return err
		case r := <-msgs:
			log.Debug().Str("component", "rabbitmq").Msg("received message")
			go e.handleMessage(r, dispatcher)
		case <-ctx.Done():
			log.Debug().Str("component", "rabbitmq").Msg("context done")
			return nil
		}
	}
}

func (e *RabbitMqEventListener) handleMessage(msg amqp.Delivery, dispatcher events.Dispatcher) {
	notification := pkg.NotificationRequest{}
	if err := json.Unmarshal(msg.Body, &notification); err != nil {
		metrics.NotificationGarbageData.WithLabelValues("rabbitmq").Inc()
		log.Error().Err(err).Str("component", "rabbitmq").Msg("could not decode message")
	}

	if err := dispatcher.Accept(notification, "rabbitmq"); err != nil {
		log.Error().Err(err).Msg("could not dispatch message")
	}

	metrics.AcceptedNotifications.WithLabelValues("rabbitmq").Inc()

	backoffImpl := backoff.NewExponentialBackOff()
	op := func() error {
		return msg.Ack(false)
	}
	notify := func(err error, duration time.Duration) {
		metrics.RabbitMqErrors.WithLabelValues("ack").Inc()
		log.Error().Err(err).Msgf("could not ack message after %v", duration)
	}
	if err := backoff.RetryNotify(op, backoffImpl, notify); err != nil {
		log.Error().Err(err).Str("event", "rabbitmq").Uint64("tag", msg.DeliveryTag).Msg("giving up on acking")
	}
}
