package ratelimiter

import (
	"context"
	"errors"
	"sync"
	"time"
)

// ErrExceedsCapacity is returned by WaitN when n is larger than capacity and
// could therefore never be satisfied.
var ErrExceedsCapacity = errors.New("ratelimiter: requested tokens exceed capacity")

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

// Wait blocks until one token is available or ctx is cancelled.
func (b *TokenBucket) Wait(ctx context.Context) error { return b.WaitN(ctx, 1) }

// WaitN blocks until n tokens are available or ctx is cancelled. It reserves
// the tokens before sleeping so concurrent waiters do not double-claim.
func (b *TokenBucket) WaitN(ctx context.Context, n float64) error {
	if n > b.capacity {
		return ErrExceedsCapacity
	}
	b.mu.Lock()
	b.refill()
	if b.tokens >= n {
		b.tokens -= n
		b.mu.Unlock()
		return nil
	}
	deficit := n - b.tokens
	delay := time.Duration(deficit / b.refillRate * float64(time.Second))
	b.tokens -= n // reserve (may go negative)
	b.last = b.clock.Now()
	b.mu.Unlock()

	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-timer.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
