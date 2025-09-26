package task

// Action types for undo/redo operations
const (
	ActionTypeAdd        = "add"
	ActionTypeDelete     = "delete"
	ActionTypeComplete   = "complete"
	ActionTypeUncomplete = "uncomplete"
	ActionTypeEdit       = "edit"
	ActionTypePriority   = "priority"
)

// Default maximum undo stack size
const DefaultMaxUndoSize = 100
