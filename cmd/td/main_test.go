package main

import (
	"os"
	"os/exec"
	"testing"
)

func TestVersionFlag(t *testing.T) {
	// Save original args and restore them after test
	originalArgs := os.Args
	defer func() { os.Args = originalArgs }()

	// Test cases
	testCases := []struct {
		name string
		args []string
	}{
		{"version flag", []string{"td", "--version"}},
		{"version shorthand", []string{"td", "-v"}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set test args
			os.Args = tc.args

			// We expect the program to exit with status 0
			// This is a simple test to ensure the flag is recognized
			if os.Getenv("TEST_VERSION_EXIT") == "1" {
				main()
				return
			}

			// Run the test in a subprocess
			executable, _ := os.Executable()
			cmd := exec.Command(executable, "-test.run=TestVersionFlag")
			cmd.Env = append(os.Environ(), "TEST_VERSION_EXIT=1")
			err := cmd.Run()

			if e, ok := err.(*exec.ExitError); ok && !e.Success() {
				t.Errorf("version flag test failed: %v", err)
			}
		})
	}
}

func newTestModel(tasks []*Task, doneTasks []*Task) model {
	m := model{
		keys:      keys,
		tasks:     tasks,
		doneTasks: doneTasks,
		cursor:    0,
		mode:      ModeNormal,
		filter:    FilterAll,
	}
	if len(tasks) > 0 {
		maxID := 0
		for _, t := range tasks {
			if t.ID > maxID {
				maxID = t.ID
			}
		}
		m.latestTaskID = maxID
	}
	return m
}

func TestModelCursorTracking(t *testing.T) {
	t.Run("cursor follows task after priority change", func(t *testing.T) {
		m := model{
			tasks: []*Task{
				{ID: 1, Name: "Task 1", Priority: PriorityNone},
				{ID: 2, Name: "Task 2", Priority: PriorityLow},
				{ID: 3, Name: "Task 3", Priority: PriorityHigh},
			},
			cursor: 2, // pointing to Task 2
		}

		// Change priority of Task 2 to High
		taskID := m.tasks[1].ID
		m.tasks[1].Priority = PriorityHigh
		m.followTask(taskID)

		// Task 2 should now be at position 2 (after Task 3)
		tasks := m.filteredTasks()
		var foundPos int
		for i, task := range tasks {
			if task.ID == taskID {
				foundPos = i + 1
				break
			}
		}

		if m.cursor != foundPos {
			t.Errorf("cursor should follow task, expected position %d, got %d", foundPos, m.cursor)
		}
	})

	t.Run("cursor adjusts after filtering", func(t *testing.T) {
		m := model{
			tasks: []*Task{
				{ID: 1, Name: "Task 1", Priority: PriorityNone},
				{ID: 2, Name: "Task 2", Priority: PriorityHigh},
				{ID: 3, Name: "Task 3", Priority: PriorityHigh},
			},
			cursor: 1,
			filter: FilterAll,
		}

		m.filter = FilterHigh
		tasks := m.filteredTasks()

		if len(tasks) != 2 {
			t.Errorf("expected 2 high priority tasks, got %d, tasks: %v",
				len(tasks),
				tasks,
			)
		}

		if m.cursor > len(tasks) {
			t.Error("cursor should adjust to valid position after filtering")
		}
	})
}

