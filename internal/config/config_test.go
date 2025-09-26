package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.DataFile == "" {
		t.Error("expected default config to have data file")
	}
	if cfg.Theme.PrimaryColor != DefaultPrimaryColor {
		t.Errorf("expected primary color to be %s, got %s", DefaultPrimaryColor, cfg.Theme.PrimaryColor)
	}
	if cfg.KeyMap.Add != "a" {
		t.Errorf("expected add key to be 'a', got %s", cfg.KeyMap.Add)
	}
}

func TestLoadConfig(t *testing.T) {
	t.Run("load default config when no file exists", func(t *testing.T) {
		cfg, err := LoadConfig("")

		if err != nil {
			t.Errorf("expected no error loading default config, got %v", err)
		}
		if cfg == nil {
			t.Error("expected config to not be nil")
		}
		if cfg.Theme.PrimaryColor != DefaultPrimaryColor {
			t.Errorf("expected default primary color, got %s", cfg.Theme.PrimaryColor)
		}
	})

	t.Run("load config from JSON file", func(t *testing.T) {
		// Create a temporary config file
		tempDir, err := os.MkdirTemp("", "td-config-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.json")

		// Create test config
		testConfig := &Config{
			DataFile: "/tmp/test.json",
			Theme: Theme{
				PrimaryColor:      "#123456",
				HighPriorityColor: "#FF0000",
			},
			KeyMap: KeyMap{
				Add: "b",
			},
		}

		err = testConfig.SaveConfig(configFile)
		if err != nil {
			t.Fatal(err)
		}

		// Load config
		loadedConfig, err := LoadConfig(configFile)
		if err != nil {
			t.Errorf("expected no error loading config, got %v", err)
		}

		if loadedConfig.DataFile != "/tmp/test.json" {
			t.Errorf("expected data file to be '/tmp/test.json', got %s", loadedConfig.DataFile)
		}
		if loadedConfig.Theme.PrimaryColor != "#123456" {
			t.Errorf("expected primary color to be '#123456', got %s", loadedConfig.Theme.PrimaryColor)
		}
		if loadedConfig.KeyMap.Add != "b" {
			t.Errorf("expected add key to be 'b', got %s", loadedConfig.KeyMap.Add)
		}
	})

	t.Run("load config from YAML file", func(t *testing.T) {
		// Create a temporary config file
		tempDir, err := os.MkdirTemp("", "td-config-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.yaml")

		// Create test config
		testConfig := &Config{
			DataFile: "/tmp/test.yaml",
			Theme: Theme{
				PrimaryColor:      "#654321",
				HighPriorityColor: "#00FF00",
			},
			KeyMap: KeyMap{
				Add: "c",
			},
		}

		err = testConfig.SaveConfig(configFile)
		if err != nil {
			t.Fatal(err)
		}

		// Load config
		loadedConfig, err := LoadConfig(configFile)
		if err != nil {
			t.Errorf("expected no error loading YAML config, got %v", err)
		}

		if loadedConfig.DataFile != "/tmp/test.yaml" {
			t.Errorf("expected data file to be '/tmp/test.yaml', got %s", loadedConfig.DataFile)
		}
		if loadedConfig.Theme.PrimaryColor != "#654321" {
			t.Errorf("expected primary color to be '#654321', got %s", loadedConfig.Theme.PrimaryColor)
		}
		if loadedConfig.KeyMap.Add != "c" {
			t.Errorf("expected add key to be 'c', got %s", loadedConfig.KeyMap.Add)
		}
	})

	t.Run("merge with defaults for missing fields", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "td-config-test")
		if err != nil {
			t.Fatal(err)
		}
		defer os.RemoveAll(tempDir)

		configFile := filepath.Join(tempDir, "config.json")

		// Create partial config
		partialConfig := `{"theme": {"primary_color": "#123456"}}`
		err = os.WriteFile(configFile, []byte(partialConfig), 0644)
		if err != nil {
			t.Fatal(err)
		}

		loadedConfig, err := LoadConfig(configFile)
		if err != nil {
			t.Errorf("expected no error loading partial config, got %v", err)
		}

		// Should have custom primary color
		if loadedConfig.Theme.PrimaryColor != "#123456" {
			t.Errorf("expected primary color to be '#123456', got %s", loadedConfig.Theme.PrimaryColor)
		}

		// Should have default values for missing fields
		if loadedConfig.Theme.HighPriorityColor != DefaultHighPriorityColor {
			t.Errorf("expected default high priority color, got %s", loadedConfig.Theme.HighPriorityColor)
		}
		if loadedConfig.KeyMap.Add != "a" {
			t.Errorf("expected default add key, got %s", loadedConfig.KeyMap.Add)
		}
	})
}

func TestSaveConfig(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "td-config-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	configFile := filepath.Join(tempDir, "config.json")

	cfg := DefaultConfig()
	cfg.Theme.PrimaryColor = "#123456"

	err = cfg.SaveConfig(configFile)
	if err != nil {
		t.Errorf("expected no error saving config, got %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Error("expected config file to be created")
	}
}
