package smtp

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
	"github.com/soerenschneider/hermes/internal/events"
	"github.com/soerenschneider/hermes/internal/metrics"
	"go.uber.org/multierr"
)

type UserAuth interface {
	Authenticate(user, password string) bool
}

type backend struct {
	auth       UserAuth
	dispatcher events.Dispatcher
}

func (b *backend) NewSession(_ *smtp.Conn) (smtp.Session, error) {
	metrics.SmtpSessions.Inc()
	return &Session{
		auth:       b.auth,
		dispatcher: b.dispatcher,
	}, nil
}

type SmtpOpts func(s *SmtpServer) error

type SmtpServer struct {
	server  *smtp.Server
	backend *backend

	certFile string
	keyFile  string
}

func NewSmtp(addr string, domain string, auth UserAuth, opts ...SmtpOpts) (*SmtpServer, error) {
	backend := &backend{
		auth: auth,
	}

	server := smtp.NewServer(backend)

	server.Addr = addr
	server.Domain = domain
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxMessageBytes = 1024 * 1024
	server.MaxRecipients = 50
	server.AllowInsecureAuth = true

	smtpServer := &SmtpServer{
		server:  server,
		backend: backend,
	}

	var errs error
	for _, opt := range opts {
		if err := opt(smtpServer); err != nil {
			errs = multierr.Append(errs, err)
		}
	}

	return smtpServer, errs
}

func (m *SmtpServer) Listen(ctx context.Context, dispatcher events.Dispatcher, wg *sync.WaitGroup) error {
	log.Info().Msgf("Starting event source smtp server")
	wg.Add(1)
	defer wg.Done()

	m.backend.dispatcher = dispatcher

	log.Info().Msgf("Starting smtp server")
	go func() {
		if err := m.server.ListenAndServe(); err != nil {
			log.Error().Err(err).Msg("could not start smtp listener")
		}
	}()

	<-ctx.Done()
	log.Info().Msg("Shutting down smtp server")
	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := m.server.Shutdown(shutdownCtx); err != nil && !errors.Is(err, context.Canceled) {
		log.Error().Err(err).Msg("")
		return err
	}
	return nil
}
