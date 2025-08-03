package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIProvider implements the Provider interface for OpenAI
type OpenAIProvider struct {
	config      *Config
	client      *http.Client
	templateMgr *TemplateManager
}

// OpenAIRequest represents the request structure for OpenAI API
type OpenAIRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
}

// OpenAIResponse represents the response structure from OpenAI API
type OpenAIResponse struct {
	Choices []Choice `json:"choices"`
	Error   *Error   `json:"error,omitempty"`
}

// Message represents a message in the OpenAI conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Choice represents a choice in the OpenAI response
type Choice struct {
	Message Message `json:"message"`
}

// Error represents an error from OpenAI API
type Error struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(config *Config) (Provider, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required")
	}

	// Set defaults if not provided
	if config.Model == "" {
		config.Model = "gpt-3.5-turbo"
	}
	if config.MaxTokens == 0 {
		config.MaxTokens = 1000
	}
	if config.Temperature == 0 {
		config.Temperature = 0.7
	}

	// Determine base URL
	baseURL := config.BaseURL
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Create template manager
	templateMgr := NewTemplateManager("./templates/ai")

	// Create default templates
	if err := templateMgr.CreateDefaultTemplates(); err != nil {
		return nil, fmt.Errorf("failed to create default templates: %v", err)
	}

	return &OpenAIProvider{
		config:      config,
		client:      client,
		templateMgr: templateMgr,
	}, nil
}

// Name returns the provider name
func (o *OpenAIProvider) Name() string {
	return "openai"
}

// Enrich enriches content using OpenAI API
func (o *OpenAIProvider) Enrich(content string, context *EnrichmentContext) (string, error) {
	// Generate prompt using template
	prompt, err := o.templateMgr.GetEnrichmentPrompt(content, context)
	if err != nil {
		return "", fmt.Errorf("failed to generate prompt: %v", err)
	}

	// Create request
	request := OpenAIRequest{
		Model:       o.config.Model,
		MaxTokens:   o.config.MaxTokens,
		Temperature: o.config.Temperature,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are an expert software development assistant. Provide clear, professional, and actionable responses.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Make API request
	response, err := o.makeRequest(request)
	if err != nil {
		return "", fmt.Errorf("OpenAI API request failed: %v", err)
	}

	// Extract response content
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}

// makeRequest makes a request to the OpenAI API
func (o *OpenAIProvider) makeRequest(request OpenAIRequest) (*OpenAIResponse, error) {
	// Determine endpoint
	endpoint := "/chat/completions"
	if o.config.BaseURL != "" {
		endpoint = o.config.BaseURL + endpoint
	} else {
		endpoint = "https://api.openai.com/v1" + endpoint
	}

	// Marshal request
	requestBody, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.config.APIKey)

	// Make request
	resp, err := o.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("OpenAI API error: %s - %s", resp.Status, string(body))
	}

	// Parse response
	var response OpenAIResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %v", err)
	}

	// Check for API errors
	if response.Error != nil {
		return nil, fmt.Errorf("OpenAI API error: %s", response.Error.Message)
	}

	return &response, nil
}

// EnrichComment enriches a comment using OpenAI
func (o *OpenAIProvider) EnrichComment(comment string, context *EnrichmentContext) (string, error) {
	// Generate comment prompt
	prompt, err := o.templateMgr.GetCommentPrompt(comment, context)
	if err != nil {
		return "", fmt.Errorf("failed to generate comment prompt: %v", err)
	}

	// Create request
	request := OpenAIRequest{
		Model:       o.config.Model,
		MaxTokens:   o.config.MaxTokens,
		Temperature: o.config.Temperature,
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are an expert software development assistant. Provide clear, professional, and helpful comments.",
			},
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	// Make API request
	response, err := o.makeRequest(request)
	if err != nil {
		return "", fmt.Errorf("OpenAI API request failed: %v", err)
	}

	// Extract response content
	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response choices from OpenAI")
	}

	return response.Choices[0].Message.Content, nil
}
