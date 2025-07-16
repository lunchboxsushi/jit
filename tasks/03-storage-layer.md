# Task 03: JSON Storage Layer

**Status:** [ready]

## Overview
Create the storage subsystem responsible for persisting tickets, context, and cache data as JSON files under `~/.local/share/jit/`.

## Objectives
- Define a `Storage` interface with basic CRUD operations
- Implement JSON file storage:
  - Atomic write (temp file + rename)
  - Safe concurrent read/write using file locks
- Helper functions for ticket path resolution (by key)
- Automatic directory creation on first use
- Backup strategy (retain last N versions)

## Deliverables
- [ready] `internal/storage/storage.go` (interface + common helpers)
- [ready] `internal/storage/json.go` (JSON implementation)
- [ready] `internal/storage/context.go` (context.json helpers)
- [ready] Unit tests covering save/load ticket and context

## Dependencies
- Task 01 structs
- Task 02 config (for data_dir path)

## Implementation Notes
- Use `os.UserCacheDir` fallback if env not set
- Use `flock` or golang `syscall` for file locking (cross-platform)
- Consider configurable backup retention in config

## Acceptance Criteria
- `go test ./internal/storage/...` passes
- Tickets and context can be saved and reloaded accurately
- Concurrent read/write tests do not corrupt data

## Next Tasks
- 04-jira-client.md 