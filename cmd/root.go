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
	Version:           "0.0.1",
	DisableAutoGenTag: true,
}

func init() {
	// Ticket Creation Commands
	epicCmd := commands.GetEpicCmd()
	epicCmd.GroupID = "ticket-creation"
	rootCmd.AddCommand(epicCmd)

	taskCmd := commands.GetTaskCmd()
	taskCmd.GroupID = "ticket-creation"
	rootCmd.AddCommand(taskCmd)

	subtaskCmd := commands.GetSubtaskCmd()
	subtaskCmd.GroupID = "ticket-creation"
	rootCmd.AddCommand(subtaskCmd)

	// Context Management Commands
	trackCmd := commands.GetTrackCmd()
	trackCmd.GroupID = "context-management"
	rootCmd.AddCommand(trackCmd)

	focusCmd := commands.GetFocusCmd()
	focusCmd.GroupID = "context-management"
	rootCmd.AddCommand(focusCmd)

	// Status & Workflow Commands
	statusCmd := commands.GetStatusCmd()
	statusCmd.GroupID = "status-workflow"
	rootCmd.AddCommand(statusCmd)

	cleanupCmd := commands.GetCleanupCmd()
	cleanupCmd.GroupID = "status-workflow"
	rootCmd.AddCommand(cleanupCmd)

	// View & Navigation Commands
	logCmd := commands.GetLogCmd()
	logCmd.GroupID = "view-navigation"
	rootCmd.AddCommand(logCmd)

	linkCmd := commands.GetLinkCmd()
	linkCmd.GroupID = "view-navigation"
	rootCmd.AddCommand(linkCmd)

	openCmd := commands.GetOpenCmd()
	openCmd.GroupID = "view-navigation"
	rootCmd.AddCommand(openCmd)

	// Collaboration Commands
	commentCmd := commands.GetCommentCmd()
	commentCmd.GroupID = "collaboration"
	rootCmd.AddCommand(commentCmd)

	// Setup & Utility Commands
	initCmd := commands.GetInitCmd()
	initCmd.GroupID = "setup-utility"
	rootCmd.AddCommand(initCmd)

	versionCmd := commands.GetVersionCmd()
	versionCmd.GroupID = "setup-utility"
	rootCmd.AddCommand(versionCmd)

	completionCmd := commands.GetCompletionCmd()
	completionCmd.GroupID = "setup-utility"
	rootCmd.AddCommand(completionCmd)

	// Set command groups
	rootCmd.AddGroup(&cobra.Group{
		ID:    "ticket-creation",
		Title: "Ticket Creation:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "context-management",
		Title: "Context Management:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "status-workflow",
		Title: "Status & Workflow:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "view-navigation",
		Title: "View & Navigation:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "collaboration",
		Title: "Collaboration:",
	})
	rootCmd.AddGroup(&cobra.Group{
		ID:    "setup-utility",
		Title: "Setup & Utility:",
	})
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
