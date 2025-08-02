package jira

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/lunchboxsushi/jit/pkg/types"
)

func TestNewClient(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)

	if client == nil {
		t.Fatal("Client should not be nil")
	}

	if client.baseURL != config.URL {
		t.Errorf("Expected baseURL %s, got %s", config.URL, client.baseURL)
	}

	if client.username != config.Username {
		t.Errorf("Expected username %s, got %s", config.Username, client.username)
	}

	if client.token != config.Token {
		t.Errorf("Expected token %s, got %s", config.Token, client.token)
	}
}

func TestGetIssue(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check authentication header
		auth := r.Header.Get("Authorization")
		if auth != "Basic dGVzdEBleGFtcGxlLmNvbTp0ZXN0LXRva2Vu" { // test@example.com:test-token
			t.Errorf("Expected Basic auth header, got %s", auth)
		}

		// Check endpoint
		if r.URL.Path != "/rest/api/3/issue/TEST-123" {
			t.Errorf("Expected path /rest/api/3/issue/TEST-123, got %s", r.URL.Path)
		}

		// Return mock issue
		issue := JiraIssue{
			Key: "TEST-123",
			ID:  "12345",
			Fields: JiraIssueFields{
				Summary:     "Test Issue",
				Description: "This is a test issue",
				Status: JiraStatus{
					ID:   "10000",
					Name: "To Do",
				},
				Priority: JiraPriority{
					ID:   "3",
					Name: "Medium",
				},
				IssueType: JiraIssueType{
					ID:   "10001",
					Name: "Task",
				},
				Project: JiraProject{
					ID:   "10000",
					Key:  "TEST",
					Name: "Test Project",
				},
				Created: time.Now(),
				Updated: time.Now(),
				Labels:  []string{"test", "example"},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(issue)
	}))
	defer server.Close()

	// Create client with mock server URL
	config := &types.JiraConfig{
		URL:           server.URL,
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}
	client := NewClient(config)

	// Test GetIssue
	ctx := context.Background()
	issue, err := client.GetIssue(ctx, "TEST-123")
	if err != nil {
		t.Fatalf("Failed to get issue: %v", err)
	}

	if issue.Key != "TEST-123" {
		t.Errorf("Expected key TEST-123, got %s", issue.Key)
	}

	if issue.Fields.Summary != "Test Issue" {
		t.Errorf("Expected summary 'Test Issue', got %s", issue.Fields.Summary)
	}

	if issue.Fields.Status.Name != "To Do" {
		t.Errorf("Expected status 'To Do', got %s", issue.Fields.Status.Name)
	}
}

func TestCreateIssue(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check endpoint
		if r.URL.Path != "/rest/api/3/issue" {
			t.Errorf("Expected path /rest/api/3/issue, got %s", r.URL.Path)
		}

		// Parse request body
		var request JiraCreateIssueRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Check request fields
		if request.Fields.Summary != "New Test Issue" {
			t.Errorf("Expected summary 'New Test Issue', got %s", request.Fields.Summary)
		}

		if request.Fields.Project.Key != "TEST" {
			t.Errorf("Expected project TEST, got %s", request.Fields.Project.Key)
		}

		// Return mock response
		response := JiraCreateIssueResponse{
			ID:   "12345",
			Key:  "TEST-456",
			Self: "https://test.atlassian.net/rest/api/3/issue/12345",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server URL
	config := &types.JiraConfig{
		URL:           server.URL,
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}
	client := NewClient(config)

	// Test CreateIssue
	ctx := context.Background()
	request := &JiraCreateIssueRequest{
		Fields: JiraCreateIssueFields{
			Project:     JiraProjectReference{Key: "TEST"},
			Summary:     "New Test Issue",
			Description: "This is a new test issue",
			IssueType:   JiraIssueTypeRef{ID: "10001"},
		},
	}

	response, err := client.CreateIssue(ctx, request)
	if err != nil {
		t.Fatalf("Failed to create issue: %v", err)
	}

	if response.Key != "TEST-456" {
		t.Errorf("Expected key TEST-456, got %s", response.Key)
	}

	if response.ID != "12345" {
		t.Errorf("Expected ID 12345, got %s", response.ID)
	}
}

func TestAddComment(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		// Check endpoint
		expectedPath := "/rest/api/3/issue/TEST-123/comment"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Parse request body
		var request JiraCreateCommentRequest
		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			t.Errorf("Failed to decode request: %v", err)
		}

		// Check comment body
		if request.Body != "Test comment" {
			t.Errorf("Expected comment 'Test comment', got %s", request.Body)
		}

		// Return mock response
		comment := JiraComment{
			ID:      "10000",
			Body:    "Test comment",
			Created: time.Now(),
			Updated: time.Now(),
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(comment)
	}))
	defer server.Close()

	// Create client with mock server URL
	config := &types.JiraConfig{
		URL:           server.URL,
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}
	client := NewClient(config)

	// Test AddComment
	ctx := context.Background()
	comment, err := client.AddComment(ctx, "TEST-123", "Test comment")
	if err != nil {
		t.Fatalf("Failed to add comment: %v", err)
	}

	if comment.Body != "Test comment" {
		t.Errorf("Expected comment body 'Test comment', got %s", comment.Body)
	}
}

