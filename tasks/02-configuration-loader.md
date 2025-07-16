# Task 02: Configuration Loader

**Status:** [ready]

## Overview
Implement the configuration subsystem that loads `~/.jit/config.yml`, applies defaults, supports environment variable overrides, and validates required fields.

## Objectives
- Parse YAML configuration into Go structs (from Task 01)
- Environment variable expansion (e.g. `${JIRA_API_TOKEN}`)
- Sensible defaults for optional fields
- Validation of required fields (Jira URL, token, etc.)
- Provide package-level accessor to obtain a singleton Config instance

## Deliverables
- [ready] `internal/config/config.go` with `Load()` and `Get()` functions
- [ready] `internal/config/defaults.go` with default constants
- [ready] `internal/config/validation.go` with validation logic
- [ready] Unit tests covering success, missing fields, bad YAML

## Dependencies
- gopkg.in/yaml.v3 for YAML parsing
- Task 01 structs

## Implementation Notes
- Respect `XDG_CONFIG_HOME`; default to `~/.jit/config.yml`
- Use `sync.Once` to cache the loaded config
- Fail fast with descriptive errors on invalid config

## Acceptance Criteria
- `go test ./internal/config/...` passes
- Running `jit` without config prompts user to run `jit init`

## Next Tasks
- 03-storage-layer.md 