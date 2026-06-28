package main

import (
	"context"
	"fmt"
	"math"
	"sync"
	"time"
)

var limiter *TokenBucket

type TokenBucket struct {
	lastRefill time.Time
	mu         sync.Mutex
	tpsRate    float64
	tpsTokens  float64
}

func NewTokenBucket(tpsRate float64) *TokenBucket {
	t := TokenBucket{
		tpsTokens:  tpsRate,
		tpsRate:    tpsRate,
		lastRefill: time.Now(),
	}

	return &t
}

// must be called inside critical section
func (t *TokenBucket) refill() {
	curr := time.Now()
	timeDiff := curr.Sub(t.lastRefill).Seconds()

	tokens := t.tpsRate * timeDiff
	//fmt.Printf("refill %f\n", tokens)

	t.tpsTokens = min(t.tpsRate, t.tpsTokens+tokens)
	t.lastRefill = curr
}

func (t *TokenBucket) Wait(ctx context.Context) {
	for {
		t.mu.Lock()
		t.refill()

		if t.tpsTokens >= 1.0 {
			t.tpsTokens--
			t.mu.Unlock()
			return
		}

		missingTokens := 1 - t.tpsTokens
		duration := time.Duration(math.Ceil(missingTokens / t.tpsRate * float64(time.Second)))

		//println(t.tpsTokens)
		//println(duration)

		t.mu.Unlock()

		select {
		case <-time.After(duration):
		case <-ctx.Done():
			println("Context cancelled")
			return
		}
	}
}

type Fetcher interface {
	Fetch(url string) (string, []string, error)
}

func Crawl(ctx context.Context, url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	limiter.Wait(ctx)

	l, links, err := fetcher.Fetch(url)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%d %s\n", depth, l)

	if depth == 1 {
		return
	}

	wg.Add(len(links))
	for _, l := range links {
		go Crawl(ctx, l, depth-1, wg)
	}
}

func main() {
	var wg sync.WaitGroup

	limiter = NewTokenBucket(1)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	wg.Add(1)

	Crawl(ctx, "http://golang.org/", 4, &wg)

	wg.Wait()
}
