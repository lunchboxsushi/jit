# Task 08: Command `jit focus`

**Status:** [ready]

## Overview
Implement the `jit focus <query>` command that switches the current context using fuzzy search on ticket keys and titles.

## Objectives
- Accept query string (partial key or title)
- Search local tickets via Storage layer (simple in-memory index for now)
- Use fuzzy matching to rank results (utils/fuzzy.go)
- If multiple matches, prompt user to select
- Update context to selected ticket

## Deliverables
- [ready] `internal/commands/focus.go` implementing command
- [ready] `internal/utils/fuzzy.go` with search helper (if not already done)
- [ready] Unit tests for matching logic

## Dependencies
- Task 03 Storage layer
- Task 05 Context management

## Implementation Notes
- Respect `--type` flag to limit search (epic|task|subtask)
- Option `--list` to just list matches without switching

## Acceptance Criteria
- `jit focus "5344"` correctly focuses SRE-5344
- Ambiguous queries prompt user

## Next Tasks
- 09-command-create.md 