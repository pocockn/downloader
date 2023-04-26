package watcher

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/pocockn/downloader/models"
	"github.com/pocockn/downloader/store"
)

// Watcher is used to watch the URLs being saved into the database. It will run a process
// function based on the interval passed in.
type Watcher struct {
	intervalDuration time.Duration
	store            store.Store

	successfulDownloads   int64
	unsuccessfulDownloads int64

	stop chan struct{}

	mu sync.RWMutex
}

// New returns a new watcher struct.
func New(i time.Duration, s store.Store) *Watcher {
	return &Watcher{
		intervalDuration: i,
		store:            s,
		mu:               sync.RWMutex{},
		stop:             make(chan struct{}),
	}
}

// Process performs the logic for the watcher. It triggers every n seconds based off the interval passed into the
// watchers constructor. It will fetch the 10 most submitted URLs then perform batch downloads of 3 URLs at a time.
// Once all URLs have been downloaded it prints the time taken and number of successful / unsuccessful downloads to stdout.
func (w *Watcher) Process() {
	fmt.Println("starting watcher...")
	ticker := time.NewTicker(w.intervalDuration)
	go func() {
		for {
			select {
			case <-ticker.C:
				results, err := w.store.GetAll()
				if err != nil {
					fmt.Println(err.Error())
					return
				}

				var urls []models.URL
				for _, d := range results {
					var url models.URL
					err := json.Unmarshal(d, &url)
					if err != nil {
						fmt.Println(err.Error())
						continue
					}
					urls = append(urls, url)
				}

				sort.Slice(urls, func(i, j int) bool {
					return urls[i].Submitted > urls[j].Submitted
				})

				// Create a wait group to wait for all the downloads to complete
				var wg sync.WaitGroup

				// Dummy channel to coordinate the number of concurrent goroutines.
				// Buffered channel, allows max 3 values
				concurrentGoroutines := make(chan struct{}, 3)

				if len(urls) > 10 {
					urls = urls[:9]
				}

				for _, url := range urls {
					wg.Add(1)
					concurrentGoroutines <- struct{}{}
					go func(url models.URL) {
						defer wg.Done()
						if err := w.downloadURL(url); err != nil {
							w.mu.Lock()
							w.unsuccessfulDownloads++
							w.mu.Unlock()
						}
						w.mu.Lock()
						w.successfulDownloads++
						w.mu.Unlock()
						// read from the channel, this will allow another URL to be processed.
						<-concurrentGoroutines
					}(url)
				}

				wg.Wait()
				fmt.Printf(
					"successfull downloads %d, unsuccessful downloads %d \n",
					w.successfulDownloads,
					w.unsuccessfulDownloads,
				)
			case <-w.stop:
				ticker.Stop()
				fmt.Println("Stopping watcher...")
				return
			}
		}
	}()
}

// Stop closes down the watcher
func (w *Watcher) Stop() {
	w.stop <- struct{}{}
}

// downloadURL performs a GET request to the URL passed in. We measure the time it takes to download the URL
// and then log the URLs stats to stdout.
func (w *Watcher) downloadURL(url models.URL) error {
	fmt.Printf("downloading %s...\n", url.URL)

	startTime := time.Now()

	resp, err := http.Get(url.URL)
	if err != nil {
		return fmt.Errorf("error downloading %s: %s\n", url.URL, err)
	}

	defer resp.Body.Close()
	elapsedTime := time.Since(startTime)
	fmt.Printf(
		"downloaded %s in %s \n",
		url.URL,
		elapsedTime,
	)

	url.UpdatedAt = time.Now()
	return nil
}
