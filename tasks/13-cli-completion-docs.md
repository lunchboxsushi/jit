# Task 13: CLI Completion & Developer Docs

**Status:** [ready]

## Overview
Add shell completion scripts for popular shells and generate developer documentation to improve usability and onboarding.

## Objectives
- Generate completion scripts via Cobra for:
  - Bash
  - Zsh
  - Fish
- Provide `jit completion <shell>` command to output/install scripts
- Add `docs/commands.md` auto-generated from Cobra command tree
- Update README with installation instructions for completion scripts

## Deliverables
- [ready] `internal/commands/completion.go` implementing completion generation
- [ready] Generated docs in `docs/commands.md` via `cobra doc` or custom
- [ready] README updates with new section

## Dependencies
- Task 00 CLI foundation
- All command files registered with root

## Implementation Notes
- Use `cmd/root.go` to add hidden `completion` command
- Provide `make docs` to regenerate command docs

## Acceptance Criteria
- `jit completion bash | source` enables tab completion
- Documentation file lists all commands with descriptions

## Done Signal
When this task is complete, all primary commands and supporting developer tooling are implemented. 🎉 