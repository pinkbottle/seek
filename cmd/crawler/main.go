package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"time"

	"github.com/pinkbottle/seek"
	"github.com/pinkbottle/seek/wiki"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/sync/errgroup"
)

var (
	publisher = flag.String("pub", "http://localhost:8080/", "publisher url")
	root      = flag.String("root", "https://en.wikipedia.org/wiki/Main_Page", "root url")
	rps       = flag.Int("rps", 10, "max requests per second")
	depth     = flag.Int("depth", 3, "max depth")
	ignore    = flag.String("ignore", "", "domains to ignore, separated by comma")
)

func main() {
	flag.Parse()

	client := http.Client{
		Timeout: 3 * time.Second,
	}

	r := make(chan *seek.Resource)
	ws := wiki.NewSink(3, 1, r)

	ctx := context.Background()
	var g errgroup.Group

	// start wiki sink crawler
	g.Go(func() error { return ws.Start(ctx, *root) })
	// start metrics server
	g.Go(func() error { return registerMetrics() })
	// flush results to elastic
	g.Go(func() error { return collect(ctx, r, &client) })

	if err := g.Wait(); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}

func collect(ctx context.Context, res <-chan *seek.Resource, client *http.Client) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case r := <-res:
			b, err := json.Marshal(r)
			if err != nil {
				log.Printf("failed to marshal result: %v", err)
				continue
			}

			req, err := http.NewRequestWithContext(context.Background(), "POST", *publisher, bytes.NewReader(b))
			if err != nil {
				log.Printf("failed to create request: %v", err)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
			res, err := client.Do(req)
			if err != nil {
				log.Printf("failed to send request: %v", err)
				continue
			}

			if res.StatusCode != http.StatusOK {
				log.Printf("failed to send request: %v", err)
				continue
			}

			log.Printf("✅ [%s] : \n%s\n", r.URL, r.Content)
		}
	}
}

func registerMetrics() error {
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":8081", nil)
}
