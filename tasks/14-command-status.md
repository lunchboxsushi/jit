# Task 14: Command `jit status` — View & Transition Ticket State

**Status:** [ready]

## Overview
Provide a unified command to display the current ticket state and move it through the Jira workflow using concise flags. Mirrors the typical start/stop/blocked flow for day-to-day development.

## Desired UX
```bash
jit status                 # Show current focus ticket, its status, assignee, etc.
jit status --start         # Move ticket to *In Progress*
jit status --complete      # Move ticket to *Done*
jit status --blocked       # Move ticket to *Blocked*
# Optional extras
jit status --review        # Move ticket to *Code Review*
jit status --custom "QA"   # Move to arbitrary status
```

## Objectives
1. **Display mode (default)**
   - Show key information about the current focus ticket:
     - Key, Title, Type, Status, Assignee, Story Points
     - Parent/children with their statuses (compact tree)
2. **Transition mode (mutually-exclusive flags)**
   - `--start`, `--complete`, `--blocked`, `--custom <status>`
   - Resolve desired Jira status → transition ID (see below)
   - Invoke Jira transitions endpoint to change status
   - Update local ticket JSON with new status and `updated` timestamp
3. **Configurable mapping**
   - Add section to config:
     ```yaml
     status_aliases:
       start: "In Progress"
       complete: "Done"
       blocked: "Blocked"
     ```
   - Fallback to literal flag name if alias not present
4. **Local-only option**
   - `--no-sync` flag updates local JSON without Jira call (offline work)

## Deliverables
- [ready] `internal/commands/status.go` implementing CLI logic
- [ready] Helper `internal/jira/transitions.go` to fetch & execute transitions
- [ready] Extend config struct with `StatusAliases` map
- [ready] Unit tests:
  - Transition resolution logic (mock Jira transitions)
  - Local JSON status update
  - Display output formatting

## Dependencies
- Task 04 Jira client (add transitions helper)
- Task 03 Storage layer (load/save ticket)
- Task 05 Context management (current focus)
- Task 02 Config loader (new field)

## Implementation Notes
- Use Jira endpoint: `GET /rest/api/3/issue/{issueIdOrKey}/transitions` then `POST` to perform.
- Cache transition list per ticket to reduce API calls (TTL 10 min).
- If status already equals desired, short-circuit with message.
- For recursive updates (future), could add `--all-children` flag.

## Acceptance Criteria
- `jit status` shows readable summary
- `jit status --start` moves ticket to *In Progress* in Jira and locally
- Configurable aliases work; `--custom` can override
- Unit tests cover edge cases (unknown status, no available transition)

## Next Tasks
- Implementation cycle starting back at Task 00 if needed. 