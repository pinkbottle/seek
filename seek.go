package seek

import "context"

type Type int

const (
	SearchFuzzy Type = iota
	SearchSentence
)

// Seeker searches indexes for resources matching a given query
type Seeker interface {
	Seek(ctx context.Context, query string, queryType Type) ([]*Result, error)
}

// Resource represents a piece of content that was found in some source
type Resource struct {
	Content string
	URL     string
}

// Result represents a search result, containing the content and the score matching the search query
// Result is indexed by the URL of the resource
type Result struct {
	Content string
	URL     string
	Score   float64
}
