package ai

import (
	"testing"
)

func TestMockProvider(t *testing.T) {
	config := &Config{
		Provider: "mock",
		Model:    "test",
	}

	provider, err := NewProvider(config)
	if err != nil {
		t.Fatalf("Failed to create mock provider: %v", err)
	}

	if provider.Name() != "mock" {
		t.Errorf("Expected provider name 'mock', got '%s'", provider.Name())
	}

	context := &EnrichmentContext{
		TicketType: "epic",
		Project:    "TEST",
	}

	enriched, err := provider.Enrich("Test content", context)
	if err != nil {
		t.Fatalf("Failed to enrich content: %v", err)
	}

	if enriched == "" {
		t.Error("Enriched content should not be empty")
	}

	// Check that the enriched content contains expected elements
	if !contains(enriched, "AI ENRICHED") {
		t.Error("Enriched content should contain 'AI ENRICHED'")
	}

	if !contains(enriched, "epic") {
		t.Error("Enriched content should contain ticket type")
	}

	if !contains(enriched, "TEST") {
		t.Error("Enriched content should contain project name")
	}
}

func TestTemplateManager(t *testing.T) {
	tm := NewTemplateManager("./templates/ai")

	// Create default templates
	err := tm.CreateDefaultTemplates()
	if err != nil {
		t.Fatalf("Failed to create default templates: %v", err)
	}

	// Test processing a template
	context := &EnrichmentContext{
		TicketType: "epic",
		Project:    "TEST",
	}

	prompt, err := tm.GetEnrichmentPrompt("Test content", context)
	if err != nil {
		t.Fatalf("Failed to generate enrichment prompt: %v", err)
	}

	if prompt == "" {
		t.Error("Generated prompt should not be empty")
	}

	// Check that the prompt contains expected elements
	if !contains(prompt, "Test content") {
		t.Error("Prompt should contain original content")
	}

	if !contains(prompt, "EPIC") {
		t.Error("Prompt should contain ticket type (EPIC)")
	}

	if !contains(prompt, "TEST") {
		t.Error("Prompt should contain project name")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
