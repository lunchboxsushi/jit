# Task 05: Context Management

**Status:** [ready]

## Overview
Implement the system that tracks the current working context (epic, task, subtask) in `context.json`, with helper functions to query and update the focus.

## Objectives
- Provide functions to get/set current context
- Persist context changes via Storage layer
- Expose CLI-friendly utilities (`PrintCurrentContext`)
- Maintain recent tickets list (most-recently-focused)

## Deliverables
- [ready] `internal/storage/context.go` (context read/write – if not already covered by Task 03)
- [ready] `internal/commands/context_util.go` helper utilities
- [ready] Unit tests for context switching logic

## Dependencies
- Task 03 Storage layer
- Task 01 Context struct

## Implementation Notes
- Ensure atomic updates to avoid context corruption
- Limit recent tickets list to N entries (configurable)

## Acceptance Criteria
- `go test ./internal/...` passes
- Context switches persist correctly across CLI invocations

## Next Tasks
- 06-command-init.md 