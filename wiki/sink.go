package wiki

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/pinkbottle/seek"
	"go.uber.org/ratelimit"

	"github.com/gocolly/colly"
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
	rl      ratelimit.Limiter
	visited map[string]struct{}
	res     chan<- *seek.Resource
}

func NewSink(maxDepth, maxRps int, res chan<- *seek.Resource) *Sink {
	return &Sink{
		c: colly.NewCollector(func(c *colly.Collector) {
			c.MaxDepth = maxDepth
			c.Async = true
		}),
		visited: map[string]struct{}{},
		res:     res,
		rl:      ratelimit.New(maxRps),
	}
}

func (s *Sink) Start(ctx context.Context, root string) error {
	err := s.setupCollector()
	if err != nil {
		return fmt.Errorf("failed to setup collector: %w", err)
	}

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

func (s *Sink) setupCollector() error {
	s.c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		_ = e.Request.Visit(e.Attr("href"))
	})

	s.c.OnHTML("p", func(h *colly.HTMLElement) {
		result := &seek.Resource{
			Content: h.Text,
			URL:     h.Request.URL.String(),
		}
		log.Printf("🔗[%s] : \n%s\n", result.URL, result.Content)
		s.res <- result
	})

	s.c.OnRequest(func(r *colly.Request) {
		s.rl.Take()
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

	s.c.OnScraped(func(_ *colly.Response) {
		crawled.Inc()
	})

	return nil
}
