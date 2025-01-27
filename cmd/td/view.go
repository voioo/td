package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) taskView(task *Task, selected bool) string {
	var sb strings.Builder

	// Priority indicator with colors
	priorityStyles := map[Priority]lipgloss.Style{
		PriorityHigh:   lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")),
		PriorityMedium: lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00")),
		PriorityLow:    lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")),
	}

	// Priority indicator
	priorityStr := "○"
	if style, ok := priorityStyles[task.Priority]; ok {
		priorityStr = style.Render("●")
	}
	sb.WriteString(priorityStr + " ")

	// Task name with selection
	taskName := task.Name
	if selected {
		taskName = m.inputStyle.Render(taskName)
	}
	sb.WriteString(taskName)

	return sb.String()
}
