package ratelimit

import (
	"context"
	"time"
)

type Limiter interface {
	Allow() (bool, time.Duration, error)
	Wait(ctx context.Context) error
}
