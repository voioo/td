package main

import (
	"testing"
	"time"
)

func TestTask(t *testing.T) {
	t.Run("create new task", func(t *testing.T) {
		task := &Task{
			CreatedAt: time.Now(),
			Name:      "Test task",
			ID:        1,
			IsDone:    false,
		}

		if task.Name != "Test task" {
			t.Errorf("expected task name to be 'Test task', got %s", task.Name)
		}

		if task.ID != 1 {
			t.Errorf("expected task ID to be 1, got %d", task.ID)
		}

		if task.IsDone {
			t.Error("expected new task to not be done")
		}
	})
}
