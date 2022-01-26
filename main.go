package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gocolly/colly"
	"github.com/pinkbottle/seek/seek"
	"github.com/pinkbottle/seek/wiki"
)

func main() {
	c := colly.NewCollector(func(c *colly.Collector) {
		c.MaxDepth = 5
	})
	res := make(chan seek.Result)
	go func() {
		for r := range res {
			fmt.Printf("[%s] : \n%s\n", r.URL, r.Content)
		}
	}()

	ws := wiki.NewSink(*c, res)
	if err := ws.Start(context.Background(), "https://en.wikipedia.org/wiki/Main_Page"); err != nil {
		log.Fatalf("failed to start: %v", err)
	}
}
