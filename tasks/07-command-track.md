# Task 07: Command `jit track`

**Status:** [ready]

## Overview
Implement the `jit track <ticket-key>` command that imports a Jira ticket (and optionally its hierarchy) into local storage and sets focus.

## Objectives
- Accept ticket key argument, validate format
- Use Jira client to fetch ticket details
- If ticket is Epic, recursively fetch child tasks & subtasks
- Convert Jira JSON to internal Ticket structs
- Persist tickets via Storage layer
- Update context to tracked ticket
- Flags:
  - `--no-children` to skip recursive fetch

## Deliverables
- [ready] `internal/commands/track.go` implementing command logic
- [ready] Unit tests with mocked Jira client and storage

## Dependencies
- Task 04 Jira client
- Task 03 Storage layer
- Task 05 Context management

## Implementation Notes
- Show progress spinner for large hierarchies
- Deduplicate tickets already in storage (update if outdated)

## Acceptance Criteria
- `jit track SRE-1234` stores ticket JSON locally and sets focus
- Recursion fetch confirmed via unit tests

## Next Tasks
- 08-command-focus.md 