package main

import (
	"fmt"
	"sync"
	"time"
)

var limiter <-chan time.Time

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

	limiter = time.Tick(time.Second)

	wg.Add(1)

	Crawl("http://golang.org/", 4, &wg)

	wg.Wait()
}
