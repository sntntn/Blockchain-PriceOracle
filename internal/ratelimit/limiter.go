package ratelimit

import (
	"context"
	"time"
)

type Limiter interface {
	// I can add a ctx context later when I add a wait time and delay for a request
	// Allow(ctx context.Context) (bool, time.Duration, error)
	Allow() (bool, time.Duration, error)
	Wait(ctx context.Context) error
}
