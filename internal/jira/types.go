package jira

import "time"

// JiraIssue represents a Jira issue from the API
type JiraIssue struct {
	Key       string          `json:"key"`
	ID        string          `json:"id"`
	Fields    JiraIssueFields `json:"fields"`
	Changelog *JiraChangelog  `json:"changelog,omitempty"`
}

// JiraIssueFields contains the main issue data
type JiraIssueFields struct {
	Summary      string                 `json:"summary"`
	Description  string                 `json:"description"`
	Status       JiraStatus             `json:"status"`
	Priority     JiraPriority           `json:"priority"`
	IssueType    JiraIssueType          `json:"issuetype"`
	Project      JiraProject            `json:"project"`
	Assignee     *JiraUser              `json:"assignee"`
	Reporter     *JiraUser              `json:"reporter"`
	Created      time.Time              `json:"created"`
	Updated      time.Time              `json:"updated"`
	Labels       []string               `json:"labels"`
	CustomFields map[string]interface{} `json:"-"`
}

// JiraStatus represents the issue status
type JiraStatus struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// JiraPriority represents the issue priority
type JiraPriority struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Icon string `json:"iconUrl"`
}

// JiraIssueType represents the issue type
type JiraIssueType struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	IconURL     string `json:"iconUrl"`
}

// JiraProject represents the project
type JiraProject struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Name string `json:"name"`
}

// JiraUser represents a Jira user
type JiraUser struct {
	AccountID   string `json:"accountId"`
	DisplayName string `json:"displayName"`
	Email       string `json:"emailAddress"`
	Active      bool   `json:"active"`
}

// JiraChangelog contains the issue changelog
type JiraChangelog struct {
	Histories []JiraChangelogHistory `json:"histories"`
}

// JiraChangelogHistory represents a changelog entry
type JiraChangelogHistory struct {
	ID      string              `json:"id"`
	Author  JiraUser            `json:"author"`
	Created time.Time           `json:"created"`
	Items   []JiraChangelogItem `json:"items"`
}

// JiraChangelogItem represents a change in the changelog
type JiraChangelogItem struct {
	Field     string `json:"field"`
	FieldType string `json:"fieldtype"`
	FieldID   string `json:"fieldId"`
	From      string `json:"fromString"`
	To        string `json:"toString"`
}

// JiraCreateIssueRequest represents the request to create an issue
type JiraCreateIssueRequest struct {
	Fields JiraCreateIssueFields `json:"fields"`
}

// JiraCreateIssueFields contains the fields for creating an issue
type JiraCreateIssueFields struct {
	Project     JiraProjectReference `json:"project"`
	Summary     string               `json:"summary"`
	Description string               `json:"description"`
	IssueType   JiraIssueTypeRef     `json:"issuetype"`
	Priority    *JiraPriorityRef     `json:"priority,omitempty"`
	Labels      []string             `json:"labels,omitempty"`
	Assignee    *JiraUserRef         `json:"assignee,omitempty"`
	Parent      *JiraParentRef       `json:"parent,omitempty"`
}

// JiraProjectReference for creating issues
type JiraProjectReference struct {
	Key string `json:"key"`
}

// JiraIssueTypeRef for creating issues
type JiraIssueTypeRef struct {
	ID string `json:"id"`
}

// JiraPriorityRef for creating issues
type JiraPriorityRef struct {
	ID string `json:"id"`
}

// JiraUserRef for creating issues
type JiraUserRef struct {
	AccountID string `json:"accountId"`
}

// JiraParentRef for creating subtasks
type JiraParentRef struct {
	Key string `json:"key"`
}

// JiraCreateIssueResponse represents the response from creating an issue
type JiraCreateIssueResponse struct {
	ID   string `json:"id"`
	Key  string `json:"key"`
	Self string `json:"self"`
}

// JiraComment represents a comment
type JiraComment struct {
	ID           string    `json:"id"`
	Author       JiraUser  `json:"author"`
	Body         string    `json:"body"`
	Created      time.Time `json:"created"`
	Updated      time.Time `json:"updated"`
	UpdateAuthor JiraUser  `json:"updateAuthor"`
}

// JiraCreateCommentRequest represents the request to create a comment
type JiraCreateCommentRequest struct {
	Body string `json:"body"`
}

// JiraSearchResponse represents the response from a JQL search
type JiraSearchResponse struct {
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []JiraIssue `json:"issues"`
}

// JiraError represents a Jira API error
type JiraError struct {
	ErrorMessages []string          `json:"errorMessages"`
	Errors        map[string]string `json:"errors"`
}

// JiraErrorResponse represents the full error response
type JiraErrorResponse struct {
	ErrorCollection JiraError `json:"errorCollection"`
}
