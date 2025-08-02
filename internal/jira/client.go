package jira

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/lunchboxsushi/jit/pkg/types"
)

// Client represents a Jira API client
type Client struct {
	baseURL    string
	username   string
	token      string
	httpClient *http.Client
	config     *types.JiraConfig
}

// NewClient creates a new Jira client
func NewClient(config *types.JiraConfig) *Client {
	return &Client{
		baseURL:    config.URL,
		username:   config.Username,
		token:      config.Token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		config:     config,
	}
}

// doRequest performs an HTTP request with authentication and error handling
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	// Construct full URL
	fullURL := c.baseURL + "/rest/api/3" + endpoint

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Set authentication
	auth := base64.StdEncoding.EncodeToString([]byte(c.username + ":" + c.token))
	req.Header.Set("Authorization", "Basic "+auth)

	// Make request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %v", err)
	}

	// Handle rate limiting
	if resp.StatusCode == 429 {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			if seconds, err := strconv.Atoi(retryAfter); err == nil {
				time.Sleep(time.Duration(seconds) * time.Second)
				// Retry the request
				return c.doRequest(ctx, method, endpoint, body)
			}
		}
		// Default retry delay
		time.Sleep(5 * time.Second)
		return c.doRequest(ctx, method, endpoint, body)
	}

	return resp, nil
}

// parseErrorResponse parses Jira error responses
func (c *Client) parseErrorResponse(resp *http.Response) error {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read error response: %v", err)
	}

	// Try to parse as Jira error
	var jiraError JiraErrorResponse
	if err := json.Unmarshal(body, &jiraError); err == nil {
		if len(jiraError.ErrorCollection.ErrorMessages) > 0 {
			return fmt.Errorf("Jira API error: %s", strings.Join(jiraError.ErrorCollection.ErrorMessages, "; "))
		}
		if len(jiraError.ErrorCollection.Errors) > 0 {
			var errors []string
			for field, message := range jiraError.ErrorCollection.Errors {
				errors = append(errors, fmt.Sprintf("%s: %s", field, message))
			}
			return fmt.Errorf("Jira API error: %s", strings.Join(errors, "; "))
		}
	}

	// Fallback to generic error
	return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
}

// GetIssue fetches a single issue by key
func (c *Client) GetIssue(ctx context.Context, key string) (*JiraIssue, error) {
	endpoint := fmt.Sprintf("/issue/%s?expand=changelog", key)

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseErrorResponse(resp)
	}

	var issue JiraIssue
	if err := json.NewDecoder(resp.Body).Decode(&issue); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &issue, nil
}

// CreateIssue creates a new issue
func (c *Client) CreateIssue(ctx context.Context, request *JiraCreateIssueRequest) (*JiraCreateIssueResponse, error) {
	// Marshal request body
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	resp, err := c.doRequest(ctx, "POST", "/issue", strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, c.parseErrorResponse(resp)
	}

	var response JiraCreateIssueResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

// AddComment adds a comment to an issue
func (c *Client) AddComment(ctx context.Context, issueKey, commentBody string) (*JiraComment, error) {
	request := JiraCreateCommentRequest{
		Body: commentBody,
	}

	// Marshal request body
	body, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %v", err)
	}

	endpoint := fmt.Sprintf("/issue/%s/comment", issueKey)
	resp, err := c.doRequest(ctx, "POST", endpoint, strings.NewReader(string(body)))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		return nil, c.parseErrorResponse(resp)
	}

	var comment JiraComment
	if err := json.NewDecoder(resp.Body).Decode(&comment); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &comment, nil
}

// SearchIssues performs a JQL search
func (c *Client) SearchIssues(ctx context.Context, jql string, maxResults int) (*JiraSearchResponse, error) {
	// Build query parameters
	params := url.Values{}
	params.Set("jql", jql)
	params.Set("maxResults", strconv.Itoa(maxResults))
	params.Set("fields", "summary,description,status,priority,issuetype,project,assignee,reporter,created,updated,labels")

	endpoint := "/search?" + params.Encode()

	resp, err := c.doRequest(ctx, "GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, c.parseErrorResponse(resp)
	}

	var response JiraSearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return &response, nil
}

// TestConnection tests the connection to Jira
func (c *Client) TestConnection(ctx context.Context) error {
	resp, err := c.doRequest(ctx, "GET", "/myself", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return c.parseErrorResponse(resp)
	}

	return nil
}
