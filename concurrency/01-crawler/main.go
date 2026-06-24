package main

import (
	"fmt"
	"sync"
	"time"
)

var limiter chan struct{}

type Fetcher interface {
	Fetch(url string) (string, []string, error)
}

func Crawl(url string, depth int, wg *sync.WaitGroup) {
	defer wg.Done()

	<-limiter

	f := fetcher

	l, links, err := f.Fetch(url)
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
		go Crawl(l, depth-1, wg)
	}
}

func main() {
	var wg sync.WaitGroup

	limiter = make(chan struct{}, 1)
	go func() {
		for {
			time.Sleep(1 * time.Second)
			limiter <- struct{}{}
		}
	}()

	wg.Add(1)

	Crawl("http://golang.org/", 4, &wg)

	wg.Wait()
}
