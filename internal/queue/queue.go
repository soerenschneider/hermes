package queue

import (
	"context"
	"math"
	"math/rand/v2"
	"time"

	"github.com/soerenschneider/hermes/pkg"
)

type Queue interface {
	Offer(ctx context.Context, item pkg.Notification) error
	Get(ctx context.Context) (pkg.Notification, error)
	IsEmpty(ctx context.Context) (bool, error)
	GetMessageCount(ctx context.Context) (int64, error)
}

func ExponentialBackoff(n int, baseDelay time.Duration, maxDelay time.Duration) time.Duration {
	backoff := float64(baseDelay) * math.Pow(2, float64(n))

	if backoff > float64(maxDelay) {
		backoff = float64(maxDelay)
	}

	jitter := rand.Float64()*0.2 + 0.9 // Generates a value between 0.9 and 1.1
	backoff *= jitter

	return time.Duration(backoff)
}
