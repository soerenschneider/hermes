package main

import (
	"github.com/soerenschneider/hermes/internal/config"
	"github.com/soerenschneider/hermes/internal/events"
	"github.com/soerenschneider/hermes/internal/events/http"
	"github.com/soerenschneider/hermes/internal/events/kafka"
	"github.com/soerenschneider/hermes/internal/events/rabbitmq"
	"github.com/soerenschneider/hermes/internal/events/smtp"
	"github.com/soerenschneider/hermes/internal/events/smtp/auth"

	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

func buildEventSources(conf *config.Config) ([]events.EventSource, error) {
	var eventSources []events.EventSource

	var errs error
	for _, eventSourceImpl := range conf.EventSourceImpl {
		var err error
		var impl events.EventSource

		switch eventSourceImpl {
		case "kafka":
			impl, err = buildKafka(conf)
		case "rabbitmq":
			impl, err = buildRabbitMq(conf)
		case "http":
			impl, err = buildHttpServer(conf)
		case "smtp":
			impl, err = buildSmtp(conf)
		default:
			log.Warn().Msgf("Unknown event source impl: %s. This should not happen", eventSourceImpl)
		}

		if err != nil {
			errs = multierr.Append(errs, err)
		} else {
			eventSources = append(eventSources, impl)
		}
	}

	return eventSources, errs
}

func buildKafka(conf *config.Config) (*kafka.KafkaReader, error) {
	var opts []kafka.KafkaReaderOpts
	if conf.Kafka.Partition > 0 {
		opts = append(opts, kafka.WithPartition(conf.Kafka.Partition))
	}

	if len(conf.Kafka.TlsCertFile) > 0 && len(conf.Kafka.TlsKeyFile) > 0 {
		opts = append(opts, kafka.WithTlsCert(conf.Kafka.TlsCertFile))
		opts = append(opts, kafka.WithTlsKey(conf.Kafka.TlsKeyFile))
	}

	return kafka.NewReader(conf.Kafka.Brokers, conf.Kafka.Topic, conf.Kafka.GroupId, opts...)
}

func buildRabbitMq(conf *config.Config) (*rabbitmq.RabbitMqEventListener, error) {
	var opts []rabbitmq.RabbitMqOpts

	// add webhook_server path
	if len(conf.RabbitMq.ConsumerName) > 0 {
		opts = append(opts, rabbitmq.WithConsumerName(conf.RabbitMq.ConsumerName))
	}

	conn := rabbitmq.RabbitMqConnection{
		BrokerHost: conf.RabbitMq.Broker,
		Port:       conf.RabbitMq.Port,
		Username:   conf.RabbitMq.Username,
		Password:   conf.RabbitMq.Password,
		Vhost:      conf.RabbitMq.Vhost,
		CertFile:   conf.RabbitMq.TlsCertFile,
		KeyFile:    conf.RabbitMq.TlsKeyFile,
		UseSsl:     conf.RabbitMq.UseSsl,
	}

	return rabbitmq.New(conn, conf.RabbitMq.QueueName, opts...)
}

func buildHttpServer(conf *config.Config) (*http.HttpServer, error) {
	var opts []http.HttpServerOpts

	// add tls keys
	if len(conf.Http.TlsCertFile) > 0 && len(conf.Http.TlsKeyFile) > 0 {
		opts = append(opts, http.WithTLS(conf.Http.TlsCertFile, conf.Http.TlsKeyFile))
	}

	if len(conf.Http.TlsClientCa) > 0 {
		opts = append(opts, http.WithClientCertificateValidation(conf.Http.TlsClientCa))
	}

	return http.New(conf.Http.Address, opts...)
}

func buildSmtp(conf *config.Config) (*smtp.SmtpServer, error) {
	var opts []smtp.SmtpOpts

	// add tls keys
	if len(conf.Smtp.TlsCertFile) > 0 && len(conf.Smtp.TlsKeyFile) > 0 {
		opts = append(opts, smtp.WithTLS(conf.Smtp.TlsCertFile, conf.Smtp.TlsKeyFile))
	}

	return smtp.NewSmtp("localhost:1025", "example.net", &auth.NoAuth{}, opts...)
}
