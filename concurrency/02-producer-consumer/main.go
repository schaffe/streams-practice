package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type producer struct {
	c  chan<- string
	id string
}

func (p *producer) produce(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	counter := 0

	for {
		select {
		case <-ctx.Done():
			fmt.Printf("producer [%s] aborted\n", p.id)
			return
		case p.c <- fmt.Sprintf("[%s]: count %d", p.id, counter):
			counter++
			time.Sleep(1000 * time.Millisecond)
		}
	}
}

type consumer struct {
	c  <-chan string
	id string
}

func (c *consumer) consume(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("consumer [%s] aborted\n", c.id)
			return

		case msg := <-c.c:
			fmt.Printf("[%s]%s\n", c.id, msg)
		}
	}
}

func main() {
	ch := make(chan string, 5)

	p := &producer{
		c:  ch,
		id: "p1",
	}

	c := &consumer{
		c:  ch,
		id: "c1",
	}

	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Go(func() {
		c.consume(ctx)
	})
	go p.produce(ctx, &wg)

	time.AfterFunc(3*time.Second, func() {
		cancel()
		close(ch)
	})

	wg.Wait()

}
