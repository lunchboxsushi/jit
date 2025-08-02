package config

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"

	"github.com/lunchboxsushi/jit/pkg/types"
)

var (
	configInstance *types.Config
	configOnce     sync.Once
	configError    error
)

// Load loads the configuration from the default location
func Load() (*types.Config, error) {
	configOnce.Do(func() {
		configInstance, configError = loadConfig()
	})
	return configInstance, configError
}

// Get returns the loaded configuration instance
func Get() *types.Config {
	config, _ := Load()
	return config
}

// loadConfig loads and validates the configuration
func loadConfig() (*types.Config, error) {
	configPath := GetDefaultConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("configuration file not found: %s. Run 'jit init' to create one", configPath)
	}

	// Read config file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	// Expand environment variables
	expandedData := ExpandEnvironmentVariables(string(data))

	// Parse YAML
	var config types.Config
	if err := yaml.Unmarshal([]byte(expandedData), &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file %s: %v", configPath, err)
	}

	// Validate configuration
	if errors := ValidateConfig(&config); len(errors) > 0 {
		var errorMsgs []string
		for _, err := range errors {
			errorMsgs = append(errorMsgs, err.Error())
		}
		return nil, fmt.Errorf("configuration validation failed:\n%s", strings.Join(errorMsgs, "\n"))
	}

	return &config, nil
}

// expandEnvironmentVariables expands ${VAR} patterns in the config content
func ExpandEnvironmentVariables(content string) string {
	// Regex to match ${VAR} patterns
	re := regexp.MustCompile(`\$\{([^}]+)\}`)

	return re.ReplaceAllStringFunc(content, func(match string) string {
		// Extract variable name from ${VAR}
		varName := match[2 : len(match)-1]

		// Get environment variable value
		value := os.Getenv(varName)

		// If not found, return the original match
		if value == "" {
			return match
		}

		return value
	})
}

// CreateDefaultConfig creates a default configuration file
func CreateDefaultConfig() error {
	configPath := GetDefaultConfigPath()

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %v", configDir, err)
	}

	// Create default config
	config := types.Config{
		Jira: types.JiraConfig{
			URL:           "https://your-company.atlassian.net",
			Username:      "your-email@company.com",
			Token:         "${JIRA_API_TOKEN}",
			Project:       "PROJ",
			EpicLinkField: "customfield_10014",
		},
		AI: types.AIConfig{
			Provider:  "openai",
			APIKey:    "${OPENAI_API_KEY}",
			Model:     "gpt-4",
			MaxTokens: 1000,
		},
		App: types.AppConfig{
			DataDir:            GetDefaultDataPath(),
			DefaultEditor:      "vim",
			ReviewBeforeCreate: true,
		},
	}

	// Marshal to YAML
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %v", configPath, err)
	}

	return nil
}

// SaveConfig saves the configuration to the default location
func SaveConfig(config *types.Config) error {
	configPath := GetDefaultConfigPath()

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory %s: %v", configDir, err)
	}

	// Marshal to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	// Write to file
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file %s: %v", configPath, err)
	}

	return nil
}
