package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/voioo/td/internal/config"
	"github.com/voioo/td/internal/logger"
	"github.com/voioo/td/internal/storage"
	"github.com/voioo/td/internal/task"
	"github.com/voioo/td/internal/ui"
	"github.com/voioo/td/internal/upgrade"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// initializeModel creates and initializes the UI model with loaded data.
func initializeModel(cfg *config.Config) (tea.Model, error) {
	logger.Info("Initializing application",
		logger.F("data_file", cfg.DataFile))

	// Load tasks from storage
	repo := storage.NewRepository(cfg.DataFile)
	activeTasks, doneTasks, nextID, err := repo.LoadTasks()
	if err != nil {
		logger.Error("Failed to load tasks", logger.F("error", err))
		return nil, fmt.Errorf("failed to load tasks: %w", err)
	}

	logger.Info("Loaded tasks from storage",
		logger.F("active_tasks", len(activeTasks)),
		logger.F("done_tasks", len(doneTasks)),
		logger.F("next_id", nextID))

	// Create task manager
	taskManager := task.NewTaskManager(activeTasks, doneTasks, nextID)

	// Create UI model
	uiModel := ui.NewModel(cfg, taskManager)

	logger.Info("Application initialized successfully")
	return uiModel, nil
}

func main() {
	versionFlag := flag.Bool("version", false, "print version information")
	flag.BoolVar(versionFlag, "v", false, "print version information (shorthand)")
	upgradeFlag := flag.Bool("upgrade", false, "upgrade td to the latest version")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("td %s (commit: %s, built at: %s)\n", version, commit, date)
		logger.Info("Version requested",
			logger.F("version", version),
			logger.F("commit", commit),
			logger.F("date", date))
		os.Exit(0)
	}

	if *upgradeFlag {
		logger.Info("Upgrade requested",
			logger.F("version", version),
			logger.F("commit", commit))
		if err := upgrade.Upgrade(version); err != nil {
			logger.Error("Upgrade failed", logger.F("error", err))
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Successfully upgraded td!")
		os.Exit(0)
	}

	logger.Info("Starting td application",
		logger.F("version", version),
		logger.F("commit", commit))

	// Cleanup old executable files from previous upgrades (Windows only)
	upgrade.CleanupOldExecutables()

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		logger.Fatal("Failed to load configuration", logger.F("error", err))
	}

	// Initialize the model
	model, err := initializeModel(cfg)
	if err != nil {
		logger.Fatal("Failed to initialize application", logger.F("error", err))
	}

	// Start the Bubble Tea program
	logger.Info("Starting Bubble Tea program")
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		logger.Fatal("Application error", logger.F("error", err))
	}

	logger.Info("Application shutdown complete")
}
