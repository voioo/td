package task

import (
	"testing"
)

func TestUndoManager(t *testing.T) {
	t.Run("new undo manager", func(t *testing.T) {
		um := NewUndoManager(10)

		if um.CanUndo() {
			t.Error("expected new undo manager to not have undo actions")
		}
		if um.CanRedo() {
			t.Error("expected new undo manager to not have redo actions")
		}
	})

	t.Run("default max size", func(t *testing.T) {
		um := NewUndoManager(0)

		// Test that it uses default size
		for i := 0; i < DefaultMaxUndoSize+10; i++ {
			um.PushUndo(Action{Type: ActionTypeAdd, Task: &Task{ID: i}})
		}

		// Should be able to undo up to the default max size
		undoCount := 0
		for um.CanUndo() {
			um.Undo(NewTaskManager([]*Task{}, []*Task{}, 1))
			undoCount++
		}

		if undoCount != DefaultMaxUndoSize {
			t.Errorf("expected to undo %d times, got %d", DefaultMaxUndoSize, undoCount)
		}
	})

	t.Run("push undo and undo", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add a task
		task := tm.AddTask("Test Task")
		um.PushUndo(Action{Type: ActionTypeAdd, Task: task})

		if !um.CanUndo() {
			t.Error("expected to be able to undo after push")
		}

		// Undo the add
		if !um.Undo(tm) {
			t.Error("expected undo to succeed")
		}

		if len(tm.GetTasks()) != 0 {
			t.Error("expected task to be removed after undo")
		}

		if !um.CanRedo() {
			t.Error("expected to be able to redo after undo")
		}
	})

	t.Run("redo", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add a task
		task := tm.AddTask("Test Task")
		um.PushUndo(Action{Type: ActionTypeAdd, Task: task})

		// Undo
		um.Undo(tm)

		// Redo
		if !um.Redo(tm) {
			t.Error("expected redo to succeed")
		}

		if len(tm.GetTasks()) != 1 {
			t.Error("expected task to be restored after redo")
		}
	})

	t.Run("undo delete", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add and then delete a task
		task := tm.AddTask("Test Task")
		deletedTask := tm.DeleteTask(task.ID)
		um.PushUndo(Action{Type: ActionTypeDelete, Task: deletedTask})

		// Undo the delete
		if !um.Undo(tm) {
			t.Error("expected undo delete to succeed")
		}

		if len(tm.GetTasks()) != 1 {
			t.Error("expected task to be restored after undo delete")
		}
	})

	t.Run("undo complete", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add and complete a task
		task := tm.AddTask("Test Task")
		completedTask := tm.CompleteTask(task.ID)
		um.PushUndo(Action{Type: ActionTypeComplete, Task: completedTask})

		// Undo the complete
		if !um.Undo(tm) {
			t.Error("expected undo complete to succeed")
		}

		if len(tm.GetTasks()) != 1 {
			t.Error("expected task to be restored to active after undo complete")
		}
		if len(tm.GetDoneTasks()) != 0 {
			t.Error("expected task to be removed from done after undo complete")
		}
	})

	t.Run("undo edit", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add a task and edit its name
		task := tm.AddTask("Original Name")
		oldName := task.Name
		tm.UpdateTaskName(task.ID, "New Name")
		um.PushUndo(Action{
			Type:     ActionTypeEdit,
			Task:     task,
			OldState: oldName,
			NewState: "New Name",
		})

		// Undo the edit
		if !um.Undo(tm) {
			t.Error("expected undo edit to succeed")
		}

		if task.Name != "Original Name" {
			t.Errorf("expected task name to be restored to 'Original Name', got %s", task.Name)
		}
	})

	t.Run("undo priority change", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add a task and change its priority
		task := tm.AddTask("Test Task")
		oldPriority := task.Priority
		tm.SetTaskPriority(task.ID, PriorityHigh)
		um.PushUndo(Action{
			Type:     ActionTypePriority,
			Task:     task,
			OldState: oldPriority,
			NewState: PriorityHigh,
		})

		// Undo the priority change
		if !um.Undo(tm) {
			t.Error("expected undo priority change to succeed")
		}

		if task.Priority != PriorityNone {
			t.Errorf("expected task priority to be restored to None, got %v", task.Priority)
		}
	})

	t.Run("clear", func(t *testing.T) {
		tm := NewTaskManager([]*Task{}, []*Task{}, 1)
		um := NewUndoManager(10)

		// Add some actions
		task := tm.AddTask("Test Task")
		um.PushUndo(Action{Type: ActionTypeAdd, Task: task})

		um.Clear()

		if um.CanUndo() {
			t.Error("expected undo stack to be cleared")
		}
		if um.CanRedo() {
			t.Error("expected redo stack to be cleared")
		}
	})
}
