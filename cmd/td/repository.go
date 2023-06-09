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
	f, err := os.OpenFile(repositoryFilePath, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			report(err)
		}
		f, err = os.Create(repositoryFilePath)
		if err != nil {
			report(err)
		}
	}
	defer f.Close()
	if err := f.Truncate(0); err != nil {
		report(err)
	}

	todos := append(m.tasks, m.doneTasks...)
	data, _ := json.Marshal(todos)

	_, err = f.Write(data)
	if err != nil {
		report(err)
	}
}
