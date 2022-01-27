package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/pinkbottle/seek"
	"github.com/pinkbottle/seek/wiki"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	publisher = flag.String("pub", "http://localhost:8080/", "publisher url")
	root      = flag.String("root", "https://en.wikipedia.org/wiki/Main_Page", "root url")
)

func main() {
	flag.Parse()

	http.Handle("/metrics", promhttp.Handler())

	c := colly.NewCollector(func(c *colly.Collector) {
		c.MaxDepth = 5
	})
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	res := make(chan seek.Resource)
	go func() {
		for r := range res {
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

			fmt.Printf("✅[%s] : \n%s\n", r.URL, r.Content)
		}
	}()

	go func() {
		log.Printf("starting metrics server on :8081")
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("metrics server: %v", err)
		}
	}()

	ws := wiki.NewSink(*c, res)
	if err := ws.Start(context.Background(), *root); err != nil {
		log.Fatalf("failed to start: %v", err)
	}

}
