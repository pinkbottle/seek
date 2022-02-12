package seekcli

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/pinkbottle/seek"
)

func printWithHighlight(results []*seek.Result) {
	for _, r := range results[0:] {
		content := r.Content
		tokens := strings.Split(content, " ")
		fmt.Printf("%s (%f)\n", r.URL, r.Score)
		for i, t := range tokens[0:] {
			if strings.Contains(t, "<em>") {
				t = strings.ReplaceAll(t, "<em>", "")
				t = strings.ReplaceAll(t, "</em>", "")
				color.Set(color.FgGreen)
			} else {
				color.Unset()
			}

			if i == len(tokens)-1 {
				color.Unset()
			}
			fmt.Printf("%s ", t)
		}
		fmt.Println("")
		fmt.Println("")
	}

}
