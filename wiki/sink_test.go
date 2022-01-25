package wiki_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/pinkbottle/seek/wiki"
)

func TestExplore(t *testing.T) {
	sink, err := wiki.NewSink(&http.Client{
		Timeout: time.Second * 10,
	})
	if err != nil {
		t.Errorf("error creating sink: %v", err)
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err = sink.Explore(ctx, "https://en.wikipedia.org/wiki/The_Marketts")
	if err != nil {
		t.Errorf("error exploring: %v", err)
	}
}
