package storage

import (
	"sync"

	"github.com/voioo/td/internal/task"
)

// MemoryRepository is an in-memory implementation of the TaskRepository interface.
// Useful for testing and temporary storage.
type MemoryRepository struct {
	mu          sync.RWMutex
	activeTasks []*task.Task
	doneTasks   []*task.Task
	nextID      int
}

// NewMemoryRepository creates a new in-memory repository.
func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		activeTasks: make([]*task.Task, 0),
		doneTasks:   make([]*task.Task, 0),
		nextID:      0,
	}
}

// Ensure MemoryRepository implements the TaskRepository interface.
var _ TaskRepository = (*MemoryRepository)(nil)

// LoadTasks returns a copy of the current tasks.
func (r *MemoryRepository) LoadTasks() ([]*task.Task, []*task.Task, int, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Return copies to prevent external modification
	activeTasks := make([]*task.Task, len(r.activeTasks))
	copy(activeTasks, r.activeTasks)

	doneTasks := make([]*task.Task, len(r.doneTasks))
	copy(doneTasks, r.doneTasks)

	return activeTasks, doneTasks, r.nextID, nil
}

// SaveTasks saves the provided tasks to memory.
func (r *MemoryRepository) SaveTasks(activeTasks []*task.Task, doneTasks []*task.Task) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Create copies to prevent external modification
	r.activeTasks = make([]*task.Task, len(activeTasks))
	copy(r.activeTasks, activeTasks)

	r.doneTasks = make([]*task.Task, len(doneTasks))
	copy(r.doneTasks, doneTasks)

	// Calculate next ID
	maxID := 0
	for _, task := range activeTasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	for _, task := range doneTasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	r.nextID = maxID + 1

	return nil
}

// Close is a no-op for memory repositories.
func (r *MemoryRepository) Close() error {
	return nil
}
