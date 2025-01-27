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

	t.Run("task priority", func(t *testing.T) {
		task := &Task{
			Name:     "Test task",
			ID:       1,
			Priority: PriorityHigh,
		}

		if task.Priority != PriorityHigh {
			t.Errorf("expected priority to be high, got %s", task.Priority)
		}

		if task.Priority.String() != "high" {
			t.Errorf("expected priority string to be 'high', got %s", task.Priority.String())
		}
	})
}

func TestTaskSorting(t *testing.T) {
	now := time.Now()
	tasks := []*Task{
		{ID: 1, Name: "Task 1", CreatedAt: now, Priority: PriorityNone},
		{ID: 2, Name: "Task 2", CreatedAt: now.Add(time.Hour), Priority: PriorityHigh},
		{ID: 3, Name: "Task 3", CreatedAt: now.Add(2 * time.Hour), Priority: PriorityMedium},
	}

	original := make([]*Task, len(tasks))
	copy(original, tasks)

	SortTasksByPriority(tasks)

	if tasks[0].ID != 2 {
		t.Error("Expected high priority task first")
	}
	if tasks[1].ID != 3 {
		t.Error("Expected medium priority task second")
	}

	for _, task := range tasks {
		found := false
		for _, orig := range original {
			if orig.ID == task.ID {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Task %d was lost during sorting", task.ID)
		}
	}
}
