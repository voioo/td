// Package task provides task management functionality for the td application.
package task

import (
	"sort"
	"time"
)

// Priority represents the priority level of a task.
type Priority int

const (
	// PriorityNone indicates no priority assigned to the task.
	PriorityNone Priority = iota
	// PriorityLow indicates low priority.
	PriorityLow
	// PriorityMedium indicates medium priority.
	PriorityMedium
	// PriorityHigh indicates high priority.
	PriorityHigh
)

// String returns the string representation of the priority.
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

// Task represents a todo task with all its properties.
type Task struct {
	CreatedAt time.Time `json:"created_at"`
	Name      string    `json:"name"`
	ID        int       `json:"id"`
	IsDone    bool      `json:"is_done"`
	Priority  Priority  `json:"priority"`
}

// TaskManager manages a collection of tasks and provides business logic operations.
type TaskManager struct {
	tasks     []*Task
	doneTasks []*Task
	nextID    int
}

// NewTaskManager creates a new task manager with the given tasks.
func NewTaskManager(tasks []*Task, doneTasks []*Task, nextID int) *TaskManager {
	tm := &TaskManager{
		tasks:     make([]*Task, len(tasks)),
		doneTasks: make([]*Task, len(doneTasks)),
		nextID:    nextID,
	}
	copy(tm.tasks, tasks)
	copy(tm.doneTasks, doneTasks)
	return tm
}

// GetTasks returns a copy of the active tasks.
func (tm *TaskManager) GetTasks() []*Task {
	tasks := make([]*Task, len(tm.tasks))
	copy(tasks, tm.tasks)
	return tasks
}

// GetDoneTasks returns a copy of the completed tasks.
func (tm *TaskManager) GetDoneTasks() []*Task {
	doneTasks := make([]*Task, len(tm.doneTasks))
	copy(doneTasks, tm.doneTasks)
	return doneTasks
}

// GetNextID returns the next available task ID.
func (tm *TaskManager) GetNextID() int {
	return tm.nextID
}

// AddTask adds a new task with the given name.
func (tm *TaskManager) AddTask(name string) *Task {
	tm.nextID++
	task := &Task{
		ID:        tm.nextID,
		Name:      name,
		CreatedAt: time.Now(),
		IsDone:    false,
		Priority:  PriorityNone,
	}
	tm.tasks = append(tm.tasks, task)
	tm.sortTasks()
	return task
}

// DeleteTask removes the task with the given ID.
func (tm *TaskManager) DeleteTask(id int) *Task {
	for i, task := range tm.tasks {
		if task.ID == id {
			deleted := task
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			return deleted
		}
	}
	for i, task := range tm.doneTasks {
		if task.ID == id {
			deleted := task
			tm.doneTasks = append(tm.doneTasks[:i], tm.doneTasks[i+1:]...)
			return deleted
		}
	}
	return nil
}

// CompleteTask marks the task with the given ID as completed.
func (tm *TaskManager) CompleteTask(id int) *Task {
	for i, task := range tm.tasks {
		if task.ID == id {
			task.IsDone = true
			tm.doneTasks = append(tm.doneTasks, task)
			tm.tasks = append(tm.tasks[:i], tm.tasks[i+1:]...)
			tm.sortTasks()
			return task
		}
	}
	return nil
}

// UncompleteTask marks the completed task with the given ID as active.
func (tm *TaskManager) UncompleteTask(id int) *Task {
	for i, task := range tm.doneTasks {
		if task.ID == id {
			task.IsDone = false
			tm.tasks = append(tm.tasks, task)
			tm.doneTasks = append(tm.doneTasks[:i], tm.doneTasks[i+1:]...)
			tm.sortTasks()
			return task
		}
	}
	return nil
}

// UpdateTaskName changes the name of the task with the given ID.
func (tm *TaskManager) UpdateTaskName(id int, newName string) *Task {
	for _, task := range tm.tasks {
		if task.ID == id {
			task.Name = newName
			return task
		}
	}
	for _, task := range tm.doneTasks {
		if task.ID == id {
			task.Name = newName
			return task
		}
	}
	return nil
}

// SetTaskPriority sets the priority of the task with the given ID.
func (tm *TaskManager) SetTaskPriority(id int, priority Priority) *Task {
	for _, task := range tm.tasks {
		if task.ID == id {
			task.Priority = priority
			tm.sortTasks()
			return task
		}
	}
	return nil
}

// FindTaskByID finds a task by its ID in both active and completed tasks.
func (tm *TaskManager) FindTaskByID(id int) *Task {
	for _, task := range tm.tasks {
		if task.ID == id {
			return task
		}
	}
	for _, task := range tm.doneTasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

// sortTasks sorts tasks by priority (high to low) and then by creation time (newest first).
func (tm *TaskManager) sortTasks() {
	sortTasks(tm.tasks)
	sortTasks(tm.doneTasks)
}

// SortTasksByPriority sorts the given tasks by priority and creation time.
func SortTasksByPriority(tasks []*Task) {
	sortTasks(tasks)
}

// sortTasks is the internal sorting function.
func sortTasks(tasks []*Task) {
	sort.Slice(tasks, func(i, j int) bool {
		// First compare priorities (higher priority first)
		if tasks[i].Priority != tasks[j].Priority {
			return tasks[i].Priority > tasks[j].Priority
		}
		// If priorities are equal, sort by creation time (newest first)
		return tasks[i].CreatedAt.After(tasks[j].CreatedAt)
	})
}
