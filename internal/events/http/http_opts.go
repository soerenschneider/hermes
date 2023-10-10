package http

import (
	"errors"
	"os"
)

func WithTLS(certFile, keyFile string) func(s *HttpServer) error {
	return func(s *HttpServer) error {
		if len(certFile) == 0 {
			return errors.New("empty certfile")
		}

		if err := canRead(certFile); err != nil {
			return err
		}

		if len(keyFile) == 0 {
			return errors.New("empty keyfile")
		}

		if err := canRead(keyFile); err != nil {
			return err
		}

		s.certFile = certFile
		s.keyFile = keyFile
		return nil
	}
}

func WithClientCertificateValidation(caFile string) func(s *HttpServer) error {
	return func(s *HttpServer) error {
		if len(caFile) == 0 {
			return errors.New("empty certfile")
		}

		if err := canRead(caFile); err != nil {
			return err
		}

		s.clientCaCertFile = caFile
		return nil
	}
}

func canRead(filePath string) error {
	_, err := os.Stat(filePath)

	if err != nil {
		return err
	}

	return nil
}
