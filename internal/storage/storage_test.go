package storage

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/voioo/td/internal/task"
)

func TestRepository(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "td-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dataFile := filepath.Join(tempDir, "test.json")
	repo := NewRepository(dataFile)

	t.Run("load empty repository", func(t *testing.T) {
		tasks, doneTasks, nextID, err := repo.LoadTasks()

		if err != nil {
			t.Errorf("expected no error for empty repository, got %v", err)
		}
		if len(tasks) != 0 {
			t.Errorf("expected 0 tasks, got %d", len(tasks))
		}
		if len(doneTasks) != 0 {
			t.Errorf("expected 0 done tasks, got %d", len(doneTasks))
		}
		if nextID != 0 {
			t.Errorf("expected next ID to be 0, got %d", nextID)
		}
	})

	t.Run("save and load tasks", func(t *testing.T) {
		// Create test tasks
		tasks := []*task.Task{
			{
				ID:        1,
				Name:      "Active Task",
				Priority:  task.PriorityHigh,
				IsDone:    false,
				CreatedAt: testTime(),
			},
		}
		doneTasks := []*task.Task{
			{
				ID:        2,
				Name:      "Done Task",
				Priority:  task.PriorityLow,
				IsDone:    true,
				CreatedAt: testTime(),
			},
		}

		// Save tasks
		err := repo.SaveTasks(tasks, doneTasks)
		if err != nil {
			t.Errorf("expected no error saving tasks, got %v", err)
		}

		// Load tasks
		loadedTasks, loadedDoneTasks, nextID, err := repo.LoadTasks()
		if err != nil {
			t.Errorf("expected no error loading tasks, got %v", err)
		}

		if len(loadedTasks) != 1 {
			t.Errorf("expected 1 active task, got %d", len(loadedTasks))
		}
		if len(loadedDoneTasks) != 1 {
			t.Errorf("expected 1 done task, got %d", len(loadedDoneTasks))
		}
		if nextID != 3 {
			t.Errorf("expected next ID to be 3, got %d", nextID)
		}

		if loadedTasks[0].Name != "Active Task" {
			t.Errorf("expected task name to be 'Active Task', got %s", loadedTasks[0].Name)
		}
		if loadedDoneTasks[0].Name != "Done Task" {
			t.Errorf("expected done task name to be 'Done Task', got %s", loadedDoneTasks[0].Name)
		}
	})

	t.Run("validate valid task", func(t *testing.T) {
		validTask := &task.Task{
			ID:        1,
			Name:      "Valid Task",
			Priority:  task.PriorityMedium,
			IsDone:    false,
			CreatedAt: testTime(),
		}

		repo := NewRepository(dataFile)
		err := repo.SaveTasks([]*task.Task{validTask}, []*task.Task{})

		if err != nil {
			t.Errorf("expected no error for valid task, got %v", err)
		}
	})

	t.Run("validate invalid task - empty name", func(t *testing.T) {
		invalidTask := &task.Task{
			ID:        1,
			Name:      "", // Empty name
			Priority:  task.PriorityMedium,
			IsDone:    false,
			CreatedAt: testTime(),
		}

		repo := NewRepository(dataFile)
		err := repo.SaveTasks([]*task.Task{invalidTask}, []*task.Task{})

		if err == nil {
			t.Error("expected error for task with empty name")
		}
	})

	t.Run("validate invalid task - zero ID", func(t *testing.T) {
		invalidTask := &task.Task{
			ID:        0, // Invalid ID
			Name:      "Task",
			Priority:  task.PriorityMedium,
			IsDone:    false,
			CreatedAt: testTime(),
		}

		repo := NewRepository(dataFile)
		err := repo.SaveTasks([]*task.Task{invalidTask}, []*task.Task{})

		if err == nil {
			t.Error("expected error for task with zero ID")
		}
	})

	t.Run("validate invalid task - invalid priority", func(t *testing.T) {
		invalidTask := &task.Task{
			ID:        1,
			Name:      "Task",
			Priority:  task.Priority(999), // Invalid priority
			IsDone:    false,
			CreatedAt: testTime(),
		}

		repo := NewRepository(dataFile)
		err := repo.SaveTasks([]*task.Task{invalidTask}, []*task.Task{})

		if err == nil {
			t.Error("expected error for task with invalid priority")
		}
	})

	t.Run("load corrupted data", func(t *testing.T) {
		// Write invalid JSON to the file
		err := os.WriteFile(dataFile, []byte("invalid json"), 0644)
		if err != nil {
			t.Fatal(err)
		}

		repo := NewRepository(dataFile)
		_, _, _, err = repo.LoadTasks()

		if err == nil {
			t.Error("expected error when loading corrupted data")
		}
	})
}

// testTime returns a consistent time for testing
func testTime() time.Time {
	return time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
}
