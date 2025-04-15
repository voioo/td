package main

import "github.com/charmbracelet/bubbles/key"

const (
	KeyAdd    = "a"
	KeyDelete = "d"
	KeyEnter  = "enter"
	KeyEscape = "esc"
	KeyUp     = "up"
	KeyDown   = "down"
	KeyLeft   = "left"
	KeyRight  = "right"
	KeyType   = "t"
	KeyHelp   = "?"
	KeyQuit   = "q"
)

const (
	ModeNormal = iota
	ModeDoneTaskList
	ModeAdditional
	ModeEdit
	ModeHelp
)

type keyMap struct {
	Add      key.Binding
	Up       key.Binding
	Down     key.Binding
	Delete   key.Binding
	Left     key.Binding
	Right    key.Binding
	Enter    key.Binding
	ListType key.Binding
	Escape   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Priority key.Binding
	Filter   key.Binding
	Undo     key.Binding
	Redo     key.Binding
}

var keys = keyMap{
	Add: key.NewBinding(
		key.WithKeys(KeyAdd),
		key.WithHelp(KeyAdd, "add new task"),
	),
	Delete: key.NewBinding(
		key.WithKeys(KeyDelete),
		key.WithHelp(KeyDelete, "delete task"),
	),
	Enter: key.NewBinding(
		key.WithKeys(KeyEnter),
		key.WithHelp(KeyEnter, "save"),
	),
	Up: key.NewBinding(
		key.WithKeys(KeyUp, "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys(KeyDown, "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys(KeyLeft, "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys(KeyRight, "l"),
		key.WithHelp("→/l", "move right"),
	),
	ListType: key.NewBinding(
		key.WithKeys("t", "tab"),
		key.WithHelp("t/tab", "list type"),
	),
	Escape: key.NewBinding(
		key.WithKeys(KeyEscape),
		key.WithHelp(KeyEscape, "back/cancel"),
	),
	Help: key.NewBinding(
		key.WithKeys(KeyHelp),
		key.WithHelp(KeyHelp, "toggle usage"),
	),
	Quit: key.NewBinding(
		key.WithKeys(KeyQuit, "ctrl+c"),
		key.WithHelp(KeyQuit, "quit"),
	),
	Priority: key.NewBinding(
		key.WithKeys("p"),
		key.WithHelp("p", "set priority"),
	),
	Filter: key.NewBinding(
		key.WithKeys("f"),
		key.WithHelp("f", "filter by priority"),
	),
	Undo: key.NewBinding(
		key.WithKeys("u"),
		key.WithHelp("u", "undo"),
	),
	Redo: key.NewBinding(
		key.WithKeys("ctrl+r"),
		key.WithHelp("ctrl+r", "redo"),
	),
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Delete, k.Up, k.Down, k.Left, k.Right},
		{k.ListType, k.Priority, k.Filter, k.Escape},
		{k.Help, k.Quit, k.Undo, k.Redo},
	}
}
