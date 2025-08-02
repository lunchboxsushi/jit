# Task 03: JSON Storage Layer

**Status:** [completed]

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
- [completed] `internal/storage/storage.go` (interface + common helpers)
- [completed] `internal/storage/file_storage.go` (JSON implementation)
- [completed] `internal/storage/context.go` (context.json helpers)
- [completed] Unit tests covering save/load ticket and context

## Dependencies
- Task 01 structs
- Task 02 config (for data_dir path)

## Implementation Notes
- Use `os.UserCacheDir` fallback if env not set
- Use `sync.RWMutex` for thread-safe concurrent access
- Atomic writes using temporary files and rename operations
- Simplified approach without per-file backups (atomic writes provide sufficient safety)
- Automatic directory creation with proper permissions

## Acceptance Criteria
- [completed] `go test ./internal/storage/...` passes
- [completed] Tickets and context can be saved and reloaded accurately
- [completed] Concurrent read/write tests do not corrupt data
- [completed] Atomic writes prevent corruption during file operations
- [completed] Context management with focus tracking and recent tickets

## Next Tasks
- 04-jira-client.md 