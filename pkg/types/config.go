package types

// Config represents the application configuration
type Config struct {
	Jira JiraConfig `yaml:"jira" json:"jira"`
	AI   AIConfig   `yaml:"ai" json:"ai"`
	App  AppConfig  `yaml:"app" json:"app"`
}

// JiraConfig contains Jira connection settings
type JiraConfig struct {
	URL           string `yaml:"url" json:"url"`
	Username      string `yaml:"username" json:"username"`
	Token         string `yaml:"token" json:"token"`
	Project       string `yaml:"project" json:"project"`
	EpicLinkField string `yaml:"epic_link_field" json:"epic_link_field"`
}

// AIConfig contains AI provider settings
type AIConfig struct {
	Provider  string `yaml:"provider" json:"provider"`
	APIKey    string `yaml:"api_key" json:"api_key"`
	Model     string `yaml:"model" json:"model"`
	MaxTokens int    `yaml:"max_tokens" json:"max_tokens"`
}

// AppConfig contains application settings
type AppConfig struct {
	DataDir            string `yaml:"data_dir" json:"data_dir"`
	DefaultEditor      string `yaml:"default_editor" json:"default_editor"`
	ReviewBeforeCreate bool   `yaml:"review_before_create" json:"review_before_create"`
}

// NewConfig creates a new config with default values
func NewConfig() *Config {
	return &Config{
		Jira: JiraConfig{
			EpicLinkField: "customfield_10014",
		},
		AI: AIConfig{
			Provider:  "openai",
			Model:     "gpt-4",
			MaxTokens: 1000,
		},
		App: AppConfig{
			DataDir:            "~/.local/share/jit",
			DefaultEditor:      "vim",
			ReviewBeforeCreate: true,
		},
	}
}
