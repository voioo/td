package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRepository(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "td-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.json")
	originalPath := repositoryFilePath
	repositoryFilePath = testFile
	defer func() {
		repositoryFilePath = originalPath
	}()

	t.Run("load empty repository", func(t *testing.T) {
		todos, doneTodos, latestID := loadTasksFromRepositoryFile()

		if len(todos) != 0 {
			t.Errorf("expected empty todos, got %d items", len(todos))
		}
		if len(doneTodos) != 0 {
			t.Errorf("expected empty doneTodos, got %d items", len(doneTodos))
		}
		if latestID != 0 {
			t.Errorf("expected latestID to be 0, got %d", latestID)
		}
	})

	t.Run("save and load tasks", func(t *testing.T) {
		m := model{
			tasks: []*Task{
				{
					CreatedAt: time.Now(),
					Name:      "Test task 1",
					ID:        1,
					IsDone:    false,
				},
			},
			doneTasks: []*Task{
				{
					CreatedAt: time.Now(),
					Name:      "Test task 2",
					ID:        2,
					IsDone:    true,
				},
			},
		}

		m.saveTasks()

		todos, doneTodos, latestID := loadTasksFromRepositoryFile()

		if len(todos) != 1 {
			t.Errorf("expected 1 todo, got %d", len(todos))
		}
		if todos[0].Name != "Test task 1" {
			t.Errorf("expected task name 'Test task 1', got %s", todos[0].Name)
		}

		if len(doneTodos) != 1 {
			t.Errorf("expected 1 done todo, got %d", len(doneTodos))
		}
		if doneTodos[0].Name != "Test task 2" {
			t.Errorf("expected task name 'Test task 2', got %s", doneTodos[0].Name)
		}

		if latestID != 2 {
			t.Errorf("expected latestID to be 2, got %d", latestID)
		}
	})

	t.Run("handle corrupted file", func(t *testing.T) {
		err := os.WriteFile(testFile, []byte("invalid json"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		defer func() {
			if r := recover(); r != nil {
				t.Errorf("loadTasksFromRepositoryFile panicked: %v", r)
			}
		}()

		loadTasksFromRepositoryFile()
	})
}
