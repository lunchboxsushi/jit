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
	rootCmd.AddCommand(commands.GetVersionCmd())
	rootCmd.AddCommand(commands.GetInitCmd())
	rootCmd.AddCommand(commands.GetTrackCmd())
	rootCmd.AddCommand(commands.GetFocusCmd())
	rootCmd.AddCommand(commands.GetEpicCmd())
	rootCmd.AddCommand(commands.GetTaskCmd())
	rootCmd.AddCommand(commands.GetSubtaskCmd())
	rootCmd.AddCommand(commands.GetLogCmd())
	rootCmd.AddCommand(commands.GetLinkCmd())
	rootCmd.AddCommand(commands.GetOpenCmd())
	rootCmd.AddCommand(commands.GetCommentCmd())
	rootCmd.AddCommand(commands.GetStatusCmd())
	rootCmd.AddCommand(commands.GetCleanupCmd())
	rootCmd.AddCommand(commands.GetCompletionCmd())
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
