package task

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

	t.Run("priority string representations", func(t *testing.T) {
		tests := []struct {
			priority Priority
			expected string
		}{
			{PriorityNone, "none"},
			{PriorityLow, "low"},
			{PriorityMedium, "medium"},
			{PriorityHigh, "high"},
		}

		for _, test := range tests {
			if test.priority.String() != test.expected {
				t.Errorf("expected %v.String() to be %s, got %s", test.priority, test.expected, test.priority.String())
			}
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

func TestTaskManager(t *testing.T) {
	t.Run("new task manager", func(t *testing.T) {
		tasks := []*Task{
			{ID: 1, Name: "Task 1", Priority: PriorityNone},
			{ID: 2, Name: "Task 2", Priority: PriorityHigh},
		}
		doneTasks := []*Task{
			{ID: 3, Name: "Done Task", Priority: PriorityLow, IsDone: true},
		}

		tm := NewTaskManager(tasks, doneTasks, 4)

		if len(tm.GetTasks()) != 2 {
			t.Errorf("expected 2 active tasks, got %d", len(tm.GetTasks()))
		}
		if len(tm.GetDoneTasks()) != 1 {
			t.Errorf("expected 1 done task, got %d", len(tm.GetDoneTasks()))
		}
		if tm.GetNextID() != 4 {
			t.Errorf("expected next ID to be 4, got %d", tm.GetNextID())
		}
	})

	t.Run("add task", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 0) // Start with 0, first task gets ID 1

		task := tm.AddTask("New Task")

		if task.Name != "New Task" {
			t.Errorf("expected task name to be 'New Task', got %s", task.Name)
		}
		if task.ID != 1 {
			t.Errorf("expected task ID to be 1, got %d", task.ID)
		}
		if tm.GetNextID() != 1 {
			t.Errorf("expected next ID to be 1, got %d", tm.GetNextID())
		}

		// Add another task to verify ID increment
		task2 := tm.AddTask("Another Task")
		if task2.ID != 2 {
			t.Errorf("expected second task ID to be 2, got %d", task2.ID)
		}
		if tm.GetNextID() != 2 {
			t.Errorf("expected next ID to be 2, got %d", tm.GetNextID())
		}
	})

	t.Run("complete task", func(t *testing.T) {
		task := &Task{ID: 1, Name: "Test", Priority: PriorityNone, IsDone: false}
		tm := NewTaskManager([]*Task{task}, []*Task{}, 2)

		completedTask := tm.CompleteTask(1)

		if completedTask == nil {
			t.Fatal("expected task to be completed")
		}
		if len(tm.GetTasks()) != 0 {
			t.Errorf("expected 0 active tasks, got %d", len(tm.GetTasks()))
		}
		if len(tm.GetDoneTasks()) != 1 {
			t.Errorf("expected 1 done task, got %d", len(tm.GetDoneTasks()))
		}
	})

	t.Run("uncomplete task", func(t *testing.T) {
		task := &Task{ID: 1, Name: "Test", Priority: PriorityNone, IsDone: true}
		tm := NewTaskManager([]*Task{}, []*Task{task}, 2)

		uncompletedTask := tm.UncompleteTask(1)

		if uncompletedTask == nil {
			t.Fatal("expected task to be uncompleted")
		}
		if len(tm.GetTasks()) != 1 {
			t.Errorf("expected 1 active task, got %d", len(tm.GetTasks()))
		}
		if len(tm.GetDoneTasks()) != 0 {
			t.Errorf("expected 0 done tasks, got %d", len(tm.GetDoneTasks()))
		}
	})

	t.Run("delete task", func(t *testing.T) {
		task := &Task{ID: 1, Name: "Test", Priority: PriorityNone}
		tm := NewTaskManager([]*Task{task}, []*Task{}, 2)

		deletedTask := tm.DeleteTask(1)

		if deletedTask == nil {
			t.Fatal("expected task to be deleted")
		}
		if len(tm.GetTasks()) != 0 {
			t.Errorf("expected 0 active tasks, got %d", len(tm.GetTasks()))
		}
	})

	t.Run("update task name", func(t *testing.T) {
		task := &Task{ID: 1, Name: "Old Name", Priority: PriorityNone}
		tm := NewTaskManager([]*Task{task}, []*Task{}, 2)

		updatedTask := tm.UpdateTaskName(1, "New Name")

		if updatedTask == nil {
			t.Fatal("expected task to be updated")
		}
		if updatedTask.Name != "New Name" {
			t.Errorf("expected task name to be 'New Name', got %s", updatedTask.Name)
		}
	})

	t.Run("set task priority", func(t *testing.T) {
		task := &Task{ID: 1, Name: "Test", Priority: PriorityNone}
		tm := NewTaskManager([]*Task{task}, []*Task{}, 2)

		updatedTask := tm.SetTaskPriority(1, PriorityHigh)

		if updatedTask == nil {
			t.Fatal("expected task priority to be set")
		}
		if updatedTask.Priority != PriorityHigh {
			t.Errorf("expected priority to be High, got %v", updatedTask.Priority)
		}
	})
}
