package wiki

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type worker struct {
	mu      sync.Mutex
	done    <-chan struct{}
	ticker  *time.Ticker
	visited map[string]bool
	jobs    chan string
	f       fetcher
	id      int
}

func newWorker(id int, jobs chan string, done <-chan struct{}, f fetcher, visited map[string]bool, interval time.Duration) *worker {
	return &worker{
		jobs:    jobs,
		done:    done,
		f:       f,
		id:      id,
		ticker:  time.NewTicker(interval),
		visited: visited,
	}
}

type fetcher struct {
	c http.Client
}

func newFetcher(c http.Client) fetcher {
	return fetcher{c: c}
}

func (f *fetcher) fetch(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for url %s: %v", url, err)
	}

	res, err := f.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("expected 200 status code, instead got: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	urls, err := findURLs(string(b))
	if err != nil {
		return nil, fmt.Errorf("error finding urls: %v", err)
	}
	return urls, nil
}

func (w *worker) run(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case job := <-w.jobs:
			<-w.ticker.C
			log.Printf("worker %d: crawling %s", w.id, job)
			links, err := w.f.fetch(ctx, job)
			if err != nil {
				log.Printf("worker %d: error crawling %s: %v", w.id, job, err)
			}
			for _, link := range links {
				link := link
				if w.seen(link) {
					continue
				}

				go func() {
					select {
					case w.jobs <- link:
						log.Printf("worker %d: enqueued %s", w.id, link)
					default:
						return
					}
				}()
			}
			// log.Printf("worker %d: finished crawling %s, added %d new links", w.id, job, len(links))
		}
	}
}

func (w *worker) seen(link string) bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	if _, ok := w.visited[link]; ok {
		return true
	}
	w.visited[link] = true
	return false
}
