package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lunchboxsushi/jit/pkg/types"
)

func TestNewJSONStorage(t *testing.T) {
	tempDir := t.TempDir()

	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	if storage == nil {
		t.Fatal("Storage should not be nil")
	}

	// Check if directories were created
	expectedDirs := []string{"tickets", "cache"}
	for _, dir := range expectedDirs {
		dirPath := filepath.Join(tempDir, dir)
		if _, err := os.Stat(dirPath); os.IsNotExist(err) {
			t.Errorf("Expected directory %s to exist", dirPath)
		}
	}
}

func TestSaveAndLoadTicket(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a test ticket
	ticket := types.NewTicket("TEST-123", "Test Ticket", types.TicketTypeTask)
	ticket.Description = "This is a test ticket"
	ticket.Status = "In Progress"
	ticket.Priority = "High"
	ticket.Metadata.Project = "TEST"
	ticket.Metadata.Assignee = "test@example.com"

	// Save ticket
	err = storage.SaveTicket(ticket)
	if err != nil {
		t.Fatalf("Failed to save ticket: %v", err)
	}

	// Check if file exists
	if !storage.Exists("TEST-123") {
		t.Error("Ticket file should exist after saving")
	}

	// Load ticket
	loadedTicket, err := storage.LoadTicket("TEST-123")
	if err != nil {
		t.Fatalf("Failed to load ticket: %v", err)
	}

	// Verify ticket data
	if loadedTicket.Key != ticket.Key {
		t.Errorf("Expected key %s, got %s", ticket.Key, loadedTicket.Key)
	}
	if loadedTicket.Title != ticket.Title {
		t.Errorf("Expected title %s, got %s", ticket.Title, loadedTicket.Title)
	}
	if loadedTicket.Description != ticket.Description {
		t.Errorf("Expected description %s, got %s", ticket.Description, loadedTicket.Description)
	}
	if loadedTicket.Status != ticket.Status {
		t.Errorf("Expected status %s, got %s", ticket.Status, loadedTicket.Status)
	}
	if loadedTicket.Metadata.Project != ticket.Metadata.Project {
		t.Errorf("Expected project %s, got %s", ticket.Metadata.Project, loadedTicket.Metadata.Project)
	}
}

func TestLoadNonExistentTicket(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Try to load non-existent ticket
	_, err = storage.LoadTicket("NONEXISTENT-123")
	if err == nil {
		t.Error("Expected error when loading non-existent ticket")
	}

	if storage.Exists("NONEXISTENT-123") {
		t.Error("Non-existent ticket should not exist")
	}
}

func TestDeleteTicket(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create and save a ticket
	ticket := types.NewTicket("TEST-456", "Test Ticket", types.TicketTypeTask)
	err = storage.SaveTicket(ticket)
	if err != nil {
		t.Fatalf("Failed to save ticket: %v", err)
	}

	// Verify it exists
	if !storage.Exists("TEST-456") {
		t.Error("Ticket should exist before deletion")
	}

	// Delete ticket
	err = storage.DeleteTicket("TEST-456")
	if err != nil {
		t.Fatalf("Failed to delete ticket: %v", err)
	}

	// Verify it's gone
	if storage.Exists("TEST-456") {
		t.Error("Ticket should not exist after deletion")
	}
}

func TestListTickets(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Initially should be empty
	tickets, err := storage.ListTickets()
	if err != nil {
		t.Fatalf("Failed to list tickets: %v", err)
	}
	if len(tickets) != 0 {
		t.Errorf("Expected 0 tickets, got %d", len(tickets))
	}

	// Save some tickets
	ticket1 := types.NewTicket("TEST-1", "Test 1", types.TicketTypeEpic)
	ticket2 := types.NewTicket("TEST-2", "Test 2", types.TicketTypeTask)
	ticket3 := types.NewTicket("TEST-3", "Test 3", types.TicketTypeSubtask)

	err = storage.SaveTicket(ticket1)
	if err != nil {
		t.Fatalf("Failed to save ticket 1: %v", err)
	}
	err = storage.SaveTicket(ticket2)
	if err != nil {
		t.Fatalf("Failed to save ticket 2: %v", err)
	}
	err = storage.SaveTicket(ticket3)
	if err != nil {
		t.Fatalf("Failed to save ticket 3: %v", err)
	}

	// List tickets
	tickets, err = storage.ListTickets()
	if err != nil {
		t.Fatalf("Failed to list tickets: %v", err)
	}

	if len(tickets) != 3 {
		t.Errorf("Expected 3 tickets, got %d", len(tickets))
	}

	// Check that all expected tickets are in the list
	expected := map[string]bool{"TEST-1": true, "TEST-2": true, "TEST-3": true}
	for _, ticket := range tickets {
		if !expected[ticket] {
			t.Errorf("Unexpected ticket in list: %s", ticket)
		}
	}
}

