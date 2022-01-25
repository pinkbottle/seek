package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

var (
	exp        = regexp.MustCompile(`<a href="(.*?)"`)
	concurrent = 10
)

type Sink struct {
	c *http.Client
}

func NewSink(c *http.Client) (*Sink, error) {
	return &Sink{c: c}, nil
}

type Article struct {
	Title string
}

type articleResponse struct {
	Query struct {
		Pages map[string]struct {
			Title string `json:"title"`
		} `json:"pages"`
	} `json:"query"`
	Continue struct {
		Grncontinue string `json:"grncontinue"`
		Continue    string `json:"continue"`
	} `json:"continue"`
	Batchcomplete string `json:"batchcomplete"`
}

func (s *Sink) GetRandomArticle(ctx context.Context) (*Article, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://en.wikipedia.org/w/api.php?action=query&generator=random&format=json", nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	res, err := s.c.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting response: %v", err)
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting response: %v", res.StatusCode)
	}
	var ar *articleResponse
	if err := json.NewDecoder(res.Body).Decode(&ar); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	var keys []string
	for k := range ar.Query.Pages {
		keys = append(keys, k)
	}
	if len(keys) == 0 {
		return nil, fmt.Errorf("no keys found")
	}

	return &Article{
		Title: ar.Query.Pages[keys[0]].Title,
	}, nil
}

// Explore takes root of a content tree(be it root link for example) and explores it
func (s *Sink) Explore(ctx context.Context, root string) error {
	//fetch root url and get all links in it
	log.Printf("bootstrapping %s", root)
	init, err := s.crawlPage(ctx, root)
	if err != nil {
		return fmt.Errorf("error bootstrapping root: %v", err)
	}

	log.Printf("found %d links", len(init))
	jobs := make(chan string, concurrent)
	var g errgroup.Group

	g.Go(func() error { return s.loop(ctx, jobs) })

	for _, url := range init {
		jobs <- url
	}
	log.Printf("starting %d workers", concurrent)

	if err := g.Wait(); err != nil {
		return fmt.Errorf("error waiting for goroutines: %v", err)
	}

	return nil
}

func (s *Sink) loop(ctx context.Context, jobs chan string) error {
	log.Printf("starting loop")
	var counter int32 = 0
	for {
		log.Println("waiting for job")
		select {
		case <-ctx.Done():
			log.Printf("loop: context done")
			return nil
		case job := <-jobs:
			log.Println(len(jobs))
			currentCount := atomic.AddInt32(&counter, 1)
			log.Printf("starting job %s, current workers %d", job, currentCount)
			go func() {
				links, err := s.crawlPage(ctx, job)
				if err != nil {
					log.Printf("loop: error crawling page %s: %v", job, err)
				}
				for _, link := range links {
					jobs <- link
					log.Printf("loop: added link %s", link)
				}
				log.Printf("loop: finished crawling page %s", job)
				atomic.AddInt32(&counter, -1)
			}()
		}
	}
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

func (s *Sink) crawlPage(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for url %s: %v", url, err)
	}

	res, err := s.c.Do(req)
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
