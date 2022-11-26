package main

import (
	"log"

	"github.com/KJone1/gophercises-task-manager/src/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal("Failed to start")
	}
}
