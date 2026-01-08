# Task-Tracker-Cli

An app to track your tasks and manage your to-do list in Golang.

I'm trying to build an app to-do list `CLI` as a foundation to my open-source project, an automation builder for Golang projects!

In this project I will be using `Cobra CLI` and `Golang 1.25.5`

## How to Run:

Do this:

```bash
git clone <...>
go mod tidy
go build -o ./bin
```

And add the bin directory to your env variables and feel free to use...

## Commands

The list of commands and their usage is given below:

```bash
# Adding a new task
task-cli add "Buy groceries"
# Output: Task added successfully (ID: 1)

# Updating and deleting tasks
task-cli update 1 "Buy groceries and cook dinner"
task-cli delete 1

# Marking a task as in progress or done
task-cli mark-in-progress 1
task-cli mark-done 1

# Listing all tasks
task-cli list

# Listing tasks by status
task-cli list done
task-cli list todo
task-cli list in-progress
```

## Contact

For questions or help, contact: `rafael.cr.carneiro@gmail.com`
