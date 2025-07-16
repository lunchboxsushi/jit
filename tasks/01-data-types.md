# Task 01: Core Data Types

**Status:** [ready]

## Overview
Define the foundational Go structs that represent tickets, relationships, context, and configuration. These types will be used across the entire application.

## Objectives
- Establish Go structs in `pkg/types/` for:
  - Ticket (Epic, Task, Subtask)
  - Relationships meta
  - LocalData / JiraData
  - Context (current focus)
  - Configuration (parsed YAML)

## Deliverables
- [ready] `pkg/types/ticket.go` with `Ticket` struct and related sub-structs
- [ready] `pkg/types/context.go` with `Context` struct
- [ready] `pkg/types/config.go` with `Config` struct mirroring `~/.jit/config.yml`
- [ready] Unit tests validating JSON (un)marshal for Ticket and Context

## Dependencies
- None (pure Go)

## Implementation Notes
- Keep structs JSON-tagged for storage; YAML tags for config where needed
- Use `omitempty` wisely to reduce file size
- Add helper constructors for ticket creation

## Acceptance Criteria
- `go test ./pkg/types/...` passes
- Tickets can round-trip JSON marshal/unmarshal without data loss
- Config can be unmarshaled from sample YAML

## Next Tasks
- 02-configuration.md 