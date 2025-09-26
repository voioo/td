// Package storage provides data persistence functionality for the td application.
package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/voioo/td/internal/logger"
	"github.com/voioo/td/internal/task"
)

var (
	// ErrFileNotFound is returned when the data file doesn't exist.
	ErrFileNotFound = errors.New("data file not found")
	// ErrInvalidData is returned when the data file contains invalid data.
	ErrInvalidData = errors.New("invalid data in file")
	// ErrPermissionDenied is returned when there's no permission to access the data file.
	ErrPermissionDenied = errors.New("permission denied")
)

// FileRepository handles file-based data persistence operations.
type FileRepository struct {
	filePath string
}

// Ensure FileRepository implements the TaskRepository interface.
var _ TaskRepository = (*FileRepository)(nil)

// NewRepository creates a new file repository with the given file path.
func NewRepository(filePath string) *FileRepository {
	return &FileRepository{filePath: filePath}
}

// LoadTasks loads tasks from the repository file.
func (r *FileRepository) LoadTasks() ([]*task.Task, []*task.Task, int, error) {
	logger.Debug("Loading tasks from repository", logger.F("file", r.filePath))

	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			logger.Debug("Data file does not exist, starting with empty repository")
			return []*task.Task{}, []*task.Task{}, 0, nil
		}
		if os.IsPermission(err) {
			logger.Error("Permission denied accessing data file", logger.F("file", r.filePath), logger.F("error", err))
			return nil, nil, 0, fmt.Errorf("%w: %v", ErrPermissionDenied, err)
		}
		logger.Error("Failed to open data file", logger.F("file", r.filePath), logger.F("error", err))
		return nil, nil, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var tasks []*task.Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		logger.Error("Failed to decode JSON data", logger.F("error", err))
		return nil, nil, 0, fmt.Errorf("%w: %v", ErrInvalidData, err)
	}

	// Validate and separate tasks
	activeTasks := []*task.Task{}
	doneTasks := []*task.Task{}
	maxID := 0

	for _, t := range tasks {
		if err := r.validateTask(t); err != nil {
			logger.Error("Invalid task data", logger.F("task_id", t.ID), logger.F("error", err))
			return nil, nil, 0, fmt.Errorf("invalid task data: %w", err)
		}

		if t.ID > maxID {
			maxID = t.ID
		}

		if t.IsDone {
			doneTasks = append(doneTasks, t)
		} else {
			activeTasks = append(activeTasks, t)
		}
	}

	logger.Info("Successfully loaded tasks",
		logger.F("active_tasks", len(activeTasks)),
		logger.F("done_tasks", len(doneTasks)),
		logger.F("next_id", maxID+1))

	return activeTasks, doneTasks, maxID + 1, nil
}

// SaveTasks saves all tasks to the repository file.
func (r *FileRepository) SaveTasks(tasks []*task.Task, doneTasks []*task.Task) error {
	logger.Debug("Saving tasks to repository",
		logger.F("file", r.filePath),
		logger.F("active_tasks", len(tasks)),
		logger.F("done_tasks", len(doneTasks)))

	// Ensure directory exists
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		logger.Error("Failed to create directory", logger.F("dir", dir), logger.F("error", err))
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Combine all tasks for saving
	allTasks := append(tasks, doneTasks...)

	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if os.IsPermission(err) {
			logger.Error("Permission denied writing to data file", logger.F("file", r.filePath), logger.F("error", err))
			return fmt.Errorf("%w: %v", ErrPermissionDenied, err)
		}
		logger.Error("Failed to open file for writing", logger.F("file", r.filePath), logger.F("error", err))
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	// Validate all tasks before saving
	for _, t := range allTasks {
		if err := r.validateTask(t); err != nil {
			logger.Error("Cannot save invalid task", logger.F("task_id", t.ID), logger.F("error", err))
			return fmt.Errorf("cannot save invalid task: %w", err)
		}
	}

	data, err := json.MarshalIndent(allTasks, "", "  ")
	if err != nil {
		logger.Error("Failed to marshal tasks to JSON", logger.F("error", err))
		return fmt.Errorf("failed to marshal tasks: %w", err)
	}

	if _, err := file.Write(data); err != nil {
		logger.Error("Failed to write data to file", logger.F("error", err))
		return fmt.Errorf("failed to write data: %w", err)
	}

	logger.Info("Successfully saved tasks to repository")
	return nil
}

// validateTask checks if a task has valid data.
func (r *FileRepository) validateTask(t *task.Task) error {
	if t == nil {
		return errors.New("task is nil")
	}
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("task name cannot be empty or only whitespace")
	}
	if len(t.Name) > 500 {
		return errors.New("task name is too long (max 500 characters)")
	}
	if t.ID <= 0 {
		return errors.New("task ID must be positive")
	}
	if t.Priority < task.PriorityNone || t.Priority > task.PriorityHigh {
		return errors.New("invalid priority value")
	}
	if t.CreatedAt.IsZero() {
		return errors.New("task must have creation time")
	}
	// Check for reasonable creation time (not too far in the future)
	if t.CreatedAt.After(time.Now().Add(24 * time.Hour)) {
		return errors.New("task creation time cannot be more than 24 hours in the future")
	}
	return nil
}

