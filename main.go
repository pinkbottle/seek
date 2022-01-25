package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/pinkbottle/seek/wiki"
)

func main() {
	c := http.Client{Timeout: time.Second * 10}
	sink, err := wiki.NewSink(&c)
	if err != nil {
		log.Fatalf("error creating sink: %v", err)
	}

	ctx := context.Background()
	if err := sink.Explore(ctx, "https://en.wikipedia.org/wiki/2007_Vodacom_Challenge"); err != nil {
		log.Fatalf("error exploring: %v", err)
	}
}
