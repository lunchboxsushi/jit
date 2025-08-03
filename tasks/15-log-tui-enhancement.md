# Task 15: Log Command TUI Enhancement

**Status:** [ready]

## Overview
Enhance the `jit log` command with a modern TUI (Terminal User Interface) library to render ticket hierarchies as interactive graphs, similar to Git's graph visualization. Implement color-coded ticket types and focus indicators for better visual hierarchy.

## Objectives
- Research and select a suitable TUI library for Go
- Redesign the log command output as an interactive graph
- Implement color-coded ticket types (epic=purple, task=blue, subtask=light-blue)
- Add focus indicator with orange color and @ symbol
- Create interactive navigation within the graph
- Support orphaned tasks with `--orphan` flag and dedicated section
- Maintain backward compatibility with existing flags

## Deliverables
- [ready] TUI library evaluation and selection
- [ready] Enhanced `internal/commands/log.go` with TUI integration
- [ready] Graph rendering engine for ticket hierarchies
- [ready] Color scheme implementation
- [ready] Interactive navigation features
- [ready] Orphaned tasks section with `== ORPHAN TASKS ==` header
- [ready] Updated documentation for new log features

## Dependencies
- Task 00-13 (all previous tasks completed)
- Existing log command functionality
- Ticket hierarchy data structures

## Implementation Notes

### TUI Library Candidates
1. **Bubble Tea** (charmbracelet/bubbletea) - Popular, well-maintained, good for interactive UIs
2. **Termui** - Feature-rich, good for dashboards and data visualization
3. **Gocui** - Lightweight, good for simple interfaces
4. **Tview** - Rich text UI library with good widget support

### Output Format
- **Focus indicator**: `@` symbol at the start of the line (orange color)
- **Ticket type**: Epic/Task/Subtask prefix
- **Ticket key**: `[PROJ-100]` format
- **Status**: `<In-Progress>` format with status colors
- **Title**: `- User Authentication` format
- **Color coding**: Purple (Epic), Blue (Task), Light Blue (Subtask)

### Color Scheme
- **Epics**: Purple (#8B5CF6) - High-level planning items
- **Tasks**: Blue (#3B82F6) - Main work items
- **Subtasks**: Light Blue (#60A5FA) - Detailed work items
- **Focus**: Orange (#F97316) with @ symbol - Current context
- **Status colors**: 
  - Green (#10B981) - Done
  - Yellow (#F59E0B) - In Progress
  - Red (#EF4444) - Blocked
  - Gray (#6B7280) - To Do

### Graph Structure
```
@─ Epic [PROJ-100]<In-Progress> - User Authentication (purple)
├─ Task [PROJ-101]<To Do> - OAuth Implementation (blue)
│  ├─ Subtask [PROJ-102]<Done> - Google OAuth (light-blue)
│  └─ Subtask [PROJ-103]<In-Progress> - GitHub OAuth (light-blue)
└─ Task [PROJ-104]<Blocked> - MFA Setup (blue)
   └─ Subtask [PROJ-105]<To Do> - TOTP Implementation (light-blue)

─ Epic [PROJ-200]<Done> - Database Migration (purple)
└─ Task [PROJ-201]<In-Progress> - Schema Updates (blue)

== ORPHAN TASKS ==
─ Task [PROJ-300]<To Do> - Standalone Feature (blue)
└─ Subtask [PROJ-301]<In-Progress> - Implementation Details (light-blue)
```

### Interactive Features
- **Navigation**: Arrow keys to move between tickets
- **Selection**: Enter to focus on selected ticket
- **Expand/Collapse**: Space to toggle ticket children visibility
- **Filtering**: Type to search within visible tickets
- **Status toggle**: 's' to cycle through status filters
- **View modes**: 't' for tree, 'l' for list, 'g' for graph

### Command Flags
- `--tui` - Enable TUI mode (default for interactive terminals)
- `--no-tui` - Force text-only output
- `--colors` - Enable/disable colors
- `--interactive` - Enable interactive navigation
- `--compact` - Compact tree view
- `--full` - Show all details
- `--orphan` - Show orphaned tasks (tasks without parent epic)

## Acceptance Criteria
- [ ] TUI library selected and integrated
- [ ] Graph renders correctly with proper hierarchy
- [ ] Color coding works for all ticket types
- [ ] Focus indicator clearly visible with @ symbol
- [ ] Interactive navigation functional
- [ ] Backward compatibility maintained
- [ ] Performance acceptable for large hierarchies
- [ ] Cross-platform compatibility (Linux, macOS, Windows)

## Technical Requirements

### Performance
- Handle hierarchies with 100+ tickets efficiently
- Smooth scrolling and navigation
- Responsive to user input

### Accessibility
- Support for non-color terminals
- Keyboard navigation for screen readers
- Clear visual hierarchy without colors

### Compatibility
- Work with existing `jit log` flags
- Maintain JSON output option
- Support for piping to other commands

## Testing Strategy
- Unit tests for graph rendering logic
- Integration tests for TUI interactions
- Performance tests with large datasets
- Cross-platform testing
- Accessibility testing

## Success Metrics
- Visual clarity improvement over current text output
- User satisfaction with interactive features
- Performance benchmarks met
- Zero regression in existing functionality

## Next Tasks
- Task 16: Advanced filtering and search
- Task 17: Export and reporting features
- Task 18: Team collaboration features

## Implementation Phases

### Phase 1: Foundation
- [ ] Research and select TUI library
- [ ] Create basic graph rendering engine
- [ ] Implement color scheme
- [ ] Add focus indicator

### Phase 2: Interactivity
- [ ] Implement keyboard navigation
- [ ] Add expand/collapse functionality
- [ ] Create search and filtering
- [ ] Add view mode switching

### Phase 3: Polish
- [ ] Performance optimization
- [ ] Cross-platform testing
- [ ] Documentation updates
- [ ] User testing and feedback

## Done Signal
When the log command provides an intuitive, interactive graph view that clearly shows ticket hierarchies with proper color coding and focus indicators, making it easy for users to navigate and understand their project structure at a glance. 