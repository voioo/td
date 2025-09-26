// Package ui provides the terminal user interface for the td application.
package ui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/voioo/td/internal/config"
	"github.com/voioo/td/internal/logger"
	"github.com/voioo/td/internal/storage"
	"github.com/voioo/td/internal/task"
)

// Mode represents different UI modes.
type Mode int

const (
	ModeNormal = iota
	ModeDoneTaskList
	ModeAdditional
	ModeEdit
	ModeHelp
)

// FilterMode represents different task filtering modes.
type FilterMode int

const (
	FilterAll FilterMode = iota
	FilterNone
	FilterLow
	FilterMedium
	FilterHigh
)

// saveAndQuitMsg is sent when the application should save and quit.
type saveAndQuitMsg struct{}

// Model represents the main UI model.
type Model struct {
	config      *config.Config
	taskManager *task.TaskManager
	undoManager *task.UndoManager

	// UI components
	help              help.Model
	inputStyle        lipgloss.Style
	keys              KeyMap
	newTaskNameInput  input.Model
	editTaskNameInput input.Model

	// UI state
	cursor     int
	mode       Mode
	filter     FilterMode
	quitting   bool
	taskCache  []*task.Task // Cache for filtered tasks
	cacheValid bool
}

// KeyMap defines the key bindings for the UI.
type KeyMap struct {
	Add      key.Binding
	Up       key.Binding
	Down     key.Binding
	Delete   key.Binding
	Left     key.Binding
	Right    key.Binding
	Edit     key.Binding
	Enter    key.Binding
	ListType key.Binding
	Escape   key.Binding
	Help     key.Binding
	Quit     key.Binding
	Filter   key.Binding
	Undo     key.Binding
	Redo     key.Binding
	// New shortcuts
	PriorityNone   key.Binding
	PriorityLow    key.Binding
	PriorityMedium key.Binding
	PriorityHigh   key.Binding
	Home           key.Binding
	End            key.Binding
	ClearCompleted key.Binding
}

// NewModel creates a new UI model with the given configuration and task manager.
func NewModel(cfg *config.Config, taskManager *task.TaskManager) *Model {
	// Create key bindings from config
	keys := KeyMap{
		Add: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Add),
			key.WithHelp(cfg.KeyMap.Add, "add new task"),
		),
		Delete: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Delete),
			key.WithHelp(cfg.KeyMap.Delete, "delete task"),
		),
		Enter: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Enter),
			key.WithHelp(cfg.KeyMap.Enter, "save"),
		),
		Up: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Up, "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Down, "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Left, "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Right, "l"),
			key.WithHelp("→/l", "move right"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		ListType: key.NewBinding(
			key.WithKeys(cfg.KeyMap.ListType, "tab"),
			key.WithHelp("t/tab", "list type"),
		),
		Escape: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Escape),
			key.WithHelp(cfg.KeyMap.Escape, "back/cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Help),
			key.WithHelp(cfg.KeyMap.Help, "toggle usage"),
		),
		Quit: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Quit, "ctrl+c"),
			key.WithHelp(cfg.KeyMap.Quit, "quit"),
		),
		Filter: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Filter),
			key.WithHelp(cfg.KeyMap.Filter, "filter by priority"),
		),
		Undo: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Undo),
			key.WithHelp(cfg.KeyMap.Undo, "undo"),
		),
		Redo: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Redo),
			key.WithHelp(cfg.KeyMap.Redo, "redo"),
		),
		// New shortcuts
		PriorityNone: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "set no priority"),
		),
		PriorityLow: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "set low priority"),
		),
		PriorityMedium: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "set medium priority"),
		),
		PriorityHigh: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "set high priority"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "go to top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "go to bottom"),
		),
		ClearCompleted: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("C", "clear completed tasks"),
		),
	}

	// Create input models
	newTaskNameModel := input.New()
	newTaskNameModel.Placeholder = "New task name..."
	editTaskNameModel := input.New()

	m := &Model{
		config:            cfg,
		taskManager:       taskManager,
		undoManager:       task.NewUndoManager(100),
		keys:              keys,
		help:              help.New(),
		inputStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.PrimaryColor)),
		newTaskNameInput:  newTaskNameModel,
		editTaskNameInput: editTaskNameModel,
		cursor:            0,
		mode:              ModeNormal,
		filter:            FilterAll,
		taskCache:         []*task.Task{},
		cacheValid:        false,
	}

	// Set initial cursor position
	m.updateTaskCache()

	return m
}

