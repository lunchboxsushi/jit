package config

import (
	"os"
	"path/filepath"
)

// GetDefaultConfigPath returns the default configuration file path
func GetDefaultConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if we can't get home dir
		return filepath.Join(".", ".jit", "config.yml")
	}

	// Check XDG_CONFIG_HOME first
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		return filepath.Join(xdgConfig, "jit", "config.yml")
	}

	// Fallback to ~/.jit/config.yml
	return filepath.Join(homeDir, ".jit", "config.yml")
}

// GetDefaultDataPath returns the default data directory path
func GetDefaultDataPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if we can't get home dir
		return filepath.Join(".", ".local", "share", "jit")
	}

	// Check XDG_DATA_HOME first
	if xdgData := os.Getenv("XDG_DATA_HOME"); xdgData != "" {
		return filepath.Join(xdgData, "jit")
	}

	// Fallback to ~/.local/share/jit
	return filepath.Join(homeDir, ".local", "share", "jit")
}
