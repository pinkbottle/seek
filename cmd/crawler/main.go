package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocolly/colly"
	"github.com/pinkbottle/seek/seek"
	"github.com/pinkbottle/seek/wiki"
)

var (
	publisher = "http://localhost:8080/publish"
)

func main() {
	c := colly.NewCollector(func(c *colly.Collector) {
		c.MaxDepth = 5
	})
	client := http.Client{
		Timeout: 3 * time.Second,
	}

	res := make(chan seek.Result)
	go func() {
		for r := range res {
			b, err := json.Marshal(r)
			if err != nil {
				log.Printf("failed to marshal result: %v", err)
				continue
			}

			req, err := http.NewRequestWithContext(context.Background(), "POST", publisher, bytes.NewReader(b))
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

	ws := wiki.NewSink(*c, res)
	if err := ws.Start(context.Background(), "https://en.wikipedia.org/wiki/Main_Page"); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
