# jit CLI Commands

`jit` is a local-first CLI tool that lets developers write tasks and sub-tasks in markdown, auto-enriches raw task descriptions into manager-optimized Jira tickets, and syncs with Jira to reflect status, updates, and structure.

## Overview

```bash
jit [command] [flags]
```

## Commands

### `init`
Initialize jit configuration and data directories.

```bash
jit init [flags]
```

**Flags:**
- `--config-path string` - Path to config file (default: ~/.jit/config.yml)
- `--data-path string` - Path to data directory (default: ~/.jit/data)

**Description:**
Creates the default configuration file and data directories. This is typically the first command you run to set up jit.

**Example:**
```bash
jit init
```

### `track <ticket-key>`
Track a Jira ticket and its hierarchy locally.

```bash
jit track <ticket-key> [flags]
```

**Flags:**
- `--recursive` - Track all children recursively (default: true)

**Description:**
Downloads a Jira ticket and its entire hierarchy (epic → tasks → subtasks) to your local workspace. This creates a local copy that you can work with offline.

**Example:**
```bash
jit track PROJ-123
```

### `focus <query>`
Set focus to a specific ticket or search for tickets.

```bash
jit focus <query> [flags]
```

**Flags:**
- `--exact` - Use exact match instead of fuzzy search

**Description:**
Searches through your locally tracked tickets and sets the current focus. You can use fuzzy search to find tickets by key, title, or description.

**Examples:**
```bash
jit focus PROJ-123          # Focus by exact ticket key
jit focus "user login"      # Fuzzy search by title/description
```

### `epic`
Create a new epic in Jira.

```bash
jit epic [flags]
```

**Flags:**
- `--no-enrich` - Skip AI enrichment of description
- `--priority string` - Set priority (low, medium, high, highest) (default: medium)

**Description:**
Opens an editor to create a new epic. The epic will be created in Jira and tracked locally.

**Example:**
```bash
jit epic
```

### `task`
Create a new task in Jira.

```bash
jit task [flags]
```

**Flags:**
- `--epic string` - Parent epic key
- `--no-enrich` - Skip AI enrichment of description
- `--orphan` - Create task without parent epic
- `--priority string` - Set priority (low, medium, high, highest) (default: medium)

**Description:**
Opens an editor to create a new task. The task will be created in Jira and tracked locally.

**Examples:**
```bash
jit task                    # Create task in current epic context
jit task --epic PROJ-100    # Create task in specific epic
jit task --orphan           # Create standalone task
```

### `subtask`
Create a new subtask in Jira.

```bash
jit subtask [flags]
```

**Flags:**
- `--task string` - Parent task key
- `--no-enrich` - Skip AI enrichment of description
- `--priority string` - Set priority (low, medium, high, highest) (default: medium)

**Description:**
Opens an editor to create a new subtask. The subtask will be created in Jira and tracked locally.

**Examples:**
```bash
jit subtask                 # Create subtask in current task context
jit subtask --task PROJ-200 # Create subtask in specific task
```

### `log`
Display the hierarchy of tracked tickets.

```bash
jit log [flags]
```

**Flags:**
- `--all` - Show all tickets, not just current focus hierarchy
- `--json` - Output in JSON format
- `--status string` - Filter by status

**Description:**
Shows a tree view of your tracked tickets, highlighting the current focus and showing the hierarchy structure.

**Examples:**
```bash
jit log                     # Show current focus hierarchy
jit log --all               # Show all tracked tickets
jit log --json              # Output as JSON
jit log --status "In Progress"
```

### `link`
Get the Jira URL for a ticket.

```bash
jit link [ticket-key] [flags]
```

**Flags:**
- `--short` - Output only the URL (no description)

**Description:**
Generates and optionally copies the Jira URL for a ticket to your clipboard.

**Examples:**
```bash
jit link                    # Get URL for current focus
jit link PROJ-123           # Get URL for specific ticket
jit link --short            # Output only the URL
```

### `open`
Open a Jira ticket in your default browser.

```bash
jit open [ticket-key]
```

**Description:**
Opens the Jira ticket in your default web browser.

**Examples:**
```bash
jit open                    # Open current focus ticket
jit open PROJ-123           # Open specific ticket
```

### `comment`
Add a comment to a Jira ticket.

```bash
jit comment [ticket-key] [comment-text] [flags]
```

**Flags:**
- `--message, -m` - Add inline comment (requires comment text)

**Description:**
Adds a comment to a Jira ticket. You can provide the comment text inline or open an editor to write a longer comment.

**Examples:**
```bash
jit comment "Updated the API endpoint"     # Inline comment
jit comment -m "Fixed the bug"             # Inline comment with flag
jit comment                                # Open editor for comment
```

### `completion`
Generate shell completion scripts.

```bash
jit completion [bash|zsh|fish|powershell]
```

**Description:**
Generates shell completion scripts for bash, zsh, fish, or PowerShell. This enables tab completion for jit commands.

**Examples:**
```bash
jit completion bash | source              # Enable bash completion
jit completion zsh > ~/.zshrc             # Install zsh completion
jit completion fish > ~/.config/fish/completions/jit.fish
```

### `version`
Show version information.

```bash
jit version
```

**Description:**
Displays the current version of jit.

## Global Flags

- `--help, -h` - Show help for the command
- `--version` - Show version information

## Configuration

jit uses a configuration file located at `~/.jit/config.yml`. You can customize:

- **Jira Settings**: URL, username, API token, project key
- **AI Settings**: Provider (openai, mock), API key, model
- **Editor Settings**: Default editor for creating tickets
- **Storage Settings**: Data directory location

## Data Storage

jit stores local data in `~/.jit/data/`:

- `tickets/` - Local copies of Jira tickets
- `context.json` - Current focus and recent tickets
- `config.yml` - Configuration file

## Examples

### Basic Workflow

```bash
# Initialize jit
jit init

# Track an epic and its children
jit track PROJ-100

# Focus on a specific task
jit focus PROJ-101

# Create a new subtask
jit subtask

# Add a comment
jit comment "Completed the implementation"

# View the hierarchy
jit log

# Open in browser
jit open
```

### Advanced Usage

```bash
# Create a standalone task
jit task --orphan

# Create task in specific epic
jit task --epic PROJ-100

# Search for tickets
jit focus "authentication"

# Get ticket URL
jit link PROJ-123

# Enable shell completion
echo 'source <(jit completion bash)' >> ~/.bashrc
```

## Troubleshooting

### Common Issues

1. **Configuration not found**: Run `jit init` to create the default configuration
2. **Jira connection failed**: Check your Jira URL and API token in the config
3. **Permission denied**: Ensure the data directory is writable
4. **Editor not found**: Set the correct editor path in your configuration

### Getting Help

- Use `jit --help` for general help
- Use `jit <command> --help` for command-specific help
- Check the configuration file at `~/.jit/config.yml` 