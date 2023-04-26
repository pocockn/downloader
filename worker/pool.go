package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/store"
)

// Now is used, so we can fix the time within our tests.
var Now = time.Now()

// Pool holds the max amount of workers, a channel that we'll send our URLs down and our store.
type Pool struct {
	maxWorker int
	urls      chan models.URL
	store     store.Store
	mu        sync.Mutex
}

// NewPool creates a new worker pool.
func NewPool(maxWorkers int, s store.Store, urls chan models.URL) *Pool {
	return &Pool{
		maxWorker: maxWorkers,
		urls:      urls,
		store:     s,
		mu:        sync.Mutex{},
	}
}

// AddURL add a url to the pool to be processed by the workers.
func (p *Pool) AddURL(url models.URL) {
	p.urls <- url
}

// Run starts our workers and listens on the channel that the URLs are sent down.
// If we encounter an error we log it and discard the URL.
func (p *Pool) Run() {
	// ensure we don't exit before all the Go routines have finished processing.
	var wg sync.WaitGroup
	wg.Add(p.maxWorker)
	seenURLs := make(map[string]bool)

	for i := 0; i < p.maxWorker; i++ {
		go func(workerID int) {
			defer wg.Done()
			fmt.Printf("starting worker %d \n", workerID)
			for url := range p.urls {
				// ensure we only process URLs once.
				p.mu.Lock()
				if seenURLs[url.URL] {
					p.mu.Unlock()
					continue
				}
				seenURLs[url.URL] = true
				p.mu.Unlock()

				if err := Process(url, p.store); err != nil {
					fmt.Printf("unable to process %s from worker %d : %+v \n", url.URL, workerID, err)
					continue
				}
				fmt.Printf("processed URL %s via worker %d \n", url.URL, workerID)
			}
		}(i)
	}

	wg.Wait()
}

// Process takes a URL and performs a GET request against the URL. If the GET request isn't successful we discard the
// URL and log the error. If it is successful we store the URL in the store.
func Process(url models.URL, store store.Store) error {
	fmt.Printf("downloading %s...\n", url.URL)
	resp, err := http.Get(url.URL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	fmt.Printf("successfully downloaded %s \n", url.URL)

	result, err := store.Get(url.URL)
	if err != nil {
		return fmt.Errorf("unable to fetch %s", url.URL)
	}

	if result == nil {
		url.Submitted = 1
		url.CreatedAt = Now.UTC()
		bytes, err := json.Marshal(url)
		if err != nil {
			return fmt.Errorf("unable to marshal URL into bytes")
		}
		fmt.Printf("first time we have seen url %s storing in the db \n", url.URL)
		return store.Set(url.URL, bytes)
	}

	if err := json.Unmarshal(result, &url); err != nil {
		return fmt.Errorf("unable to unmarshal bytes into URL")
	}

	url.Submitted++
	url.UpdatedAt = Now.UTC()
	fmt.Printf("seen url %s %d times, updating \n", url.URL, url.Submitted)
	bytes, err := json.Marshal(url)
	if err != nil {
		return fmt.Errorf("unable to marshal URL into bytes")
	}

	return store.Set(url.URL, bytes)
}
