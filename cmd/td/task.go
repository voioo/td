package main

import (
	"sort"
	"time"
)

type Priority int

const (
	PriorityNone   Priority = iota // 0
	PriorityLow                    // 1
	PriorityMedium                 // 2
	PriorityHigh                   // 3
)

type Task struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	ID        int       `json:"id"`
	IsDone    bool      `json:"is_done"`
	Priority  Priority  `json:"priority"`
}

// Helper method to get priority as string
func (p Priority) String() string {
	switch p {
	case PriorityHigh:
		return "high"
	case PriorityMedium:
		return "medium"
	case PriorityLow:
		return "low"
	default:
		return "none"
	}
}

// Sort tasks by priority (high to low) and then by creation time (newest first)
func SortTasksByPriority(tasks []*Task) {
	sort.Slice(tasks, func(i, j int) bool {
		// First compare priorities
		if tasks[i].Priority != tasks[j].Priority {
			return tasks[i].Priority > tasks[j].Priority
		}
		// If priorities are equal, sort by creation time (newest first)
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
}
