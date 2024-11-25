package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"time"
)

var tasks []task

func main() {
	defer panicRecover()
	err := checkFile()
	if err != nil {
		panic(err)
	}

	err = readJsonFile()
	if err != nil {
		panic(err)
	}

	args := os.Args
	args = args[1:]

	switch args[0] {
	case "add":
		addTask(&args[1])
		err := writeJsonFile()
		if err != nil {
			panic(err)
		}
		return
	case "update":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		changeDescriptionTask(id, args[2])

		err = writeJsonFile()
		if err != nil {
			panic(err)
		}
		return
	case "delete":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		deleteTask(id)

		err = writeJsonFile()
		if err != nil {
			panic(err)
		}
		return
	case "list":
		if len(args) > 1 {
			listTask(status(args[1]))
			return
		} else {
			listTask("")
		}
	case "mark-in-progress":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		changeStatusTask(StatusInProgress, id)

		err = writeJsonFile()
		if err != nil {
			panic(err)
		}
		return
	case "mark-done":
		id, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}

		changeStatusTask(StatusDone, id)

		err = writeJsonFile()
		if err != nil {
			panic(err)
		}
		return
	default:
		fmt.Println("command not found")
		return
	}
}

func panicRecover() {
	if r := recover(); r != nil {
		fmt.Printf("Error recovered: %v\n", r)
	}
}

func getFile() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return path.Join(wd, "tasks.json"), nil
}

func checkFile() error {
	filepath, err := getFile()
	if err != nil {
		return err
	}

	_, err = os.Stat(filepath)

	if os.IsNotExist(err) {
		file, err := os.Create("tasks.json")
		if err != nil {
			return err
		}
		defer file.Close()

		err = os.WriteFile(filepath, []byte("[]"), os.ModeAppend)
		if err != nil {
			return err
		}
	}

	return nil
}

type status string

const (
	StatusTodo       status = "todo"
	StatusInProgress status = "in-progress"
	StatusDone       status = "done"
)

type task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Status      status    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func readJsonFile() error {
	filepath, err := getFile()
	if err != nil {
		return err
	}

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&tasks)
	if err != nil {
		return err
	}

	return nil
}

func writeJsonFile() error {
	filepath, err := getFile()
	if err != nil {
		return err
	}

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(tasks)
	if err != nil {
		return err
	}

	return nil
}

func addTask(description *string) {
	id := 1

	if len(tasks) > 0 {
		id = tasks[len(tasks)-1].ID + 1
	}

	createdAt := time.Now()

	task := task{
		ID:          id,
		Description: *description,
		Status:      StatusTodo,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	tasks = append(tasks, task)
}

func listTask(status status) {
	filteredTasks := []task{}

	if status == "" {
		filteredTasks = tasks
	} else {
		switch status {
		case StatusTodo:
			for _, task := range tasks {
				if task.Status == StatusTodo {
					filteredTasks = append(filteredTasks, task)
				}
			}
		case StatusInProgress:
			for _, task := range tasks {
				if task.Status == StatusInProgress {
					filteredTasks = append(filteredTasks, task)
				}
			}
		case StatusDone:
			for _, task := range tasks {
				if task.Status == StatusDone {
					filteredTasks = append(filteredTasks, task)
				}
			}
		default:
			fmt.Println("command not found")
			return
		}
	}

	for _, task := range filteredTasks {
		fmt.Printf(
			"Task: %d\nStatus: %s\nDescription: %s\nCreated at: %s\nUpdate at: %s\n\n",
			task.ID, task.Status, task.Description, task.CreatedAt.Format(time.RFC822), task.UpdatedAt.Format(time.RFC822),
		)
	}
}

func changeStatusTask(status status, id int) {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == id {
			tasks[i].Status = status
			tasks[i].UpdatedAt = time.Now()
			return
		}
	}

	fmt.Println("task not found")
}

func changeDescriptionTask(id int, description string) {
	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID == id {
			tasks[i].Description = description
			tasks[i].UpdatedAt = time.Now()
			return
		}
	}

	fmt.Println("task not found")
}

func deleteTask(id int) {
	undeletedTask := []task{}
	for i := 0; i < len(tasks); i++ {
		if tasks[i].ID != id {
			undeletedTask = append(undeletedTask, tasks[i])
		}
	}

	tasks = undeletedTask
}
