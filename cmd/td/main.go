package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	input "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
	todo "github.com/voioo/td"
)

type keyMap struct {
	Add    key.Binding
	Up     key.Binding
	Down   key.Binding
	Delete key.Binding
	Left   key.Binding
	Right  key.Binding
	Enter  key.Binding
	Type   key.Binding
	Help   key.Binding
	Quit   key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

var keys = keyMap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "add new task"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete task"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "save"),
	),
	Up: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "move down"),
	),
	Left: key.NewBinding(
		key.WithKeys("left", "h"),
		key.WithHelp("←/h", "move left"),
	),
	Right: key.NewBinding(
		key.WithKeys("right", "l"),
		key.WithHelp("→/l", "move right"),
	),
	Type: key.NewBinding(
		key.WithKeys("t"),
		key.WithHelp("t", "list type"),
	),
	Help: key.NewBinding(
		key.WithKeys("?", "?"),
		key.WithHelp("?", "toggle usage"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
}

// FullHelp returns keybindings for the expanded help view. It's part of the
// key.Map interface.
func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Delete, k.Up, k.Down, k.Left, k.Right}, // first column
		{k.Help, k.Quit}, // second column
	}
}

const (
	normalMode = iota
	doneTaskListMode
	additionalMode
	editMode
	helpMode

	usage = `

--Normal Mode--

j - move cursor one line down
k - move cursor one line up
a - add a new task(move to additional mode)
d - remove a task
e - edit the task name(mode to edit mode)
h - help(switch to help mode)
x, enter - mark as done
t - switch to done tasks list mode
q - save tasks and close this app

--Done Tasks List Mode--

j - move cursor one line down
k - move cursor one line up
d - remove a task
t - switch to normal mode
x, enter - mark as not done
q - save tasks and close this app

--Additional Mode--

ctrl+q - switch to normal mode
enter - submit

--Edit Mode--

left - go back
ctrl+q - switch to normal mode
enter - submit

--Help Mode--
left - back
q - switch to normal mode
`
)

type model struct {
	cursor       int
	mode         int
	latestTaskID int
	tasks        []*todo.Task
	doneTasks    []*todo.Task
	keys         keyMap
	help         help.Model
	inputStyle   lipgloss.Style
	quitting     bool

	newTaskNameInput  input.Model
	editTaskNameInput input.Model
}

func initializeModel() tea.Model {
	tasks, doneTasks, ltID := loadTasksFromRepositoryFile()

	cursor := 0
	if len(tasks) != 0 {
		cursor = 1
	}

	newTaskNameModel := input.New()
	newTaskNameModel.Placeholder = "New task name..."
	newTaskNameModel.Focus()
	editTaskNameModel := input.New()
	editTaskNameModel.Focus()

	return model{
		cursor:            cursor,
		mode:              normalMode,
		latestTaskID:      ltID,
		tasks:             tasks,
		doneTasks:         doneTasks,
		newTaskNameInput:  newTaskNameModel,
		editTaskNameInput: editTaskNameModel,
		keys:              keys,
		help:              help.New(),
		inputStyle:        lipgloss.NewStyle().Foreground(lipgloss.Color("#FF75B7")),
	}
}

func (model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case normalMode:
		return m.normalUpdate(msg)
	case doneTaskListMode:
		return m.doneTaskListUpdate(msg)
	case additionalMode:
		return m.addingTaskUpdate(msg)
	case editMode:
		return m.editTaskUpdate(msg)
	case helpMode:
		return m.helpUpdate(msg)
	}

	return m, nil
}

