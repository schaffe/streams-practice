package ratelimiter

import (
	"context"
	"errors"
	"sync"
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

func TestWaitNCancelledRefundsTokens(t *testing.T) {
	b, _ := newTestBucket(5, 1)
	// Drain all tokens so WaitN must block.
	b.AllowN(5)
	// Cancelled context: select takes ctx.Done() branch immediately.
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := b.WaitN(ctx, 3); !errors.Is(err, context.Canceled) {
		t.Fatalf("want context.Canceled, got %v", err)
	}
	// Tokens should be refunded: balance should be back to 0 (not -3).
	b.mu.Lock()
	got := b.tokens
	b.mu.Unlock()
	if got < 0 {
		t.Fatalf("tokens not refunded after cancellation: got %v", got)
	}
}

func TestWaitNConcurrentWaitersNoDoubleClaim(t *testing.T) {
	// High rate so goroutines finish quickly.
	b := New(1, 1000) // 1 capacity, 1000 tokens/sec
	b.AllowN(1)       // drain

	const n = 5
	errs := make([]error, n)
	var wg sync.WaitGroup
	wg.Add(n)
	for i := 0; i < n; i++ {
		i := i
		go func() {
			defer wg.Done()
			errs[i] = b.WaitN(context.Background(), 1)
		}()
	}
	done := make(chan struct{})
	go func() { wg.Wait(); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("timed out waiting for concurrent waiters")
	}
	for i, err := range errs {
		if err != nil {
			t.Errorf("goroutine %d got unexpected error: %v", i, err)
		}
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
