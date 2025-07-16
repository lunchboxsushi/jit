# Task 09: Commands `jit epic`, `jit task`, `jit subtask`

**Status:** [ready]

## Overview
Implement the ticket creation commands that open an editor, capture user input, optionally enhance with AI, create the ticket in Jira, store locally, and set focus. Support orphan task creation via `jit task -o`.

## Objectives
- Common helper to open editor with template (from `templates/`)
- AI enrichment (optional `--no-enrich` flag) – placeholder until Task 11
- Jira client call to create ticket
- Update relationships (link to parent epic/task unless orphan `-o`)
- Save new ticket via Storage layer
- Update context to new ticket
- Flags:
  - `--no-enrich`, `--no-create` (local only), `--orphan/-o` for tasks

## Deliverables
- [ready] `internal/commands/epic.go`
- [ready] `internal/commands/task.go`
- [ready] `internal/commands/subtask.go`
- [ready] Editor helper `internal/ui/editor.go`
- [ready] Unit tests mocking editor, Jira, storage

## Dependencies
- Task 04 Jira client
- Task 03 Storage layer
- Task 02 Config (templates path)

## Implementation Notes
- Use `$EDITOR` env var fallback to `vim`
- For orphan tasks (`task -o`), ensure `epic_key` is empty

## Acceptance Criteria
- Creating epic/task/subtask produces valid local JSON and Jira ticket
- Orphan task creation does not set `epic_key`

## Next Tasks
- 10-command-aux.md 