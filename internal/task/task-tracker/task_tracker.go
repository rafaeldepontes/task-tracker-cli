package tasktracker

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"os"
	"time"

	"github.com/rafaeldepontes/task-tracker-cli/internal/task/model"
	"github.com/spf13/cobra"
)

const (
	InitStatus  model.StatusTask = "init"
	DoingStatus model.StatusTask = "doing"
	DoneStatus  model.StatusTask = "done"
)

type Executable struct {
	incID uint64
	cmd   *cobra.Command
	file  *os.File
}

func (exec *Executable) Execute() error {
	return exec.cmd.Execute()
}

// CreateTask adds the task to the json file...
func (exec *Executable) CreateTask() *cobra.Command {
	return &cobra.Command{
		Use:   "add [description]",
		Short: "Create a new task with the description",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			task := model.Task{
				ID:          exec.incID,
				Description: args[0],
				Status:      InitStatus,
				CreatedAt:   time.Now(),
			}

			// read the whole file
			data, err := readFile(exec)
			if err != nil {
				return
			}

			var tasks []model.Task
			if len(bytes.TrimSpace(data)) > 0 {
				if err := json.Unmarshal(data, &tasks); err != nil {
					log.Println("Could not create a new task,", err.Error())
					return
				}
			} else {
				tasks = make([]model.Task, 0)
			}

			// Append and marshal
			tasks = append(tasks, task)
			out, err := json.MarshalIndent(tasks, "", "  ")
			if err != nil {
				log.Println("could not marshal tasks:", err)
				return
			}

			// Truncate and write back
			if err := exec.file.Truncate(0); err != nil {
				log.Println("truncate failed:", err)
				return
			}

			if _, err := exec.file.Seek(0, io.SeekStart); err != nil {
				log.Println("seek failed:", err)
				return
			}

			if _, err := exec.file.Write(out); err != nil {
				log.Println("[ERROR] write failed:", err)
				return
			}

			exec.incID++
		},
	}
}

func (exec *Executable) UpdateTask() *cobra.Command {
	return &cobra.Command{
		Use:   "update [id] [description]",
		Short: "Update a task description based on the id",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: Implement this
		},
	}
}

func (exec *Executable) DeleteTask() *cobra.Command {
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
func (exec *Executable) MarkInProgressTask() *cobra.Command {
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
func (exec *Executable) ListTasks() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: `List all the task, using some flags to define if you want some filters`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			//TODO: Implement this
		},
	}
}

func NewCommand(f *os.File) *Executable {
	var id uint64 = 1

	if _, err := f.Seek(0, io.SeekStart); err == nil {
		data, err := io.ReadAll(f)
		if err == nil && len(bytes.TrimSpace(data)) > 0 {
			var tasks []model.Task
			if err := json.Unmarshal(data, &tasks); err == nil {
				for _, t := range tasks {
					if t.ID >= id {
						id = t.ID + 1
					}
				}
			}
		}
	}

	exec := &Executable{
		incID: id,
		file:  f,
	}

	cmd := &cobra.Command{}
	cmd.AddCommand(exec.CreateTask())

	exec.cmd = cmd
	return exec
}

func readFile(exec *Executable) ([]byte, error) {
	if _, err := exec.file.Seek(0, io.SeekStart); err != nil {
		log.Println("seek failed:", err)
		return nil, err
	}
	data, err := io.ReadAll(exec.file)
	if err != nil {
		log.Println("read failed:", err)
		return nil, err
	}
	return data, nil
}
