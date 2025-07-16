# Task 06: Command `jit init`

**Status:** [ready]

## Overview
Implement the interactive setup wizard that generates a new configuration file and initializes the data directory structure.

## Objectives
- CLI command `jit init`
- Prompt user for Jira URL, email, token, project key, AI provider, etc.
- Create `~/.jit/config.yml` with provided info and defaults
- Create data directories (`~/.local/share/jit/`)
- Validate Jira credentials by pinging the API (optional flag `--no-verify`)

## Deliverables
- [ready] `internal/commands/init.go` implementing the command
- [ready] Template config file for guided editing
- [ready] Unit tests using pexpect-like integration (optional)

## Dependencies
- Task 02 Configuration loader
- Cobra CLI foundation (Task 00)

## Implementation Notes
- Use survey/v2 or simple stdin prompts
- Allow non-interactive flags (`--url`, `--token`, etc.) for scripting
- Respect `--force` to overwrite existing config

## Acceptance Criteria
- Running `jit init` creates config and directories
- Running again without `--force` aborts with warning
- Config passes validation (Task 02)

## Next Tasks
- 07-command-track.md 