package main

import (
	"log"
	"os"
)

var build = "develop"

func run() error {
	return nil
}

func main() {
	if err := run(); err != nil {
		log.Println("error: ", err)
		os.Exit(1)
	}
}
