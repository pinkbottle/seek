package wiki

import (
	"context"
	"fmt"
	"sync"

	"github.com/gocolly/colly"
	"github.com/pinkbottle/seek"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	crawled = promauto.NewCounter(prometheus.CounterOpts{
		Name: "crawler_crawled",
		Help: "The total number of cralwed pages",
	})
)

type Sink struct {
	mu      sync.Mutex
	c       *colly.Collector
	visited map[string]struct{}
}

func NewSink(c colly.Collector, res chan<- seek.Resource) *Sink {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("p", func(h *colly.HTMLElement) {
		result := seek.Resource{
			Content: h.Text,
			URL:     h.Request.URL.String(),
		}
		res <- result
		crawled.Inc()
	})

	return &Sink{
		c:       &c,
		visited: map[string]struct{}{},
	}
}

func (s *Sink) Start(ctx context.Context, root string) error {
	s.c.OnRequest(func(r *colly.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		url := r.URL.String()
		if _, ok := s.visited[url]; ok {
			fmt.Println("skipping", url)
			r.Abort()
			return
		}
		s.visited[url] = struct{}{}
	})
	if err := s.c.Visit(root); err != nil {
		return fmt.Errorf("failed to start visitor: %w", err)
	}
	done := make(chan struct{}, 1)
	go func() {
		s.c.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil
	case <-done:
		return nil
	}
}
