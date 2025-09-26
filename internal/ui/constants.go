package ui

// Default UI settings
const (
	DefaultMaxUndoSize = 100
	DefaultDataFile    = "~/.td.json"
)

// UI mode names for display
const (
	ModeNameNormal   = "Normal"
	ModeNameDoneList = "Completed Tasks"
	ModeNameAdd      = "Add Task"
	ModeNameEdit     = "Edit Task"
	ModeNameHelp     = "Help"
)

// Filter mode names
const (
	FilterNameAll    = "All"
	FilterNameNone   = "None"
	FilterNameLow    = "Low Priority"
	FilterNameMedium = "Medium Priority"
	FilterNameHigh   = "High Priority"
)
