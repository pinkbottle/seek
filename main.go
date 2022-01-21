package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

func main() {
	for {
		req, err := http.NewRequestWithContext(context.Background(), "GET", "https://en.wikipedia.org/wiki/Special:Random", nil)
		if err != nil {
			log.Fatalf("error creating request: %v", err)
		}

		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatalf("error sending request: %v", err)
		}

		b, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalf("error reading response: %v", err)
		}

		b, _ = json.Marshal(map[string]string{"message": string(b), "source": "https://en.wikipedia.org/wiki/Special:Random"})
		res, err = http.Post("http://localhost:8081/publish", "application/json", bytes.NewReader(b))
		if err != nil {
			log.Fatalf("Error: %v", err)
		}
		log.Printf("%d", res.StatusCode)
		res.Body.Close()
		time.Sleep(time.Second * 10)
	}
}

//write a regex expression to match text inside <p> tags
var exp = regexp.MustCompile(`<p>(.+?)</p>`)

func findSentence(source string) (string, error) {
	match := exp.FindString(source)
	if strings.EqualFold(match, "") {
		return "", fmt.Errorf("no match found")
	}
	return match, nil
}
