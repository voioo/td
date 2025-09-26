package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/voioo/td/internal/config"
	"github.com/voioo/td/internal/storage"
	"github.com/voioo/td/internal/task"
)

func TestFullApplicationFlow(t *testing.T) {
	// Create a temporary directory for the test
	tempDir, err := os.MkdirTemp("", "td-integration-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	dataFile := filepath.Join(tempDir, "tasks.json")

	t.Run("create, save, and load tasks", func(t *testing.T) {
		// Create config
		cfg := config.DefaultConfig()
		cfg.DataFile = dataFile

		// Create repository and task manager
		repo := storage.NewRepository(cfg.DataFile)
		tm := task.NewTaskManager([]*task.Task{}, []*task.Task{}, 0)

		// Add some tasks
		task1 := tm.AddTask("Write integration tests")
		task2 := tm.AddTask("Refactor codebase")
		task3 := tm.AddTask("Add documentation")

		// Set priorities
		tm.SetTaskPriority(task2.ID, task.PriorityHigh)
		tm.SetTaskPriority(task3.ID, task.PriorityMedium)

		// Complete one task
		tm.CompleteTask(task1.ID)

		// Save tasks
		err := repo.SaveTasks(tm.GetTasks(), tm.GetDoneTasks())
		if err != nil {
			t.Errorf("failed to save tasks: %v", err)
		}

		// Load tasks in a new task manager
		loadedTasks, loadedDoneTasks, nextID, err := repo.LoadTasks()
		if err != nil {
			t.Errorf("failed to load tasks: %v", err)
		}

		newTM := task.NewTaskManager(loadedTasks, loadedDoneTasks, nextID)

		// Verify tasks were loaded correctly
		if len(newTM.GetTasks()) != 2 {
			t.Errorf("expected 2 active tasks, got %d", len(newTM.GetTasks()))
		}
		if len(newTM.GetDoneTasks()) != 1 {
			t.Errorf("expected 1 done task, got %d", len(newTM.GetDoneTasks()))
		}

		// Check task details
		activeTasks := newTM.GetTasks()
		doneTasks := newTM.GetDoneTasks()

		// Find the high priority task
		var highPriorityTask *task.Task
		var mediumPriorityTask *task.Task
		for _, tsk := range activeTasks {
			if tsk.Priority == task.PriorityHigh {
				highPriorityTask = tsk
			}
			if tsk.Priority == task.PriorityMedium {
				mediumPriorityTask = tsk
			}
		}

		if highPriorityTask == nil {
			t.Error("expected to find high priority task")
		} else if highPriorityTask.Name != "Refactor codebase" {
			t.Errorf("expected high priority task name to be 'Refactor codebase', got %s", highPriorityTask.Name)
		}

		if mediumPriorityTask == nil {
			t.Error("expected to find medium priority task")
		} else if mediumPriorityTask.Name != "Add documentation" {
			t.Errorf("expected medium priority task name to be 'Add documentation', got %s", mediumPriorityTask.Name)
		}

		if doneTasks[0].Name != "Write integration tests" {
			t.Errorf("expected done task name to be 'Write integration tests', got %s", doneTasks[0].Name)
		}
	})

	t.Run("undo and redo operations", func(t *testing.T) {
		tm := task.NewTaskManager([]*task.Task{}, []*task.Task{}, 0)
		undoManager := task.NewUndoManager(10)

		// Add a task
		addedTask := tm.AddTask("Test undo/redo")
		undoManager.PushUndo(task.Action{
			Type: task.ActionTypeAdd,
			Task: addedTask,
		})

		// Change priority
		oldPriority := addedTask.Priority
		tm.SetTaskPriority(addedTask.ID, task.PriorityHigh)
		undoManager.PushUndo(task.Action{
			Type:     task.ActionTypePriority,
			Task:     addedTask,
			OldState: oldPriority,
			NewState: task.PriorityHigh,
		})

		// Complete the task
		completedTask := tm.CompleteTask(addedTask.ID)
		undoManager.PushUndo(task.Action{
			Type:     task.ActionTypeComplete,
			Task:     completedTask,
			OldState: false,
			NewState: true,
		})

		// Verify initial state
		if len(tm.GetTasks()) != 0 {
			t.Error("expected no active tasks after completion")
		}
		if len(tm.GetDoneTasks()) != 1 {
			t.Error("expected 1 done task")
		}
		if tm.GetDoneTasks()[0].Priority != task.PriorityHigh {
			t.Error("expected done task to have high priority")
		}

		// Undo completion
		if !undoManager.Undo(tm) {
			t.Error("expected undo to succeed")
		}

		if len(tm.GetTasks()) != 1 {
			t.Error("expected task to be restored to active")
		}
		if len(tm.GetDoneTasks()) != 0 {
			t.Error("expected no done tasks after undo")
		}

		// Undo priority change
		if !undoManager.Undo(tm) {
			t.Error("expected undo priority change to succeed")
		}

		if tm.GetTasks()[0].Priority != task.PriorityNone {
			t.Error("expected task priority to be restored to none")
		}

		// Undo add
		if !undoManager.Undo(tm) {
			t.Error("expected undo add to succeed")
		}

		if len(tm.GetTasks()) != 0 {
			t.Error("expected no tasks after undoing add")
		}

		// Redo operations
		if !undoManager.Redo(tm) {
			t.Error("expected redo to succeed")
		}

		if len(tm.GetTasks()) != 1 {
			t.Error("expected task to be restored after redo")
		}
	})

	t.Run("config persistence", func(t *testing.T) {
		configFile := filepath.Join(tempDir, "config.json")

		// Create custom config
		customConfig := config.DefaultConfig()
		customConfig.DataFile = dataFile
		customConfig.Theme.PrimaryColor = "#ABCDEF"
		customConfig.KeyMap.Add = "x"

		// Save config
		err := customConfig.SaveConfig(configFile)
		if err != nil {
			t.Errorf("failed to save config: %v", err)
		}

		// Load config
		loadedConfig, err := config.LoadConfig(configFile)
		if err != nil {
			t.Errorf("failed to load config: %v", err)
		}

		// Verify config values
		if loadedConfig.Theme.PrimaryColor != "#ABCDEF" {
			t.Errorf("expected primary color to be '#ABCDEF', got %s", loadedConfig.Theme.PrimaryColor)
		}
		if loadedConfig.KeyMap.Add != "x" {
			t.Errorf("expected add key to be 'x', got %s", loadedConfig.KeyMap.Add)
		}
		// Check that defaults are preserved for missing fields
		if loadedConfig.KeyMap.Delete != "d" {
			t.Errorf("expected default delete key 'd', got %s", loadedConfig.KeyMap.Delete)
		}
	})
}
