package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/lunchboxsushi/jit/internal/storage"
	"github.com/lunchboxsushi/jit/pkg/types"
)

// CommandContext holds the common context for commands
type CommandContext struct {
	Config  *types.Config
	Storage storage.Storage
}

// InitializeCommand sets up the common context for commands
func InitializeCommand() (*CommandContext, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %v\nRun 'jit init' to create a configuration file", err)
	}

	// Initialize storage
	storageInstance, err := storage.NewJSONStorage(cfg.App.DataDir)
	if err != nil {
		return nil, fmt.Errorf("storage error: %v", err)
	}

	return &CommandContext{
		Config:  cfg,
		Storage: storageInstance,
	}, nil
}

// HandleError provides consistent error handling across commands
func HandleError(err error, message string) {
	if err != nil {
		fmt.Printf("Error: %s: %v\n", message, err)
	}
}

// PrintSuccess provides consistent success messaging
func PrintSuccess(message string) {
	fmt.Printf("Success: %s\n", message)
}

// PrintInfo provides consistent info messaging
func PrintInfo(message string) {
	fmt.Printf("Info: %s\n", message)
}

// PrintWarning provides consistent warning messaging
func PrintWarning(message string) {
	fmt.Printf("Warning: %s\n", message)
}
