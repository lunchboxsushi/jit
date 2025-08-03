package commands

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/lunchboxsushi/jit/internal/jira"
	"github.com/lunchboxsushi/jit/internal/storage"
	"github.com/lunchboxsushi/jit/internal/ui"
	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

// CreateOptions defines the configuration for ticket creation
type CreateOptions struct {
	TicketType       string
	TemplateName     string
	ValidateContext  func(contextManager *storage.ContextManager) (string, error)
	SetRelationships func(ticket *types.Ticket, context string)
	SuccessMessage   string
	ParentInfo       string
}

// CreateTicket is the shared pipeline for creating tickets
func CreateTicket(cmd *cobra.Command, options CreateOptions, flags CreateFlags) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("configuration error: %v\n💡 Run 'jit init' to create a configuration file", err)
	}

	// Initialize storage
	storageInstance, err := storage.NewJSONStorage(cfg.App.DataDir)
	if err != nil {
		return fmt.Errorf("storage error: %v", err)
	}

	// Validate context if required
	contextManager := storage.NewContextManager(storageInstance)
	contextInfo := ""
	if options.ValidateContext != nil {
		contextInfo, err = options.ValidateContext(contextManager)
		if err != nil {
			return err
		}
	}

	// Initialize editor
	editor := ui.NewEditor()

	// Create temporary file for editing
	tempDir := filepath.Join(cfg.App.DataDir, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %v", err)
	}

	tempFile := filepath.Join(tempDir, fmt.Sprintf("%s_%d.md", options.TicketType, time.Now().Unix()))

	// Open template in editor
	templatePath := filepath.Join("templates", options.TemplateName)
	if err := editor.EditTemplate(templatePath, tempFile); err != nil {
		return fmt.Errorf("failed to open editor: %v", err)
	}

	// Read the edited content
	content, err := editor.ReadFile(tempFile)
	if err != nil {
		return fmt.Errorf("failed to read edited content: %v", err)
	}

	// Parse the markdown content
	title, description, err := editor.ParseMarkdownTicket(content)
	if err != nil {
		return fmt.Errorf("failed to parse content: %v\n💡 Please ensure you have a title and description", err)
	}

	// AI enrichment (placeholder for now)
	if !flags.NoEnrich {
		fmt.Println("🤖 AI enrichment would happen here (Task 11)")
		// TODO: Implement AI enrichment in Task 11
	}

	// Create ticket object
	ticket := types.NewTicket("", title, options.TicketType)
	ticket.Description = description
	ticket.Metadata.Project = cfg.Jira.Project
	ticket.Metadata.Assignee = cfg.Jira.Username

	// Set relationships
	if options.SetRelationships != nil {
		options.SetRelationships(ticket, contextInfo)
	}

	// Create in Jira if requested
	if !flags.NoCreate {
		fmt.Printf("🚀 Creating %s in Jira...\n", options.TicketType)

		// Initialize Jira client
		jiraClient := jira.NewClient(&cfg.Jira)
		ticketService := jira.NewTicketService(jiraClient)

		// Create the ticket
		createdTicket, err := ticketService.CreateTicket(cmd.Context(), ticket)
		if err != nil {
			return fmt.Errorf("failed to create %s in Jira: %v\n💡 Use --no-create to save locally only", options.TicketType, err)
		}
		ticket = createdTicket
		fmt.Printf("✅ Created %s %s in Jira\n", options.TicketType, ticket.Key)
	} else {
		// Generate a temporary key for local-only tickets
		ticket.Key = fmt.Sprintf("LOCAL-%s-%d", options.TicketType, time.Now().Unix())
		fmt.Println("💾 Saving locally only")
	}

	// Save locally
	if err := storageInstance.SaveTicket(ticket); err != nil {
		return fmt.Errorf("failed to save %s locally: %v", options.TicketType, err)
	}

	// Update context
	if err := contextManager.SetFocus(ticket.Key, ticket.Type); err != nil {
		fmt.Printf("⚠️  Warning: Failed to set focus: %v\n", err)
	}

	// Add to recent tickets
	if err := contextManager.AddToRecent(ticket.Key); err != nil {
		fmt.Printf("⚠️  Warning: Failed to add to recent tickets: %v\n", err)
	}

	// Clean up temp file
	os.Remove(tempFile)

	// Success message
	fmt.Printf("✅ %s\n", options.SuccessMessage)
	fmt.Printf("🎯 Focused on: %s\n", ticket.Key)
	if options.ParentInfo != "" {
		fmt.Printf("%s\n", options.ParentInfo)
	}

	return nil
}

// CreateFlags holds the common flags for ticket creation
type CreateFlags struct {
	NoEnrich bool
	NoCreate bool
	Orphan   bool
}

// ValidateEpicContext validates context for epic creation (no validation needed)
func ValidateEpicContext(contextManager *storage.ContextManager) (string, error) {
	return "", nil // Epics don't need context validation
}

// ValidateTaskContext validates context for task creation
func ValidateTaskContext(contextManager *storage.ContextManager) (string, error) {
	currentEpic, err := contextManager.GetCurrentEpic()
	if err != nil {
		return "", fmt.Errorf("failed to get current context: %v", err)
	}

	if currentEpic == "" {
		return "", fmt.Errorf("no epic context found\n💡 Use 'jit track <epic>' to set an epic context, or use --orphan to create an orphan task")
	}

	return currentEpic, nil
}

// ValidateSubtaskContext validates context for subtask creation
func ValidateSubtaskContext(contextManager *storage.ContextManager) (string, error) {
	currentTask, err := contextManager.GetCurrentTask()
	if err != nil {
		return "", fmt.Errorf("failed to get current context: %v", err)
	}

	if currentTask == "" {
		return "", fmt.Errorf("no task context found\n💡 Use 'jit focus <task>' to set a task context")
	}

	return currentTask, nil
}

// SetEpicRelationships sets relationships for epics (no parent)
func SetEpicRelationships(ticket *types.Ticket, context string) {
	// Epics have no parent relationships
}

// SetTaskRelationships sets relationships for tasks
func SetTaskRelationships(ticket *types.Ticket, context string) {
	if context != "" {
		ticket.Relationships.ParentKey = context
	}
}

// SetSubtaskRelationships sets relationships for subtasks
func SetSubtaskRelationships(ticket *types.Ticket, context string) {
	if context != "" {
		ticket.Relationships.ParentKey = context
	}
}
