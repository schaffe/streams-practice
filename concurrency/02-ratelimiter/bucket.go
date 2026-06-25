package ratelimiter

import (
	"sync"
	"time"
)

// TokenBucket is a lazily-refilled token bucket rate limiter.
type TokenBucket struct {
	mu         sync.Mutex
	capacity   float64
	tokens     float64
	refillRate float64 // tokens per second
	last       time.Time
	clock      Clock
}

// New returns a bucket that starts full, refilling at refillRate tokens/sec.
func New(capacity, refillRate float64) *TokenBucket {
	return newWithClock(capacity, refillRate, realClock{})
}

func newWithClock(capacity, refillRate float64, c Clock) *TokenBucket {
	return &TokenBucket{
		capacity:   capacity,
		tokens:     capacity,
		refillRate: refillRate,
		last:       c.Now(),
		clock:      c,
	}
}

// refill credits tokens accrued since last, clamped to capacity.
// Caller must hold mu.
func (b *TokenBucket) refill() {
	now := b.clock.Now()
	elapsed := now.Sub(b.last).Seconds()
	b.tokens = min(b.capacity, b.tokens+elapsed*b.refillRate)
	b.last = now
}

// Allow reports whether one token is available, consuming it if so.
func (b *TokenBucket) Allow() bool { return b.AllowN(1) }

// AllowN consumes n tokens if available; non-blocking.
func (b *TokenBucket) AllowN(n float64) bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.refill()
	if b.tokens >= n {
		b.tokens -= n
		return true
	}
	return false
}
