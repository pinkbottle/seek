package wiki

import (
	"context"
	"fmt"

	"github.com/gocolly/colly"
	"github.com/pinkbottle/seek/seek"
)

type Sink struct {
	c *colly.Collector
}

func NewSink(c colly.Collector, res chan<- seek.Result) *Sink {
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnHTML("p", func(h *colly.HTMLElement) {
		result := seek.Result{
			Content: h.Text,
			URL:     h.Request.URL.String(),
		}
		res <- result
	})

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })

	return &Sink{
		c: &c,
	}
}

func (s *Sink) Start(ctx context.Context, root string) error {
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
