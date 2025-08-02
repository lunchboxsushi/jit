package jira

import (
	"context"
	"fmt"
	"strings"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// TicketService provides high-level ticket operations
type TicketService struct {
	client *Client
}

// NewTicketService creates a new ticket service
func NewTicketService(client *Client) *TicketService {
	return &TicketService{
		client: client,
	}
}

// GetTicket fetches a ticket and converts it to our internal format
func (ts *TicketService) GetTicket(ctx context.Context, key string) (*types.Ticket, error) {
	jiraIssue, err := ts.client.GetIssue(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch issue %s: %v", key, err)
	}

	return ts.convertJiraIssueToTicket(jiraIssue), nil
}

// CreateTicket creates a new ticket from our internal format
func (ts *TicketService) CreateTicket(ctx context.Context, ticket *types.Ticket) (*types.Ticket, error) {
	// Determine issue type ID based on ticket type
	issueTypeID, err := ts.getIssueTypeID(ticket.Type)
	if err != nil {
		return nil, err
	}

	// Build create request
	request := &JiraCreateIssueRequest{
		Fields: JiraCreateIssueFields{
			Project: JiraProjectReference{
				Key: ticket.Metadata.Project,
			},
			Summary:     ticket.Title,
			Description: ticket.Description,
			IssueType: JiraIssueTypeRef{
				ID: issueTypeID,
			},
			Labels: ticket.Metadata.Labels,
		},
	}

	// Add priority if specified
	if ticket.Priority != "" {
		priorityID, err := ts.getPriorityID(ticket.Priority)
		if err == nil {
			request.Fields.Priority = &JiraPriorityRef{ID: priorityID}
		}
	}

	// Add assignee if specified
	if ticket.Metadata.Assignee != "" {
		request.Fields.Assignee = &JiraUserRef{AccountID: ticket.Metadata.Assignee}
	}

	// Add parent for subtasks
	if ticket.Type == types.TicketTypeSubtask && ticket.Relationships.ParentKey != "" {
		request.Fields.Parent = &JiraParentRef{Key: ticket.Relationships.ParentKey}
	}

	// Create the issue
	response, err := ts.client.CreateIssue(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to create issue: %v", err)
	}

	// Fetch the created issue to get full details
	return ts.GetTicket(ctx, response.Key)
}

// AddComment adds a comment to a ticket
func (ts *TicketService) AddComment(ctx context.Context, ticketKey, commentBody string) error {
	_, err := ts.client.AddComment(ctx, ticketKey, commentBody)
	if err != nil {
		return fmt.Errorf("failed to add comment: %v", err)
	}

	return nil
}

// SearchTickets performs a JQL search and returns tickets
func (ts *TicketService) SearchTickets(ctx context.Context, jql string, maxResults int) ([]*types.Ticket, error) {
	response, err := ts.client.SearchIssues(ctx, jql, maxResults)
	if err != nil {
		return nil, fmt.Errorf("failed to search issues: %v", err)
	}

	var tickets []*types.Ticket
	for _, issue := range response.Issues {
		ticket := ts.convertJiraIssueToTicket(&issue)
		tickets = append(tickets, ticket)
	}

	return tickets, nil
}

// convertJiraIssueToTicket converts a Jira API issue to our internal ticket format
func (ts *TicketService) convertJiraIssueToTicket(jiraIssue *JiraIssue) *types.Ticket {
	ticket := &types.Ticket{
		Key:         jiraIssue.Key,
		Title:       jiraIssue.Fields.Summary,
		Type:        ts.convertIssueType(jiraIssue.Fields.IssueType.Name),
		Status:      jiraIssue.Fields.Status.Name,
		Priority:    jiraIssue.Fields.Priority.Name,
		Description: jiraIssue.Fields.Description,
		Metadata: types.TicketMetadata{
			Project: jiraIssue.Fields.Project.Key,
			Created: jiraIssue.Fields.Created,
			Updated: jiraIssue.Fields.Updated,
			Labels:  jiraIssue.Fields.Labels,
		},
		Relationships: types.TicketRelationships{
			Children: []string{},
		},
		JiraData: types.JiraData{
			URL:          fmt.Sprintf("%s/browse/%s", ts.client.baseURL, jiraIssue.Key),
			CustomFields: make(map[string]interface{}),
		},
		LocalData: types.LocalData{
			LastSync:     jiraIssue.Fields.Updated,
			LocalChanges: false,
			AIEnhanced:   false,
		},
	}

	// Set assignee if available
	if jiraIssue.Fields.Assignee != nil {
		ticket.Metadata.Assignee = jiraIssue.Fields.Assignee.Email
	}

	// Extract parent relationships from changelog
	if jiraIssue.Changelog != nil {
		ts.extractParentRelationships(ticket, jiraIssue.Changelog)
	}

	return ticket
}

// convertIssueType converts Jira issue type names to our internal types
func (ts *TicketService) convertIssueType(jiraType string) string {
	switch strings.ToLower(jiraType) {
	case "epic":
		return types.TicketTypeEpic
	case "story", "task":
		return types.TicketTypeTask
	case "sub-task", "subtask":
		return types.TicketTypeSubtask
	default:
		return types.TicketTypeTask // Default to task
	}
}

// extractParentRelationships extracts parent/child relationships from changelog
func (ts *TicketService) extractParentRelationships(ticket *types.Ticket, changelog *JiraChangelog) {
	// This is a simplified implementation
	// In a real implementation, you'd parse the changelog to find parent/child relationships
	// For now, we'll leave this empty and let the calling code handle relationships
}

// getIssueTypeID returns the Jira issue type ID for our ticket type
func (ts *TicketService) getIssueTypeID(ticketType string) (string, error) {
	// This is a simplified implementation
	// In a real implementation, you'd fetch issue types from Jira and cache them
	switch ticketType {
	case types.TicketTypeEpic:
		return "10000", nil // Epic issue type ID
	case types.TicketTypeTask:
		return "10001", nil // Story/Task issue type ID
	case types.TicketTypeSubtask:
		return "10003", nil // Sub-task issue type ID
	default:
		return "", fmt.Errorf("unknown ticket type: %s", ticketType)
	}
}

// getPriorityID returns the Jira priority ID for our priority name
func (ts *TicketService) getPriorityID(priority string) (string, error) {
	// This is a simplified implementation
	// In a real implementation, you'd fetch priorities from Jira and cache them
	switch strings.ToLower(priority) {
	case "highest":
		return "1", nil
	case "high":
		return "2", nil
	case "medium":
		return "3", nil
	case "low":
		return "4", nil
	case "lowest":
		return "5", nil
	default:
		return "3", nil // Default to medium
	}
}

// GetEpicChildren fetches all children of an epic
func (ts *TicketService) GetEpicChildren(ctx context.Context, epicKey string) ([]*types.Ticket, error) {
	// Search for issues that are linked to this epic
	jql := fmt.Sprintf("'Epic Link' = %s ORDER BY created DESC", epicKey)
	return ts.SearchTickets(ctx, jql, 100)
}

// GetTaskSubtasks fetches all subtasks of a task
func (ts *TicketService) GetTaskSubtasks(ctx context.Context, taskKey string) ([]*types.Ticket, error) {
	// Search for subtasks that have this task as parent
	jql := fmt.Sprintf("parent = %s ORDER BY created DESC", taskKey)
	return ts.SearchTickets(ctx, jql, 100)
}
