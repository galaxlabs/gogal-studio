package main

import (
	"log"

	"github.com/galaxylabs/gogal-studio/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}
}
