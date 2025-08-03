package config

import (
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("config validation error in %s: %s", e.Field, e.Message)
}

// ValidateConfig validates the configuration and returns any errors
func ValidateConfig(config *types.Config) []error {
	var errors []error

	// Validate Jira configuration
	if err := validateJiraConfig(config.Jira); err != nil {
		errors = append(errors, err)
	}

	// Validate AI configuration (optional)
	if err := validateAIConfig(config.AI); err != nil {
		errors = append(errors, err)
	}

	// Validate App configuration
	if err := validateAppConfig(config.App); err != nil {
		errors = append(errors, err)
	}

	return errors
}

// validateJiraConfig validates Jira-specific configuration
func validateJiraConfig(jira types.JiraConfig) error {
	// URL is required
	if jira.URL == "" {
		return ValidationError{Field: "jira.url", Message: "Jira URL is required"}
	}

	// Validate URL format
	if _, err := url.ParseRequestURI(jira.URL); err != nil {
		return ValidationError{Field: "jira.url", Message: fmt.Sprintf("Invalid URL format: %v", err)}
	}

	// Username is required
	if jira.Username == "" {
		return ValidationError{Field: "jira.username", Message: "Jira username is required"}
	}

	// Token is required
	if jira.Token == "" {
		return ValidationError{Field: "jira.token", Message: "Jira API token is required"}
	}

	// Project is required
	if jira.Project == "" {
		return ValidationError{Field: "jira.project", Message: "Jira project key is required"}
	}

	// Project should be uppercase
	if jira.Project != strings.ToUpper(jira.Project) {
		return ValidationError{Field: "jira.project", Message: "Project key should be uppercase"}
	}

	return nil
}

// validateAIConfig validates AI-specific configuration (optional)
func validateAIConfig(ai types.AIConfig) error {
	// If provider is empty, AI is disabled (which is fine)
	if ai.Provider == "" {
		return nil
	}

	// Validate provider
	validProviders := []string{"openai", "anthropic", "local", "mock", "test"}
	valid := false
	for _, provider := range validProviders {
		if ai.Provider == provider {
			valid = true
			break
		}
	}
	if !valid {
		return ValidationError{Field: "ai.provider", Message: fmt.Sprintf("Invalid AI provider: %s. Valid providers: %v", ai.Provider, validProviders)}
	}

	// If provider is set, API key is required
	if ai.APIKey == "" {
		return ValidationError{Field: "ai.api_key", Message: "AI API key is required when AI provider is configured"}
	}

	// Validate model (optional, but if provided should be non-empty)
	if ai.Model != "" && strings.TrimSpace(ai.Model) == "" {
		return ValidationError{Field: "ai.model", Message: "AI model cannot be empty if provided"}
	}

	// Validate max tokens
	if ai.MaxTokens <= 0 {
		return ValidationError{Field: "ai.max_tokens", Message: "AI max tokens must be greater than 0"}
	}

	return nil
}

// validateAppConfig validates application-specific configuration
func validateAppConfig(app types.AppConfig) error {
	// Data directory is required
	if app.DataDir == "" {
		return ValidationError{Field: "app.data_dir", Message: "Data directory is required"}
	}

	// Editor is required
	if app.DefaultEditor == "" {
		return ValidationError{Field: "app.default_editor", Message: "Default editor is required"}
	}

	return nil
}

// IsConfigMissing checks if the configuration file is missing
func IsConfigMissing() bool {
	configPath := GetDefaultConfigPath()
	_, err := os.Stat(configPath)
	return os.IsNotExist(err)
}
