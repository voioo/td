package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/voioo/td/internal/config"
	"github.com/voioo/td/internal/task"
)

// normalView renders the normal view (active tasks or completed tasks).
func (m *Model) normalView() string {
	if m.quitting {
		return "Bye!\n"
	}

	var s strings.Builder
	var title termenv.Style
	var tasksToDisplay []*task.Task

	switch m.mode {
	case ModeNormal:
		if len(m.taskManager.GetTasks()) == 0 {
			helpView := m.help.FullHelpView(m.keys.FullHelp())
			return "You have no tasks.\n" + helpView
		}
		title = termenv.String("YOUR TASKS")
		m.updateTaskCache()
		tasksToDisplay = m.taskCache
	case ModeDoneTaskList:
		doneTasks := m.taskManager.GetDoneTasks()
		if len(doneTasks) == 0 {
			return "You have no completed tasks.\n"
		}
		title = termenv.String("YOUR COMPLETED TASKS")
		tasksToDisplay = doneTasks
	}

	title = title.Bold().Underline()
	s.WriteString(fmt.Sprintf("%v\n\n", title))

	for i, task := range tasksToDisplay {
		cursor := termenv.String(" ")
		if m.cursor == i+1 {
			cursor = termenv.String(">").Foreground(termenv.ANSIYellow)
		}

		taskStr := m.taskView(task, m.cursor == i+1)
		timeLayout := "2006-01-02 15:04"
		s.WriteString(fmt.Sprintf("%v #%d: %s (%s)\n", cursor, task.ID, taskStr, task.CreatedAt.Format(timeLayout)))
	}

	helpView := m.help.FullHelpView(m.keys.FullHelp())
	height := 1
	s.WriteString("\n" + strings.Repeat("\n", height))
	s.WriteString(helpView)

	return s.String()
}

// addingTaskView renders the task adding view.
func (m *Model) addingTaskView() string {
	title := termenv.String("Additional Mode").Bold().Underline()
	return fmt.Sprintf("%v\n\nInput the new task name\n\n%s\n", title, m.newTaskNameInput.View())
}

// editTaskView renders the task editing view.
func (m *Model) editTaskView() string {
	title := termenv.String("Edit Mode").Bold().Underline()
	return fmt.Sprintf("%v\n\nInput the new task name\n\n%s\n", title, m.editTaskNameInput.View())
}

// helpView renders the help view.
func (m *Model) helpView() string {
	title := termenv.String("USAGE").Bold().Underline()
	return fmt.Sprintf("%v\n\n%s", title, getUsageView(m.config))
}

// taskView renders a single task with priority indicators and colors.
func (m *Model) taskView(t *task.Task, selected bool) string {
	var sb strings.Builder

	// Priority indicator with colors
	priorityStyles := map[task.Priority]lipgloss.Style{
		task.PriorityHigh:   lipgloss.NewStyle().Foreground(lipgloss.Color(m.config.Theme.HighPriorityColor)),
		task.PriorityMedium: lipgloss.NewStyle().Foreground(lipgloss.Color(m.config.Theme.MediumPriorityColor)),
		task.PriorityLow:    lipgloss.NewStyle().Foreground(lipgloss.Color(m.config.Theme.LowPriorityColor)),
	}

	// Priority indicator
	priorityStr := "○"
	if style, ok := priorityStyles[t.Priority]; ok {
		priorityStr = style.Render("●")
	}
	sb.WriteString(priorityStr + " ")

	// Task name with selection
	taskName := t.Name
	if selected {
		taskName = m.inputStyle.Render(taskName)
	}
	sb.WriteString(taskName)

	return sb.String()
}

// getUsageView returns the usage help text.
func getUsageView(config *config.Config) string {
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50")).
			Render("Task Management"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("───────────────"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Add+"     ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" add new task"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Delete+"     ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" delete task"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Enter+" ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" mark done/undone"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • →/l   ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" edit task name"),
		"",
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50")).
			Render("Navigation"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("──────────"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • ↑/k   ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" move up"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • ↓/j   ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" move down"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.ListType+"     ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" toggle tasks view"),
		"",
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50")).
			Render("General"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("───────"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Help+"     ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" show/hide help"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Quit+"   ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" quit"),
	)

	rightColumn := lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50")).
			Render("Priority Management"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("──────────────────"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Priority+"    ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" cycle priority"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • 1-4   ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" set priority directly"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • "+config.KeyMap.Filter+"    ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" filter by priority"),
		"",
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50")).
			Render("Priority Levels"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("──────────────"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(config.Theme.HighPriorityColor)).Render("  • ● ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render("4    high priority"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(config.Theme.MediumPriorityColor)).Render("  • ● ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render("3    medium priority"),
		lipgloss.NewStyle().Foreground(lipgloss.Color(config.Theme.LowPriorityColor)).Render("  • ● ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render("2    low priority"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("#666666")).Render("  • ○ ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render("1    no priority"),
		"",
		lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50")).
			Render("Navigation"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")).
			Render("──────────"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • home/g ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" go to top"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • end/G  ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" go to bottom"),
		lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9")).
			Render("  • C      ")+lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Render(" clear completed"),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftColumn,
		"     ", // gap between columns
		rightColumn,
	)
}
