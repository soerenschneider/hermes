package smtp

import (
	"errors"
	"io"
	"net/mail"
	"strings"

	"github.com/soerenschneider/hermes/internal/metrics"
	"github.com/soerenschneider/hermes/internal/notification"
	"github.com/soerenschneider/hermes/pkg"

	"github.com/emersion/go-smtp"
	"github.com/rs/zerolog/log"
)

type Session struct {
	auth            UserAuth
	cortex          *notification.Dispatcher
	isAuthenticated bool

	user    string
	from    string
	to      []string
	content *mail.Message
}

func (s *Session) AuthPlain(username, password string) error {
	if !s.auth.Authenticate(username, password) {
		metrics.SmtpAuthFailure.Inc()
		return smtp.ErrAuthFailed
	}

	s.isAuthenticated = true
	s.user = username
	return nil
}

func (s *Session) Mail(from string, opts *smtp.MailOptions) error {
	if !s.isAuthenticated {
		return smtp.ErrAuthRequired
	}
	s.from = from
	return nil
}

func (s *Session) Rcpt(to string, opts *smtp.RcptOptions) error {
	if !s.isAuthenticated {
		return smtp.ErrAuthRequired
	}
	s.to = append(s.to, to)
	return nil
}

func (s *Session) Data(r io.Reader) error {
	if !s.isAuthenticated {
		return smtp.ErrAuthRequired
	}

	var err error
	s.content, err = mail.ReadMessage(r)
	if err != nil {
		metrics.NotificationGarbageData.WithLabelValues("smtp").Inc()
	}
	return err
}

func (s *Session) Reset() {}

func (s *Session) Logout() error {
	if s.isAuthenticated {
		log.Info().Str("from", s.from).Any("to", s.to).Msgf("Accepted event via smtp")
		msg, err := s.toNotification()
		if err != nil {
			metrics.NotificationGarbageData.WithLabelValues("smtp").Inc()
			log.Error().Err(err).Str("source", "smtp").Msg("can not create message")
			return err
		}

		if err := s.cortex.Accept(*msg, "smtp"); err != nil {
			log.Error().Err(err).Msg("could not accept notification")
		}
	}

	s.isAuthenticated = false
	return nil
}

func (s *Session) toNotification() (*pkg.NotificationRequest, error) {
	if s.content == nil {
		return nil, errors.New("invalid data")
	}
	bodyBytes, err := io.ReadAll(s.content.Body)
	if err != nil {
		return nil, err
	}

	body := string(bodyBytes)

	return &pkg.NotificationRequest{
		ServiceId: extractServiceId(extractServiceId(s.to[0])),
		Subject:   s.content.Header.Get("Subject"),
		Message:   body,
	}, nil
}

func extractServiceId(to string) string {
	if !strings.Contains(to, "@") {
		return to
	}

	split := strings.Split(to, "@")
	return split[0]
}
