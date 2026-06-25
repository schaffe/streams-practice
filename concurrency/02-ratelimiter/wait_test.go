package ratelimiter

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWaitNFastPath(t *testing.T) {
	b, _ := newTestBucket(5, 2)
	if err := b.WaitN(context.Background(), 3); err != nil {
		t.Fatalf("fast path should not error: %v", err)
	}
}

func TestWaitNExceedsCapacity(t *testing.T) {
	b, _ := newTestBucket(5, 2)
	if err := b.WaitN(context.Background(), 6); !errors.Is(err, ErrExceedsCapacity) {
		t.Fatalf("want ErrExceedsCapacity, got %v", err)
	}
}

func TestWaitNCancelled(t *testing.T) {
	b, _ := newTestBucket(1, 1) // real-clock bucket via newTestBucket uses fakeClock;
	// drain it so Wait must block.
	b.AllowN(1)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.WaitN(ctx, 1); !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
}

func TestWaitNSucceedsAfterAccrual(t *testing.T) {
	// Use a real-clock bucket so time.After actually fires.
	b := New(1, 100) // 100 tokens/sec => ~10ms per token
	b.AllowN(1)      // drain
	start := time.Now()
	if err := b.WaitN(context.Background(), 1); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if elapsed := time.Since(start); elapsed < 5*time.Millisecond {
		t.Fatalf("expected to wait for accrual, only waited %v", elapsed)
	}
}
