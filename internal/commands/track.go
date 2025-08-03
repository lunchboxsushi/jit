package commands

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var (
	noChildrenFlag bool
)

var trackCmd = &cobra.Command{
	Use:   "track <ticket-key>",
	Short: "Track a Jira ticket and its children",
	Long: `Download a Jira ticket and optionally its entire hierarchy (children) into local storage.
	
Examples:
  jit track SRE-1234          # Track a single ticket
  jit track SRE-1234 --no-children  # Track without children
  jit track EPIC-567          # Track epic and all its tasks/subtasks`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		ticketKey := strings.ToUpper(args[0])

		// Validate ticket key format
		if err := validateTicketKey(ticketKey); err != nil {
			fmt.Printf("Invalid ticket key: %v\n", err)
			return
		}

		// Initialize command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "Failed to initialize")
			return
		}

		// Track the ticket using common method
		if err := ctx.TrackTicketWithChildren(ticketKey, !noChildrenFlag); err != nil {
			HandleError(err, "Failed to track ticket")
			return
		}

		PrintSuccess(fmt.Sprintf("Successfully tracked %s", ticketKey))
	},
}

func init() {
	trackCmd.Flags().BoolVar(&noChildrenFlag, "no-children", false, "Skip fetching child tickets")
}

// validateTicketKey validates the format of a Jira ticket key
func validateTicketKey(key string) error {
	// Jira ticket keys are typically PROJECT-NUMBER format
	// e.g., SRE-1234, PROJ-567, etc.
	pattern := `^[A-Z]+-\d+$`
	matched, err := regexp.MatchString(pattern, key)
	if err != nil {
		return fmt.Errorf("failed to validate ticket key: %v", err)
	}

	if !matched {
		return fmt.Errorf("ticket key must be in format PROJECT-NUMBER (e.g., SRE-1234)")
	}

	return nil
}

// GetTrackCmd returns the track command
func GetTrackCmd() *cobra.Command {
	return trackCmd
}
