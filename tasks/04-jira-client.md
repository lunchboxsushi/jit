# Task 04: Jira API Client

**Status:** [ready]

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
- [ready] `internal/jira/client.go` with HTTP client and auth
- [ready] `internal/jira/tickets.go` with high-level helpers (`GetIssue`, `CreateIssue`, `AddComment`)
- [ready] `internal/jira/types.go` with minimal Jira JSON structs
- [ready] Unit tests with mocked HTTP server

## Dependencies
- Task 02 config (Jira creds)
- net/http, encoding/json

## Implementation Notes
- Use `context.Context` for timeout/cancellation
- Backoff on 429 responses (Respect `Retry-After`)
- Keep external structs minimal to reduce maintenance

## Acceptance Criteria
- `go test ./internal/jira/...` passes with mocked server
- Able to fetch real ticket when correct creds in env (manual test)

## Next Tasks
- 05-context-management.md 