func (m model) normalUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.tasks) {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 1 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Add):
			m.mode = additionalMode
			m.newTaskNameInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.Delete):
			if m.cursor == 0 {
				break
			}
			m.tasks = append(m.tasks[:m.cursor-1], m.tasks[m.cursor:]...)
			if len(m.tasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
		case key.Matches(msg, m.keys.Right):
			if m.cursor == 0 {
				break
			}
			m.mode = editMode
			m.editTaskNameInput.Placeholder = m.tasks[m.cursor-1].Name
			m.editTaskNameInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.mode = helpMode
		case key.Matches(msg, m.keys.Enter):
			if m.cursor == 0 {
				break
			}

			t := m.tasks[m.cursor-1]
			t.IsDone = true
			m.doneTasks = append(m.doneTasks, t)
			m.tasks = append(m.tasks[:m.cursor-1], m.tasks[m.cursor:]...)

			if len(m.tasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
		case key.Matches(msg, m.keys.Type):
			if m.mode == doneTaskListMode {
				if len(m.tasks) == 0 {
					m.cursor = 0
				} else {
					m.cursor = 1
				}
				m.mode = normalMode
				return m, nil
			}
			if len(m.doneTasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
			m.mode = doneTaskListMode
		case key.Matches(msg, m.keys.Quit):
			m.saveTasks()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) doneTaskListUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.doneTasks) {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 1 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Delete):
			if m.cursor == 0 {
				break
			}
			m.doneTasks = append(m.tasks[:m.cursor-1], m.tasks[m.cursor:]...)
			if len(m.doneTasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
		case key.Matches(msg, m.keys.Type):
			if len(m.tasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
			m.mode = normalMode
		case key.Matches(msg, m.keys.Enter):
			if m.cursor == 0 {
				break
			}
			t := m.doneTasks[m.cursor-1]
			t.IsDone = false
			m.tasks = append(m.tasks, t)
			m.doneTasks = append(m.doneTasks[:m.cursor-1], m.doneTasks[m.cursor:]...)
			if len(m.doneTasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
		case key.Matches(msg, m.keys.Quit):
			m.saveTasks()
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) addingTaskUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			m.mode = normalMode
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.mode = normalMode
			m.newTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.newTaskNameInput.Value() == "" {
				return m, nil
			}

			m.latestTaskID++
			m.tasks = append(m.tasks, &todo.Task{
				ID:        m.latestTaskID,
				Name:      m.newTaskNameInput.Value(),
				CreatedAt: time.Now(),
			})

			m.cursor++
			m.mode = normalMode
			m.newTaskNameInput.Reset()
			return m, nil
		}
	}

	m.newTaskNameInput, cmd = m.newTaskNameInput.Update(msg)

	return m, cmd
}

func (m model) editTaskUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keys.Left):
			m.mode = normalMode
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.mode = normalMode
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.editTaskNameInput.Value() == "" {
				return m, nil
			}
			m.tasks[m.cursor-1].Name = m.editTaskNameInput.Value()

			m.mode = normalMode
			m.editTaskNameInput.Reset()
			return m, nil
		}
	}

	m.editTaskNameInput, cmd = m.editTaskNameInput.Update(msg)

	return m, cmd
}

func (m model) helpUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.saveTasks()
			return m, tea.Quit
		case "q":
			m.mode = normalMode
		}
	}

	return m, nil
}

func (m model) View() string {
	switch m.mode {
	case normalMode, doneTaskListMode:
		return m.normalView()
	case additionalMode:
		return m.addingTaskView()
	case editMode:
		return m.editTaskView()
	case helpMode:
		return m.helpView()
	}

	return ""
}

func (m model) normalView() string {
	if m.quitting {
		return "Bye!\n"
	}
	var s string
	var title termenv.Style
	var tasksToDisplay []*todo.Task
	switch m.mode {
	case normalMode:
		if len(m.tasks) == 0 {
			helpView := m.help.FullHelpView(m.keys.FullHelp())
			return "You have no tasks.\n" + s + helpView
		}
		title = termenv.String("YOUR TASKS")
		tasksToDisplay = m.tasks
	case doneTaskListMode:
		if len(m.doneTasks) == 0 {
			return "You have no completed tasks.\n"
		}
		title = termenv.String("YOUR COMPLETED TASKS")
		tasksToDisplay = m.doneTasks
	}
	title = title.Bold().Underline()
	s = fmt.Sprintf("%v\n\n", title)

	for i, v := range tasksToDisplay {
		cursor := termenv.String(" ")
		if m.cursor == i+1 {
			cursor = termenv.String(">").Foreground(termenv.ANSIYellow)
		}
		taskName := termenv.String(v.Name)
		taskName = taskName.Bold()
		timeLayout := "2006-01-02 15:04"

		s += fmt.Sprintf("%v #%d: %s (%s)\n", cursor, v.ID, taskName, v.CreatedAt.Format(timeLayout))
	}

	helpView := m.help.FullHelpView(m.keys.FullHelp())
	height := 1

	return "\n" + s + strings.Repeat("\n", height) + helpView
}

func (m model) addingTaskView() string {
	title := termenv.String("Additional Mode").Bold().Underline()
	return fmt.Sprintf("%v\n\nInput the new task name\n\n%s\n", title, m.newTaskNameInput.View())
}

func (m model) editTaskView() string {
	title := termenv.String("Edit Mode").Bold().Underline()
	return fmt.Sprintf("%v\n\nInput the new task name\n\n%s\n", title, m.editTaskNameInput.View())
}

func (m model) helpView() string {
	title := termenv.String("USAGE").Bold().Underline()
	return fmt.Sprintf("%v"+usage, title)
}

func main() {
	p := tea.NewProgram(initializeModel(), tea.WithAltScreen())
	p.Run()
}

func report(err error) {
	fmt.Printf("td: %s\n", err.Error())
	os.Exit(1)
}
