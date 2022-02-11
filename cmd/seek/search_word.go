package main

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pinkbottle/seek/elastic"
)

type SearchWordCommand struct {
}

func NewSearchWordCommand() (cli.Command, error) {
	return &SearchWordCommand{}, nil
}

func (s *SearchWordCommand) Help() string {
	return ""
}

func (s *SearchWordCommand) Run(args []string) int {
	search, err := elastic.NewSearch(index)
	if err != nil {
		fmt.Printf("error creating the client: %s", err)
		return -1
	}

	phrase := args[0]
	results, err := search.SearchFuzzy(phrase)
	if err != nil {
		fmt.Printf("error searching: %s", err)
		return -1
	}

	for _, r := range results[0:] {
		fmt.Printf("%s (%f)\n\n%s\n%s\n\n", r.URL, r.Score, r.Content, strings.Repeat(".", 37))
	}

	return 0
}
func (s *SearchWordCommand) Synopsis() string {
	return ""
}
