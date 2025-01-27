package main

import (
	"flag"
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
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help, k.Quit}
}

type filterMode int

const (
	FilterAll filterMode = iota
	FilterNone
	FilterLow
	FilterMedium
	FilterHigh
)

type model struct {
	help              help.Model
	inputStyle        lipgloss.Style
	keys              keyMap
	tasks             []*Task
	doneTasks         []*Task
	newTaskNameInput  input.Model
	editTaskNameInput input.Model
	cursor            int
	mode              int
	latestTaskID      int
	quitting          bool
	filter            filterMode
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
		mode:              ModeNormal,
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

func (m model) normalUpdate(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		tasksToDisplay := m.filteredTasks()
		switch {
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(tasksToDisplay) {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 1 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Right):
			if m.cursor == 0 || m.cursor > len(tasksToDisplay) {
				break
			}
			taskToEdit := tasksToDisplay[m.cursor-1]
			m.mode = ModeEdit
			m.editTaskNameInput.Placeholder = taskToEdit.Name
			m.editTaskNameInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.cursor == 0 || m.cursor > len(tasksToDisplay) {
				break
			}
			taskToComplete := tasksToDisplay[m.cursor-1]

			// Find and move task to done tasks
			for i, task := range m.tasks {
				if task.ID == taskToComplete.ID {
					task.IsDone = true
					m.doneTasks = append(m.doneTasks, task)
					m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
					break
				}
			}

			// Adjust cursor
			if len(m.filteredTasks()) == 0 {
				m.cursor = 0
			} else if m.cursor > len(m.filteredTasks()) {
				m.cursor = len(m.filteredTasks())
			}
		case key.Matches(msg, m.keys.Priority):
			if len(tasksToDisplay) > 0 && m.cursor > 0 && m.cursor <= len(tasksToDisplay) {
				taskToUpdate := tasksToDisplay[m.cursor-1]
				// Find and update priority in original slice
				for _, task := range m.tasks {
					if task.ID == taskToUpdate.ID {
						task.Priority = (task.Priority + 1) % 4
						m.saveTasks()
						m.followTask(task.ID) // Track task after priority change
						break
					}
				}
			}
		case key.Matches(msg, m.keys.Add):
			m.mode = ModeAdditional
			m.newTaskNameInput.Focus()
			return m, nil
		case key.Matches(msg, m.keys.Delete):
			if m.cursor == 0 {
				break
			}

			// Get the actual task from filtered view
			tasksToDisplay := m.filteredTasks()
			if m.cursor > len(tasksToDisplay) {
				break
			}
			taskToDelete := tasksToDisplay[m.cursor-1]

			// Find and delete the task from original slice
			for i, task := range m.tasks {
				if task.ID == taskToDelete.ID {
					m.tasks = append(m.tasks[:i], m.tasks[i+1:]...)
					break
				}
			}

			// Adjust cursor
			if len(m.filteredTasks()) == 0 {
				m.cursor = 0
			} else if m.cursor > len(m.filteredTasks()) {
				m.cursor = len(m.filteredTasks())
			}
		case key.Matches(msg, m.keys.Type):
			if m.mode == ModeDoneTaskList {
				if len(m.tasks) == 0 {
					m.cursor = 0
				} else {
					m.cursor = 1
				}
				m.mode = ModeNormal
				return m, nil
			}
			if len(m.doneTasks) == 0 {
				m.cursor = 0
			} else {
				m.cursor = 1
			}
			m.mode = ModeDoneTaskList
		case key.Matches(msg, m.keys.Quit):
			m.saveTasks()
			return m, tea.Quit
		case key.Matches(msg, m.keys.Filter):
			m.filter = (m.filter + 1) % 5 // 5 is the number of filter modes
			return m, nil
		case key.Matches(msg, m.keys.Help):
			m.mode = ModeHelp
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
			m.doneTasks = append(m.doneTasks[:m.cursor-1], m.doneTasks[m.cursor:]...)
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
			m.mode = ModeNormal
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
		case key.Matches(msg, m.keys.Escape):
			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.mode = ModeNormal
			m.newTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.newTaskNameInput.Value() == "" {
				return m, nil
			}

			m.latestTaskID++
			m.tasks = append(m.tasks, &Task{
				ID:        m.latestTaskID,
				Name:      m.newTaskNameInput.Value(),
				CreatedAt: time.Now(),
			})

			m.cursor++
			m.mode = ModeNormal
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
		case key.Matches(msg, m.keys.Escape):
			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Quit):
			m.mode = ModeNormal
			m.editTaskNameInput.Reset()
			return m, nil
		case key.Matches(msg, m.keys.Enter):
			if m.editTaskNameInput.Value() == "" {
				return m, nil
			}

			tasksToDisplay := m.filteredTasks()
			if m.cursor > 0 && m.cursor <= len(tasksToDisplay) {
				taskToEdit := tasksToDisplay[m.cursor-1]
				// Find and update name in original slice
				for _, task := range m.tasks {
					if task.ID == taskToEdit.ID {
						task.Name = m.editTaskNameInput.Value()
						break
					}
				}
			}

			m.mode = ModeNormal
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			m.mode = ModeNormal
		}
	}

	return m, nil
}

