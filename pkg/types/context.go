package types

import "time"

// Context represents the current working context
type Context struct {
	CurrentEpic    string    `json:"current_epic"`
	CurrentTask    string    `json:"current_task"`
	CurrentSubtask string    `json:"current_subtask"`
	LastUpdated    time.Time `json:"last_updated"`
	RecentTickets  []string  `json:"recent_tickets"`
}

// NewContext creates a new context with default values
func NewContext() *Context {
	return &Context{
		CurrentEpic:    "",
		CurrentTask:    "",
		CurrentSubtask: "",
		LastUpdated:    time.Now(),
		RecentTickets:  []string{},
	}
}

// SetFocus updates the context focus based on ticket type
func (c *Context) SetFocus(ticketKey, ticketType string) {
	c.LastUpdated = time.Now()

	switch ticketType {
	case TicketTypeEpic:
		c.CurrentEpic = ticketKey
		c.CurrentTask = ""
		c.CurrentSubtask = ""
	case TicketTypeTask:
		c.CurrentTask = ticketKey
		c.CurrentSubtask = ""
	case TicketTypeSubtask:
		c.CurrentSubtask = ticketKey
	}

	// Add to recent tickets if not already present
	c.addToRecent(ticketKey)
}

// GetCurrentFocus returns the most specific current focus
func (c *Context) GetCurrentFocus() string {
	if c.CurrentSubtask != "" {
		return c.CurrentSubtask
	}
	if c.CurrentTask != "" {
		return c.CurrentTask
	}
	return c.CurrentEpic
}

// addToRecent adds a ticket to recent tickets list
func (c *Context) addToRecent(ticketKey string) {
	// Remove if already present
	for i, key := range c.RecentTickets {
		if key == ticketKey {
			c.RecentTickets = append(c.RecentTickets[:i], c.RecentTickets[i+1:]...)
			break
		}
	}

	// Add to beginning
	c.RecentTickets = append([]string{ticketKey}, c.RecentTickets...)

	// Keep only last 10
	if len(c.RecentTickets) > 10 {
		c.RecentTickets = c.RecentTickets[:10]
	}
}
