package main

import (
	"encoding/base64"
	"flag"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/cespare/xxhash"
)

func main() {
	var numWorkers int
	var numChecks int

	flag.IntVar(&numWorkers, "w", 2, "Number of workers")
	flag.IntVar(&numChecks, "n", 0, "Number of uuid's to check")
	flag.IntVar(&numChecks, "b", 0, "Size of buffer")
	flag.Parse()

	c := make(chan string, numChecks)
	wg := sync.WaitGroup {}

	wg.Add(numWorkers)

	for i := 0; i < numWorkers; i++ {
		go worker(numChecks, c, &wg)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	m := make(map[string]int)

	for s := range c {
		m[s]++

		if m[s] > 1 {
			log.Printf("Found %d collisions for %s", m[s], s)
		}
	}
}

func worker(size int, c chan<- string, wg *sync.WaitGroup) {
	defer wg.Done()

	for i := 0; i < size; i++ {
		u := uuid.New()
		b, err := u.MarshalBinary()

		if err != nil {
			log.Printf("Failed to marshal UUID %s: %s", u.String(), err)
			continue
		}

		h := xxhash.New()
		_, err = h.Write(b)

		if err != nil {
			log.Printf("Failed to hash UUID %s: %s", u.String(), err)
			continue
		}

		c <- base64.URLEncoding.EncodeToString(h.Sum(nil))
	}
}
