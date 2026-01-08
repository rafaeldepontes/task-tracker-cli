package tasktracker

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
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
	InitStatus  model.StatusTask = "TODO"
	DoingStatus model.StatusTask = "In progress"
	DoneStatus  model.StatusTask = "Done"

	// Text Format
	TaskInfo              = "Task %d - %s\nProgress: %s\nCreated: %s\nUpdated last: %s\n\n"
	BrazilianFormatLayout = "02/01/2006 15:04:05"
)

type RootCmd struct {
	cmd *cobra.Command
}

func (exec *RootCmd) Execute() error {
	return exec.cmd.Execute()
}

func (exec *RootCmd) CreateTask() *cobra.Command {
	return &cobra.Command{
		Use:   "add [description]",
		Short: "Create a new task with the description",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := readFile()
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
			tasks = append(tasks, task)

			err = writeFile(tasks)
			if err != nil {
				return
			}

			log.Printf("Task added successfully (ID: %d)\n", id)
		},
	}
}

func (exec *RootCmd) UpdateTask() *cobra.Command {
	return &cobra.Command{
		Use:   "update [id] [description]",
		Short: "Update a task description based on the id",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := readFile()
			if err != nil {
				return
			}

			pos, err := searchTaks(args[0], tasks)
			if err != nil {
				return
			}

			tasks[pos].Description = args[1]
			err = writeFile(tasks)
			if err != nil {
				return
			}

			log.Printf("Task updated successfully (ID: %d)\n", tasks[pos].ID)
		},
	}
}

func (exec *RootCmd) DeleteTask() *cobra.Command {
	return &cobra.Command{
		Use:   "delete [id]",
		Short: "Delete a task based on the id",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := readFile()
			if err != nil {
				log.Println("couldn't read from the storage file:", err.Error())
				return
			}

			pos, err := searchTaks(args[0], tasks)
			if err != nil {
				return
			}

			tasks = append(tasks[:pos], tasks[pos+1:]...)
			err = writeFile(tasks)
			if err != nil {
				return
			}

			log.Println("Task removed successfully")
		},
	}
}

func (exec *RootCmd) MarkInProgressTask() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-in-progress [id]",
		Short: `Set a task as "in progress" based on the id`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := readFile()
			if err != nil {
				return
			}

			pos, err := searchTaks(args[0], tasks)
			if err != nil {
				return
			}
			now := time.Now()

			tasks[pos].Status = DoingStatus
			tasks[pos].UpdatedAt = &now
			err = writeFile(tasks)
			if err != nil {
				return
			}

			log.Println("Task marked as in progress!")
		},
	}
}

func (exec *RootCmd) MarkDoneTask() *cobra.Command {
	return &cobra.Command{
		Use:   "mark-done [id]",
		Short: `Set a task as "done" based on the id`,
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			tasks, err := readFile()
			if err != nil {
				return
			}

			pos, err := searchTaks(args[0], tasks)
			if err != nil {
				return
			}
			now := time.Now()

			tasks[pos].Status = DoneStatus
			tasks[pos].UpdatedAt = &now
			err = writeFile(tasks)
			if err != nil {
				return
			}

			log.Println("Task marked as done!")
		},
	}
}

// TODO: Check if it's possible to use a switch case instead of multiple functions doing basically the same thing
func (exec *RootCmd) ListTasks() *cobra.Command {
	return &cobra.Command{
		Use:   "list [filter]",
		Short: `List all the task, using some flags to define if you want some filters`,
		Long: `List all the tasks, if a filter is specified than the filter is considered.

		Else it will show all the tasks no matter the status.`,
		Run: func(cmd *cobra.Command, args []string) {
			action := "all"
			if len(args) > 0 {
				action = args[0]
			}

			tasks, err := readFile()
			if err != nil {
				return
			}

			println()
			switch action {
			case "all":
				var updatedAt string
				var createdAt string
				for _, t := range tasks {
					printTask(&t, &updatedAt, &createdAt)
				}

			case "done":
				var updatedAt string
				var createdAt string
				for _, t := range tasks {
					if t.Status == DoneStatus {
						printTask(&t, &updatedAt, &createdAt)
					}
				}

			case "todo":
				var updatedAt string
				var createdAt string
				for _, t := range tasks {
					if t.Status == InitStatus {
						printTask(&t, &updatedAt, &createdAt)
					}
				}

			case "in-progress":
				var updatedAt string
				var createdAt string
				for _, t := range tasks {
					if t.Status == DoingStatus {
						printTask(&t, &updatedAt, &createdAt)
					}
				}
			}
		},
	}
}

func printTask(t *model.Task, updatedAt *string, createdAt *string) {
	if t.UpdatedAt == nil {
		*updatedAt = "-"
	} else {
		*updatedAt = t.UpdatedAt.Format(BrazilianFormatLayout)
	}
	*createdAt = t.CreatedAt.Format(BrazilianFormatLayout)

	fmt.Printf(TaskInfo, t.ID, t.Description, t.Status, *createdAt, *updatedAt)
}

func NewCommand() *RootCmd {
	_, err := os.Stat(Path)
	if err != nil {
		if err := os.WriteFile(Path, []byte("[]"), OwnerPropertyMode); err != nil {
			log.Fatal("Couldn't create the needed file,", err.Error())
		}
	}

	rootCmd := &RootCmd{}

	cmd := &cobra.Command{}
	cmd.AddCommand(rootCmd.ListTasks())
	cmd.AddCommand(rootCmd.CreateTask())
	cmd.AddCommand(rootCmd.UpdateTask())
	cmd.AddCommand(rootCmd.DeleteTask())
	cmd.AddCommand(rootCmd.MarkDoneTask())
	cmd.AddCommand(rootCmd.MarkInProgressTask())

	rootCmd.cmd = cmd
	return rootCmd
}

func readFile() ([]model.Task, error) {
	f, err := os.OpenFile(Path, os.O_RDWR|os.O_CREATE, OwnerPropertyMode)
	if err != nil {
		log.Fatal("Couldn't open the storage file,", err.Error())
		return nil, err
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

// searchTaks returns the task position in the underlying array and an error if any.
func searchTaks(strID string, tasks []model.Task) (int, error) {
	id_, err := strconv.ParseInt(strID, 10, 64)
	if err != nil {
		log.Println(`[ERROR] didn't parsed the id:`, err.Error())
		return 0, err
	}
	id := uint64(id_)

	pos := -1
	for i, t := range tasks {
		if t.ID == id {
			pos = i
			break
		}
	}

	if pos == -1 {
		log.Println(`Task not found with id:`, strID)
		return 0, errors.New("Task not found")
	}
	return pos, nil
}
