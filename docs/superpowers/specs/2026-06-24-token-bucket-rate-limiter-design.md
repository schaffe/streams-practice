# Token Bucket Rate Limiter (Go) — Design

**Goal:** Interview-prep exercise. Build a token bucket rate limiter that is
correct, clearly explainable, and not over-engineered. Optimize for being able
to whiteboard and defend every design decision.

**Location:** `concurrency/02-ratelimiter/`, package `ratelimiter` (sibling to
the existing crawler exercise).

## The mental model

The bucket holds tokens. Two forces act on it:

- **Tokens drip in** at `refillRate` per second, but the bucket overflows at
  `capacity` — extra tokens are discarded.
- **Each request removes** a token. No token → request denied.

Two independent knobs:

- `capacity` → **burst size**: how many requests can fire back-to-back after an
  idle period.
- `refillRate` → **sustained rate**: the long-run average throughput.

Separating burst from sustained rate is the whole point of token bucket
(contrast: leaky bucket enforces a strictly constant output with no bursts).

## Lazy refill (the core trick)

Do **not** run a goroutine adding tokens on a ticker. Instead, recompute tokens
only when someone asks. Store `last` (instant of last recompute) and `tokens`
(value at that instant). On any call:

```
elapsed = now - last
tokens  = min(capacity, tokens + elapsed * refillRate)
last    = now
```

Time itself is the source of truth; `tokens` is a cached snapshot. Idle for an
hour? The next call adds an hour's worth of tokens in one line (clamped to
capacity). Zero goroutines, exact accuracy.

**Ticker alternative (discuss, don't build):** a background goroutine ticking
tokens wastes a goroutine per limiter, has granularity bounded by the tick
interval, and still needs the same clamp. Lazy refill is strictly better for a
single limiter.

## Worked example

`capacity = 5`, `refillRate = 2`/sec (one token every 500ms), starts full.

| Time | Event | Refill adds | tokens before | tokens after | Result |
|------|-------|------------|---------------|--------------|--------|
| 0.0s | Allow() | 0 | 5 | 4 | ✅ |
| 0.0s | Allow() ×4 | 0 | 4→1 | 0 | ✅✅✅✅ (burst of 5) |
| 0.0s | Allow() | 0 | 0 | 0 | ❌ |
| 0.3s | Allow() | 0.6 | 0 | 0.6 | ❌ (0.6 < 1) |
| 0.5s | Allow() | 0.4 | 0.6 | 1.0 → 0 | ✅ |
| 10s  | Allow() | 19, clamp to 5 | 0 | 5 → 4 | ✅ |

Notes:
- **Fractional tokens** (0.6) are why `tokens` is a `float64`: partial accrual
  must carry over, but a request is granted only at a whole token.
- **Clamp at 10s**: idle time does not bank unlimited burst; capacity is the ceiling.

## Architecture / components

```
Clock interface        ← seam for deterministic tests; one method Now() time.Time
  realClock{}          ← time.Now in prod
  fakeClock{}          ← test-controlled, advanceable

TokenBucket
  capacity   float64   (max tokens = burst size)
  tokens     float64   (current; fractional)
  refillRate float64   (tokens/sec)
  last       time.Time (last refill instant)
  clock      Clock
  mu         sync.Mutex
```

The `Clock` interface is the reason fake-clock tests work: prod injects
`realClock`, tests inject a `fakeClock` they advance by exact durations.

## Public API

- `New(capacity, refillRate float64) *TokenBucket` — starts full.
- `Allow() bool` → `AllowN(1)`.
- `AllowN(n float64) bool` — lazily refill, then consume `n` if available;
  non-blocking core.
- `Wait(ctx context.Context) error` → `WaitN(ctx, 1)`.
- `WaitN(ctx, n)` — same refill math; compute deficit, sleep until enough
  tokens accrue (or `ctx` cancels), then consume.

### Allow (non-blocking)

```
Allow():
  lock; defer unlock
  refill()
  if tokens >= 1: tokens -= 1; return true
  return false
```

### Wait (blocking, built on the same math)

```
WaitN(ctx, n):
  if n > capacity: return error            // impossible forever
  lock
  refill()
  if tokens >= n: tokens -= n; unlock; return nil   // fast path
  deficit = n - tokens
  delay   = deficit / refillRate
  tokens -= n                              // reserve now (may go negative)
  last    = now
  unlock                                   // release BEFORE sleeping
  select:
    case <-time.After(delay): return nil
    case <-ctx.Done():        return ctx.Err()
```

Three subtleties to say out loud:

1. **Release the lock before sleeping** — holding it across the sleep would
   serialize all waiters. Compute delay under lock, sleep lock-free.
2. **Reserve by going negative** — subtracting `n` immediately stops concurrent
   waiters from claiming the same future token; each reservation pushes
   `last`/`tokens` forward so the next caller computes a later delay. `refill`
   works off `last` regardless of sign, so the math stays consistent.
3. **`ctx` cancellation** lets a caller bail out — essential for request timeouts.

## Thread safety

One `sync.Mutex` guards `tokens` + `last`. Every public method refills-then-acts
under the lock (except the sleep in `WaitN`). A lock-free atomic version is a
"push further" answer, not the lead.

## Error handling / edge cases

- `WaitN` with `n > capacity` → error immediately (never satisfiable).
- `ctx` already cancelled, or cancels during sleep → return `ctx.Err()`,
  consume nothing on the cancellation path.
- Mutex never held across the sleep.

## Testing

- **Table tests + fake clock** driving `AllowN`: full-bucket burst, exhaustion
  returns false, advance time → tokens return, clamp at capacity, fractional
  carry-over.
- **WaitN tests:** success after accrual, `ctx` cancellation, `n > capacity`
  error.
- **Concurrency test:** hammer `Allow()` from many goroutines under `-race`.
