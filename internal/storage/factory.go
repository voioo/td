package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// StorageType represents the type of storage backend.
type StorageType string

const (
	// StorageTypeFile represents file-based storage.
	StorageTypeFile StorageType = "file"
	// StorageTypeMemory represents in-memory storage (for testing).
	StorageTypeMemory StorageType = "memory"
)

// DefaultFactory is the default repository factory.
type DefaultFactory struct{}

// NewDefaultFactory creates a new default factory.
func NewDefaultFactory() *DefaultFactory {
	return &DefaultFactory{}
}

// CreateRepository creates a repository based on the given configuration.
func (f *DefaultFactory) CreateRepository(config map[string]interface{}) (TaskRepository, error) {
	// Extract storage type
	storageTypeRaw, ok := config["type"]
	if !ok {
		// Default to file storage
		storageTypeRaw = StorageTypeFile
	}

	storageType, ok := storageTypeRaw.(StorageType)
	if !ok {
		storageTypeStr, ok := storageTypeRaw.(string)
		if !ok {
			return nil, fmt.Errorf("invalid storage type: %v", storageTypeRaw)
		}
		storageType = StorageType(storageTypeStr)
	}

	switch storageType {
	case StorageTypeFile:
		return f.createFileRepository(config)
	case StorageTypeMemory:
		return f.createMemoryRepository(config)
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", storageType)
	}
}

// createFileRepository creates a file-based repository.
func (f *DefaultFactory) createFileRepository(config map[string]interface{}) (TaskRepository, error) {
	// Extract file path
	filePathRaw, ok := config["file_path"]
	if !ok {
		return nil, fmt.Errorf("file_path is required for file storage")
	}

	filePath, ok := filePathRaw.(string)
	if !ok {
		return nil, fmt.Errorf("file_path must be a string")
	}

	// Expand home directory if needed
	if strings.HasPrefix(filePath, "~") {
		homeDir, err := getHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get home directory: %w", err)
		}
		filePath = filepath.Join(homeDir, filePath[1:])
	}

	return NewRepository(filePath), nil
}

// createMemoryRepository creates an in-memory repository.
func (f *DefaultFactory) createMemoryRepository(_ map[string]interface{}) (TaskRepository, error) {
	return NewMemoryRepository(), nil
}

// getHomeDir gets the user's home directory.
func getHomeDir() (string, error) {
	// This would normally use os.UserHomeDir(), but we'll keep it simple
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir, nil
}
