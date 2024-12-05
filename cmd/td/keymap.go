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
	KeyQuit   = "esc"
)

const (
	ModeNormal = iota
	ModeDoneTaskList
	ModeAdditional
	ModeEdit
	ModeHelp
)

type keyMap struct {
	Add    key.Binding
	Up     key.Binding
	Down   key.Binding
	Delete key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Escape key.Binding
	Type   key.Binding
	Help   key.Binding
	Quit   key.Binding
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
	Escape: key.NewBinding(
		key.WithKeys(KeyEscape),
		key.WithHelp(KeyEscape, "go back/quit"),
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
	Type: key.NewBinding(
		key.WithKeys(KeyType),
		key.WithHelp(KeyType, "list type"),
	),
	Help: key.NewBinding(
		key.WithKeys(KeyHelp, KeyHelp),
		key.WithHelp(KeyHelp, "toggle usage"),
	),
	Quit: key.NewBinding(
		key.WithKeys(KeyQuit, "ctrl+c"),
		key.WithHelp(KeyQuit, "quit"),
	),
}
