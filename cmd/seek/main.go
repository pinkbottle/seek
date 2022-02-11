package main

import (
	"log"
	"os"

	"github.com/mitchellh/cli"
)

var (
	searchAddr = "http://localhost:9200/seek/_search"
	index      = "seek"
)

func main() {
	c := cli.NewCLI("seek", "0.1.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"word":     NewSearchWordCommand,
		"sentence": NewSearchSentenceCommand,
	}

	exitStatus, err := c.Run()

	if err != nil {
		log.Fatal(err)
	}

	os.Exit(exitStatus)
}
