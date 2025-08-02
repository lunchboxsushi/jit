package cmd

import (
	"github.com/lunchboxsushi/jit/internal/commands"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "jit",
	Short: "JIT - A jira experience like git, focus ticket work like branches",
	Long: `jit is a local-first CLI tool that lets developers write tasks and sub-tasks in markdown,
auto-enriches raw task descriptions into manager-optimized Jira tickets, and syncs with Jira
to reflect status, updates, and structure.`,
	Version: "0.0.1",
}

func init() {
	// Add version command
	rootCmd.AddCommand(commands.GetVersionCmd())

	// Add init command
	rootCmd.AddCommand(commands.GetInitCmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
