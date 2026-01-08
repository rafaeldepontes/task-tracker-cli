package main

import (
	"log"

	taskTracker "github.com/rafaeldepontes/task-tracker-cli/internal/task/task-tracker"
)

func main() {
	rootCmd := taskTracker.NewCommand()
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
	}
}
