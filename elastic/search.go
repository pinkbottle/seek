package elastic

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"

	"github.com/elastic/go-elasticsearch"
	"github.com/pinkbottle/seek"
)

type searchType int

const (
	SearchFuzzy searchType = iota
	SearchSentence
)

type Search struct {
	client *elasticsearch.Client
	index  string
}

func NewSearch(index string) (*Search, error) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil, fmt.Errorf("error creating the client: %s", err)
	}

	return &Search{
		client: es,
		index:  index,
	}, nil

}

func (s *Search) search(st searchType, phrase string) []*seek.Result {
	var buf bytes.Buffer
	query, err := getQuery(st, phrase)
	if err != nil {
		log.Printf("error creating the query: %s", err)
		return nil
	}

	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	// Perform the search request.
	res, err := s.client.Search(
		s.client.Search.WithContext(context.Background()),
		s.client.Search.WithIndex(s.index),
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
		highlight := hit.(map[string]interface{})["highlight"]

		url := source.(map[string]interface{})["URL"]
		// content := source.(map[string]interface{})["Content"]
		score := hit.(map[string]interface{})["_score"]

		highlightedContent := highlight.(map[string]interface{})["Content"]

		results = append(results, &seek.Result{
			URL:     url.(string),
			Content: highlightedContent.([]interface{})[0].(string),
			Score:   score.(float64),
		})
	}

	return results
}

func (s *Search) SearchFuzzy(word string) ([]*seek.Result, error) {
	results := make([]*seek.Result, 0)
	results = s.search(SearchFuzzy, word)
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	//sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func (s *Search) SearchSentence(sentence string) ([]*seek.Result, error) {
	results := make([]*seek.Result, 0)
	results = s.search(SearchSentence, sentence)
	if len(results) == 0 {
		return nil, fmt.Errorf("no results found")
	}

	//sort results by score
	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}

func getQuery(st searchType, phrase string) (map[string]interface{}, error) {
	switch st {
	case SearchFuzzy:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"fuzzy": map[string]interface{}{
					"Content": map[string]interface{}{
						"value": phrase,
					},
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					"Content": map[string]interface{}{},
				},
			},
		}, nil
	case SearchSentence:
		return map[string]interface{}{
			"query": map[string]interface{}{
				"query_string": map[string]interface{}{
					"query": phrase,
				},
			},
			"highlight": map[string]interface{}{
				"fields": map[string]interface{}{
					"Content": map[string]interface{}{},
				},
			},
		}, nil
	default:
		return nil, fmt.Errorf("invalid search type")
	}
}
