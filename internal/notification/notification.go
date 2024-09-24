package notification

import (
	"context"
	"errors"
	"net/http"
)

var ErrServiceNotFound = errors.New("notification service not found")

type NotificationProvider interface {
	Send(ctx context.Context, subject, message string) error
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}
