# JIT - Jira In Terminal

A CLI-first, markdown-native workflow tool for managing Jira tickets with speed, clarity, and flow.

## 🔥 Overview

`jit` is a local-first CLI tool that provides a Git-like workflow for Jira tickets. It's designed to complement, not replace, Jira by offering efficient terminal-based workflows for developers.

**Key Features:**
* Git-like commands for Jira ticket management
* AI-powered ticket enhancement with customizable templates
* Local-first architecture with offline capability
* Context-aware workflow that maintains focus
* Markdown-native editing experience
* Seamless Jira integration with bi-directional sync

## 🚀 Quick Start

### Installation

```bash
# Install directly from source
go install -buildvcs=false github.com/lunchboxsushi/jit/cmd/jit@latest

# Or clone and build
git clone https://github.com/lunchboxsushi/jit.git
cd jit
./build.sh  # Uses build script to handle VCS issues
# Or manually: go install -buildvcs=false ./cmd/jit
```
> **Note:** Ensure your Go binary path (`$GOPATH/bin` or `$HOME/go/bin`) is in your system's `PATH`.

### Shell Completion

Enable tab completion for faster command entry:

**Bash:**
```bash
# Temporary (current session only)
source <(jit completion bash)

# Permanent (add to ~/.bashrc)
echo 'source <(jit completion bash)' >> ~/.bashrc
```

**Zsh:**
```bash
# Temporary (current session only)
source <(jit completion zsh)

# Permanent (add to ~/.zshrc)
echo 'source <(jit completion zsh)' >> ~/.zshrc
```

**Fish:**
```bash
# Temporary (current session only)
jit completion fish | source

# Permanent
jit completion fish > ~/.config/fish/completions/jit.fish
```

**PowerShell:**
```powershell
# Temporary (current session only)
jit completion powershell | Out-String | Invoke-Expression

# Permanent
jit completion powershell > jit.ps1
# Add to your PowerShell profile
```

### Initial Setup

```bash
# Initialize configuration
jit init

# Follow the interactive setup to configure Jira connection and AI provider
```

### Basic Workflow

```bash
# 1. Track an existing epic and its children
jit track SRE-5344

# 2. Focus on a specific ticket
jit focus "5433"           # Fuzzy search by key or title

# 3. Create a new task under the current epic
jit task

# 4. Create a new subtask under the current task
jit subtask

# 5. View your ticket tree
jit log

# 6. Add a comment to the current focus ticket
jit comment

# 7. Quick actions
jit link                   # Get Jira URL
jit link -s               # Copy ticket number
jit open                  # Open in browser
```

## 📖 Commands Reference

### Core Workflow Commands

#### `jit track <ticket-key>`
Downloads a Jira ticket and all its children (Epic → Tasks → Subtasks) for local work.
```bash
jit track SRE-5344        # Track epic and all children
```

#### `jit focus <query>`
Switch your working context using fuzzy search.
```bash
jit focus "5433"          # Match SRE-5433 by partial key
jit focus "tracing"       # Match by title content
```

#### `jit log`
Display the hierarchical view of all tracked tickets with current focus marked.
```bash
jit log                   # Show ticket tree with focus indicator
```

### Ticket Creation Commands

#### `jit epic`
Create a new epic. Opens editor for content, optionally enhances with AI, and creates Jira ticket.
```bash
jit epic                  # Create new epic
jit epic --no-enrich     # Skip AI enhancement
jit epic --no-create     # Create locally only
```

#### `jit task`
Create a new task under the current epic.
```bash
jit task                  # Create task under current epic
```

#### `jit subtask`
Create a new subtask under the current task.
```bash
jit subtask               # Create subtask under current task
```

### Context-Aware Operations

#### `jit link`
Get Jira URL for the current focus ticket.
```bash
jit link                  # Display full Jira URL
jit link -s               # Copy ticket key to clipboard
jit link -f               # Copy full URL to clipboard
```

#### `jit open`
Open the current focus ticket in your default browser.
```bash
jit open                  # Open current ticket
jit open -e               # Open in edit mode
jit open -c               # Open comments section
```

#### `jit comment`
Add a comment to the current focus ticket.
```bash
jit comment               # Open editor for comment
jit comment "Quick update"  # Add inline comment
jit comment --enrich      # AI-enhanced comment
jit comment --no-create   # Draft without posting
```

### Configuration Commands

#### `jit init`
Interactive setup wizard for first-time configuration.
```bash
jit init                  # Set up Jira connection, AI provider, etc.
```

## 🗂️ Data Structure

### Storage Locations
- **Configuration**: `~/.jit/config.yml`
- **Data**: `~/.local/share/jit/`
- **Templates**: `~/.jit/templates/`

### Local Data Organization
```text
~/.local/share/jit/
├── tickets/
│   ├── SRE-5344.json                  # Epic ticket data
│   ├── SRE-5345.json                  # Task ticket data
│   └── SRE-5346.json                  # Subtask ticket data
├── context.json                       # Current focus context
├── cache/                             # Cached project data
└── logs/                              # Application logs
```

