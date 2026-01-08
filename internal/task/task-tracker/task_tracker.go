package tasktracker

import (
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/rafaeldepontes/task-tracker-cli/internal/task/model"
	"github.com/spf13/cobra"
)

const (
	// Path
	Path = "storage/tasks.json"

	// File Mode
	OwnerPropertyMode = 0644

	// Enum
	InitStatus  model.StatusTask = "init"
	DoingStatus model.StatusTask = "doing"
	DoneStatus  model.StatusTask = "done"
)

type RootCmd struct {
	cmd *cobra.Command
}

func (exec *RootCmd) Execute() error {
	return exec.cmd.Execute()
}

// CreateTask adds the task to the json file...
func (exec *RootCmd) CreateTask() *cobra.Command {
	return &cobra.Command{
		Use:   "add [description]",
		Short: "Create a new task with the description",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := getTasks()
			if err != nil {
				return
			}

			var id uint64 = 1
			if len(tasks) > 0 {
				id = tasks[len(tasks)-1].ID + 1
			}

			task := model.Task{
				ID:          id,
				Description: args[0],
				Status:      InitStatus,
				CreatedAt:   time.Now(),
			}

			log.Println("Tasks:", tasks)

			tasks = append(tasks, task)

			log.Println("Tasks:", tasks)
			err = writeFile(tasks)
			if err != nil {
				return
			}
		},
	}
}

func (exec *RootCmd) UpdateTask() *cobra.Command {
	return &cobra.Command{
		Use:   "update [id] [description]",
		Short: "Update a task description based on the id",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: Implement this
		},
	}
}

func (exec *RootCmd) DeleteTask() *cobra.Command {
	return &cobra.Command{
		Use:   "detele [id]",
		Short: "Delete a task based on the id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: Implement this
		},
	}
}

// TODO: Check if it's possible to use a switch case instead of multiple functions doing basically the same thing
func (exec *RootCmd) MarkInProgressTask() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-in-progress [id]",
		Short: `Set a task as "in progress" based on the id`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: Implement this
		},
	}
}

// TODO: Check if it's possible to use a switch case instead of multiple functions doing basically the same thing
func (exec *RootCmd) ListTasks() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: `List all the task, using some flags to define if you want some filters`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: Implement this
		},
	}
}

func NewCommand() *RootCmd {
	_, err := os.Stat(Path)
	if err != nil {
		if err := os.WriteFile(Path, []byte("[]"), OwnerPropertyMode); err != nil {
			log.Fatal("Couldn't create the needed file,", err.Error())
		}
	}

	exec := &RootCmd{}

	cmd := &cobra.Command{}
	cmd.AddCommand(exec.CreateTask())

	exec.cmd = cmd
	return exec
}

func getTasks() ([]model.Task, error) {
	f, err := os.OpenFile(Path, os.O_RDWR|os.O_CREATE, OwnerPropertyMode)
	if err != nil {
		log.Fatal("Couldn't open the storage file,", err.Error())
	}
	defer f.Close()

	buffer := make([]byte, 4096)
	length := 0
	for {
		n, err := f.Read(buffer[length:])
		if err != nil {
			break
		}
		length += n
	}

	var tasks []model.Task
	if err := json.Unmarshal(buffer[:length], &tasks); err != nil {
		log.Println("[ERROR] unmarshiling went wrong:", err.Error())
		return nil, err
	}

	return tasks, nil
}

func writeFile(tasks []model.Task) error {
	data, err := json.Marshal(tasks)
	if err != nil {
		log.Println("didn't marshal data:", err.Error())
		return err
	}

	err = os.WriteFile(Path, data, OwnerPropertyMode)
	if err != nil {
		log.Println("Somehow I coded poorly, couldn't write file:", err.Error())
		return err
	}
	return nil
}
