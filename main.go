package main

import (
	"log"
	"os"

	taskTracker "github.com/rafaeldepontes/task-tracker-cli/internal/task/task-tracker"
)

const (
	// Path
	Path = "storage/tasks.json"

	// File Mode
	OwnerPropertyMode = 0644
)

func main() {
	_, err := os.Stat(Path)
	if err != nil {
		if err := os.WriteFile(Path, []byte("[]"), OwnerPropertyMode); err != nil {
			log.Fatal("Couldn't create the needed file,", err.Error())
		}
	}

	f, err := os.OpenFile(Path, os.O_RDWR|os.O_CREATE, OwnerPropertyMode)
	if err != nil {
		log.Fatal("Couldn't open the storage file,", err.Error())
	}
	defer f.Close()

	rootCmd := taskTracker.NewCommand(f)
	if err := rootCmd.Execute(); err != nil {
		log.Println(err.Error())
	}
}