func TestSaveAndLoadContext(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a test context
	context := types.NewContext()
	context.CurrentEpic = "TEST-100"
	context.CurrentTask = "TEST-101"
	context.RecentTickets = []string{"TEST-100", "TEST-101", "TEST-102"}

	// Save context
	err = storage.SaveContext(context)
	if err != nil {
		t.Fatalf("Failed to save context: %v", err)
	}

	// Load context
	loadedContext, err := storage.LoadContext()
	if err != nil {
		t.Fatalf("Failed to load context: %v", err)
	}

	// Verify context data
	if loadedContext.CurrentEpic != context.CurrentEpic {
		t.Errorf("Expected current epic %s, got %s", context.CurrentEpic, loadedContext.CurrentEpic)
	}
	if loadedContext.CurrentTask != context.CurrentTask {
		t.Errorf("Expected current task %s, got %s", context.CurrentTask, loadedContext.CurrentTask)
	}
	if len(loadedContext.RecentTickets) != len(context.RecentTickets) {
		t.Errorf("Expected %d recent tickets, got %d", len(context.RecentTickets), len(loadedContext.RecentTickets))
	}
}

func TestLoadDefaultContext(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Load context when no file exists (should return default)
	context, err := storage.LoadContext()
	if err != nil {
		t.Fatalf("Failed to load default context: %v", err)
	}

	// Should be a new context
	if context.CurrentEpic != "" {
		t.Errorf("Expected empty current epic, got %s", context.CurrentEpic)
	}
	if context.CurrentTask != "" {
		t.Errorf("Expected empty current task, got %s", context.CurrentTask)
	}
	if context.CurrentSubtask != "" {
		t.Errorf("Expected empty current subtask, got %s", context.CurrentSubtask)
	}
	if len(context.RecentTickets) != 0 {
		t.Errorf("Expected empty recent tickets, got %d", len(context.RecentTickets))
	}
}

func TestContextManager(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	cm := NewContextManager(storage)

	// Test initial state
	focus, err := cm.GetCurrentFocus()
	if err != nil {
		t.Fatalf("Failed to get current focus: %v", err)
	}
	if focus != "" {
		t.Errorf("Expected empty focus, got %s", focus)
	}

	// Test setting focus
	err = cm.SetFocus("TEST-100", types.TicketTypeEpic)
	if err != nil {
		t.Fatalf("Failed to set focus: %v", err)
	}

	focus, err = cm.GetCurrentFocus()
	if err != nil {
		t.Fatalf("Failed to get current focus: %v", err)
	}
	if focus != "TEST-100" {
		t.Errorf("Expected focus TEST-100, got %s", focus)
	}

	// Test adding to recent
	err = cm.AddToRecent("TEST-200")
	if err != nil {
		t.Fatalf("Failed to add to recent: %v", err)
	}

	recent, err := cm.GetRecentTickets()
	if err != nil {
		t.Fatalf("Failed to get recent tickets: %v", err)
	}
	if len(recent) != 2 { // TEST-200 and TEST-100
		t.Errorf("Expected 2 recent tickets, got %d", len(recent))
	}
	if recent[0] != "TEST-200" {
		t.Errorf("Expected first recent ticket to be TEST-200, got %s", recent[0])
	}
}

func TestConcurrentAccess(t *testing.T) {
	tempDir := t.TempDir()
	storage, err := NewJSONStorage(tempDir)
	if err != nil {
		t.Fatalf("Failed to create storage: %v", err)
	}

	// Create a ticket
	ticket := types.NewTicket("CONCURRENT-123", "Concurrent Test", types.TicketTypeTask)
	err = storage.SaveTicket(ticket)
	if err != nil {
		t.Fatalf("Failed to save ticket: %v", err)
	}

	// Test concurrent reads
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			defer func() { done <- true }()
			_, err := storage.LoadTicket("CONCURRENT-123")
			if err != nil {
				t.Errorf("Failed to load ticket concurrently: %v", err)
			}
		}()
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify ticket is still intact
	loadedTicket, err := storage.LoadTicket("CONCURRENT-123")
	if err != nil {
		t.Fatalf("Failed to load ticket after concurrent access: %v", err)
	}
	if loadedTicket.Key != "CONCURRENT-123" {
		t.Errorf("Ticket corrupted after concurrent access")
	}
}
