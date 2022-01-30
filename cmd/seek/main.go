package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/elastic/go-elasticsearch"
	"github.com/pinkbottle/seek"
)

var (
	searchAddr = "http://localhost:9200/seek/_search"
	index      = "seek"
)

func main() {
	args := os.Args
	if len(args) == 1 {
		fmt.Println("Usage: seek <keyword/sentence>")
		return
	}

	input := args[1:]
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}
	s := &Search{
		client: es,
	}
	results := make([]*seek.Result, 0)

	for _, phrase := range input {
		results = append(results, s.search(phrase)...)
	}
	resLen := len(results)
	if resLen == 0 {
		fmt.Println("No results found")
		return
	}

	show := 3
	if resLen < show {
		show = resLen
	}

	//sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	for _, r := range results[0:show] {
		fmt.Printf("%s (%f)\n\n%s\n%s\n\n", r.URL, r.Score, r.Content, strings.Repeat(".", 37))
	}
}

type Search struct {
	client *elasticsearch.Client
}

func (s *Search) search(phrase string) []*seek.Result {
	var buf bytes.Buffer
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"fuzzy": map[string]interface{}{
				"Content": map[string]interface{}{
					"value": phrase,
				},
			},
		},
	}
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(index),
		s.client.Search.WithBody(&buf),
		s.client.Search.WithTrackTotalHits(true),
		s.client.Search.WithPretty(),
	)
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			log.Fatalf("Error parsing the response body: %s", err)
		} else {
			// Print the response status and error information.
			log.Fatalf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	var results []*seek.Result
	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {
		source := hit.(map[string]interface{})["_source"]
		url := source.(map[string]interface{})["URL"]
		content := source.(map[string]interface{})["Content"]
		score := hit.(map[string]interface{})["_score"]

		results = append(results, &seek.Result{
			URL:     url.(string),
			Content: content.(string),
			Score:   score.(float64),
		})
	}

	return results
}
