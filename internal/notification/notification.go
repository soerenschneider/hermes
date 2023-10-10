package notification

import (
	"context"
	"errors"
)

var ErrServiceNotFound = errors.New("notification service not found")

type NotificationProvider interface {
	Send(ctx context.Context, subject, message string) error
}
