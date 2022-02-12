package elastic

import (
	"fmt"

	"github.com/pinkbottle/seek"
)

func getQuery(st seek.Type, phrase string) (map[string]interface{}, error) {
	switch st {
	case seek.SearchFuzzy:
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
	case seek.SearchSentence:
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
