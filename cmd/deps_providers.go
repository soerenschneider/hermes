package main

import (
	"fmt"
	"os"

	"github.com/soerenschneider/hermes/internal/config"
	"github.com/soerenschneider/hermes/internal/notification"

	"github.com/nikoksr/notify/service/mail"
	"github.com/nikoksr/notify/service/telegram"
	"go.uber.org/multierr"
)

func buildProviders(conf *config.Config) (map[string]notification.NotificationProvider, error) {
	ret := make(map[string]notification.NotificationProvider)

	var errs error
	errs = multierr.Append(errs, buildTelegram(conf, ret))
	errs = multierr.Append(errs, buildEmail(conf, ret))
	return ret, errs
}

func buildTelegram(conf *config.Config, n map[string]notification.NotificationProvider) error {
	var errs error
	for _, t := range conf.Telegram {
		token, err := returnValueOrFileContent(t.Token, t.TokenFile)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		telegram, err := telegram.New(token)
		if err != nil {
			errs = multierr.Append(errs, fmt.Errorf("can not build telegram notififer: %w", err))
			continue
		}
		telegram.AddReceivers(t.Receivers...)

		_, ok := n[t.ServiceUri]
		if ok {
			errs = multierr.Append(errs, fmt.Errorf("not adding telegram notification: serviceUri %s already registed", t.ServiceUri))
			continue
		}
		n[t.ServiceUri] = telegram
	}

	return errs
}

func buildEmail(conf *config.Config, n map[string]notification.NotificationProvider) error {
	var errs error
	for _, e := range conf.Email {
		mail := mail.New(e.Sender, e.Host)
		password, err := returnValueOrFileContent(e.Password, e.PasswordFile)
		if err != nil {
			errs = multierr.Append(errs, err)
			continue
		}

		mail.AuthenticateSMTP("", e.UserName, password, e.Host)
		mail.AddReceivers(e.Receivers...)
		_, ok := n[e.ServiceUri]
		if ok {
			errs = multierr.Append(errs, fmt.Errorf("not adding email notification: serviceUri %s already registed", e.ServiceUri))
			continue
		}
		n[e.ServiceUri] = mail
	}

	return errs
}

func returnValueOrFileContent(val, file string) (string, error) {
	if len(val) == 0 {
		data, err := os.ReadFile(file)
		return string(data), err
	}

	return val, nil
}
