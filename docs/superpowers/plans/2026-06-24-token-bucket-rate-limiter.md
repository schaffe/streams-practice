# Token Bucket Rate Limiter Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build an interview-ready token bucket rate limiter in Go with a non-blocking `Allow`, a blocking `Wait`, lazy refill, and deterministic fake-clock tests.

**Architecture:** Single `TokenBucket` struct guarded by a `sync.Mutex`. Tokens are refilled lazily — recomputed from elapsed time on each call rather than by a background goroutine. A small `Clock` interface is injected so tests advance time deterministically.

**Tech Stack:** Go 1.26 (module `streams-practice`), standard library only (`sync`, `time`, `context`, `testing`).

## Global Constraints

- Package `ratelimiter` lives in `concurrency/02-ratelimiter/`.
- Standard library only — no third-party dependencies.
- `tokens` is `float64` to carry fractional accrual.
- The mutex is never held across a sleep.
- Match existing repo style (tabs, gofmt).

---

### Task 1: Clock seam

**Files:**
- Create: `concurrency/02-ratelimiter/clock.go`
- Test: `concurrency/02-ratelimiter/clock_test.go`

**Interfaces:**
- Produces: `Clock interface { Now() time.Time }`; `realClock struct{}` with `func (realClock) Now() time.Time`; `fakeClock struct{ t time.Time }` with `func (c *fakeClock) Now() time.Time` and `func (c *fakeClock) advance(d time.Duration)`.

- [ ] **Step 1: Write the failing test**

```go
package ratelimiter

import (
	"testing"
	"time"
)

func TestFakeClockAdvance(t *testing.T) {
	start := time.Unix(0, 0)
	c := &fakeClock{t: start}
	if got := c.Now(); !got.Equal(start) {
		t.Fatalf("Now() = %v, want %v", got, start)
	}
	c.advance(500 * time.Millisecond)
	if got := c.Now(); !got.Equal(start.Add(500 * time.Millisecond)) {
		t.Fatalf("after advance Now() = %v, want %v", got, start.Add(500*time.Millisecond))
	}
}

func TestRealClockMovesForward(t *testing.T) {
	c := realClock{}
	a := c.Now()
	b := c.Now()
	if b.Before(a) {
		t.Fatalf("real clock went backwards: %v then %v", a, b)
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./concurrency/02-ratelimiter/ -run TestFakeClock -v`
Expected: FAIL — `undefined: fakeClock`.

- [ ] **Step 3: Write minimal implementation**

```go
package ratelimiter

import "time"

// Clock is the time source. Production uses realClock; tests use fakeClock.
type Clock interface {
	Now() time.Time
}

type realClock struct{}

func (realClock) Now() time.Time { return time.Now() }

type fakeClock struct{ t time.Time }

func (c *fakeClock) Now() time.Time { return c.t }

func (c *fakeClock) advance(d time.Duration) { c.t = c.t.Add(d) }
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./concurrency/02-ratelimiter/ -run 'TestFakeClock|TestRealClock' -v`
Expected: PASS.

- [ ] **Step 5: Commit**

```bash
git add concurrency/02-ratelimiter/clock.go concurrency/02-ratelimiter/clock_test.go
git commit -m "feat(ratelimiter): add injectable Clock seam"
```

---

### Task 2: TokenBucket with lazy refill + AllowN

**Files:**
- Create: `concurrency/02-ratelimiter/bucket.go`
- Test: `concurrency/02-ratelimiter/bucket_test.go`

