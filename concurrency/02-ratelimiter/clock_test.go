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
