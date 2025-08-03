package commands

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status [status]",
	Short: "Change status of current focus item",
	Long: `Change the status of the currently focused ticket.
	
Available statuses:
  todo, to-do, todo     - Set status to "To Do"
  progress, in-progress - Set status to "In Progress" 
  done, completed       - Set status to "Done"
  blocked              - Set status to "Blocked"

Examples:
  jit status done       # Mark current focus as done
  jit status progress   # Start working on current focus
  jit status blocked    # Mark current focus as blocked`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		status := normalizeStatus(args[0])
		if status == "" {
			fmt.Println("Invalid status. Use: todo, progress, done, or blocked")
			return
		}

		// TODO: Implement status change logic
		// For now, just show what would happen
		fmt.Printf("Would change current focus status to: %s\n", status)
		fmt.Println("(Status command not yet implemented)")
	},
}

// normalizeStatus converts various status inputs to standard status values
func normalizeStatus(input string) string {
	status := strings.ToLower(strings.TrimSpace(input))

	switch status {
	case "todo", "to-do", "to do":
		return "To Do"
	case "progress", "in-progress", "in progress":
		return "In Progress"
	case "done", "completed", "complete":
		return "Done"
	case "blocked", "block":
		return "Blocked"
	default:
		return ""
	}
}

// GetStatusCmd returns the status command
func GetStatusCmd() *cobra.Command {
	return statusCmd
}
