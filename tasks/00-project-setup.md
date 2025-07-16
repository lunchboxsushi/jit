# Task 00: Project Setup & Foundation

**Status:** [ready]

## Overview
Set up the basic Go project structure with proper module initialization and core dependencies.

## Objectives
- Initialize Go module
- Set up basic project structure
- Add core dependencies (Cobra, YAML, etc.)
- Create basic main.go entry point

## Deliverables
- [ready] Initialize `go.mod` with proper module name
- [ready] Create basic CLI structure with Cobra
- [ready] Set up cmd/jit/main.go entry point
- [ready] Create cmd/root.go with basic command setup
- [ready] Add core dependencies to go.mod

## Dependencies
- Go 1.21+
- github.com/spf13/cobra (CLI framework)
- gopkg.in/yaml.v3 (YAML parsing)

## Implementation Notes
- Use `github.com/lunchboxsushi/jit` as module name
- Set up basic command structure that can be extended
- Ensure proper error handling from the start
- Add basic version command

## Acceptance Criteria
- `go build` succeeds without errors
- `jit --help` shows basic help output
- `jit version` shows version information
- Project structure matches architecture document

## Next Tasks
- 01-data-types.md
- 02-configuration.md 