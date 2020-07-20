package main

import (
	"io"
	"log"

	"github.com/rdnply/backend-trainee-assignment/internal/app"
)

func main() {
	app, closers, err := app.New(":9000")
	if err != nil {
		log.Fatal(err)
	}
	defer handleClosers(closers)

	app.RunServer()
}

func handleClosers(m map[string]io.Closer) {
	for n, c := range m {
		if err := c.Close(); err != nil {
			log.Printf("can't close %q: %s", n, err)
		}
	}
}
