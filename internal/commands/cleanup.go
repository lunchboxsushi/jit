package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cleanupCmd = &cobra.Command{
	Use:   "cleanup [ticket-key]",
	Short: "Remove done tickets",
	Long: `Remove tickets that are marked as done.
	
Examples:
  jit cleanup              # Remove all done tickets
  jit cleanup PROJ-100     # Remove specific done ticket
  jit cleanup --dry-run    # Show what would be removed without doing it`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if len(args) == 0 {
			// Remove all done tickets
			if dryRun {
				fmt.Println("Would remove all done tickets")
			} else {
				fmt.Println("Removing all done tickets...")
			}
		} else {
			// Remove specific ticket
			ticketKey := args[0]
			if dryRun {
				fmt.Printf("Would remove ticket: %s\n", ticketKey)
			} else {
				fmt.Printf("Removing ticket: %s\n", ticketKey)
			}
		}

		// TODO: Implement actual cleanup logic
		fmt.Println("(Cleanup command not yet implemented)")
	},
}

func init() {
	cleanupCmd.Flags().Bool("dry-run", false, "Show what would be removed without doing it")
}

// GetCleanupCmd returns the cleanup command
func GetCleanupCmd() *cobra.Command {
	return cleanupCmd
}
