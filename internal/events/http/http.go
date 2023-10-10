package http

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/internal/notification"
	"github.com/soerenschneider/hermes/pkg"

	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog/log"
	"go.uber.org/multierr"
)

type HttpServer struct {
	address string
	cortex  *notification.Dispatcher

	// optional
	certFile string
	keyFile  string

	clientCaCertFile string
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

func (s *HttpServer) notifyHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Warn().Err(err).Msg("could not close http body")
	}

	msg := pkg.NotificationRequest{}
	if err := json.Unmarshal(data, &msg); err != nil {
		metrics.NotificationGarbageData.WithLabelValues("http").Inc()
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.cortex.Accept(msg, "http"); err != nil {
		_, isValidationErr := err.(validator.ValidationErrors)
		if errors.Is(err, notification.ErrServiceNotFound) || isValidationErr {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (s *HttpServer) Listen(ctx context.Context, cortex *notification.Dispatcher, wg *sync.WaitGroup) error {
	log.Info().Msgf("Starting http server event source, listening on %q", s.address)
	wg.Add(1)

	s.cortex = cortex

	mux := http.NewServeMux()
	mux.HandleFunc("/notify", s.notifyHandler)

	tlsConfig, err := s.getTlsConf()
	if err != nil {
		return err
	}

	server := http.Server{
		Addr:              s.address,
		Handler:           mux,
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
