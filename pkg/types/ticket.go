package types

import "time"

// Ticket types as constants
const (
	TicketTypeEpic    = "Epic"
	TicketTypeTask    = "Task"
	TicketTypeSubtask = "Subtask"
)

// Ticket represents a Jira ticket (Epic, Task, or Subtask)
type Ticket struct {
	Key           string              `json:"key"`
	Title         string              `json:"title"`
	Type          string              `json:"type"` // Use constants: TicketTypeEpic, TicketTypeTask, TicketTypeSubtask
	Status        string              `json:"status"`
	Priority      string              `json:"priority"`
	Description   string              `json:"description"`
	Metadata      TicketMetadata      `json:"metadata"`
	Relationships TicketRelationships `json:"relationships"`
	JiraData      JiraData            `json:"jira_data"`
	LocalData     LocalData           `json:"local_data"`
}

// TicketMetadata contains metadata about the ticket
type TicketMetadata struct {
	Project  string    `json:"project"`
	Assignee string    `json:"assignee"`
	Created  time.Time `json:"created"`
	Updated  time.Time `json:"updated"`
	Labels   []string  `json:"labels"`
}

// TicketRelationships defines parent/child relationships
type TicketRelationships struct {
	ParentKey string   `json:"parent_key"` // Parent task for subtasks, empty for epics and orphan tasks
	Children  []string `json:"children"`   // Child tickets
}

// JiraData contains Jira-specific information
type JiraData struct {
	URL          string                 `json:"url"`
	CustomFields map[string]interface{} `json:"custom_fields"`
}

// LocalData contains local-only information
type LocalData struct {
	LastSync     time.Time `json:"last_sync"`
	LocalChanges bool      `json:"local_changes"`
	AIEnhanced   bool      `json:"ai_enhanced"`
}

// NewTicket creates a new ticket with default values
func NewTicket(key, title, ticketType string) *Ticket {
	now := time.Now()
	return &Ticket{
		Key:         key,
		Title:       title,
		Type:        ticketType,
		Status:      "To Do",
		Priority:    "Medium",
		Description: "",
		Metadata: TicketMetadata{
			Created: now,
			Updated: now,
			Labels:  []string{},
		},
		Relationships: TicketRelationships{
			Children: []string{},
		},
		JiraData: JiraData{
			CustomFields: make(map[string]interface{}),
		},
		LocalData: LocalData{
			LastSync:     now,
			LocalChanges: false,
			AIEnhanced:   false,
		},
	}
}

// IsEpic returns true if the ticket is an epic
func (t *Ticket) IsEpic() bool {
	return t.Type == TicketTypeEpic
}

// IsTask returns true if the ticket is a task
func (t *Ticket) IsTask() bool {
	return t.Type == TicketTypeTask
}

// IsSubtask returns true if the ticket is a subtask
func (t *Ticket) IsSubtask() bool {
	return t.Type == TicketTypeSubtask
}

// IsOrphanTask returns true if the task has no parent (not a subtask)
func (t *Ticket) IsOrphanTask() bool {
	return t.Type == TicketTypeTask && t.Relationships.ParentKey == ""
}