// NewTestModel creates a minimal UI model for testing purposes.
// It skips initializing Bubble Tea input components that require a TTY.
func NewTestModel(cfg *config.Config, taskManager *task.TaskManager) (*Model, error) {
	// Create key bindings from config
	keys := KeyMap{
		Add: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Add),
			key.WithHelp(cfg.KeyMap.Add, "add new task"),
		),
		Delete: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Delete),
			key.WithHelp(cfg.KeyMap.Delete, "delete task"),
		),
		Enter: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Enter),
			key.WithHelp(cfg.KeyMap.Enter, "save"),
		),
		Up: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Up, "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Down, "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Left, "h"),
			key.WithHelp("←/h", "move left"),
		),
		Right: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Right, "l"),
			key.WithHelp("→/l", "move right"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit task"),
		),
		ListType: key.NewBinding(
			key.WithKeys(cfg.KeyMap.ListType, "tab"),
			key.WithHelp("t/tab", "list type"),
		),
		Escape: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Escape),
			key.WithHelp(cfg.KeyMap.Escape, "back/cancel"),
		),
		Help: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Help),
			key.WithHelp(cfg.KeyMap.Help, "toggle usage"),
		),
		Quit: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Quit, "ctrl+c"),
			key.WithHelp(cfg.KeyMap.Quit, "quit"),
		),
		Filter: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Filter),
			key.WithHelp(cfg.KeyMap.Filter, "filter by priority"),
		),
		Undo: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Undo),
			key.WithHelp(cfg.KeyMap.Undo, "undo"),
		),
		Redo: key.NewBinding(
			key.WithKeys(cfg.KeyMap.Redo),
			key.WithHelp(cfg.KeyMap.Redo, "redo"),
		),
		PriorityNone: key.NewBinding(
			key.WithKeys("1"),
			key.WithHelp("1", "set no priority"),
		),
		PriorityLow: key.NewBinding(
			key.WithKeys("2"),
			key.WithHelp("2", "set low priority"),
		),
		PriorityMedium: key.NewBinding(
			key.WithKeys("3"),
			key.WithHelp("3", "set medium priority"),
		),
		PriorityHigh: key.NewBinding(
			key.WithKeys("4"),
			key.WithHelp("4", "set high priority"),
		),
		Home: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("home/g", "go to top"),
		),
		End: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("end/G", "go to bottom"),
		),
		ClearCompleted: key.NewBinding(
			key.WithKeys("C"),
			key.WithHelp("C", "clear completed tasks"),
		),
	}

	m := &Model{
		config:      cfg,
		taskManager: taskManager,
		undoManager: task.NewUndoManager(100),
		keys:        keys,
		help:        help.New(),
		inputStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color(cfg.Theme.PrimaryColor)),
		// Skip input models for testing - they're not needed for quit save testing
		cursor:     0,
		mode:       ModeNormal,
		filter:     FilterAll,
		taskCache:  []*task.Task{},
		cacheValid: false,
	}

	// Set initial cursor position
	m.updateTaskCache()

	return m, nil
}