func (m model) View() string {
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

func (m model) normalView() string {
	if m.quitting {
		return "Bye!\n"
	}
	var s strings.Builder
	var title termenv.Style
	var tasksToDisplay []*Task
	switch m.mode {
	case ModeNormal:
		if len(m.tasks) == 0 {
			helpView := m.help.FullHelpView(m.keys.FullHelp())
			return "You have no tasks.\n" + helpView
		}
		title = termenv.String("YOUR TASKS")
		tasksToDisplay = m.filteredTasks()
	case ModeDoneTaskList:
		if len(m.doneTasks) == 0 {
			return "You have no completed tasks.\n"
		}
		title = termenv.String("YOUR COMPLETED TASKS")
		tasksToDisplay = m.doneTasks
	}
	title = title.Bold().Underline()
	s.WriteString(fmt.Sprintf("%v\n\n", title))

	for i, v := range tasksToDisplay {
		cursor := termenv.String(" ")
		if m.cursor == i+1 {
			cursor = termenv.String(">").Foreground(termenv.ANSIYellow)
		}

		taskStr := m.taskView(v, m.cursor == i+1)
		timeLayout := "2006-01-02 15:04"
		s.WriteString(fmt.Sprintf("%v #%d: %s (%s)\n", cursor, v.ID, taskStr, v.CreatedAt.Format(timeLayout)))
	}

	helpView := m.help.FullHelpView(m.keys.FullHelp())
	height := 1

	s.WriteString("\n" + strings.Repeat("\n", height))
	s.WriteString(helpView)

	return s.String()
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
	// Setup flags
	versionFlag := flag.Bool("version", false, "print version information")
	flag.BoolVar(versionFlag, "v", false, "print version information (shorthand)")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("td %s (commit: %s, built at: %s)\n", version, commit, date)
		os.Exit(0)
	}

	p := tea.NewProgram(initializeModel(), tea.WithAltScreen())
	p.Run()
}

func report(err error) {
	fmt.Printf("td: %s\n", err.Error())
	os.Exit(1)
}

func (m model) filteredTasks() []*Task {
	var tasks []*Task

	if m.filter == FilterAll {
		tasks = append(tasks, m.tasks...)
	} else {
		for _, task := range m.tasks {
			if filterToPriority(m.filter) == task.Priority {
				tasks = append(tasks, task)
			}
		}
	}

	SortTasksByPriority(tasks)

	return tasks
}

func (m *model) followTask(taskID int) {
	// Find task's new position in filtered/sorted list
	tasks := m.filteredTasks()
	for i, task := range tasks {
		if task.ID == taskID {
			m.cursor = i + 1
			return
		}
	}
	// If task not found, adjust cursor to valid position
	if len(tasks) == 0 {
		m.cursor = 0
	} else if m.cursor > len(tasks) {
		m.cursor = len(tasks)
	}
}

func filterToPriority(f filterMode) Priority {
	switch f {
	case FilterHigh:
		return PriorityHigh
	case FilterMedium:
		return PriorityMedium
	case FilterLow:
		return PriorityLow
	case FilterNone:
		return PriorityNone
	default:
		return PriorityNone
	}
}
