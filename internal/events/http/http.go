package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/soerenschneider/hermes/internal/domain"
	"github.com/soerenschneider/hermes/internal/events"
	"github.com/soerenschneider/hermes/internal/notification"
	"gitlab.com/tanna.dev/openapi-doc-http-handler/elements"

	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

type HttpServer struct {
	address    string
	dispatcher events.Dispatcher

	// optional
	certFile string
	keyFile  string

	clientCaCertFile string
}

func (s *HttpServer) SendNotification(ctx context.Context, request SendNotificationRequestObject) (SendNotificationResponseObject, error) {
	msg := domain.NotificationRequest{
		ServiceId: request.Body.ServiceId,
		Subject:   request.Body.Subject,
		Message:   request.Body.Message,
	}

	if err := s.dispatcher.Accept(msg, "http"); err != nil {
		_, isValidationErr := err.(validator.ValidationErrors)
		if isValidationErr {
			return SendNotification400JSONResponse{Error: "validation error"}, nil
		}

		if errors.Is(err, notification.ErrServiceNotFound) {
			return SendNotification400JSONResponse{Error: "service_id not found"}, nil
		}

		return SendNotification500JSONResponse{}, nil
	}

	return SendNotification200Response{}, nil
}

type HttpServerOpts func(*HttpServer) error

func New(address string, opts ...HttpServerOpts) (*HttpServer, error) {
	if len(address) == 0 {
		return nil, errors.New("empty address provided")
	}

	w := &HttpServer{
		address: address,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(w); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return w, errs
}

func (s *HttpServer) IsTLSConfigured() bool {
	return len(s.certFile) > 0 && len(s.keyFile) > 0
}

func (s *HttpServer) Listen(ctx context.Context, dispatcher events.Dispatcher, wg *sync.WaitGroup) error {
	log.Info().Msgf("Starting http server event source, listening on %q", s.address)
	wg.Add(1)

	s.dispatcher = dispatcher

	handler, err := s.getOpenApiHandler()
	if err != nil {
		return err
	}

	tlsConfig, err := s.getTlsConf()
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:              s.address,
		Handler:           handler,
		ReadTimeout:       3 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       30 * time.Second,
		TLSConfig:         tlsConfig,
	}

	errChan := make(chan error)
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- fmt.Errorf("can not start webhook server: %w", err)
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("Stopping webhook server")
		err := server.Shutdown(ctx)
		wg.Done()
		return err
	case err := <-errChan:
		return err
	}
}

func (s *HttpServer) getTlsConf() (*tls.Config, error) {
	if !s.IsTLSConfigured() {
		return nil, nil
	}

	tlsConfig := &tls.Config{
		GetCertificate: s.getCertificate,
		MinVersion:     tls.VersionTLS13,
	}

	if len(s.clientCaCertFile) > 0 {
		caPool, err := generateClientCaCertPool(s.clientCaCertFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.ClientCAs = caPool
		tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
	}

	return tlsConfig, nil
}

func generateClientCaCertPool(caFile string) (*x509.CertPool, error) {
	data, err := os.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(data) {
		return nil, fmt.Errorf("could not read valid cert data from %q", caFile)
	}

	return certPool, nil
}

func (s *HttpServer) getCertificate(info *tls.ClientHelloInfo) (*tls.Certificate, error) {
	log.Info().Msg("Reading TLS certs")
	if len(s.certFile) == 0 || len(s.keyFile) == 0 {
		return nil, errors.New("no client certificates defined")
	}

	certificate, err := tls.LoadX509KeyPair(s.certFile, s.keyFile)
	if err != nil {
		log.Error().Err(err).Msg("user-defined client certificates could not be loaded")
	}
	return &certificate, err
}

func (s *HttpServer) getOpenApiHandler() (http.Handler, error) {
	// add a mux that serves /docs
	swagger, err := GetSwagger()
	if err != nil {
		return nil, err
	}

	docs, err := elements.NewHandler(swagger, err)
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/docs", docs)
	mux.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	options := StdHTTPServerOptions{
		Middlewares: []MiddlewareFunc{},
		BaseRouter:  mux,
	}

	strictHandler := NewStrictHandler(s, nil)
	return HandlerWithOptions(strictHandler, options), nil
}
