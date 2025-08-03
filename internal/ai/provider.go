package ai

import (
	"fmt"
	"strings"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// Provider defines the interface for AI enrichment services
type Provider interface {
	// Enrich takes raw content and context, returns enriched content
	Enrich(content string, context *EnrichmentContext) (string, error)

	// Name returns the provider name
	Name() string
}

// EnrichmentContext provides context for AI enrichment
type EnrichmentContext struct {
	TicketType   string
	Project      string
	CurrentEpic  string
	CurrentTask  string
	UserEmail    string
	CustomFields map[string]interface{}
}

// Config holds AI provider configuration
type Config struct {
	Provider    string
	Model       string
	MaxTokens   int
	Temperature float64
	APIKey      string
	BaseURL     string
}

// NewProvider creates a new AI provider based on configuration
func NewProvider(config *Config) (Provider, error) {
	switch strings.ToLower(config.Provider) {
	case "openai":
		return NewOpenAIProvider(config)
	case "mock", "test":
		return NewMockProvider(config)
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", config.Provider)
	}
}

// MockProvider is a test provider that returns predictable responses
type MockProvider struct {
	config *Config
}

// NewMockProvider creates a new mock provider for testing
func NewMockProvider(config *Config) (Provider, error) {
	return &MockProvider{
		config: config,
	}, nil
}

// Name returns the provider name
func (m *MockProvider) Name() string {
	return "mock"
}

// Enrich returns a mock enriched response
func (m *MockProvider) Enrich(content string, context *EnrichmentContext) (string, error) {
	// Return a predictable mock response for testing
	enriched := fmt.Sprintf("[AI ENRICHED] %s\n\nEnhanced with AI assistance for %s ticket in project %s.",
		content, context.TicketType, context.Project)

	if context.CurrentEpic != "" {
		enriched += fmt.Sprintf("\n\nRelated to epic: %s", context.CurrentEpic)
	}

	return enriched, nil
}

// EnrichTicket enriches a ticket's description using the provided AI provider
func EnrichTicket(provider Provider, ticket *types.Ticket, context *EnrichmentContext) error {
	if provider == nil {
		return fmt.Errorf("no AI provider configured")
	}

	// Enrich the description
	enrichedDescription, err := provider.Enrich(ticket.Description, context)
	if err != nil {
		return fmt.Errorf("failed to enrich ticket description: %v", err)
	}

	// Update the ticket
	ticket.Description = enrichedDescription
	ticket.LocalData.AIEnhanced = true

	return nil
}

// EnrichComment enriches a comment using the provided AI provider
func EnrichComment(provider Provider, comment string, context *EnrichmentContext) (string, error) {
	if provider == nil {
		return comment, nil // Return original comment if no provider
	}

	enrichedComment, err := provider.Enrich(comment, context)
	if err != nil {
		return comment, fmt.Errorf("failed to enrich comment: %v", err)
	}

	return enrichedComment, nil
}
