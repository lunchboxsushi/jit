# Task 02: Configuration Loader

**Status:** [completed]

## Overview
Implement the configuration subsystem that loads `~/.jit/config.yml`, applies defaults, supports environment variable overrides, and validates required fields.

## Objectives
- Parse YAML configuration into Go structs (from Task 01)
- Environment variable expansion (e.g. `${JIRA_API_TOKEN}`)
- Sensible defaults for optional fields
- Validation of required fields (Jira URL, token, etc.)
- Provide package-level accessor to obtain a singleton Config instance

## Deliverables
- [completed] `internal/config/config.go` with `Load()` and `Get()` functions
- [completed] `internal/config/defaults.go` with path resolution functions
- [completed] `internal/config/validation.go` with validation logic
- [completed] Unit tests covering success, missing fields, bad YAML, environment variables

## Dependencies
- gopkg.in/yaml.v3 for YAML parsing
- Task 01 structs

## Implementation Notes
- Respect `XDG_CONFIG_HOME`; default to `~/.jit/config.yml`
- Use `sync.Once` to cache the loaded config
- Fail fast with descriptive errors on invalid config
- No default values applied - all values must be explicitly set in config
- Environment variable expansion for sensitive data
- Comprehensive validation with detailed error messages

## Acceptance Criteria
- [completed] `go test ./internal/config/...` passes
- [completed] Running `jit` without config prompts user to run `jit init`
- [completed] Environment variable expansion works correctly
- [completed] Validation catches all required field errors
- [completed] Singleton pattern ensures config is loaded only once

## Next Tasks
- 03-storage-layer.md 