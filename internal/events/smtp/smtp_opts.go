package smtp

import "errors"

func WithTLS(certFile, keyFile string) SmtpOpts {
	return func(s *SmtpServer) error {
		if len(certFile) == 0 {
			return errors.New("empty certfile")
		}

		if len(keyFile) == 0 {
			return errors.New("empty keyfile")
		}

		s.certFile = certFile
		s.keyFile = keyFile
		return nil
	}
}
