package wiki

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	exp = regexp.MustCompile(`<a href="(.*?)"`)
)

type Sink struct {
	c          *http.Client
	maxWorkers int
}

func NewSink(c *http.Client, maxWorkers int) (*Sink, error) {
	return &Sink{c: c, maxWorkers: maxWorkers}, nil
}

// Explore takes root of a content tree(be it root link for example) and explores it
func (s *Sink) Explore(ctx context.Context, root string) error {
	//fetch root url and get all links in it
	log.Printf("bootstrapping %d workers for %s", s.maxWorkers, root)
	workers := make([]*worker, s.maxWorkers)
	jobs := make(chan string)
	visited := make(map[string]bool)
	var g errgroup.Group
	for i := 0; i < s.maxWorkers; i++ {
		i := i
		f := newFetcher(*s.c)
		workers[i] = newWorker(i, jobs, ctx.Done(), f, visited, time.Second)
		g.Go(func() error { return workers[i].run(ctx) })
	}

	jobs <- root

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error waiting for goroutines: %v", err)
	}

	return nil
}

func findURLs(content string) ([]string, error) {
	matches := exp.FindAllStringSubmatch(content, -1)

	var links []string
	for _, m := range matches {
		link := m[1]
		if !strings.HasPrefix(m[1], "http") {
			link = "https://en.wikipedia.org" + link
		}
		links = append(links, link)
	}
	return links, nil
}
