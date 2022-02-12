package seekcli

import (
	"context"
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/pinkbottle/seek"
	"github.com/pinkbottle/seek/elastic"
)

type SearchWordCommand struct {
	index string
}

func NewSearchWordCommand() (cli.Command, error) {
	return &SearchWordCommand{
		index: "seek",
	}, nil
}

func (s *SearchWordCommand) Help() string {
	return ""
}

func (s *SearchWordCommand) Run(args []string) int {
	search, err := elastic.NewSeeker(s.index)
	if err != nil {
		fmt.Printf("error creating the client: %s", err)
		return -1
	}

	phrase := args[0]
	results, err := search.Seek(context.Background(), phrase, seek.SearchFuzzy)
	if err != nil {
		fmt.Printf("error searching: %s", err)
		return -1
	}

	printWithHighlight(results)

	return 0
}
func (s *SearchWordCommand) Synopsis() string {
	return ""
}
