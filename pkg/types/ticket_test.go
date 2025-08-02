package types

import (
	"encoding/json"
	"testing"
)

func TestTicketJSONMarshalUnmarshal(t *testing.T) {
	// Create a test ticket
	ticket := NewTicket("TEST-123", "Test Ticket", TicketTypeTask)
	ticket.Description = "This is a test ticket"
	ticket.Status = "In Progress"
	ticket.Priority = "High"
	ticket.Metadata.Project = "TEST"
	ticket.Metadata.Assignee = "test@example.com"
	ticket.Metadata.Labels = []string{"test", "example"}
	ticket.Relationships.ParentKey = "TEST-100"
	ticket.JiraData.URL = "https://example.atlassian.net/browse/TEST-123"
	ticket.JiraData.CustomFields["story_points"] = 5

	// Marshal to JSON
	jsonData, err := json.Marshal(ticket)
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}

	// Unmarshal back
	var unmarshaledTicket Ticket
	err = json.Unmarshal(jsonData, &unmarshaledTicket)
	if err != nil {
		t.Fatalf("Failed to unmarshal ticket: %v", err)
	}

	// Verify key fields
	if unmarshaledTicket.Key != ticket.Key {
		t.Errorf("Key mismatch: got %s, want %s", unmarshaledTicket.Key, ticket.Key)
	}
	if unmarshaledTicket.Title != ticket.Title {
		t.Errorf("Title mismatch: got %s, want %s", unmarshaledTicket.Title, ticket.Title)
	}
	if unmarshaledTicket.Type != ticket.Type {
		t.Errorf("Type mismatch: got %s, want %s", unmarshaledTicket.Type, ticket.Type)
	}
	if unmarshaledTicket.Status != ticket.Status {
		t.Errorf("Status mismatch: got %s, want %s", unmarshaledTicket.Status, ticket.Status)
	}
	if unmarshaledTicket.Description != ticket.Description {
		t.Errorf("Description mismatch: got %s, want %s", unmarshaledTicket.Description, ticket.Description)
	}
	if unmarshaledTicket.Metadata.Project != ticket.Metadata.Project {
		t.Errorf("Project mismatch: got %s, want %s", unmarshaledTicket.Metadata.Project, ticket.Metadata.Project)
	}
	if unmarshaledTicket.Relationships.ParentKey != ticket.Relationships.ParentKey {
		t.Errorf("ParentKey mismatch: got %s, want %s", unmarshaledTicket.Relationships.ParentKey, ticket.Relationships.ParentKey)
	}
	if unmarshaledTicket.JiraData.URL != ticket.JiraData.URL {
		t.Errorf("URL mismatch: got %s, want %s", unmarshaledTicket.JiraData.URL, ticket.JiraData.URL)
	}
}

func TestNewTicket(t *testing.T) {
	ticket := NewTicket("TEST-456", "New Test Ticket", TicketTypeEpic)

	if ticket.Key != "TEST-456" {
		t.Errorf("Expected key TEST-456, got %s", ticket.Key)
	}
	if ticket.Title != "New Test Ticket" {
		t.Errorf("Expected title 'New Test Ticket', got %s", ticket.Title)
	}
	if ticket.Type != TicketTypeEpic {
		t.Errorf("Expected type %s, got %s", TicketTypeEpic, ticket.Type)
	}
	if ticket.Status != "To Do" {
		t.Errorf("Expected status 'To Do', got %s", ticket.Status)
	}
	if ticket.Priority != "Medium" {
		t.Errorf("Expected priority Medium, got %s", ticket.Priority)
	}
	if len(ticket.Metadata.Labels) != 0 {
		t.Errorf("Expected empty labels, got %v", ticket.Metadata.Labels)
	}
	if len(ticket.Relationships.Children) != 0 {
		t.Errorf("Expected empty children, got %v", ticket.Relationships.Children)
	}
}

func TestTicketTypeHelpers(t *testing.T) {
	epic := NewTicket("TEST-100", "Test Epic", TicketTypeEpic)
	task := NewTicket("TEST-101", "Test Task", TicketTypeTask)
	subtask := NewTicket("TEST-102", "Test Subtask", TicketTypeSubtask)

	// Test IsEpic
	if !epic.IsEpic() {
		t.Errorf("Expected epic.IsEpic() to be true")
	}
	if task.IsEpic() {
		t.Errorf("Expected task.IsEpic() to be false")
	}
	if subtask.IsEpic() {
		t.Errorf("Expected subtask.IsEpic() to be false")
	}

	// Test IsTask
	if epic.IsTask() {
		t.Errorf("Expected epic.IsTask() to be false")
	}
	if !task.IsTask() {
		t.Errorf("Expected task.IsTask() to be true")
	}
	if subtask.IsTask() {
		t.Errorf("Expected subtask.IsTask() to be false")
	}

	// Test IsSubtask
	if epic.IsSubtask() {
		t.Errorf("Expected epic.IsSubtask() to be false")
	}
	if task.IsSubtask() {
		t.Errorf("Expected task.IsSubtask() to be false")
	}
	if !subtask.IsSubtask() {
		t.Errorf("Expected subtask.IsSubtask() to be true")
	}

	// Test IsOrphanTask
	if epic.IsOrphanTask() {
		t.Errorf("Expected epic.IsOrphanTask() to be false")
	}
	if !task.IsOrphanTask() {
		t.Errorf("Expected task.IsOrphanTask() to be true")
	}
	if subtask.IsOrphanTask() {
		t.Errorf("Expected subtask.IsOrphanTask() to be false")
	}

	// Test orphan task with parent
	taskWithParent := NewTicket("TEST-103", "Task with Parent", TicketTypeTask)
	taskWithParent.Relationships.ParentKey = "TEST-100"
	if taskWithParent.IsOrphanTask() {
		t.Errorf("Expected task with parent to not be orphan")
	}
}
