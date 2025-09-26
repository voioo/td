// Package config provides configuration management for the td application.
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/voioo/td/internal/storage"
	"gopkg.in/yaml.v3"
)

// Default theme colors
const (
	DefaultPrimaryColor        = "#FF75B7"
	DefaultHighPriorityColor   = "#FF0000"
	DefaultMediumPriorityColor = "#FFFF00"
	DefaultLowPriorityColor    = "#00FF00"
)

// Config holds all configuration options for the td application.
type Config struct {
	// DataFile is the path to the data file.
	DataFile string `json:"data_file"`
	// Theme controls the UI appearance.
	Theme Theme `json:"theme"`
	// KeyMap defines keyboard shortcuts.
	KeyMap KeyMap `json:"keymap"`
}

// Theme defines the visual appearance settings.
type Theme struct {
	// PrimaryColor is the main accent color.
	PrimaryColor string `json:"primary_color"`
	// HighPriorityColor for high priority tasks.
	HighPriorityColor string `json:"high_priority_color"`
	// MediumPriorityColor for medium priority tasks.
	MediumPriorityColor string `json:"medium_priority_color"`
	// LowPriorityColor for low priority tasks.
	LowPriorityColor string `json:"low_priority_color"`
}

// KeyMap defines keyboard shortcuts.
type KeyMap struct {
	Add      string `json:"add"`
	Delete   string `json:"delete"`
	Enter    string `json:"enter"`
	Escape   string `json:"escape"`
	Up       string `json:"up"`
	Down     string `json:"down"`
	Left     string `json:"left"`
	Right    string `json:"right"`
	ListType string `json:"list_type"`
	Help     string `json:"help"`
	Quit     string `json:"quit"`
	Priority string `json:"priority"`
	Filter   string `json:"filter"`
	Undo     string `json:"undo"`
	Redo     string `json:"redo"`
}

// DefaultConfig returns a configuration with default values.
func DefaultConfig() *Config {
	return &Config{
		DataFile: storage.GetDefaultRepositoryPath(),
		Theme: Theme{
			PrimaryColor:        DefaultPrimaryColor,
			HighPriorityColor:   DefaultHighPriorityColor,
			MediumPriorityColor: DefaultMediumPriorityColor,
			LowPriorityColor:    DefaultLowPriorityColor,
		},
		KeyMap: KeyMap{
			Add:      "a",
			Delete:   "d",
			Enter:    "enter",
			Escape:   "esc",
			Up:       "up",
			Down:     "down",
			Left:     "left",
			Right:    "right",
			ListType: "t",
			Help:     "?",
			Quit:     "q",
			Priority: "p",
			Filter:   "f",
			Undo:     "ctrl+u",
			Redo:     "ctrl+r",
		},
	}
}

// ConfigFormat represents the format of a configuration file.
type ConfigFormat int

const (
	// ConfigFormatJSON represents JSON format.
	ConfigFormatJSON ConfigFormat = iota
	// ConfigFormatYAML represents YAML format.
	ConfigFormatYAML
)

// detectConfigFormat detects the format of a configuration file based on its extension.
func detectConfigFormat(filename string) ConfigFormat {
	if strings.HasSuffix(strings.ToLower(filename), ".yaml") || strings.HasSuffix(strings.ToLower(filename), ".yml") {
		return ConfigFormatYAML
	}
	return ConfigFormatJSON
}

// LoadConfig loads configuration from the specified file path.
// If the file doesn't exist, it returns the default configuration.
// Supports both JSON and YAML formats based on file extension.
func LoadConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return DefaultConfig(), nil
	}

	file, err := os.Open(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var config Config
	format := detectConfigFormat(configPath)

	if format == ConfigFormatYAML {
		if err := yaml.NewDecoder(file).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to parse YAML config file: %w", err)
		}
	} else {
		if err := json.NewDecoder(file).Decode(&config); err != nil {
			return nil, fmt.Errorf("failed to parse JSON config file: %w", err)
		}
	}

	// Merge with defaults for any missing fields
	defaults := DefaultConfig()
	mergeWithDefaults(&config, defaults)

	return &config, nil
}

// mergeWithDefaults merges missing fields from defaults into the config.
func mergeWithDefaults(config, defaults *Config) {
	if config.DataFile == "" {
		config.DataFile = defaults.DataFile
	}
	if config.Theme.PrimaryColor == "" {
		config.Theme.PrimaryColor = defaults.Theme.PrimaryColor
	}
	if config.Theme.HighPriorityColor == "" {
		config.Theme.HighPriorityColor = defaults.Theme.HighPriorityColor
	}
	if config.Theme.MediumPriorityColor == "" {
		config.Theme.MediumPriorityColor = defaults.Theme.MediumPriorityColor
	}
	if config.Theme.LowPriorityColor == "" {
		config.Theme.LowPriorityColor = defaults.Theme.LowPriorityColor
	}

	// Fill in missing keymap entries
	if config.KeyMap.Add == "" {
		config.KeyMap.Add = defaults.KeyMap.Add
	}
	if config.KeyMap.Delete == "" {
		config.KeyMap.Delete = defaults.KeyMap.Delete
	}
	if config.KeyMap.Enter == "" {
		config.KeyMap.Enter = defaults.KeyMap.Enter
	}
	if config.KeyMap.Escape == "" {
		config.KeyMap.Escape = defaults.KeyMap.Escape
	}
	if config.KeyMap.Up == "" {
		config.KeyMap.Up = defaults.KeyMap.Up
	}
	if config.KeyMap.Down == "" {
		config.KeyMap.Down = defaults.KeyMap.Down
	}
	if config.KeyMap.Left == "" {
		config.KeyMap.Left = defaults.KeyMap.Left
	}
	if config.KeyMap.Right == "" {
		config.KeyMap.Right = defaults.KeyMap.Right
	}
	if config.KeyMap.ListType == "" {
		config.KeyMap.ListType = defaults.KeyMap.ListType
	}
	if config.KeyMap.Help == "" {
		config.KeyMap.Help = defaults.KeyMap.Help
	}
	if config.KeyMap.Quit == "" {
		config.KeyMap.Quit = defaults.KeyMap.Quit
	}
	if config.KeyMap.Priority == "" {
		config.KeyMap.Priority = defaults.KeyMap.Priority
	}
	if config.KeyMap.Filter == "" {
		config.KeyMap.Filter = defaults.KeyMap.Filter
	}
	if config.KeyMap.Undo == "" {
		config.KeyMap.Undo = defaults.KeyMap.Undo
	}
	if config.KeyMap.Redo == "" {
		config.KeyMap.Redo = defaults.KeyMap.Redo
	}
}

// SaveConfig saves the configuration to the specified file path.
// Supports both JSON and YAML formats based on file extension.
func (c *Config) SaveConfig(configPath string) error {
	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	file, err := os.OpenFile(configPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open config file for writing: %w", err)
	}
	defer file.Close()

	format := detectConfigFormat(configPath)
	var data []byte

	if format == ConfigFormatYAML {
		data, err = yaml.Marshal(c)
		if err != nil {
			return fmt.Errorf("failed to marshal config to YAML: %w", err)
		}
	} else {
		data, err = json.MarshalIndent(c, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to marshal config to JSON: %w", err)
		}
	}

	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetConfigPath returns the default configuration file path.
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ".td-config.json"
	}
	return filepath.Join(homeDir, ".config", "td", "config.json")
}
