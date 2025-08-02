package storage

import (
	"fmt"
	"time"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// ContextManager provides high-level context management operations
type ContextManager struct {
	storage Storage
}

// NewContextManager creates a new context manager
func NewContextManager(storage Storage) *ContextManager {
	return &ContextManager{
		storage: storage,
	}
}

// GetCurrentContext loads the current context or creates a new one
func (cm *ContextManager) GetCurrentContext() (*types.Context, error) {
	context, err := cm.storage.LoadContext()
	if err != nil {
		return nil, fmt.Errorf("failed to load context: %v", err)
	}
	return context, nil
}

// SetFocus updates the context focus and saves it
func (cm *ContextManager) SetFocus(ticketKey, ticketType string) error {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return err
	}

	// Update focus
	context.SetFocus(ticketKey, ticketType)

	// Save context
	return cm.storage.SaveContext(context)
}

// GetCurrentFocus returns the current focus ticket key
func (cm *ContextManager) GetCurrentFocus() (string, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return "", err
	}

	return context.GetCurrentFocus(), nil
}

// GetCurrentEpic returns the current epic key
func (cm *ContextManager) GetCurrentEpic() (string, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return "", err
	}

	return context.CurrentEpic, nil
}

// GetCurrentTask returns the current task key
func (cm *ContextManager) GetCurrentTask() (string, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return "", err
	}

	return context.CurrentTask, nil
}

// GetCurrentSubtask returns the current subtask key
func (cm *ContextManager) GetCurrentSubtask() (string, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return "", err
	}

	return context.CurrentSubtask, nil
}

// AddToRecent adds a ticket to the recent list
func (cm *ContextManager) AddToRecent(ticketKey string) error {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return err
	}

	// Add to recent tickets
	context.LastUpdated = time.Now()

	// Remove if already present
	for i, key := range context.RecentTickets {
		if key == ticketKey {
			context.RecentTickets = append(context.RecentTickets[:i], context.RecentTickets[i+1:]...)
			break
		}
	}

	// Add to beginning
	context.RecentTickets = append([]string{ticketKey}, context.RecentTickets...)

	// Keep only last 10
	if len(context.RecentTickets) > 10 {
		context.RecentTickets = context.RecentTickets[:10]
	}

	// Save context
	return cm.storage.SaveContext(context)
}

// GetRecentTickets returns the list of recent tickets
func (cm *ContextManager) GetRecentTickets() ([]string, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return nil, err
	}

	return context.RecentTickets, nil
}

// ClearContext clears the current context
func (cm *ContextManager) ClearContext() error {
	context := types.NewContext()
	return cm.storage.SaveContext(context)
}

// IsInEpic checks if we're currently focused on an epic
func (cm *ContextManager) IsInEpic() (bool, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return false, err
	}

	return context.CurrentEpic != "", nil
}

// IsInTask checks if we're currently focused on a task
func (cm *ContextManager) IsInTask() (bool, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return false, err
	}

	return context.CurrentTask != "", nil
}

// IsInSubtask checks if we're currently focused on a subtask
func (cm *ContextManager) IsInSubtask() (bool, error) {
	context, err := cm.GetCurrentContext()
	if err != nil {
		return false, err
	}

	return context.CurrentSubtask != "", nil
}
