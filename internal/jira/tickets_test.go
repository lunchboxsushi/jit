package jira

import (
	"context"
	"testing"
	"time"

	"github.com/lunchboxsushi/jit/pkg/types"
)

func TestNewTicketService(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	if service == nil {
		t.Fatal("TicketService should not be nil")
	}

	if service.client != client {
		t.Error("TicketService should have the correct client")
	}
}

func TestConvertIssueType(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	tests := []struct {
		jiraType string
		expected string
	}{
		{"Epic", types.TicketTypeEpic},
		{"epic", types.TicketTypeEpic},
		{"Story", types.TicketTypeTask},
		{"Task", types.TicketTypeTask},
		{"Sub-task", types.TicketTypeSubtask},
		{"Subtask", types.TicketTypeSubtask},
		{"Bug", types.TicketTypeTask},     // Default to task
		{"Unknown", types.TicketTypeTask}, // Default to task
	}

	for _, tt := range tests {
		result := service.convertIssueType(tt.jiraType)
		if result != tt.expected {
			t.Errorf("convertIssueType(%s) = %s, expected %s", tt.jiraType, result, tt.expected)
		}
	}
}

func TestGetIssueTypeID(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	tests := []struct {
		ticketType string
		expected   string
		shouldErr  bool
	}{
		{types.TicketTypeEpic, "10000", false},
		{types.TicketTypeTask, "10001", false},
		{types.TicketTypeSubtask, "10003", false},
		{"Unknown", "", true},
	}

	for _, tt := range tests {
		result, err := service.getIssueTypeID(tt.ticketType)
		if tt.shouldErr && err == nil {
			t.Errorf("getIssueTypeID(%s) should have returned error", tt.ticketType)
		}
		if !tt.shouldErr && err != nil {
			t.Errorf("getIssueTypeID(%s) returned unexpected error: %v", tt.ticketType, err)
		}
		if !tt.shouldErr && result != tt.expected {
			t.Errorf("getIssueTypeID(%s) = %s, expected %s", tt.ticketType, result, tt.expected)
		}
	}
}

func TestGetPriorityID(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	tests := []struct {
		priority string
		expected string
	}{
		{"Highest", "1"},
		{"High", "2"},
		{"Medium", "3"},
		{"Low", "4"},
		{"Lowest", "5"},
		{"Unknown", "3"}, // Default to medium
		{"", "3"},        // Default to medium
	}

	for _, tt := range tests {
		result, err := service.getPriorityID(tt.priority)
		if err != nil {
			t.Errorf("getPriorityID(%s) returned unexpected error: %v", tt.priority, err)
		}
		if result != tt.expected {
			t.Errorf("getPriorityID(%s) = %s, expected %s", tt.priority, result, tt.expected)
		}
	}
}

func TestConvertJiraIssueToTicket(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	// Create a test Jira issue
	jiraIssue := &JiraIssue{
		Key: "TEST-123",
		ID:  "12345",
		Fields: JiraIssueFields{
			Summary:     "Test Issue",
			Description: "This is a test issue",
			Status: JiraStatus{
				ID:   "10000",
				Name: "In Progress",
			},
			Priority: JiraPriority{
				ID:   "2",
				Name: "High",
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
			Assignee: &JiraUser{
				AccountID:   "12345",
				DisplayName: "Test User",
				Email:       "test@example.com",
			},
			Created: time.Now(),
			Updated: time.Now(),
			Labels:  []string{"test", "example"},
		},
	}

	// Convert to our ticket format
	ticket := service.convertJiraIssueToTicket(jiraIssue)

	// Verify conversion
	if ticket.Key != "TEST-123" {
		t.Errorf("Expected key TEST-123, got %s", ticket.Key)
	}

	if ticket.Title != "Test Issue" {
		t.Errorf("Expected title 'Test Issue', got %s", ticket.Title)
	}

	if ticket.Description != "This is a test issue" {
		t.Errorf("Expected description 'This is a test issue', got %s", ticket.Description)
	}

	if ticket.Status != "In Progress" {
		t.Errorf("Expected status 'In Progress', got %s", ticket.Status)
	}

	if ticket.Priority != "High" {
		t.Errorf("Expected priority 'High', got %s", ticket.Priority)
	}

	if ticket.Type != types.TicketTypeTask {
		t.Errorf("Expected type %s, got %s", types.TicketTypeTask, ticket.Type)
	}

	if ticket.Metadata.Project != "TEST" {
		t.Errorf("Expected project TEST, got %s", ticket.Metadata.Project)
	}

	if ticket.Metadata.Assignee != "test@example.com" {
		t.Errorf("Expected assignee test@example.com, got %s", ticket.Metadata.Assignee)
	}

	if len(ticket.Metadata.Labels) != 2 {
		t.Errorf("Expected 2 labels, got %d", len(ticket.Metadata.Labels))
	}

	if ticket.JiraData.URL != "https://test.atlassian.net/browse/TEST-123" {
		t.Errorf("Expected URL https://test.atlassian.net/browse/TEST-123, got %s", ticket.JiraData.URL)
	}
}

func TestConvertJiraIssueToTicketWithEpic(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	// Create a test Jira epic
	jiraIssue := &JiraIssue{
		Key: "TEST-100",
		Fields: JiraIssueFields{
			Summary: "Test Epic",
			IssueType: JiraIssueType{
				Name: "Epic",
			},
			Project: JiraProject{
				Key: "TEST",
			},
		},
	}

	// Convert to our ticket format
	ticket := service.convertJiraIssueToTicket(jiraIssue)

	// Verify it's converted to an epic
	if ticket.Type != types.TicketTypeEpic {
		t.Errorf("Expected type %s, got %s", types.TicketTypeEpic, ticket.Type)
	}
}

func TestConvertJiraIssueToTicketWithSubtask(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	// Create a test Jira subtask
	jiraIssue := &JiraIssue{
		Key: "TEST-101",
		Fields: JiraIssueFields{
			Summary: "Test Subtask",
			IssueType: JiraIssueType{
				Name: "Sub-task",
			},
			Project: JiraProject{
				Key: "TEST",
			},
		},
	}

	// Convert to our ticket format
	ticket := service.convertJiraIssueToTicket(jiraIssue)

	// Verify it's converted to a subtask
	if ticket.Type != types.TicketTypeSubtask {
		t.Errorf("Expected type %s, got %s", types.TicketTypeSubtask, ticket.Type)
	}
}

func TestGetEpicChildren(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	// This test would require mocking the SearchTickets method
	// For now, we'll just test that the method exists and doesn't panic
	ctx := context.Background()
	_, err := service.GetEpicChildren(ctx, "TEST-100")
	// We expect an error since we're not mocking the HTTP client
	if err == nil {
		t.Error("Expected error when calling GetEpicChildren without mocked client")
	}
}

func TestGetTaskSubtasks(t *testing.T) {
	config := &types.JiraConfig{
		URL:           "https://test.atlassian.net",
		Username:      "test@example.com",
		Token:         "test-token",
		Project:       "TEST",
		EpicLinkField: "customfield_10014",
	}

	client := NewClient(config)
	service := NewTicketService(client)

	// This test would require mocking the SearchTickets method
	// For now, we'll just test that the method exists and doesn't panic
	ctx := context.Background()
	_, err := service.GetTaskSubtasks(ctx, "TEST-101")
	// We expect an error since we're not mocking the HTTP client
	if err == nil {
		t.Error("Expected error when calling GetTaskSubtasks without mocked client")
	}
}