func TestModelUndoRedo(t *testing.T) {
	t.Run("add task", func(t *testing.T) {
		m := newTestModel([]*Task{
			{ID: 1, Name: "Task 1", Priority: PriorityNone},
		}, []*Task{})
		m.latestTaskID = 1

		// Simulate adding a task
		newTask := &Task{ID: 2, Name: "Task 2", Priority: PriorityNone}
		m.tasks = append(m.tasks, newTask)
		m.pushUndo(action{Type: "add", Task: newTask})

		if len(m.tasks) != 2 {
			t.Fatalf("Expected 2 tasks after add, got %d", len(m.tasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}
		if len(m.redoStack) != 0 {
			t.Fatalf("Expected 0 redo actions, got %d", len(m.redoStack))
		}

		// Perform Undo
		m.performUndo()

		if len(m.tasks) != 1 {
			t.Fatalf("Expected 1 task after undo, got %d", len(m.tasks))
		}
		if len(m.undoStack) != 0 {
			t.Fatalf("Expected 0 undo actions after undo, got %d", len(m.undoStack))
		}
		if len(m.redoStack) != 1 {
			t.Fatalf("Expected 1 redo action after undo, got %d", len(m.redoStack))
		}
		if m.tasks[0].ID != 1 {
			t.Errorf("Expected remaining task ID to be 1, got %d", m.tasks[0].ID)
		}

		// Perform Redo
		m.performRedo()

		if len(m.tasks) != 2 {
			t.Fatalf("Expected 2 tasks after redo, got %d", len(m.tasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action after redo, got %d", len(m.undoStack))
		}
		if len(m.redoStack) != 0 {
			t.Fatalf("Expected 0 redo actions after redo, got %d", len(m.redoStack))
		}
		if m.tasks[0].ID != 1 || m.tasks[1].ID != 2 {
			t.Errorf("Expected task IDs 1 and 2 after redo, got %d, %d", m.tasks[0].ID, m.tasks[1].ID)
		}
	})

	t.Run("delete task", func(t *testing.T) {
		taskToDelete := &Task{ID: 2, Name: "Task 2", Priority: PriorityNone}
		m := newTestModel([]*Task{
			{ID: 1, Name: "Task 1", Priority: PriorityNone},
			taskToDelete,
		}, []*Task{})

		// Simulate deleting a task
		m.tasks = m.tasks[:1] // Keep only the first task
		m.pushUndo(action{Type: "delete", Task: taskToDelete})

		if len(m.tasks) != 1 {
			t.Fatalf("Expected 1 task after delete, got %d", len(m.tasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}
		if len(m.redoStack) != 0 {
			t.Fatalf("Expected 0 redo actions, got %d", len(m.redoStack))
		}

		// Perform Undo
		m.performUndo()

		if len(m.tasks) != 2 {
			t.Fatalf("Expected 2 tasks after undo, got %d", len(m.tasks))
		}
		if len(m.undoStack) != 0 {
			t.Fatalf("Expected 0 undo actions after undo, got %d", len(m.undoStack))
		}
		if len(m.redoStack) != 1 {
			t.Fatalf("Expected 1 redo action after undo, got %d", len(m.redoStack))
		}
		if m.tasks[0].ID != 1 || m.tasks[1].ID != 2 {
			t.Errorf("Expected task IDs 1 and 2 after undo, got %d, %d", m.tasks[0].ID, m.tasks[1].ID)
		}

		// Perform Redo
		m.performRedo()

		if len(m.tasks) != 1 {
			t.Fatalf("Expected 1 task after redo, got %d", len(m.tasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action after redo, got %d", len(m.undoStack))
		}
		if len(m.redoStack) != 0 {
			t.Fatalf("Expected 0 redo actions after redo, got %d", len(m.redoStack))
		}
		if m.tasks[0].ID != 1 {
			t.Errorf("Expected remaining task ID to be 1 after redo, got %d", m.tasks[0].ID)
		}
	})

	t.Run("complete task", func(t *testing.T) {
		taskToComplete := &Task{ID: 1, Name: "Task 1", Priority: PriorityNone, IsDone: false}
		m := newTestModel([]*Task{taskToComplete}, []*Task{})

		// Simulate completing a task
		m.tasks = []*Task{}
		m.doneTasks = []*Task{taskToComplete}
		taskToComplete.IsDone = true
		m.pushUndo(action{Type: "complete", Task: taskToComplete})

		if len(m.tasks) != 0 {
			t.Fatalf("Expected 0 active tasks after complete, got %d", len(m.tasks))
		}
		if len(m.doneTasks) != 1 {
			t.Fatalf("Expected 1 done task after complete, got %d", len(m.doneTasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}

		// Perform Undo
		m.performUndo()

		if len(m.tasks) != 1 {
			t.Fatalf("Expected 1 active task after undo, got %d", len(m.tasks))
		}
		if len(m.doneTasks) != 0 {
			t.Fatalf("Expected 0 done tasks after undo, got %d", len(m.doneTasks))
		}
		if len(m.redoStack) != 1 {
			t.Fatalf("Expected 1 redo action, got %d", len(m.redoStack))
		}
		if m.tasks[0].IsDone {
			t.Error("Task should not be done after undo")
		}

		// Perform Redo
		m.performRedo()

		if len(m.tasks) != 0 {
			t.Fatalf("Expected 0 active tasks after redo, got %d", len(m.tasks))
		}
		if len(m.doneTasks) != 1 {
			t.Fatalf("Expected 1 done task after redo, got %d", len(m.doneTasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}
		if !m.doneTasks[0].IsDone {
			t.Error("Task should be done after redo")
		}
	})

	t.Run("uncomplete task", func(t *testing.T) {
		taskToUncomplete := &Task{ID: 1, Name: "Task 1", Priority: PriorityNone, IsDone: true}
		m := newTestModel([]*Task{}, []*Task{taskToUncomplete})

		// Simulate uncompleting a task
		m.tasks = []*Task{taskToUncomplete}
		m.doneTasks = []*Task{}
		taskToUncomplete.IsDone = false
		m.pushUndo(action{Type: "uncomplete", Task: taskToUncomplete})

		if len(m.tasks) != 1 {
			t.Fatalf("Expected 1 active task after uncomplete, got %d", len(m.tasks))
		}
		if len(m.doneTasks) != 0 {
			t.Fatalf("Expected 0 done tasks after uncomplete, got %d", len(m.doneTasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}

		// Perform Undo
		m.performUndo()

		if len(m.tasks) != 0 {
			t.Fatalf("Expected 0 active tasks after undo, got %d", len(m.tasks))
		}
		if len(m.doneTasks) != 1 {
			t.Fatalf("Expected 1 done task after undo, got %d", len(m.doneTasks))
		}
		if len(m.redoStack) != 1 {
			t.Fatalf("Expected 1 redo action, got %d", len(m.redoStack))
		}
		if !m.doneTasks[0].IsDone {
			t.Error("Task should be done after undo")
		}

		// Perform Redo
		m.performRedo()

		if len(m.tasks) != 1 {
			t.Fatalf("Expected 1 active task after redo, got %d", len(m.tasks))
		}
		if len(m.doneTasks) != 0 {
			t.Fatalf("Expected 0 done tasks after redo, got %d", len(m.doneTasks))
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}
		if m.tasks[0].IsDone {
			t.Error("Task should not be done after redo")
		}
	})

	t.Run("edit task", func(t *testing.T) {
		originalName := "Task Original"
		newName := "Task Edited"
		taskToEdit := &Task{ID: 1, Name: originalName, Priority: PriorityNone}
		m := newTestModel([]*Task{taskToEdit}, []*Task{})

		// Simulate editing a task
		m.tasks[0].Name = newName
		m.pushUndo(action{Type: "edit", Task: taskToEdit, OldState: originalName, NewState: newName})

		if m.tasks[0].Name != newName {
			t.Fatalf("Expected task name to be '%s', got '%s'", newName, m.tasks[0].Name)
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}

		// Perform Undo
		m.performUndo()

		if m.tasks[0].Name != originalName {
			t.Fatalf("Expected task name to be '%s' after undo, got '%s'", originalName, m.tasks[0].Name)
		}
		if len(m.redoStack) != 1 {
			t.Fatalf("Expected 1 redo action, got %d", len(m.redoStack))
		}

		// Perform Redo
		m.performRedo()

		if m.tasks[0].Name != newName {
			t.Fatalf("Expected task name to be '%s' after redo, got '%s'", newName, m.tasks[0].Name)
		}
		if len(m.undoStack) != 1 {
			t.Fatalf("Expected 1 undo action, got %d", len(m.undoStack))
		}
	})
}
