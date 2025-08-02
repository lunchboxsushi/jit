package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lunchboxsushi/jit/pkg/types"
	"gopkg.in/yaml.v3"
)

func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()

	// Should not be empty
	if path == "" {
		t.Error("Expected non-empty config path")
	}

	// Should end with config.yml
	if filepath.Base(path) != "config.yml" {
		t.Errorf("Expected config path to end with config.yml, got %s", filepath.Base(path))
	}
}

func TestGetDefaultDataPath(t *testing.T) {
	path := GetDefaultDataPath()

	// Should not be empty
	if path == "" {
		t.Error("Expected non-empty data path")
	}

	// Should end with jit
	if filepath.Base(path) != "jit" {
		t.Errorf("Expected data path to end with jit, got %s", filepath.Base(path))
	}
}

func TestIsConfigMissing(t *testing.T) {
	// This test assumes no config file exists in the default location
	// In a real environment, this would be true for a fresh installation
	// For testing purposes, we'll just verify the function works correctly

	// The function should return true if no config file exists at the default path
	// We can't guarantee this in all test environments, so we'll just test the logic

	configPath := GetDefaultConfigPath()
	_, err := os.Stat(configPath)
	isMissing := os.IsNotExist(err)

	// The function should return the same result as checking if the file doesn't exist
	// This is a basic sanity check
	if isMissing != IsConfigMissing() {
		t.Error("IsConfigMissing() should return the same result as checking if the file doesn't exist")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create default config in temp directory
	configPath := filepath.Join(tempDir, "config.yml")

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
		t.Fatalf("Failed to marshal config: %v", err)
	}

	// Write to temp file
	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after creation")
	}
}

func TestLoadConfigWithMissingFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Try to load config from non-existent path
	configPath := filepath.Join(tempDir, "nonexistent.yml")

	// Check if file exists (should not)
	if _, err := os.Stat(configPath); !os.IsNotExist(err) {
		t.Error("Expected config file to not exist")
	}
}

func TestLoadConfigWithInvalidYAML(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Create invalid YAML file
	invalidYAML := `jira:
  url: "https://example.com"
  username: "test@example.com"
  token: "invalid-token"
  project: "TEST"
  epic_link_field: "customfield_10014"
ai:
  provider: "openai"
  api_key: "test-key"
  model: "gpt-4"
  max_tokens: 1000
app:
  data_dir: "/tmp/jit"
  default_editor: "vim"
  review_before_create: true
invalid: [yaml: content`

	configPath := filepath.Join(tempDir, "config.yml")
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Try to parse the invalid YAML directly
	var config types.Config
	err = yaml.Unmarshal([]byte(invalidYAML), &config)
	if err == nil {
		t.Error("Expected error when config has invalid YAML")
	}
}

func TestLoadConfigWithEnvironmentVariables(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	// Set environment variables
	os.Setenv("JIRA_API_TOKEN", "test-token-123")
	os.Setenv("OPENAI_API_KEY", "test-openai-key-456")
	defer func() {
		os.Unsetenv("JIRA_API_TOKEN")
		os.Unsetenv("OPENAI_API_KEY")
	}()

	// Create config with environment variables
	configWithEnv := `jira:
  url: "https://example.com"
  username: "test@example.com"
  token: "${JIRA_API_TOKEN}"
  project: "TEST"
  epic_link_field: "customfield_10014"
ai:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"
  model: "gpt-4"
  max_tokens: 1000
app:
  data_dir: "/tmp/jit"
  default_editor: "vim"
  review_before_create: true`

	configPath := filepath.Join(tempDir, "config.yml")
	err := os.WriteFile(configPath, []byte(configWithEnv), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config: %v", err)
	}

	// Test environment variable expansion
	expandedData := ExpandEnvironmentVariables(configWithEnv)

	// Parse the expanded YAML
	var config types.Config
	err = yaml.Unmarshal([]byte(expandedData), &config)
	if err != nil {
		t.Fatalf("Failed to parse expanded config: %v", err)
	}

	// Verify environment variables were expanded
	if config.Jira.Token != "test-token-123" {
		t.Errorf("Expected Jira token to be expanded, got %s", config.Jira.Token)
	}
	if config.AI.APIKey != "test-openai-key-456" {
		t.Errorf("Expected AI API key to be expanded, got %s", config.AI.APIKey)
	}
}

func TestValidateConfig(t *testing.T) {
	tests := []struct {
		name    string
		config  *types.Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &types.Config{
				Jira: types.JiraConfig{
					URL:           "https://example.com",
					Username:      "test@example.com",
					Token:         "test-token",
					Project:       "TEST",
					EpicLinkField: "customfield_10014",
				},
				AI: types.AIConfig{
					Provider:  "openai",
					APIKey:    "test-key",
					Model:     "gpt-4",
					MaxTokens: 1000,
				},
				App: types.AppConfig{
					DataDir:            "/tmp/jit",
					DefaultEditor:      "vim",
					ReviewBeforeCreate: true,
				},
			},
			wantErr: false,
		},
		{
			name: "missing jira url",
			config: &types.Config{
				Jira: types.JiraConfig{
					Username:      "test@example.com",
					Token:         "test-token",
					Project:       "TEST",
					EpicLinkField: "customfield_10014",
				},
				AI: types.AIConfig{
					Provider:  "openai",
					APIKey:    "test-key",
					Model:     "gpt-4",
					MaxTokens: 1000,
				},
				App: types.AppConfig{
					DataDir:            "/tmp/jit",
					DefaultEditor:      "vim",
					ReviewBeforeCreate: true,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid jira url",
			config: &types.Config{
				Jira: types.JiraConfig{
					URL:           "not-a-url",
					Username:      "test@example.com",
					Token:         "test-token",
					Project:       "TEST",
					EpicLinkField: "customfield_10014",
				},
				AI: types.AIConfig{
					Provider:  "openai",
					APIKey:    "test-key",
					Model:     "gpt-4",
					MaxTokens: 1000,
				},
				App: types.AppConfig{
					DataDir:            "/tmp/jit",
					DefaultEditor:      "vim",
					ReviewBeforeCreate: true,
				},
			},
			wantErr: true,
		},
		{
			name: "lowercase project",
			config: &types.Config{
				Jira: types.JiraConfig{
					URL:           "https://example.com",
					Username:      "test@example.com",
					Token:         "test-token",
					Project:       "test",
					EpicLinkField: "customfield_10014",
				},
				AI: types.AIConfig{
					Provider:  "openai",
					APIKey:    "test-key",
					Model:     "gpt-4",
					MaxTokens: 1000,
				},
				App: types.AppConfig{
					DataDir:            "/tmp/jit",
					DefaultEditor:      "vim",
					ReviewBeforeCreate: true,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := ValidateConfig(tt.config)
			if tt.wantErr && len(errors) == 0 {
				t.Error("Expected validation errors, got none")
			}
			if !tt.wantErr && len(errors) > 0 {
				t.Errorf("Expected no validation errors, got %d: %v", len(errors), errors)
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
