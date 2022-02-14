package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"
)

var (
	sources = flag.String("sources", "http://localhost:9200/seek", "what sources to use for searching")
	port    = flag.Int("port", 8080, "port to listen on")
)

func main() {
	flag.Parse()
	c := &http.Client{
		Timeout: time.Second * 5,
	}

	server := newServer(c)
	mux := http.NewServeMux()
	s := http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}
	server.registerRoutes(mux)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("error running server: %v", err)
	}
}

type server struct {
	c *http.Client
}

func newServer(c *http.Client) *server {
	return &server{c: c}
}

func (s *server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprintln(w, "hello world")
	})
}
