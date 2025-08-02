package types

import (
	"encoding/json"
	"testing"
)

func TestContextJSONMarshalUnmarshal(t *testing.T) {
	// Create a test context
	context := NewContext()
	context.CurrentEpic = "TEST-100"
	context.CurrentTask = "TEST-101"
	context.RecentTickets = []string{"TEST-100", "TEST-101", "TEST-102"}

	// Marshal to JSON
	jsonData, err := json.Marshal(context)
	if err != nil {
		t.Fatalf("Failed to marshal context: %v", err)
	}

	// Unmarshal back
	var unmarshaledContext Context
	err = json.Unmarshal(jsonData, &unmarshaledContext)
	if err != nil {
		t.Fatalf("Failed to unmarshal context: %v", err)
	}

	// Verify key fields
	if unmarshaledContext.CurrentEpic != context.CurrentEpic {
		t.Errorf("CurrentEpic mismatch: got %s, want %s", unmarshaledContext.CurrentEpic, context.CurrentEpic)
	}
	if unmarshaledContext.CurrentTask != context.CurrentTask {
		t.Errorf("CurrentTask mismatch: got %s, want %s", unmarshaledContext.CurrentTask, context.CurrentTask)
	}
	if len(unmarshaledContext.RecentTickets) != len(context.RecentTickets) {
		t.Errorf("RecentTickets length mismatch: got %d, want %d", len(unmarshaledContext.RecentTickets), len(context.RecentTickets))
	}
}

func TestContextSetFocus(t *testing.T) {
	context := NewContext()

	// Test setting epic focus
	context.SetFocus("TEST-100", TicketTypeEpic)
	if context.CurrentEpic != "TEST-100" {
		t.Errorf("Expected CurrentEpic TEST-100, got %s", context.CurrentEpic)
	}
	if context.CurrentTask != "" {
		t.Errorf("Expected empty CurrentTask, got %s", context.CurrentTask)
	}
	if context.CurrentSubtask != "" {
		t.Errorf("Expected empty CurrentSubtask, got %s", context.CurrentSubtask)
	}

	// Test setting task focus
	context.SetFocus("TEST-101", TicketTypeTask)
	if context.CurrentEpic != "TEST-100" {
		t.Errorf("Expected CurrentEpic to remain TEST-100, got %s", context.CurrentEpic)
	}
	if context.CurrentTask != "TEST-101" {
		t.Errorf("Expected CurrentTask TEST-101, got %s", context.CurrentTask)
	}
	if context.CurrentSubtask != "" {
		t.Errorf("Expected empty CurrentSubtask, got %s", context.CurrentSubtask)
	}

	// Test setting subtask focus
	context.SetFocus("TEST-102", TicketTypeSubtask)
	if context.CurrentEpic != "TEST-100" {
		t.Errorf("Expected CurrentEpic to remain TEST-100, got %s", context.CurrentEpic)
	}
	if context.CurrentTask != "TEST-101" {
		t.Errorf("Expected CurrentTask to remain TEST-101, got %s", context.CurrentTask)
	}
	if context.CurrentSubtask != "TEST-102" {
		t.Errorf("Expected CurrentSubtask TEST-102, got %s", context.CurrentSubtask)
	}
}

func TestContextGetCurrentFocus(t *testing.T) {
	context := NewContext()

	// Test with no focus
	if focus := context.GetCurrentFocus(); focus != "" {
		t.Errorf("Expected empty focus, got %s", focus)
	}

	// Test with epic focus
	context.CurrentEpic = "TEST-100"
	if focus := context.GetCurrentFocus(); focus != "TEST-100" {
		t.Errorf("Expected focus TEST-100, got %s", focus)
	}

	// Test with task focus (should override epic)
	context.CurrentTask = "TEST-101"
	if focus := context.GetCurrentFocus(); focus != "TEST-101" {
		t.Errorf("Expected focus TEST-101, got %s", focus)
	}

	// Test with subtask focus (should override task)
	context.CurrentSubtask = "TEST-102"
	if focus := context.GetCurrentFocus(); focus != "TEST-102" {
		t.Errorf("Expected focus TEST-102, got %s", focus)
	}
}
