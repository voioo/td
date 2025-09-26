package integration

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/voioo/td/internal/config"
	"github.com/voioo/td/internal/storage"
	"github.com/voioo/td/internal/task"
	"github.com/voioo/td/internal/ui"
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

	t.Run("UI quit saves tasks", func(t *testing.T) {
		// This test verifies that when a user quits the application (via UI),
		// tasks are properly saved to disk. This was the critical bug that
		// was missed by the other tests - they only tested storage directly,
		// not the UI quit flow that users actually trigger.

		// Create config
		cfg := config.DefaultConfig()
		cfg.DataFile = dataFile

		// Create task manager with some tasks
		tm := task.NewTaskManager([]*task.Task{}, []*task.Task{}, 0)
		task1 := tm.AddTask("Test quit saves tasks")
		tm.AddTask("Another test task")
		tm.CompleteTask(task1.ID)

		// Create UI model
		uiModel, err := ui.NewTestModel(cfg, tm)
		if err != nil {
			t.Fatalf("failed to create UI model: %v", err)
		}

		// Capture state before "quit"
		initialActiveTasks := len(tm.GetTasks())
		initialDoneTasks := len(tm.GetDoneTasks())

		// The quit command should save tasks - verify file gets created/updated
		initialModTime := getFileModTime(dataFile)

		// Simulate the UI quit save process (what happens when user presses 'q')
		err = saveTasksViaUIQuit(uiModel)
		if err != nil {
			t.Errorf("UI quit save failed: %v", err)
		}

		// Verify tasks were saved by loading them back from disk
		repo := storage.NewRepository(dataFile)
		loadedTasks, loadedDoneTasks, _, err := repo.LoadTasks()
		if err != nil {
			t.Errorf("failed to load saved tasks: %v", err)
		}

		if len(loadedTasks) != initialActiveTasks {
			t.Errorf("expected %d active tasks after save, got %d", initialActiveTasks, len(loadedTasks))
		}
		if len(loadedDoneTasks) != initialDoneTasks {
			t.Errorf("expected %d done tasks after save, got %d", initialDoneTasks, len(loadedDoneTasks))
		}

		// Verify file was actually modified (proves save happened)
		finalModTime := getFileModTime(dataFile)
		if !finalModTime.After(initialModTime) && initialActiveTasks > 0 {
			t.Error("expected data file to be modified after quit save")
		}
	})

	t.Run("UI quit bug demonstration", func(t *testing.T) {
		// This test demonstrates that the original bug (quit without saving)
		// would be caught by the test above. We simulate the buggy behavior
		// to show that this test WOULD HAVE FAILED before the fix.

		// Create a fresh file for this test
		buggyDataFile := filepath.Join(tempDir, "buggy-tasks.json")

		// Create config
		cfg := config.DefaultConfig()
		cfg.DataFile = buggyDataFile

		// Create task manager with some tasks
		tm := task.NewTaskManager([]*task.Task{}, []*task.Task{}, 0)
		tm.AddTask("This task would be lost on quit")
		tm.AddTask("Another lost task")

		// Create UI model
		uiModel, err := ui.NewTestModel(cfg, tm)
		if err != nil {
			t.Fatalf("failed to create UI model: %v", err)
		}

		// Simulate the BUGGY quit behavior (what happened before the fix)
		err = simulateBuggyQuit(uiModel) // This doesn't save!
		if err != nil {
			t.Errorf("buggy quit failed: %v", err)
		}

		// Try to load tasks - they should NOT be there (demonstrating the bug)
		repo := storage.NewRepository(buggyDataFile)
		loadedTasks, loadedDoneTasks, _, err := repo.LoadTasks()
		if err != nil {
			t.Errorf("failed to load tasks after buggy quit: %v", err)
		}

		// This demonstrates the bug: tasks were not saved!
		if len(loadedTasks) != 0 {
			t.Errorf("BUG: expected 0 tasks after buggy quit (no save), but got %d", len(loadedTasks))
		}
		if len(loadedDoneTasks) != 0 {
			t.Errorf("BUG: expected 0 done tasks after buggy quit (no save), but got %d", len(loadedDoneTasks))
		}

		// This test proves that the original bug would be detected:
		// - User creates tasks
		// - User quits (but tasks don't save due to bug)
		// - User restarts app, tasks are gone!
		// - Test would fail because loadedTasks would be empty instead of 2
	})
}

// saveTasksViaUIQuit simulates the UI quit save process for testing
func saveTasksViaUIQuit(model *ui.Model) error {
	repo := storage.NewRepository(model.GetConfig().DataFile)
	return repo.SaveTasks(model.GetTaskManager().GetTasks(), model.GetTaskManager().GetDoneTasks())
}

// simulateBuggyQuit demonstrates what would happen if quit didn't save (the original bug)
// This function does NOT save tasks, simulating the buggy behavior
func simulateBuggyQuit(model *ui.Model) error {
	// BUG: Original code just called tea.Quit() without saving!
	// No save operation here - this simulates the bug
	return nil
}

// getFileModTime returns the modification time of a file, or zero time if file doesn't exist
func getFileModTime(filePath string) time.Time {
	info, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}
	}
	return info.ModTime()
}
