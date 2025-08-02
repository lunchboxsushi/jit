package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// SaveTicket saves a ticket to JSON file
func (s *JSONStorage) SaveTicket(ticket *types.Ticket) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if ticket == nil {
		return fmt.Errorf("ticket cannot be nil")
	}

	if ticket.Key == "" {
		return fmt.Errorf("ticket key cannot be empty")
	}

	path := s.GetTicketPath(ticket.Key)

	// Marshal ticket to JSON
	data, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal ticket: %v", err)
	}

	// Write atomically
	if err := s.atomicWrite(path, data); err != nil {
		return fmt.Errorf("failed to write ticket: %v", err)
	}

	return nil
}

// LoadTicket loads a ticket from JSON file
func (s *JSONStorage) LoadTicket(key string) (*types.Ticket, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if key == "" {
		return nil, fmt.Errorf("ticket key cannot be empty")
	}

	path := s.GetTicketPath(key)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("ticket %s not found", key)
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read ticket: %v", err)
	}

	// Unmarshal JSON
	var ticket types.Ticket
	if err := json.Unmarshal(data, &ticket); err != nil {
		return nil, fmt.Errorf("failed to unmarshal ticket: %v", err)
	}

	return &ticket, nil
}

// DeleteTicket deletes a ticket file
func (s *JSONStorage) DeleteTicket(key string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if key == "" {
		return fmt.Errorf("ticket key cannot be empty")
	}

	path := s.GetTicketPath(key)

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("ticket %s not found", key)
	}

	// Delete file
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete ticket: %v", err)
	}

	return nil
}

// ListTickets returns a list of all ticket keys
func (s *JSONStorage) ListTickets() ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ticketsDir := filepath.Join(s.dataDir, "tickets")

	// Check if directory exists
	if _, err := os.Stat(ticketsDir); os.IsNotExist(err) {
		return []string{}, nil // Return empty list if directory doesn't exist
	}

	// Read directory
	entries, err := os.ReadDir(ticketsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read tickets directory: %v", err)
	}

	var keys []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			// Remove .json extension to get the key
			key := strings.TrimSuffix(entry.Name(), ".json")
			keys = append(keys, key)
		}
	}

	return keys, nil
}

// SaveContext saves context to JSON file
func (s *JSONStorage) SaveContext(context *types.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if context == nil {
		return fmt.Errorf("context cannot be nil")
	}

	path := s.GetContextPath()

	// Marshal context to JSON
	data, err := json.MarshalIndent(context, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %v", err)
	}

	// Write atomically
	if err := s.atomicWrite(path, data); err != nil {
		return fmt.Errorf("failed to write context: %v", err)
	}

	return nil
}

// LoadContext loads context from JSON file
func (s *JSONStorage) LoadContext() (*types.Context, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	path := s.GetContextPath()

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Return default context if file doesn't exist
		return types.NewContext(), nil
	}

	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read context: %v", err)
	}

	// Unmarshal JSON
	var context types.Context
	if err := json.Unmarshal(data, &context); err != nil {
		return nil, fmt.Errorf("failed to unmarshal context: %v", err)
	}

	return &context, nil
}