### Ticket Hierarchy
```
Epic (SRE-5344)
├── Task (SRE-5345)
│   ├── Subtask (SRE-5346)
│   └── Subtask (SRE-5347)
└── Task (SRE-5348)
    └── Subtask (SRE-5349)
```

## ⚙️ Configuration

### Configuration File (`~/.jit/config.yml`)
```yaml
# Jira Configuration
jira:
  url: "https://company.atlassian.net"
  username: "your-email@company.com"
  token: "${JIRA_API_TOKEN}"
  project: "SRE"
  epic_link_field: "customfield_10014"

# AI Configuration
ai:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"
  model: "gpt-4"
  max_tokens: 1000
  temperature: 0.7

# Application Settings
app:
  data_dir: "~/.local/share/jit"
  default_editor: "vim"
  review_before_create: true
  auto_sync: false
```

### Environment Variables
```bash
export JIRA_API_TOKEN="your-jira-token"
export OPENAI_API_KEY="your-openai-key"
export JIT_EDITOR="code"              # Override default editor
```

### Finding Your Epic Link Field ID
1. Go to **Jira Administration → Issues → Custom Fields**
2. Search for **Epic Link**
3. Click the gear icon and select **View field information**
4. Note the field ID (e.g., `customfield_10014`)

## 🤖 AI Enhancement

### Custom Templates
Create custom AI prompt templates in `~/.jit/templates/prompts/`:

**Example Epic Template:**
```markdown
You are helping create a well-structured Jira Epic for a software development team.

Title: {{TITLE}}
Raw Content: {{RAW_CONTENT}}

Please enhance this epic with:
1. Clear problem statement
2. Success criteria
3. High-level approach
4. Dependencies and risks: {{What are the main technical risks for this epic?}}
5. Estimated timeline

Keep the tone professional but concise.
```

### AI Expressions
Use `{{expression}}` for dynamic AI-generated content:
- `{{What are the main risks for this change?}}`
- `{{List potential dependencies for this feature}}`
- `{{Suggest 3 acceptance criteria for this task}}`

## 🔄 Typical Workflows

### Starting New Work
```bash
# Create and work on a new epic
jit epic                  # Create epic, AI enhances, sets focus
jit task                  # Create task under epic
jit subtask               # Create subtask under task
jit comment "Started implementation"
```

### Continuing Existing Work
```bash
# Resume work on existing tickets
jit track PROJ-1234      # Track epic and children
jit focus "login"        # Switch to login-related task
jit subtask              # Add new subtask
jit open                 # Review in browser
```

### Quick Operations
```bash
# Fast ticket operations
jit link -s              # Copy ticket number
jit comment "Status update"
jit log                  # Check current state
```

## 🛠️ Development

### Prerequisites
- Go 1.21+
- Jira Cloud instance with API access
- OpenAI API key (optional, for AI features)

### Building from Source
```bash
git clone https://github.com/lunchboxsushi/jit.git
cd jit
go build -o jit ./cmd/jit
```

### Running Tests
```bash
go test ./...
```

## 📚 Examples

### Epic Creation Example
```bash
$ jit epic
# Opens editor with template:
# ---
# Title: User Authentication Refactor
# 
# Description:
# Modernize our authentication system to improve security and user experience.
# 
# Goals:
# - Implement OAuth 2.0
# - Add multi-factor authentication
# - Improve password policies
# ---
# Save and close editor
✓ Epic created: AUTH-123 "User Authentication Refactor"
✓ Focus set to AUTH-123
```

### Task Creation Example
```bash
$ jit task
# (must be focused on an epic first)
# Opens editor for task description
# AI enhances content with acceptance criteria
✓ Task created: AUTH-124 "Implement OAuth 2.0 Integration"
✓ Focus set to AUTH-124
✓ Linked to epic AUTH-123
```

## 🎯 Philosophy

`jit` follows these principles:

1. **Terminal-First**: Optimized for keyboard-driven workflows
2. **Context-Aware**: Maintains focus on current work
3. **Git-Like**: Familiar commands for developers
4. **Local-First**: Work offline, sync when ready
5. **AI-Enhanced**: Intelligent assistance without replacement
6. **Markdown-Native**: Natural editing experience

## 🗺️ Roadmap

- [ ] Real-time sync with Jira webhooks
- [ ] Plugin system for extensibility
- [ ] Advanced search and filtering
- [ ] Team collaboration features
- [ ] Mobile companion app
- [ ] Integration with popular IDEs
- [ ] Custom field mapping
- [ ] Bulk operations support

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by Git's elegant command structure
- Built with the excellent [Cobra](https://github.com/spf13/cobra) CLI framework
- Powered by AI for intelligent enhancement
- Integrates with Atlassian Jira Cloud API 