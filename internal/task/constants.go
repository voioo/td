package task

// Action types for undo/redo operations
const (
	ActionTypeAdd        ActionType = "add"
	ActionTypeDelete     ActionType = "delete"
	ActionTypeComplete   ActionType = "complete"
	ActionTypeUncomplete ActionType = "uncomplete"
	ActionTypeEdit       ActionType = "edit"
	ActionTypePriority   ActionType = "priority"
)

// Default maximum undo stack size
const DefaultMaxUndoSize = 100