// saveAndQuitCmd returns a command that saves tasks and then quits.
func (m *Model) saveAndQuitCmd() tea.Cmd {
	return tea.Sequence(
		func() tea.Msg {
			logger.Info("Saving tasks before quit",
				logger.F("active_tasks", len(m.taskManager.GetTasks())),
				logger.F("done_tasks", len(m.taskManager.GetDoneTasks())))

			repo := storage.NewRepository(m.config.DataFile)
			err := repo.SaveTasks(m.taskManager.GetTasks(), m.taskManager.GetDoneTasks())
			if err != nil {
				logger.Error("Failed to save tasks", logger.F("error", err))
			} else {
				logger.Info("Tasks saved successfully")
			}
			return saveAndQuitMsg{}
		},
		tea.Quit,
	)
}

// Init initializes the Bubble Tea model.
func (m *Model) Init() tea.Cmd {
	return nil
}

// Update handles UI updates based on messages.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case saveAndQuitMsg:
		m.quitting = true
		return m, tea.Quit
	default:
		switch m.mode {
		case ModeNormal:
			return m.normalUpdate(msg)
		case ModeDoneTaskList:
			return m.doneTaskListUpdate(msg)
		case ModeAdditional:
			return m.addingTaskUpdate(msg)
		case ModeEdit:
			return m.editTaskUpdate(msg)
		case ModeHelp:
			return m.helpUpdate(msg)
		default:
			return m, nil
		}
	}
}

// View renders the current UI state.
func (m *Model) View() string {
	switch m.mode {
	case ModeNormal, ModeDoneTaskList:
		return m.normalView()
	case ModeAdditional:
		return m.addingTaskView()
	case ModeEdit:
		return m.editTaskView()
	case ModeHelp:
		return m.helpView()
	}
	return ""
}

// GetTaskManager returns the task manager.
func (m *Model) GetTaskManager() *task.TaskManager {
	return m.taskManager
}

// GetConfig returns the configuration.
func (m *Model) GetConfig() *config.Config {
	return m.config
}

// updateTaskCache updates the cached filtered tasks.
func (m *Model) updateTaskCache() {
	if m.cacheValid {
		return
	}

	tasks := m.taskManager.GetTasks()
	if m.filter == FilterAll {
		m.taskCache = tasks
	} else {
		m.taskCache = []*task.Task{}
		for _, task := range tasks {
			if filterToPriority(m.filter) == task.Priority {
				m.taskCache = append(m.taskCache, task)
			}
		}
	}

	task.SortTasksByPriority(m.taskCache)
	m.cacheValid = true
}

// invalidateCache marks the task cache as invalid.
func (m *Model) invalidateCache() {
	m.cacheValid = false
}

// getCurrentTask returns the currently selected task.
func (m *Model) getCurrentTask() *task.Task {
	m.updateTaskCache()
	if m.cursor > 0 && m.cursor <= len(m.taskCache) {
		return m.taskCache[m.cursor-1]
	}
	return nil
}

// followTask adjusts cursor position to follow a specific task.
func (m *Model) followTask(taskID int) {
	m.updateTaskCache()
	for i, task := range m.taskCache {
		if task.ID == taskID {
			m.cursor = i + 1
			return
		}
	}
	if len(m.taskCache) == 0 {
		m.cursor = 0
	} else if m.cursor > len(m.taskCache) {
		m.cursor = len(m.taskCache)
	}
}

// filterToPriority converts filter mode to priority.
func filterToPriority(f FilterMode) task.Priority {
	switch f {
	case FilterHigh:
		return task.PriorityHigh
	case FilterMedium:
		return task.PriorityMedium
	case FilterLow:
		return task.PriorityLow
	case FilterNone:
		return task.PriorityNone
	default:
		return task.PriorityNone
	}
}

// filterModeName returns a human-readable name for the current filter mode.
func (m *Model) filterModeName() string {
	switch m.filter {
	case FilterNone:
		return "no priority"
	case FilterLow:
		return "low priority"
	case FilterMedium:
		return "medium priority"
	case FilterHigh:
		return "high priority"
	default:
		return "all"
	}
}
