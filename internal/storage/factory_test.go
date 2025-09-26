package storage

import (
	"testing"
)

func TestDefaultFactory(t *testing.T) {
	factory := NewDefaultFactory()

	t.Run("create file repository", func(t *testing.T) {
		config := map[string]interface{}{
			"type":      StorageTypeFile,
			"file_path": "/tmp/test.json",
		}

		repo, err := factory.CreateRepository(config)
		if err != nil {
			t.Errorf("expected no error creating file repository, got %v", err)
		}

		if _, ok := repo.(*FileRepository); !ok {
			t.Errorf("expected *FileRepository, got %T", repo)
		}

		// Test that it implements the interface
		if err := repo.Close(); err != nil {
			t.Errorf("expected no error closing repository, got %v", err)
		}
	})

	t.Run("create memory repository", func(t *testing.T) {
		config := map[string]interface{}{
			"type": StorageTypeMemory,
		}

		repo, err := factory.CreateRepository(config)
		if err != nil {
			t.Errorf("expected no error creating memory repository, got %v", err)
		}

		if _, ok := repo.(*MemoryRepository); !ok {
			t.Errorf("expected *MemoryRepository, got %T", repo)
		}
	})

	t.Run("default to file storage", func(t *testing.T) {
		config := map[string]interface{}{
			"file_path": "/tmp/test.json",
		}

		repo, err := factory.CreateRepository(config)
		if err != nil {
			t.Errorf("expected no error creating repository, got %v", err)
		}

		if _, ok := repo.(*FileRepository); !ok {
			t.Errorf("expected *FileRepository, got %T", repo)
		}
	})

	t.Run("invalid storage type", func(t *testing.T) {
		config := map[string]interface{}{
			"type": "invalid",
		}

		_, err := factory.CreateRepository(config)
		if err == nil {
			t.Error("expected error for invalid storage type")
		}
	})

	t.Run("missing file path for file storage", func(t *testing.T) {
		config := map[string]interface{}{
			"type": StorageTypeFile,
		}

		_, err := factory.CreateRepository(config)
		if err == nil {
			t.Error("expected error for missing file path")
		}
	})
}
