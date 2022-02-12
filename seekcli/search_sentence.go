package seekcli

import (
	"fmt"
	"strings"

	"github.com/mitchellh/cli"
	"github.com/pinkbottle/seek/elastic"
)

type SearchSentenceCommand struct {
	index string
}

func NewSearchSentenceCommand() (cli.Command, error) {
	return &SearchSentenceCommand{
		index: "seek",
	}, nil
}

func (s *SearchSentenceCommand) Help() string {
	return ""
}

func (s *SearchSentenceCommand) Run(args []string) int {
	search, err := elastic.NewSearch(s.index)
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

	printWithHighlight(results)
	return 0
}

func (s *SearchSentenceCommand) Synopsis() string {
	return ""
}