**Interfaces:**
- Consumes: `Clock`, `realClock`, `fakeClock` from Task 1.
- Produces:
  - `type TokenBucket struct { ... }` with unexported fields `capacity, tokens, refillRate float64`, `last time.Time`, `clock Clock`, `mu sync.Mutex`.
  - `func New(capacity, refillRate float64) *TokenBucket` — starts full, uses `realClock`.
  - `func newWithClock(capacity, refillRate float64, c Clock) *TokenBucket` — test constructor, starts full.
  - `func (b *TokenBucket) refill()` — caller holds `mu`.
  - `func (b *TokenBucket) Allow() bool`.
  - `func (b *TokenBucket) AllowN(n float64) bool`.

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./concurrency/02-ratelimiter/ -run TestAllow -v`
Expected: FAIL — `undefined: newWithClock` / `TokenBucket`.

- [ ] **Step 3: Write minimal implementation**

```go
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
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./concurrency/02-ratelimiter/ -run TestAllow -v`
Expected: PASS (all four tests).

- [ ] **Step 5: Commit**

```bash
git add concurrency/02-ratelimiter/bucket.go concurrency/02-ratelimiter/bucket_test.go
git commit -m "feat(ratelimiter): add TokenBucket with lazy refill and AllowN"
```

---

### Task 3: Blocking WaitN / Wait

**Files:**
- Modify: `concurrency/02-ratelimiter/bucket.go`
- Test: `concurrency/02-ratelimiter/wait_test.go`

**Interfaces:**
- Consumes: `TokenBucket`, `refill` from Task 2.
- Produces:
  - `func (b *TokenBucket) Wait(ctx context.Context) error`.
  - `func (b *TokenBucket) WaitN(ctx context.Context, n float64) error`.
  - Sentinel: `var ErrExceedsCapacity = errors.New("ratelimiter: requested tokens exceed capacity")`.

- [ ] **Step 1: Write the failing test**

```go
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
```

- [ ] **Step 2: Run test to verify it fails**

Run: `go test ./concurrency/02-ratelimiter/ -run TestWaitN -v`
Expected: FAIL — `undefined: WaitN` / `ErrExceedsCapacity`.

- [ ] **Step 3: Write minimal implementation**

Add to the top of `bucket.go` imports: `"context"` and `"errors"`. Add the sentinel and methods:

```go
// ErrExceedsCapacity is returned by WaitN when n is larger than capacity and
// could therefore never be satisfied.
var ErrExceedsCapacity = errors.New("ratelimiter: requested tokens exceed capacity")

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
	b.tokens -= n      // reserve (may go negative)
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
```

Note: the `import` block in `bucket.go` becomes:

```go
import (
	"context"
	"errors"
	"sync"
	"time"
)
```

- [ ] **Step 4: Run test to verify it passes**

Run: `go test ./concurrency/02-ratelimiter/ -run TestWaitN -v`
Expected: PASS (all four tests).

- [ ] **Step 5: Commit**

```bash
git add concurrency/02-ratelimiter/bucket.go concurrency/02-ratelimiter/wait_test.go
git commit -m "feat(ratelimiter): add blocking WaitN with ctx cancellation"
```

---

### Task 4: Concurrency safety test under -race

**Files:**
- Test: `concurrency/02-ratelimiter/race_test.go`

**Interfaces:**
- Consumes: `New`, `Allow` from Task 2.

- [ ] **Step 1: Write the failing test**

```go
package ratelimiter

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestAllowConcurrent(t *testing.T) {
	b := New(100, 1_000_000) // generous rate; we only assert no race / no over-grant beyond reason
	const goroutines = 50
	const each = 20
	var granted int64
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < each; j++ {
				if b.Allow() {
					atomic.AddInt64(&granted, 1)
				}
			}
		}()
	}
	wg.Wait()
	if granted == 0 {
		t.Fatal("expected some grants")
	}
	if granted > goroutines*each {
		t.Fatalf("granted %d exceeds total attempts", granted)
	}
}
```

- [ ] **Step 2: Run test to verify it fails (or passes) cleanly under -race**

Run: `go test ./concurrency/02-ratelimiter/ -run TestAllowConcurrent -race -v`
Expected: PASS with no race detected. (If the data race detector reports anything, the mutex usage in Task 2 must be fixed.)

- [ ] **Step 3: Run the full package under -race**

Run: `go test ./concurrency/02-ratelimiter/ -race`
Expected: `ok` — all tests pass, no races.

- [ ] **Step 4: Commit**

```bash
git add concurrency/02-ratelimiter/race_test.go
git commit -m "test(ratelimiter): add concurrent Allow test under -race"
```

---

### Task 5: Package doc + runnable example

**Files:**
- Create: `concurrency/02-ratelimiter/doc.go`
- Test: `concurrency/02-ratelimiter/example_test.go`

**Interfaces:**
- Consumes: `New`, `Allow` from Task 2.

- [ ] **Step 1: Write the example test**

```go
package ratelimiter_test

import (
	"fmt"

	ratelimiter "streams-practice/concurrency/02-ratelimiter"
)

func ExampleTokenBucket_Allow() {
	b := ratelimiter.New(2, 1) // burst 2, 1 token/sec
	fmt.Println(b.Allow())
	fmt.Println(b.Allow())
	fmt.Println(b.Allow()) // empty
	// Output:
	// true
	// true
	// false
}
```

- [ ] **Step 2: Run the example to verify output matches**

Run: `go test ./concurrency/02-ratelimiter/ -run ExampleTokenBucket_Allow -v`
Expected: PASS (output matches).

- [ ] **Step 3: Write the package doc**

```go
// Package ratelimiter implements a token bucket rate limiter.
//
// A bucket holds up to capacity tokens and refills at refillRate tokens per
// second. Capacity bounds burst size; refillRate bounds the sustained average.
// Tokens are refilled lazily — recomputed from elapsed time on each call —
// so no background goroutine is needed.
//
// Use Allow/AllowN for non-blocking checks and Wait/WaitN to block until
// tokens are available.
package ratelimiter
```

- [ ] **Step 4: Run the whole package once more**

Run: `go test ./concurrency/02-ratelimiter/ -race`
Expected: `ok`.

- [ ] **Step 5: Commit**

```bash
git add concurrency/02-ratelimiter/doc.go concurrency/02-ratelimiter/example_test.go
git commit -m "docs(ratelimiter): add package doc and runnable example"
```
