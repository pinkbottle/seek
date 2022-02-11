package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pinkbottle/seek/elastic"
)

type SearchSentenceCommand struct {
}

func NewSearchSentenceCommand() (cli.Command, error) {
	return &SearchSentenceCommand{}, nil
}

func (s *SearchSentenceCommand) Help() string {
	return ""
}

func (s *SearchSentenceCommand) Run(args []string) int {
	search, err := elastic.NewSearch(index)
	if err != nil {
		fmt.Printf("error creating the client: %s", err)
		return -1
	}

	phrase := args[0:]
	query := strings.Join(phrase, " ")

	results, err := search.SearchSentence(query)
	if err != nil {
		fmt.Printf("error searching: %s", err)
		return -1
	}

	for _, r := range results[0:] {
		content := r.Content
		fmt.Printf("%s (%f)\n\n%s\n\n\n", r.URL, r.Score, content)
	}

	return 0
}

func (s *SearchSentenceCommand) Synopsis() string {
	return ""
}
