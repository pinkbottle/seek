package crawler

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

type HTTPCrawler struct {
	mu      sync.Mutex
	c       *colly.Collector
	rl      ratelimit.Limiter
	visited map[string]struct{}
}

func NewCrawler(maxDepth, maxRps int) *HTTPCrawler {
	return &HTTPCrawler{
		c: colly.NewCollector(func(c *colly.Collector) {
			c.MaxDepth = maxDepth
			c.Async = true
		}),
		visited: map[string]struct{}{},
		rl:      ratelimit.New(maxRps),
	}
}

func (s *HTTPCrawler) Crawl(ctx context.Context, root string) (<-chan *seek.Resource, error) {
	res := make(chan *seek.Resource)
	err := s.setupCollector(res)
	if err != nil {
		return nil, fmt.Errorf("failed to setup collector: %w", err)
	}
	if err := s.c.Visit(root); err != nil {
		return nil, fmt.Errorf("failed to start visitor: %w", err)
	}
	return res, nil
}

func (s *HTTPCrawler) setupCollector(res chan<- *seek.Resource) error {
	s.c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		_ = e.Request.Visit(e.Attr("href"))
	})

	s.c.OnHTML("p", func(h *colly.HTMLElement) {
		result := &seek.Resource{
			Content: h.Text,
			URL:     h.Request.URL.String(),
		}
		log.Printf("🔗[%s] : \n%s\n", result.URL, result.Content)
		res <- result
		crawled.Inc()
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

	return nil
}
