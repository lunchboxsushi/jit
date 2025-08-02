# Task 04: Jira API Client

**Status:** [completed]

## Overview
Implement a lightweight Jira REST API client to authenticate, fetch, create, and update tickets, as well as add comments.

## Objectives
- Wrap Jira Cloud REST API v3 endpoints we need:
  - GET issue (with expand=changelog to fetch children)
  - POST issue (create)
  - POST comment
  - GET search (JQL)
- Token-based authentication via Basic Auth or PAT
- Rate limit handling with retry/backoff
- Typed response structs for deserialization
- Error wrapping and surfacing meaningful messages

## Deliverables
- [completed] `internal/jira/client.go` with HTTP client and auth
- [completed] `internal/jira/tickets.go` with high-level helpers (`GetTicket`, `CreateTicket`, `AddComment`)
- [completed] `internal/jira/types.go` with minimal Jira JSON structs
- [completed] Unit tests with mocked HTTP server

## Dependencies
- Task 02 config (Jira creds)
- net/http, encoding/json

## Implementation Notes
- Use `context.Context` for timeout/cancellation
- Backoff on 429 responses (Respect `Retry-After`)
- Keep external structs minimal to reduce maintenance
- Basic authentication with username:token
- Comprehensive error handling with Jira error response parsing
- Rate limiting with automatic retry and backoff

## Acceptance Criteria
- [completed] `go test ./internal/jira/...` passes with mocked server
- [completed] All API endpoints properly implemented and tested
- [completed] Authentication and error handling working correctly
- [completed] Type conversion between Jira API and internal types
- [completed] Rate limiting and retry logic implemented

## Next Tasks
- 05-context-management.md 