package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/pinkbottle/seek/wiki"
)

func main() {
	for {
		sink, err := wiki.NewSink(http.DefaultClient)
		if err != nil {
			log.Fatalf("error creating sink: %v", err)
		}

		ctx := context.Background()
		article, err := sink.GetRandomArticle(ctx)
		if err != nil {
			log.Fatalf("error getting random article: %v", err)
		}
		log.Printf("got article: %v", article)
		// b, _ := json.Marshal(article)
		// res, err := http.Post("http://localhost:8081/publish", "application/json", bytes.NewReader(b))
		// if err != nil {
		// 	log.Fatalf("Error: %v", err)
		// }
		// log.Printf("%d", res.StatusCode)
		// res.Body.Close()
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
