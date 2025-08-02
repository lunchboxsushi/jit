package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/spf13/cobra"
)

var testConfigCmd = &cobra.Command{
	Use:   "test-config",
	Short: "Test configuration loading",
	Long:  `Test command to verify configuration loading works correctly`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Testing configuration loading...")

		// Try to load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ Configuration error: %v\n", err)
			fmt.Println("\n💡 Run 'jit init' to create a configuration file")
			return
		}

		fmt.Println("✅ Configuration loaded successfully!")
		fmt.Printf("Jira URL: %s\n", cfg.Jira.URL)
		fmt.Printf("Jira Project: %s\n", cfg.Jira.Project)
		fmt.Printf("AI Provider: %s\n", cfg.AI.Provider)
		fmt.Printf("Data Directory: %s\n", cfg.App.DataDir)
	},
}

// GetTestConfigCmd returns the test config command
func GetTestConfigCmd() *cobra.Command {
	return testConfigCmd
}