// DataIntegrity represents data integrity information.
type DataIntegrity struct {
	Version   int       `json:"version"`
	Checksum  string    `json:"checksum"`
	CreatedAt time.Time `json:"created_at"`
}

// saveTasksWithIntegrity saves tasks with integrity checks.
func (r *FileRepository) saveTasksWithIntegrity(tasks []*task.Task, doneTasks []*task.Task) error {
	// Ensure directory exists
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Combine all tasks for saving
	allTasks := append(tasks, doneTasks...)

	// Create integrity data
	integrity := DataIntegrity{
		Version:   1,
		Checksum:  r.calculateChecksum(allTasks),
		CreatedAt: time.Now(),
	}

	// Create wrapper structure
	data := struct {
		Integrity DataIntegrity `json:"integrity"`
		Tasks     []*task.Task  `json:"tasks"`
		DoneTasks []*task.Task  `json:"done_tasks"`
	}{
		Integrity: integrity,
		Tasks:     tasks,
		DoneTasks: doneTasks,
	}

	file, err := os.OpenFile(r.filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		if os.IsPermission(err) {
			return fmt.Errorf("%w: %v", ErrPermissionDenied, err)
		}
		return fmt.Errorf("failed to open file for writing: %w", err)
	}
	defer file.Close()

	encodedData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	if _, err := file.Write(encodedData); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	return nil
}

// loadTasksWithIntegrity loads tasks with integrity verification.
func (r *FileRepository) loadTasksWithIntegrity() ([]*task.Task, []*task.Task, int, error) {
	file, err := os.Open(r.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []*task.Task{}, []*task.Task{}, 0, nil
		}
		if os.IsPermission(err) {
			return nil, nil, 0, fmt.Errorf("%w: %v", ErrPermissionDenied, err)
		}
		return nil, nil, 0, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var data struct {
		Integrity DataIntegrity `json:"integrity"`
		Tasks     []*task.Task  `json:"tasks"`
		DoneTasks []*task.Task  `json:"done_tasks"`
	}

	if err := json.NewDecoder(file).Decode(&data); err != nil {
		// Try loading old format for backward compatibility
		return r.loadLegacyFormat(file)
	}

	// Verify integrity
	allTasks := append(data.Tasks, data.DoneTasks...)
	expectedChecksum := r.calculateChecksum(allTasks)
	if data.Integrity.Checksum != expectedChecksum {
		return nil, nil, 0, fmt.Errorf("data integrity check failed: checksum mismatch")
	}

	// Validate all tasks
	maxID := 0
	for _, t := range allTasks {
		if err := r.validateTask(t); err != nil {
			return nil, nil, 0, fmt.Errorf("invalid task data: %w", err)
		}
		if t.ID > maxID {
			maxID = t.ID
		}
	}

	return data.Tasks, data.DoneTasks, maxID, nil
}

// loadLegacyFormat loads the old JSON format for backward compatibility.
func (r *FileRepository) loadLegacyFormat(file *os.File) ([]*task.Task, []*task.Task, int, error) {
	// Reset file position
	if _, err := file.Seek(0, 0); err != nil {
		return nil, nil, 0, err
	}

	var tasks []*task.Task
	if err := json.NewDecoder(file).Decode(&tasks); err != nil {
		return nil, nil, 0, fmt.Errorf("%w: %v", ErrInvalidData, err)
	}

	// Separate tasks and validate
	activeTasks := []*task.Task{}
	doneTasks := []*task.Task{}
	maxID := 0

	for _, t := range tasks {
		if err := r.validateTask(t); err != nil {
			return nil, nil, 0, fmt.Errorf("invalid task data: %w", err)
		}

		if t.ID > maxID {
			maxID = t.ID
		}

		if t.IsDone {
			doneTasks = append(doneTasks, t)
		} else {
			activeTasks = append(activeTasks, t)
		}
	}

	return activeTasks, doneTasks, maxID, nil
}

// calculateChecksum calculates a simple checksum for data integrity.
func (r *FileRepository) calculateChecksum(tasks []*task.Task) string {
	// Simple checksum based on task count and IDs
	sum := len(tasks)
	for _, t := range tasks {
		sum += t.ID
		sum += int(t.Priority)
		sum += len(t.Name)
	}
	return fmt.Sprintf("%x", sum)
}

// Close closes the repository. For file-based repositories, this is a no-op.
func (r *FileRepository) Close() error {
	// File-based repository doesn't need explicit closing
	return nil
}

// GetDefaultRepositoryPath returns the default repository file path.
func GetDefaultRepositoryPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not available
		return ".td.json"
	}
	return filepath.Join(homeDir, ".td.json")
}
