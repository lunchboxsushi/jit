package commands

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/lunchboxsushi/jit/internal/jira"
	"github.com/lunchboxsushi/jit/internal/storage"
	"github.com/lunchboxsushi/jit/pkg/types"
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
			fmt.Printf("❌ Invalid ticket key: %v\n", err)
			return
		}

		// Load configuration
		cfg, err := config.Load()
		if err != nil {
			fmt.Printf("❌ Configuration error: %v\n", err)
			fmt.Println("💡 Run 'jit init' to create a configuration file")
			return
		}

		// Initialize storage
		storageInstance, err := storage.NewJSONStorage(cfg.App.DataDir)
		if err != nil {
			fmt.Printf("❌ Storage error: %v\n", err)
			return
		}

		// Initialize Jira client
		jiraClient := jira.NewClient(&cfg.Jira)

		ticketService := jira.NewTicketService(jiraClient)
		contextManager := storage.NewContextManager(storageInstance)

		// Track the ticket
		if err := trackTicket(context.Background(), ticketKey, ticketService, storageInstance, contextManager, !noChildrenFlag); err != nil {
			fmt.Printf("❌ Failed to track ticket: %v\n", err)
			return
		}

		fmt.Printf("✅ Successfully tracked %s\n", ticketKey)
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

// trackTicket tracks a ticket and optionally its children
func trackTicket(ctx context.Context, ticketKey string, ticketService *jira.TicketService, storage storage.Storage, contextManager *storage.ContextManager, fetchChildren bool) error {
	fmt.Printf("🔍 Fetching ticket %s...\n", ticketKey)

	// Check if ticket already exists
	if storage.Exists(ticketKey) {
		fmt.Printf("📝 Ticket %s already exists, updating...\n", ticketKey)
	}

	// Fetch the main ticket
	ticket, err := ticketService.GetTicket(ctx, ticketKey)
	if err != nil {
		return fmt.Errorf("failed to fetch ticket %s: %v", ticketKey, err)
	}

	// Save the main ticket
	if err := storage.SaveTicket(ticket); err != nil {
		return fmt.Errorf("failed to save ticket %s: %v", ticketKey, err)
	}

	fmt.Printf("✅ Saved %s (%s)\n", ticketKey, ticket.Title)

	// Set focus to the tracked ticket
	if err := contextManager.SetFocus(ticketKey, ticket.Type); err != nil {
		return fmt.Errorf("failed to set focus: %v", err)
	}

	// Add to recent tickets
	if err := contextManager.AddToRecent(ticketKey); err != nil {
		return fmt.Errorf("failed to add to recent tickets: %v", err)
	}

	// Fetch children if requested and ticket is an epic
	if fetchChildren && ticket.Type == types.TicketTypeEpic {
		if err := trackEpicChildren(ctx, ticketKey, ticketService, storage); err != nil {
			return fmt.Errorf("failed to track epic children: %v", err)
		}
	}

	// Fetch subtasks if requested and ticket is a task
	if fetchChildren && ticket.Type == types.TicketTypeTask {
		if err := trackTaskSubtasks(ctx, ticketKey, ticketService, storage); err != nil {
			return fmt.Errorf("failed to track task subtasks: %v", err)
		}
	}

	return nil
}

// trackEpicChildren tracks all children of an epic
func trackEpicChildren(ctx context.Context, epicKey string, ticketService *jira.TicketService, storage storage.Storage) error {
	fmt.Printf("🌳 Fetching children of epic %s...\n", epicKey)

	children, err := ticketService.GetEpicChildren(ctx, epicKey)
	if err != nil {
		return fmt.Errorf("failed to fetch epic children: %v", err)
	}

	if len(children) == 0 {
		fmt.Println("   No children found")
		return nil
	}

	fmt.Printf("   Found %d children\n", len(children))

	// Save each child
	for i, child := range children {
		fmt.Printf("   [%d/%d] Saving %s (%s)...\n", i+1, len(children), child.Key, child.Title)
		
		if err := storage.SaveTicket(child); err != nil {
			return fmt.Errorf("failed to save child ticket %s: %v", child.Key, err)
		}

		// If this child is a task, fetch its subtasks too
		if child.Type == types.TicketTypeTask {
			if err := trackTaskSubtasks(ctx, child.Key, ticketService, storage); err != nil {
				return fmt.Errorf("failed to track subtasks of %s: %v", child.Key, err)
			}
		}

		// Small delay to be respectful to the API
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// trackTaskSubtasks tracks all subtasks of a task
func trackTaskSubtasks(ctx context.Context, taskKey string, ticketService *jira.TicketService, storage storage.Storage) error {
	fmt.Printf("   📋 Fetching subtasks of task %s...\n", taskKey)

	subtasks, err := ticketService.GetTaskSubtasks(ctx, taskKey)
	if err != nil {
		return fmt.Errorf("failed to fetch task subtasks: %v", err)
	}

	if len(subtasks) == 0 {
		fmt.Println("      No subtasks found")
		return nil
	}

	fmt.Printf("      Found %d subtasks\n", len(subtasks))

	// Save each subtask
	for i, subtask := range subtasks {
		fmt.Printf("      [%d/%d] Saving %s (%s)...\n", i+1, len(subtasks), subtask.Key, subtask.Title)
		
		if err := storage.SaveTicket(subtask); err != nil {
			return fmt.Errorf("failed to save subtask %s: %v", subtask.Key, err)
		}

		// Small delay to be respectful to the API
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// GetTrackCmd returns the track command
func GetTrackCmd() *cobra.Command {
	return trackCmd
}
