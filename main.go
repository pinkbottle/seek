package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
)

func main() {
	b, _ := json.Marshal(map[string]string{"hello": "world"})
	res, err := http.Post("http://localhost:8081/publish", "application/json", bytes.NewReader(b))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("%d", res.StatusCode)
}
