package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"sync"
	"time"

	"github.com/soerenschneider/hermes/internal/events"
	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/pkg"

	"github.com/rs/zerolog/log"
	"github.com/segmentio/kafka-go"
	"go.uber.org/multierr"
)

const defaultTimeout = 10 * time.Second

type KafkaReader struct {
	brokers   []string
	topic     string
	groupId   string
	partition int

	reader *kafka.Reader

	tlsKey  string
	tlsCert string

	once sync.Once
}

type KafkaReaderOpts func(*KafkaReader) error

func NewReader(brokers []string, topic string, groupId string, opts ...KafkaReaderOpts) (*KafkaReader, error) {
	if len(brokers) == 0 {
		return nil, errors.New("empty list of kafka brokers supplied")
	}

	if len(topic) == 0 {
		return nil, errors.New("empty topic supplied")
	}

	if len(groupId) == 0 {
		return nil, errors.New("empty groupId supplied")
	}

	kafka := &KafkaReader{
		topic:   topic,
		brokers: brokers,
		groupId: groupId,
	}

	var errs error
	for _, opt := range opts {
		err := opt(kafka)
		errs = multierr.Append(errs, err)
	}

	return kafka, errs
}

func (k *KafkaReader) initReader() {
	initReader := func() {
		dialer := &kafka.Dialer{
			Timeout:   defaultTimeout,
			DualStack: true,
			TLS: &tls.Config{
				GetClientCertificate: k.loadTlsClientCerts,
				MinVersion:           tls.VersionTLS12,
			},
		}

		k.reader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:   k.brokers,
			Topic:     k.topic,
			Partition: k.partition,
			GroupID:   k.groupId,
			MaxBytes:  10e6,
			Dialer:    dialer,
		})
	}

	k.once.Do(initReader)
}

func (k *KafkaReader) Listen(ctx context.Context, dispatcher events.Dispatcher, wg *sync.WaitGroup) error {
	log.Info().Msgf("Starting kafka reader event source")
	if ctx == nil {
		return errors.New("empty context supplied")
	}

	if dispatcher == nil {
		return errors.New("closed channel supplied")
	}

	if wg == nil {
		return errors.New("empty waitgroup supplied")
	}

	wg.Add(1)
	defer wg.Done()

	k.initReader()

	continueConsuming := true
	for continueConsuming {
		select {
		case <-ctx.Done():
			log.Info().Msg("Kafka received signal")
			continueConsuming = false
		default:
			k.readMessage(ctx, dispatcher)
		}
	}

	return k.reader.Close()
}

func (k *KafkaReader) readMessage(ctx context.Context, dispatcher events.Dispatcher) {
	msg, err := k.reader.FetchMessage(ctx)
	if err != nil && !errors.Is(err, context.Canceled) {
		log.Error().Err(err).Msg("Error while reading kafka message")
		return
	}

	// TODO: error handling
	if err := k.handleMessage(msg, dispatcher); err != nil {
		log.Error().Err(err).Msg("error handling message")
	}

	if err := k.reader.CommitMessages(ctx, msg); err != nil {
		log.Error().Err(err).Msg("could not commit message")
	}
}

func (k *KafkaReader) handleMessage(msg kafka.Message, dispatcher events.Dispatcher) error {
	notification := pkg.NotificationRequest{}
	if err := json.Unmarshal(msg.Value, &notification); err != nil {
		metrics.NotificationGarbageData.WithLabelValues("kafka").Inc()
		return err
	}

	if err := dispatcher.Accept(notification, "kafka"); err != nil {
		log.Error().Err(err).Msg("could not dispatch message")
		return err
	}

	metrics.AcceptedNotifications.WithLabelValues("kafka").Inc()
	return nil
}

func (k *KafkaReader) loadTlsClientCerts(info *tls.CertificateRequestInfo) (*tls.Certificate, error) {
	if len(k.tlsCert) == 0 || len(k.tlsKey) == 0 {
		return nil, errors.New("no client certificates defined")
	}

	certificate, err := tls.LoadX509KeyPair(k.tlsCert, k.tlsKey)
	if err != nil {
		log.Error().Err(err).Msg("user-defined client certificates could not be loaded")
	}
	return &certificate, err
}
