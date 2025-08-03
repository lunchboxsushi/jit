package commands

import (
	"fmt"
	"os"
	"time"

	"github.com/lunchboxsushi/jit/internal/storage"
	"github.com/lunchboxsushi/jit/internal/ui"
	"github.com/lunchboxsushi/jit/pkg/types"
	"github.com/spf13/cobra"
)

// CreateOptions defines the configuration for ticket creation
type CreateOptions struct {
	TicketType       string
	TemplateName     string
	ValidateContext  func(contextManager *storage.ContextManager, flags CreateFlags) (string, error)
	SetRelationships func(ticket *types.Ticket, context string)
	SuccessMessage   string
	ParentInfo       string
}

// CreateTicket creates a new ticket using the shared pipeline
func CreateTicket(cmd *cobra.Command, options CreateOptions, flags CreateFlags) error {
	// Initialize command context
	ctx, err := InitializeCommand()
	if err != nil {
		return fmt.Errorf("failed to initialize: %v", err)
	}

	// Validate context
	contextInfo, err := options.ValidateContext(ctx.ContextManager, flags)
	if err != nil {
		return fmt.Errorf("context validation failed: %v", err)
	}

	// Create temporary file for editing
	tempFile, err := os.CreateTemp("", "jit-*.md")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Open editor with template
	editor := ui.NewEditor()
	if err := editor.EditTemplate(options.TemplateName, tempFile.Name()); err != nil {
		return fmt.Errorf("failed to open editor: %v", err)
	}

	// Read and parse the content
	content, err := editor.ReadFile(tempFile.Name())
	if err != nil {
		return fmt.Errorf("failed to read file: %v", err)
	}

	// Parse markdown content
	title, description, err := editor.ParseMarkdownTicket(content)
	if err != nil {
		return fmt.Errorf("failed to parse markdown: %v", err)
	}

	// Create ticket object
	ticket := &types.Ticket{
		Title:       title,
		Description: description,
		Type:        options.TicketType,
		Status:      "To Do",
		Priority:    "Medium",
		Metadata: types.TicketMetadata{
			Project: ctx.Config.Jira.Project,
			Created: time.Now(),
			Updated: time.Now(),
			Labels:  []string{},
		},
		Relationships: types.TicketRelationships{
			Children: []string{},
		},
		JiraData: types.JiraData{
			CustomFields: make(map[string]interface{}),
		},
		LocalData: types.LocalData{
			LastSync:     time.Now(),
			LocalChanges: false,
			AIEnhanced:   false,
		},
	}

	// Set assignee
	ticket.Metadata.Assignee = ctx.Config.Jira.Username

	// Set relationships
	if options.SetRelationships != nil {
		options.SetRelationships(ticket, contextInfo)
	}

	// Create in Jira if requested
	if !flags.NoCreate {
		if err := ctx.CreateTicketInJira(ticket, options.TicketType); err != nil {
			return err
		}
	} else {
		if err := ctx.SaveTicketLocally(ticket, options.TicketType); err != nil {
			return err
		}
	}

	// Update context and recent tickets
	if err := ctx.UpdateContextAndRecent(ticket.Key, ticket.Type); err != nil {
		// Context update errors are non-fatal
	}

	// Success message
	PrintSuccess(options.SuccessMessage)
	fmt.Printf("Focused on: %s\n", ticket.Key)
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
func ValidateEpicContext(contextManager *storage.ContextManager, flags CreateFlags) (string, error) {
	return "", nil // Epics don't need context validation
}

// ValidateTaskContext validates context for task creation
func ValidateTaskContext(contextManager *storage.ContextManager, flags CreateFlags) (string, error) {
	// If orphan flag is set, skip context validation
	if flags.Orphan {
		return "", nil
	}

	currentEpic, err := contextManager.GetCurrentEpic()
	if err != nil {
		return "", fmt.Errorf("failed to get current context: %v", err)
	}

	if currentEpic == "" {
		return "", fmt.Errorf("no epic context found\nTip: Use 'jit focus <epic>' to set an epic context, or use --orphan to create an orphan task")
	}

	return currentEpic, nil
}

// ValidateSubtaskContext validates context for subtask creation
func ValidateSubtaskContext(contextManager *storage.ContextManager, flags CreateFlags) (string, error) {
	currentTask, err := contextManager.GetCurrentTask()
	if err != nil {
		return "", fmt.Errorf("failed to get current context: %v", err)
	}

	if currentTask == "" {
		return "", fmt.Errorf("no task context found\nTip: Use 'jit focus <task>' to set a task context")
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
