package ratelimiter

import (
	"testing"
	"time"
)

func newTestBucket(cap, rate float64) (*TokenBucket, *fakeClock) {
	c := &fakeClock{t: time.Unix(0, 0)}
	return newWithClock(cap, rate, c), c
}

func TestAllowBurstThenExhaust(t *testing.T) {
	b, _ := newTestBucket(5, 2) // full = 5
	for i := 0; i < 5; i++ {
		if !b.Allow() {
			t.Fatalf("burst call %d should be allowed", i)
		}
	}
	if b.Allow() {
		t.Fatal("6th call should be denied (bucket empty)")
	}
}

func TestRefillOverTime(t *testing.T) {
	b, clock := newTestBucket(5, 2) // 2 tokens/sec
	for i := 0; i < 5; i++ {
		b.Allow()
	}
	clock.advance(300 * time.Millisecond) // 0.6 tokens
	if b.Allow() {
		t.Fatal("0.6 tokens should not satisfy a request")
	}
	clock.advance(200 * time.Millisecond) // now 1.0 token total
	if !b.Allow() {
		t.Fatal("1.0 token should satisfy a request")
	}
}

func TestRefillClampsAtCapacity(t *testing.T) {
	b, clock := newTestBucket(5, 2)
	for i := 0; i < 5; i++ {
		b.Allow()
	}
	clock.advance(10 * time.Second) // would be 20 tokens, clamp to 5
	for i := 0; i < 5; i++ {
		if !b.Allow() {
			t.Fatalf("post-idle call %d should be allowed", i)
		}
	}
	if b.Allow() {
		t.Fatal("idle time must not bank more than capacity")
	}
}

func TestAllowNConsumesMultiple(t *testing.T) {
	b, _ := newTestBucket(5, 2)
	if !b.AllowN(5) {
		t.Fatal("AllowN(5) on full bucket should succeed")
	}
	if b.AllowN(1) {
		t.Fatal("bucket should be empty after AllowN(5)")
	}
}
