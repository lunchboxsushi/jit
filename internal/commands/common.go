package commands

import (
	"context"
	"fmt"
	"time"

	"github.com/lunchboxsushi/jit/internal/ai"
	"github.com/lunchboxsushi/jit/internal/config"
	"github.com/lunchboxsushi/jit/internal/jira"
	"github.com/lunchboxsushi/jit/internal/storage"
	"github.com/lunchboxsushi/jit/pkg/types"
)

// CommandContext holds the common context for commands
type CommandContext struct {
	Config         *types.Config
	Storage        storage.Storage
	JiraClient     *jira.Client
	TicketService  *jira.TicketService
	ContextManager *storage.ContextManager
	AIProvider     ai.Provider
}

// InitializeCommand sets up the common context for commands
func InitializeCommand() (*CommandContext, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("configuration error: %v\nRun 'jit init' to create a configuration file", err)
	}

	// Initialize storage
	storageInstance, err := storage.NewJSONStorage(cfg.App.DataDir)
	if err != nil {
		return nil, fmt.Errorf("storage error: %v", err)
	}

	// Initialize Jira services
	jiraClient := jira.NewClient(&cfg.Jira)
	ticketService := jira.NewTicketService(jiraClient)
	contextManager := storage.NewContextManager(storageInstance)

	// Initialize AI provider
	var aiProvider ai.Provider
	if cfg.AI.Provider != "" && cfg.AI.APIKey != "" {
		aiConfig := &ai.Config{
			Provider:    cfg.AI.Provider,
			Model:       cfg.AI.Model,
			MaxTokens:   cfg.AI.MaxTokens,
			Temperature: 0.7, // Default temperature
			APIKey:      cfg.AI.APIKey,
			BaseURL:     "", // Use default OpenAI URL
		}

		aiProvider, err = ai.NewProvider(aiConfig)
		if err != nil {
			// Non-fatal error - log warning but continue
			PrintWarning(fmt.Sprintf("Failed to initialize AI provider: %v", err))
		}
	}

	return &CommandContext{
		Config:         cfg,
		Storage:        storageInstance,
		JiraClient:     jiraClient,
		TicketService:  ticketService,
		ContextManager: contextManager,
		AIProvider:     aiProvider,
	}, nil
}

// HandleError provides consistent error handling across commands
func HandleError(err error, message string) {
	if err != nil {
		fmt.Printf("Error: %s: %v\n", message, err)
	}
}

// PrintSuccess provides consistent success messaging
func PrintSuccess(message string) {
	fmt.Printf("Success: %s\n", message)
}

// PrintInfo provides consistent info messaging
func PrintInfo(message string) {
	fmt.Printf("Info: %s\n", message)
}

// PrintWarning provides consistent warning messaging
func PrintWarning(message string) {
	fmt.Printf("Warning: %s\n", message)
}

// Common Jira Operations

// CreateTicketInJira creates a ticket in Jira and saves it locally
func (ctx *CommandContext) CreateTicketInJira(ticket *types.Ticket, ticketType string) error {
	fmt.Printf("Creating %s in Jira...\n", ticketType)

	createdTicket, err := ctx.TicketService.CreateTicket(context.Background(), ticket)
	if err != nil {
		return fmt.Errorf("failed to create %s in Jira: %v\nTip: Use --no-create to save locally only", ticketType, err)
	}

	// Save the created ticket locally
	if err := ctx.Storage.SaveTicket(createdTicket); err != nil {
		return fmt.Errorf("failed to save %s locally: %v", ticketType, err)
	}

	fmt.Printf("Created %s %s in Jira\n", ticketType, createdTicket.Key)
	return nil
}

// SaveTicketLocally saves a ticket locally with a temporary key
func (ctx *CommandContext) SaveTicketLocally(ticket *types.Ticket, ticketType string) error {
	// Generate a temporary key for local-only tickets
	ticket.Key = fmt.Sprintf("LOCAL-%s-%d", ticketType, time.Now().Unix())
	fmt.Println("Saving locally only")

	if err := ctx.Storage.SaveTicket(ticket); err != nil {
		return fmt.Errorf("failed to save %s locally: %v", ticketType, err)
	}

	return nil
}

// UpdateContextAndRecent updates the working context and adds to recent tickets
func (ctx *CommandContext) UpdateContextAndRecent(ticketKey, ticketType string) error {
	// Update context
	if err := ctx.ContextManager.SetFocus(ticketKey, ticketType); err != nil {
		PrintWarning(fmt.Sprintf("Failed to set focus: %v", err))
	}

	// Add to recent tickets
	if err := ctx.ContextManager.AddToRecent(ticketKey); err != nil {
		PrintWarning(fmt.Sprintf("Failed to add to recent tickets: %v", err))
	}

	return nil
}

