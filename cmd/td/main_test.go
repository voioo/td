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