func TestSearchIssues(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		// Check endpoint
		if r.URL.Path != "/rest/api/3/search" {
			t.Errorf("Expected path /rest/api/3/search, got %s", r.URL.Path)
		}

		// Check query parameters
		jql := r.URL.Query().Get("jql")
		if jql != "project = TEST" {
			t.Errorf("Expected JQL 'project = TEST', got %s", jql)
		}

		maxResults := r.URL.Query().Get("maxResults")
		if maxResults != "10" {
			t.Errorf("Expected maxResults 10, got %s", maxResults)
		}

		// Return mock response
		response := JiraSearchResponse{
			StartAt:    0,
			MaxResults: 10,
			Total:      1,
			Issues: []JiraIssue{
				{
					Key: "TEST-123",
					Fields: JiraIssueFields{
						Summary:   "Test Issue",
						Status:    JiraStatus{Name: "To Do"},
						IssueType: JiraIssueType{Name: "Task"},
						Project:   JiraProject{Key: "TEST"},
					},
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create client with mock server URL
	config := &types.JiraConfig{
		URL:           server.URL,
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}
	client := NewClient(config)

	// Test SearchIssues
	ctx := context.Background()
	response, err := client.SearchIssues(ctx, "project = TEST", 10)
	if err != nil {
		t.Fatalf("Failed to search issues: %v", err)
	}

	if response.Total != 1 {
		t.Errorf("Expected 1 issue, got %d", response.Total)
	}

	if len(response.Issues) != 1 {
		t.Errorf("Expected 1 issue in response, got %d", len(response.Issues))
	}

	if response.Issues[0].Key != "TEST-123" {
		t.Errorf("Expected issue key TEST-123, got %s", response.Issues[0].Key)
	}
}

func TestTestConnection(t *testing.T) {
	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check endpoint
		if r.URL.Path != "/rest/api/3/myself" {
			t.Errorf("Expected path /rest/api/3/myself, got %s", r.URL.Path)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}))
	defer server.Close()

	// Create client with mock server URL
	config := &types.JiraConfig{
		URL:           server.URL,
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}
	client := NewClient(config)

	// Test TestConnection
	ctx := context.Background()
	err := client.TestConnection(ctx)
	if err != nil {
		t.Fatalf("Failed to test connection: %v", err)
	}
}

func TestErrorHandling(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Return Jira error response
		errorResponse := JiraErrorResponse{
			ErrorCollection: JiraError{
				ErrorMessages: []string{"Issue does not exist"},
				Errors: map[string]string{
					"issueKey": "TEST-999 does not exist",
				},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errorResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	config := &types.JiraConfig{
		URL:           server.URL,
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}
	client := NewClient(config)

	// Test error handling
	ctx := context.Background()
	_, err := client.GetIssue(ctx, "TEST-999")
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	// Check error message
	if !contains(err.Error(), "Issue does not exist") {
		t.Errorf("Expected error message to contain 'Issue does not exist', got: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