// TrackTicketWithChildren tracks a ticket and optionally its children
func (ctx *CommandContext) TrackTicketWithChildren(ticketKey string, fetchChildren bool) error {
	fmt.Printf("Fetching ticket %s...\n", ticketKey)

	// Check if ticket already exists
	if ctx.Storage.Exists(ticketKey) {
		fmt.Printf("Ticket %s already exists, updating...\n", ticketKey)
	}

	// Fetch the main ticket
	ticket, err := ctx.TicketService.GetTicket(context.Background(), ticketKey)
	if err != nil {
		return fmt.Errorf("failed to fetch ticket %s: %v", ticketKey, err)
	}

	// Save the main ticket
	if err := ctx.Storage.SaveTicket(ticket); err != nil {
		return fmt.Errorf("failed to save ticket %s: %v", ticketKey, err)
	}

	fmt.Printf("Saved %s (%s)\n", ticketKey, ticket.Title)

	// Set focus to the tracked ticket
	if err := ctx.ContextManager.SetFocus(ticketKey, ticket.Type); err != nil {
		return fmt.Errorf("failed to set focus: %v", err)
	}

	// Add to recent tickets
	if err := ctx.ContextManager.AddToRecent(ticketKey); err != nil {
		return fmt.Errorf("failed to add to recent tickets: %v", err)
	}

	// Fetch children if requested and ticket is an epic
	if fetchChildren && ticket.Type == types.TicketTypeEpic {
		if err := ctx.trackEpicChildren(ticketKey); err != nil {
			return fmt.Errorf("failed to track epic children: %v", err)
		}
	}

	// Fetch subtasks if requested and ticket is a task
	if fetchChildren && ticket.Type == types.TicketTypeTask {
		if err := ctx.trackTaskSubtasks(ticketKey); err != nil {
			return fmt.Errorf("failed to track task subtasks: %v", err)
		}
	}

	return nil
}

// trackEpicChildren tracks all children of an epic
func (ctx *CommandContext) trackEpicChildren(epicKey string) error {
	fmt.Printf("Fetching children of epic %s...\n", epicKey)

	children, err := ctx.TicketService.GetEpicChildren(context.Background(), epicKey)
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

		if err := ctx.Storage.SaveTicket(child); err != nil {
			return fmt.Errorf("failed to save child ticket %s: %v", child.Key, err)
		}

		// If this child is a task, fetch its subtasks too
		if child.Type == types.TicketTypeTask {
			if err := ctx.trackTaskSubtasks(child.Key); err != nil {
				return fmt.Errorf("failed to track subtasks of %s: %v", child.Key, err)
			}
		}

		// Small delay to be respectful to the API
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// trackTaskSubtasks tracks all subtasks of a task
func (ctx *CommandContext) trackTaskSubtasks(taskKey string) error {
	fmt.Printf("   Fetching subtasks of task %s...\n", taskKey)

	subtasks, err := ctx.TicketService.GetTaskSubtasks(context.Background(), taskKey)
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

		if err := ctx.Storage.SaveTicket(subtask); err != nil {
			return fmt.Errorf("failed to save subtask %s: %v", subtask.Key, err)
		}

		// Small delay to be respectful to the API
		time.Sleep(100 * time.Millisecond)
	}

	return nil
}

// EnrichTicketWithAI enriches a ticket's description using AI
func (ctx *CommandContext) EnrichTicketWithAI(ticket *types.Ticket) error {
	if ctx.AIProvider == nil {
		return fmt.Errorf("no AI provider configured")
	}

	// Get current context
	currentEpic, _ := ctx.ContextManager.GetCurrentEpic()
	currentTask, _ := ctx.ContextManager.GetCurrentTask()

	// Create enrichment context
	context := &ai.EnrichmentContext{
		TicketType:   ticket.Type,
		Project:      ctx.Config.Jira.Project,
		CurrentEpic:  currentEpic,
		CurrentTask:  currentTask,
		UserEmail:    ctx.Config.Jira.Username,
		CustomFields: make(map[string]interface{}),
	}

	// Enrich the ticket
	if err := ai.EnrichTicket(ctx.AIProvider, ticket, context); err != nil {
		return fmt.Errorf("failed to enrich ticket: %v", err)
	}

	PrintInfo("Ticket enriched with AI assistance")
	return nil
}

// EnrichCommentWithAI enriches a comment using AI
func (ctx *CommandContext) EnrichCommentWithAI(comment string, ticketKey string) (string, error) {
	if ctx.AIProvider == nil {
		return comment, nil // Return original comment if no AI provider
	}

	// Load ticket for context
	ticket, err := ctx.Storage.LoadTicket(ticketKey)
	if err != nil {
		return comment, fmt.Errorf("failed to load ticket for context: %v", err)
	}

	// Get current context
	currentEpic, _ := ctx.ContextManager.GetCurrentEpic()
	currentTask, _ := ctx.ContextManager.GetCurrentTask()

	// Create enrichment context
	context := &ai.EnrichmentContext{
		TicketType:   ticket.Type,
		Project:      ctx.Config.Jira.Project,
		CurrentEpic:  currentEpic,
		CurrentTask:  currentTask,
		UserEmail:    ctx.Config.Jira.Username,
		CustomFields: make(map[string]interface{}),
	}

	// Enrich the comment
	enrichedComment, err := ai.EnrichComment(ctx.AIProvider, comment, context)
	if err != nil {
		return comment, fmt.Errorf("failed to enrich comment: %v", err)
	}

	PrintInfo("Comment enriched with AI assistance")
	return enrichedComment, nil
}
