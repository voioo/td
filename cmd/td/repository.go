package main

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

var repositoryFilePath = func() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		report(err)
	}
	return filepath.Join(homeDir, ".td.json")
}()

func loadTasksFromRepositoryFile() (todos []*Task, doneTodos []*Task, latestTaskID int) {
	f, err := os.Open(repositoryFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return todos, doneTodos, latestTaskID
		}
		report(err)
	}
	defer f.Close()

	var t []*Task
	if err = json.NewDecoder(f).Decode(&t); err != nil {
		return todos, doneTodos, latestTaskID
	}

	for _, v := range t {
		if v.ID >= latestTaskID {
			latestTaskID = v.ID
		}

		if v.IsDone {
			doneTodos = append(doneTodos, v)
		} else {
			todos = append(todos, v)
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
	data, err := json.Marshal(todos)
	if err != nil {
		report(err)
	}

	if err := f.Truncate(0); err != nil {
		report(err)
	}

	_, err = f.Write(data)
	if err != nil {
		report(err)
	}
}
