package wiki

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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
	Batchcomplete string `json:"batchcomplete"`
	Continue      struct {
		Grncontinue string `json:"grncontinue"`
		Continue    string `json:"continue"`
	} `json:"continue"`
	Query struct {
		Pages map[string]struct {
			Title string `json:"title"`
		} `json:"pages"`
	} `json:"query"`
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
