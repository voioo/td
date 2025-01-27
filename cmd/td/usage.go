package main

import "github.com/charmbracelet/lipgloss"

var (
	sectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#4CAF50"))

	dividerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666"))

	commandStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#90CAF9"))

	descStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	priorityHighStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000"))
	priorityMedStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFF00"))
	priorityLowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	priorityNoneStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
)

func getUsageView() string {
	leftColumn := lipgloss.JoinVertical(lipgloss.Left,
		sectionStyle.Render("Task Management"),
		dividerStyle.Render("───────────────"),
		commandStyle.Render("  • a     ")+descStyle.Render(" add new task"),
		commandStyle.Render("  • d     ")+descStyle.Render(" delete task"),
		commandStyle.Render("  • enter ")+descStyle.Render(" mark done/undone"),
		commandStyle.Render("  • →/l   ")+descStyle.Render(" edit task name"),
		"",
		sectionStyle.Render("Navigation"),
		dividerStyle.Render("──────────"),
		commandStyle.Render("  • ↑/k   ")+descStyle.Render(" move up"),
		commandStyle.Render("  • ↓/j   ")+descStyle.Render(" move down"),
		commandStyle.Render("  • t     ")+descStyle.Render(" toggle tasks view"),
		"",
		sectionStyle.Render("General"),
		dividerStyle.Render("───────"),
		commandStyle.Render("  • ?     ")+descStyle.Render(" show/hide help"),
		commandStyle.Render("  • esc   ")+descStyle.Render(" go back/quit"),
		commandStyle.Render("  • ctrl+c")+descStyle.Render(" quit"),
	)

	rightColumn := lipgloss.JoinVertical(lipgloss.Left,
		sectionStyle.Render("Priority Management"),
		dividerStyle.Render("──────────────────"),
		commandStyle.Render("  • p    ")+descStyle.Render(" set priority"),
		commandStyle.Render("  • f    ")+descStyle.Render(" filter by priority"),
		"",
		sectionStyle.Render("Priority Levels"),
		dividerStyle.Render("──────────────"),
		priorityHighStyle.Render("  • ● ")+descStyle.Render("     high priority"),
		priorityMedStyle.Render("  • ● ")+descStyle.Render("     medium priority"),
		priorityLowStyle.Render("  • ● ")+descStyle.Render("     low priority"),
		priorityNoneStyle.Render("  • ○ ")+descStyle.Render("     no priority"),
		"",
		sectionStyle.Render("Tips"),
		dividerStyle.Render("────"),
		descStyle.Render("  • Tasks auto-sort by priority"),
		descStyle.Render("  • Use filters to focus"),
		descStyle.Render("  • Newest tasks first"),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top,
		leftColumn,
		"     ", // gap between columns
		rightColumn,
	)
}
