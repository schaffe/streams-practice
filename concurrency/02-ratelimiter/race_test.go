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
