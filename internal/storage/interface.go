// Package storage provides data persistence functionality for the td application.
package storage

import "github.com/voioo/td/internal/task"

// TaskRepository defines the interface for task data persistence.
type TaskRepository interface {
	// LoadTasks loads all tasks from the repository.
	// Returns active tasks, completed tasks, and the next available task ID.
	LoadTasks() ([]*task.Task, []*task.Task, int, error)

	// SaveTasks saves all tasks to the repository.
	SaveTasks(activeTasks []*task.Task, completedTasks []*task.Task) error

	// Close closes the repository and releases any resources.
	Close() error
}

// RepositoryFactory creates repositories based on configuration.
type RepositoryFactory interface {
	// CreateRepository creates a repository based on the given configuration.
	CreateRepository(config map[string]interface{}) (TaskRepository, error)
}

// FileRepositoryConfig holds configuration for file-based repositories.
type FileRepositoryConfig struct {
	// FilePath is the path to the data file.
	FilePath string `json:"file_path"`
	// Format specifies the file format (json, yaml, etc.).
	Format string `json:"format,omitempty"`
}
