package commands

import (
	"fmt"

	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize jit configuration",
	Long:  `Create a default configuration file for jit`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing jit configuration...")

		// Check if config already exists
		if !config.IsConfigMissing() {
			fmt.Println("Warning: Configuration file already exists!")
			fmt.Println("Tip: If you want to recreate it, delete the existing file first.")
			return
		}

		// Create default configuration
		err := config.CreateDefaultConfig()
		if err != nil {
			fmt.Printf("Error: Failed to create configuration: %v\n", err)
			return
		}

		fmt.Println("Success: Configuration file created successfully!")
		fmt.Printf("Location: %s\n", config.GetDefaultConfigPath())
		fmt.Println("\nNext steps:")
		fmt.Println("1. Edit the configuration file with your Jira and AI settings")
		fmt.Println("2. Set the required environment variables:")
		fmt.Println("   - JIRA_API_TOKEN")
		fmt.Println("   - OPENAI_API_KEY (if using AI features)")
		fmt.Println("3. Run 'jit test-config' to verify your configuration")
	},
}

// GetInitCmd returns the init command
func GetInitCmd() *cobra.Command {
	return initCmd
}
