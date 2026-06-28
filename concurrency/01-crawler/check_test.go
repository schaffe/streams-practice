package main

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	fetchSig := fetchSignalInstance()
	errCh := make(chan struct{}, 1)

	start := time.Unix(0, 0)
	go func(start time.Time) {
		for {
			switch {
			case <-fetchSig:
				if time.Since(start).Nanoseconds() < 950000000 {
					errCh <- struct{}{}
					return
				}
				start = time.Now()
			}
		}
	}(start)

	main()

	select {
	case <-errCh:
		t.Log("There exists a two crawls that were executed less than 1 second apart.")
		t.Log("Solution is incorrect.")
		t.FailNow()
	default:
	}
}
