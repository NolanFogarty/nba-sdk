package httpx

import (
	"context"
	"sync"
	"time"
)

// Limiter is a simple token-bucket rate limiter. It is safe for concurrent use.
type Limiter struct {
	mu       sync.Mutex
	tokens   float64
	capacity float64
	rate     float64
	last     time.Time
}

// NewLimiter returns a limiter that refills at `rate` tokens per second up to `burst` tokens.
// A nil receiver is treated as unlimited, which is useful for tests.
func NewLimiter(rate float64, burst int) *Limiter {
	return &Limiter{
		tokens:   float64(burst),
		capacity: float64(burst),
		rate:     rate,
		last:     time.Now(),
	}
}

// Wait blocks until a token is available or ctx is cancelled.
func (l *Limiter) Wait(ctx context.Context) error {
	if l == nil {
		return nil
	}
	for {
		l.mu.Lock()
		now := time.Now()
		elapsed := now.Sub(l.last).Seconds()
		l.tokens = min(l.capacity, l.tokens+elapsed*l.rate)
		l.last = now
		if l.tokens >= 1 {
			l.tokens--
			l.mu.Unlock()
			return nil
		}
		needed := (1 - l.tokens) / l.rate
		l.mu.Unlock()

		t := time.NewTimer(time.Duration(needed * float64(time.Second)))
		select {
		case <-ctx.Done():
			t.Stop()
			return ctx.Err()
		case <-t.C:
		}
	}
}
