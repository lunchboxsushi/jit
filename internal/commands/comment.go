package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/lunchboxsushi/jit/internal/ui"
	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

var (
	commentMessageFlag bool
)

var commentCmd = &cobra.Command{
	Use:   "comment [ticket-key]",
	Short: "Add comment to Jira ticket",
	Long: `Add a comment to a Jira ticket. If no ticket is specified, uses current focus.
	
Examples:
  jit comment                    # Add comment to current focus
  jit comment SRE-1234          # Add comment to specific ticket
  jit comment -m "Quick note"   # Inline comment`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// Initialize command context
		ctx, err := InitializeCommand()
		if err != nil {
			HandleError(err, "Failed to initialize")
			return
		}

		var ticketKey string

		// Get ticket key from args or current focus
		if len(args) > 0 {
			ticketKey = args[0]
		} else {
			// Try to get from current focus
			currentEpic, _ := ctx.ContextManager.GetCurrentEpic()
			currentTask, _ := ctx.ContextManager.GetCurrentTask()
			currentSubtask, _ := ctx.ContextManager.GetCurrentSubtask()

			// Prefer subtask > task > epic
			if currentSubtask != "" {
				ticketKey = currentSubtask
			} else if currentTask != "" {
				ticketKey = currentTask
			} else if currentEpic != "" {
				ticketKey = currentEpic
			} else {
				fmt.Println("No ticket specified and no current focus.")
				fmt.Println("Use 'jit focus <ticket>' to set focus or specify a ticket key.")
				return
			}
		}

		// Validate ticket exists
		if !ctx.Storage.Exists(ticketKey) {
			fmt.Printf("Ticket %s not found in local storage.\n", ticketKey)
			fmt.Printf("Use 'jit track %s' to track it first.\n", ticketKey)
			return
		}

		// Load ticket for display
		ticket, err := ctx.Storage.LoadTicket(ticketKey)
		if err != nil {
			HandleError(err, "Failed to load ticket")
			return
		}

		// Get comment content
		var commentBody string
		if commentMessageFlag {
			// Get comment from remaining args
			if len(args) < 2 {
				fmt.Println("Error: Inline comment requires comment text.")
				fmt.Println("Usage: jit comment -m \"Your comment text\"")
				return
			}
			commentBody = strings.Join(args[1:], " ")
		} else {
			// Open editor for comment
			commentBody, err = getCommentFromEditor(ticket)
			if err != nil {
				HandleError(err, "Failed to get comment from editor")
				return
			}
		}

		// Validate comment is not empty
		if strings.TrimSpace(commentBody) == "" {
			fmt.Println("Comment is empty. Aborting.")
			return
		}

		// AI enrichment (if provider available)
		if ctx.AIProvider != nil {
			enrichedComment, err := ctx.EnrichCommentWithAI(commentBody, ticketKey)
			if err != nil {
				PrintWarning(fmt.Sprintf("AI enrichment failed: %v", err))
			} else {
				commentBody = enrichedComment
			}
		}

		// Add comment to Jira
		fmt.Printf("Adding comment to %s...\n", ticketKey)

		if err := ctx.TicketService.AddComment(cmd.Context(), ticketKey, commentBody); err != nil {
			HandleError(err, "Failed to add comment to Jira")
			return
		}

		PrintSuccess(fmt.Sprintf("Comment added to %s", ticketKey))
	},
}

func init() {
	commentCmd.Flags().BoolVarP(&commentMessageFlag, "message", "m", false, "Add inline comment (requires comment text)")
}

// getCommentFromEditor opens an editor to get comment content
func getCommentFromEditor(ticket *types.Ticket) (string, error) {
	// Create temporary file for editing
	tempFile, err := os.CreateTemp("", "jit-comment-*.md")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Create comment template
	template := fmt.Sprintf(`# Comment for %s: %s

Add your comment below this line:
---

`)

	// Write template to temp file
	if err := os.WriteFile(tempFile.Name(), []byte(template), 0644); err != nil {
		return "", fmt.Errorf("failed to write template: %v", err)
	}

	// Open editor
	editor := ui.NewEditor()
	if err := editor.EditFile(tempFile.Name()); err != nil {
		return "", fmt.Errorf("failed to open editor: %v", err)
	}

	// Read the edited content
	content, err := editor.ReadFile(tempFile.Name())
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	// Extract comment from content (everything after the separator)
	lines := strings.Split(content, "\n")
	var commentLines []string
	foundSeparator := false

	for _, line := range lines {
		if strings.TrimSpace(line) == "---" {
			foundSeparator = true
			continue
		}
		if foundSeparator {
			commentLines = append(commentLines, line)
		}
	}

	comment := strings.Join(commentLines, "\n")
	return strings.TrimSpace(comment), nil
}

// GetCommentCmd returns the comment command
func GetCommentCmd() *cobra.Command {
	return commentCmd
}
