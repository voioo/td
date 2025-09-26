package ui

import "github.com/charmbracelet/bubbles/key"

// FullHelp returns the full help view for all key bindings.
func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Delete, k.Up, k.Down, k.Left, k.Right},
		{k.ListType, k.Priority, k.Filter, k.Escape},
		{k.Help, k.Quit, k.Undo, k.Redo},
		{k.PriorityNone, k.PriorityLow, k.PriorityMedium, k.PriorityHigh},
		{k.Home, k.End, k.ClearCompleted},
	}
}
