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
	results, err := search.SearchSentence(strings.Join(phrase, " "))
	if err != nil {
		fmt.Printf("error searching: %s", err)
		return -1
	}

	for _, r := range results[0:] {
		fmt.Printf("%s (%f)\n\n%s\n%s\n\n", r.URL, r.Score, r.Content, strings.Repeat(".", 37))
	}

	return 0
}

func (s *SearchSentenceCommand) Synopsis() string {
	return ""
}
