package ratelimit

import (
	"time"

	"golang.org/x/time/rate"
)

type LocalLimiter struct {
	limiter *rate.Limiter
}

func NewLocalLimiter(r rate.Limit, burst int) *LocalLimiter {
	return &LocalLimiter{
		limiter: rate.NewLimiter(r, burst),
	}
}

func (l *LocalLimiter) Allow() (bool, time.Duration, error) {
	if l.limiter.Allow() {
		return true, 0, nil
	}

	r := l.limiter.Reserve()
	delay := r.DelayFrom(time.Now())
	r.Cancel()

	return false, delay, nil
}
