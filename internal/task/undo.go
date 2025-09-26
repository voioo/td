// Package task provides undo/redo functionality for task operations.
package task

// ActionType represents the type of action that can be undone/redone.
type ActionType string

// Action represents a reversible action on tasks.
type Action struct {
	Type     ActionType
	Task     *Task
	OldState interface{}
	NewState interface{}
}

// UndoManager manages undo and redo operations.
type UndoManager struct {
	undoStack []Action
	redoStack []Action
	maxSize   int
}

// NewUndoManager creates a new undo manager with the specified maximum stack size.
func NewUndoManager(maxSize int) *UndoManager {
	if maxSize <= 0 {
		maxSize = DefaultMaxUndoSize
	}
	return &UndoManager{
		undoStack: make([]Action, 0, maxSize),
		redoStack: make([]Action, 0, maxSize),
		maxSize:   maxSize,
	}
}

// PushUndo adds an action to the undo stack and clears the redo stack.
func (um *UndoManager) PushUndo(action Action) {
	um.undoStack = append(um.undoStack, action)
	um.redoStack = nil // Clear redo stack when new action is performed

	// Maintain maximum stack size
	if len(um.undoStack) > um.maxSize {
		um.undoStack = um.undoStack[1:]
	}
}

// CanUndo returns true if there are actions that can be undone.
func (um *UndoManager) CanUndo() bool {
	return len(um.undoStack) > 0
}

// CanRedo returns true if there are actions that can be redone.
func (um *UndoManager) CanRedo() bool {
	return len(um.redoStack) > 0
}

// Undo performs the last undone action.
func (um *UndoManager) Undo(taskManager *TaskManager) bool {
	if len(um.undoStack) == 0 {
		return false
	}

	lastAction := um.undoStack[len(um.undoStack)-1]
	um.undoStack = um.undoStack[:len(um.undoStack)-1]

	switch lastAction.Type {
	case ActionTypeAdd:
		// Remove the added task
		taskManager.DeleteTask(lastAction.Task.ID)
	case ActionTypeDelete:
		// Restore the deleted task
		if lastAction.Task.IsDone {
			taskManager.doneTasks = append(taskManager.doneTasks, lastAction.Task)
		} else {
			taskManager.tasks = append(taskManager.tasks, lastAction.Task)
			taskManager.sortTasks()
		}
	case ActionTypeComplete:
		// Mark as incomplete and move back to active tasks
		lastAction.Task.IsDone = false
		taskManager.tasks = append(taskManager.tasks, lastAction.Task)
		// Remove from done tasks
		for i, task := range taskManager.doneTasks {
			if task.ID == lastAction.Task.ID {
				taskManager.doneTasks = append(taskManager.doneTasks[:i], taskManager.doneTasks[i+1:]...)
				break
			}
		}
		taskManager.sortTasks()
	case ActionTypeUncomplete:
		// Mark as complete and move to done tasks
		lastAction.Task.IsDone = true
		taskManager.doneTasks = append(taskManager.doneTasks, lastAction.Task)
		// Remove from active tasks
		for i, task := range taskManager.tasks {
			if task.ID == lastAction.Task.ID {
				taskManager.tasks = append(taskManager.tasks[:i], taskManager.tasks[i+1:]...)
				break
			}
		}
	case ActionTypeEdit:
		// Restore old name
		if oldName, ok := lastAction.OldState.(string); ok {
			lastAction.Task.Name = oldName
		}
	case ActionTypePriority:
		// Restore old priority
		if oldPriority, ok := lastAction.OldState.(Priority); ok {
			lastAction.Task.Priority = oldPriority
			taskManager.sortTasks()
		}
	}

	um.redoStack = append(um.redoStack, lastAction)
	return true
}

// Redo performs the last undone action.
func (um *UndoManager) Redo(taskManager *TaskManager) bool {
	if len(um.redoStack) == 0 {
		return false
	}

	lastAction := um.redoStack[len(um.redoStack)-1]
	um.redoStack = um.redoStack[:len(um.redoStack)-1]

	var correspondingUndoAction Action

	switch lastAction.Type {
	case ActionTypeAdd:
		// Re-add the task
		taskManager.tasks = append(taskManager.tasks, lastAction.Task)
		taskManager.sortTasks()
		correspondingUndoAction = Action{Type: ActionTypeAdd, Task: lastAction.Task}
	case ActionTypeDelete:
		// Re-delete the task
		taskManager.DeleteTask(lastAction.Task.ID)
		correspondingUndoAction = Action{Type: ActionTypeDelete, Task: lastAction.Task}
	case ActionTypeComplete:
		// Re-complete the task
		taskManager.CompleteTask(lastAction.Task.ID)
		correspondingUndoAction = Action{Type: ActionTypeComplete, Task: lastAction.Task}
	case ActionTypeUncomplete:
		// Re-uncomplete the task
		taskManager.UncompleteTask(lastAction.Task.ID)
		correspondingUndoAction = Action{Type: ActionTypeUncomplete, Task: lastAction.Task}
	case ActionTypeEdit:
		// Re-apply the edit
		if newName, ok := lastAction.NewState.(string); ok {
			currentName := lastAction.Task.Name
			lastAction.Task.Name = newName
			correspondingUndoAction = Action{
				Type:     ActionTypeEdit,
				Task:     lastAction.Task,
				OldState: currentName,
				NewState: newName,
			}
		}
	case ActionTypePriority:
		// Re-apply the priority change
		if newPriority, ok := lastAction.NewState.(Priority); ok {
			currentPriority := lastAction.Task.Priority
			lastAction.Task.Priority = newPriority
			taskManager.sortTasks()
			correspondingUndoAction = Action{
				Type:     ActionTypePriority,
				Task:     lastAction.Task,
				OldState: currentPriority,
				NewState: newPriority,
			}
		}
	}

	if correspondingUndoAction.Type != "" {
		um.undoStack = append(um.undoStack, correspondingUndoAction)
	}

	return true
}

// Clear clears both undo and redo stacks.
func (um *UndoManager) Clear() {
	um.undoStack = nil
	um.redoStack = nil
}
