package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// Storage defines the interface for ticket and context storage
type Storage interface {
	// Ticket operations
	SaveTicket(ticket *types.Ticket) error
	LoadTicket(key string) (*types.Ticket, error)
	DeleteTicket(key string) error
	ListTickets() ([]string, error)

	// Context operations
	SaveContext(context *types.Context) error
	LoadContext() (*types.Context, error)

	// Utility operations
	Exists(key string) bool
	GetTicketPath(key string) string
}

// JSONStorage implements Storage interface using JSON files
type JSONStorage struct {
	dataDir string
	mu      sync.RWMutex
}

// NewJSONStorage creates a new JSON storage instance
func NewJSONStorage(dataDir string) (*JSONStorage, error) {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory %s: %v", dataDir, err)
	}

	// Create subdirectories
	subdirs := []string{"tickets", "cache"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(dataDir, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			return nil, fmt.Errorf("failed to create subdirectory %s: %v", subdirPath, err)
		}
	}

	return &JSONStorage{
		dataDir: dataDir,
	}, nil
}

// GetTicketPath returns the file path for a ticket
func (s *JSONStorage) GetTicketPath(key string) string {
	return filepath.Join(s.dataDir, "tickets", key+".json")
}

// GetContextPath returns the file path for context
func (s *JSONStorage) GetContextPath() string {
	return filepath.Join(s.dataDir, "context.json")
}

// Exists checks if a ticket exists
func (s *JSONStorage) Exists(key string) bool {
	path := s.GetTicketPath(key)
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// atomicWrite writes data to a file atomically using a temporary file
func (s *JSONStorage) atomicWrite(path string, data []byte) error {
	// Create temporary file in the same directory
	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, "*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name()) // Clean up temp file on error

	// Write data to temp file
	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write to temp file: %v", err)
	}

	// Sync to disk
	if err := tmpFile.Sync(); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to sync temp file: %v", err)
	}

	// Close temp file
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %v", err)
	}

	// Atomic rename
	if err := os.Rename(tmpFile.Name(), path); err != nil {
		return fmt.Errorf("failed to rename temp file: %v", err)
	}

	return nil
}
