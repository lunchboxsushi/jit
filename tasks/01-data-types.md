# Task 01: Core Data Types

**Status:** [completed]

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
- [completed] `pkg/types/ticket.go` with `Ticket` struct and related sub-structs
- [completed] `pkg/types/context.go` with `Context` struct
- [completed] `pkg/types/config.go` with `Config` struct mirroring `~/.jit/config.yml`
- [completed] Unit tests validating JSON (un)marshal for Ticket and Context
- [completed] Ticket type constants (Epic, Task, Subtask)
- [completed] Helper methods for ticket type checking

## Dependencies
- None (pure Go)

## Implementation Notes
- Keep structs JSON-tagged for storage; YAML tags for config where needed
- Use `omitempty` wisely to reduce file size
- Add helper constructors for ticket creation
- Use constants for ticket types to prevent typos
- Simplified relationships to use ParentKey and Children only

## Acceptance Criteria
- [completed] `go test ./pkg/types/...` passes
- [completed] Tickets can round-trip JSON marshal/unmarshal without data loss
- [completed] Config can be unmarshaled from sample YAML
- [completed] Ticket type constants work correctly
- [completed] Helper methods (IsEpic, IsTask, IsSubtask, IsOrphanTask) work correctly

## Next Tasks
- 02-configuration.md 