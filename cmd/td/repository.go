package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	todo "github.com/voioo/td"
)

var repositoryFilePath = func() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		report(err)
	}
	return filepath.Join(homeDir, ".td.json")
}()

func loadTasksFromRepositoryFile() (todos []*todo.Task, doneTodos []*todo.Task, latestTaskID int) {
	f, err := os.Open(repositoryFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return todos, doneTodos, latestTaskID
		}
		report(err)
	}
	defer f.Close()

	var t []*todo.Task
	if err = json.NewDecoder(f).Decode(&t); err != nil {
		report(err)
	}

	for _, v := range t {
		if v.IsDone {
			doneTodos = append(doneTodos, v)
			continue
		}
		todos = append(todos, v)

		if v.ID >= latestTaskID {
			latestTaskID = v.ID
		}
	}

	return todos, doneTodos, latestTaskID
}

func (m model) saveTasks() {
	f, err := os.OpenFile(repositoryFilePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		report(err)
	}
	defer f.Close()

	todos := append(m.tasks, m.doneTasks...)
	data, _ := json.Marshal(todos)

	if err := f.Truncate(0); err != nil {
		report(err)
	}

	_, err = f.Write(data)
	if err != nil {
		report(err)
	}
}
