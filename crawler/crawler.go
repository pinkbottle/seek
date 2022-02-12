package crawler

import (
	"context"

	"github.com/pinkbottle/seek"
)

// Crawler crawls a given root URL and returns the results as a channel of seek.Resource
// This interface is useless now, might be deleted in the future
type Crawler interface {
	Crawl(ctx context.Context, root string) (<-chan *seek.Resource, error)
}
