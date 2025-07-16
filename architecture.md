# JIT Architecture

## Overview
`jit` is a local-first CLI application that provides a Git-like workflow for Jira tickets. The architecture emphasizes simplicity, performance, and offline capability while maintaining seamless integration with Jira.

## Project Structure

```
jit/
├── cmd/
│   ├── jit/
│   │   └── main.go                    # CLI entry point
│   └── root.go                        # Root command setup
├── internal/
│   ├── commands/                      # Command implementations
│   │   ├── init.go                    # jit init
│   │   ├── track.go                   # jit track
│   │   ├── focus.go                   # jit focus
│   │   ├── epic.go                    # jit epic
│   │   ├── task.go                    # jit task (includes orphan tasks)
│   │   ├── subtask.go                 # jit subtask
│   │   ├── log.go                     # jit log
│   │   ├── link.go                    # jit link
│   │   ├── open.go                    # jit open
│   │   └── comment.go                 # jit comment
│   ├── config/                        # Configuration management
│   │   ├── config.go                  # Config struct and loading
│   │   ├── validation.go              # Config validation
│   │   └── defaults.go                # Default configuration
│   ├── jira/                          # Jira API integration
│   │   ├── client.go                  # Jira API client
│   │   ├── auth.go                    # Authentication handling
│   │   ├── tickets.go                 # Ticket operations
│   │   └── types.go                   # Jira data types
│   ├── ai/                            # AI integration
│   │   ├── provider.go                # AI provider interface
│   │   ├── openai.go                  # OpenAI implementation
│   │   └── templates.go               # Template processing
│   ├── storage/                       # Data persistence
│   │   ├── storage.go                 # Storage interface
│   │   ├── json.go                    # JSON file storage
│   │   ├── context.go                 # Context management
│   │   └── search.go                  # Search functionality
│   ├── ui/                            # User interface
│   │   ├── editor.go                  # Editor integration
│   │   ├── display.go                 # Display formatting
│   │   └── prompts.go                 # Interactive prompts
│   └── utils/                         # Utility functions
│       ├── fuzzy.go                   # Fuzzy search
│       ├── clipboard.go               # Clipboard operations
│       └── filesystem.go              # File operations
├── pkg/                               # Public packages
│   └── types/                         # Shared types
│       ├── ticket.go                  # Ticket data structures
│       ├── context.go                 # Context types
│       └── config.go                  # Configuration types
├── templates/                         # Default templates
│   ├── epic.md                        # Epic template
│   ├── task.md                        # Task template
│   ├── subtask.md                     # Subtask template
│   └── prompts/                       # AI prompt templates
│       ├── epic_enrichment.txt
│       ├── task_enrichment.txt
│       └── subtask_enrichment.txt
├── go.mod                             # Go module definition
├── go.sum                             # Go dependencies
├── README.md                          # Project documentation
├── features.md                        # Feature specification
└── architecture.md                    # This file
```

## Data Architecture

### Storage Structure
- **Configuration**: `~/.jit/config.yml`
- **Data Directory**: `~/.local/share/jit/`
- **Templates**: `~/.jit/templates/`

```
~/.local/share/jit/
├── tickets/                           # Individual ticket JSON files
│   ├── SRE-5344.json                  # Epic
│   ├── SRE-5345.json                  # Task (with epic_key)
│   ├── SRE-5346.json                  # Subtask (with parent_key)
│   └── PROJ-1234.json                 # Orphan task (no epic_key)
├── context.json                       # Current focus context
└── cache/                             # Project metadata cache
    ├── projects.json
    └── users.json
```

### Ticket Data Format (JSON)
```json
{
  "key": "SRE-5344",
  "title": "Implement distributed tracing",
  "type": "Epic",                      # Epic, Task, Subtask
  "status": "In Progress",
  "priority": "High",
  "description": "Full description with markdown formatting preserved",
  "metadata": {
    "project": "SRE",
    "assignee": "john.doe@company.com",
    "created": "2024-01-15T10:30:00Z",
    "updated": "2024-01-20T14:45:00Z",
    "labels": ["infrastructure", "observability"]
  },
  "relationships": {
    "epic_key": "",                    # Empty for epics and orphan tasks
    "parent_key": "",                  # Parent task for subtasks
    "children": ["SRE-5345", "SRE-5346"] # Child tickets
  },
  "jira_data": {
    "url": "https://company.atlassian.net/browse/SRE-5344",
    "custom_fields": {
      "story_points": 8,
      "epic_link": "SRE-5344"
    }
  },
  "local_data": {
    "last_sync": "2024-01-20T14:45:00Z",
    "local_changes": false,
    "ai_enhanced": true
  }
}
```

### Context Management (`context.json`)
```json
{
  "current_epic": "SRE-5344",
  "current_task": "SRE-5345",
  "current_subtask": "SRE-5346",
  "last_updated": "2024-01-20T14:45:00Z",
  "recent_tickets": ["SRE-5344", "SRE-5345", "PROJ-1234"]
}
```

### Configuration (`~/.jit/config.yml`)
```yaml
jira:
  url: "https://company.atlassian.net"
  username: "user@company.com"
  token: "${JIRA_API_TOKEN}"
  project: "SRE"
  epic_link_field: "customfield_10014"
  
ai:
  provider: "openai"
  api_key: "${OPENAI_API_KEY}"
  model: "gpt-4"
  max_tokens: 1000
  
app:
  data_dir: "~/.local/share/jit"
  default_editor: "vim"
  review_before_create: true
```

## Core Components

### 1. Command Layer
- **Cobra CLI Framework**: Handles command parsing and routing
- **Command Handlers**: Individual command implementations
- **Middleware**: Common functionality like authentication, logging
- **Validation**: Input validation and error handling

### 2. Configuration Management
- **Config Loading**: YAML parsing with environment variable support
- **Validation**: Comprehensive configuration validation
- **Defaults**: Sensible default values for all settings
- **Migration**: Configuration version management

### 3. Jira Integration
- **REST API Client**: Jira Cloud API integration
- **Authentication**: Token-based authentication with refresh
- **Rate Limiting**: Respect Jira API rate limits
- **Error Handling**: Graceful handling of API errors

### 4. AI Integration
- **Provider Interface**: Abstract AI provider interface
- **Multiple Providers**: OpenAI, Anthropic, local models
- **Template Engine**: Dynamic prompt generation
- **Context Injection**: Project-specific context in prompts

### 5. Storage Layer
- **JSON Storage**: Human-readable ticket storage
- **Atomic Operations**: Safe concurrent file operations
- **Indexing**: Fast search and retrieval
- **Backup**: Automatic backup of critical data

### 6. Search & Discovery
- **Fuzzy Search**: Intelligent ticket finding
- **Full-text Search**: Search within ticket content
- **Filtering**: Multi-criteria filtering
- **Caching**: Performance optimization

## Data Flow

### 1. Ticket Tracking Flow
```
User: jit track SRE-5344
    ↓
1. Validate ticket exists in Jira
2. Fetch ticket and all children
3. Store in local JSON format
4. Set context to tracked ticket
5. Display confirmation
```

### 2. Ticket Creation Flow
```
User: jit task
    ↓
1. Validate current context (must be in epic)
2. Open editor for ticket description
3. Process content through AI (if enabled)
4. Show review screen (if enabled)
5. Create ticket in Jira
6. Store locally and update context
```

### 3. Context Switching Flow
```
User: jit focus "5433"
    ↓
1. Perform fuzzy search across local tickets
2. Display matches if multiple found
3. Update context.json with selection
4. Confirm context switch
```

## Error Handling

### Error Categories
1. **Configuration Errors**: Invalid config, missing credentials
2. **Network Errors**: Jira API failures, connectivity issues
3. **Data Errors**: Corrupt files, invalid ticket data
4. **User Errors**: Invalid commands, missing context
5. **System Errors**: File permissions, disk space

### Error Recovery
- **Graceful Degradation**: Offline mode when Jira unavailable
- **Auto-retry**: Exponential backoff for transient failures
- **User Guidance**: Clear error messages with suggested actions
- **Logging**: Comprehensive error logging for debugging

## Performance Considerations

### Optimization Strategies
1. **Lazy Loading**: Load ticket data on demand
2. **Caching**: Intelligent caching of frequently accessed data
3. **Batch Operations**: Minimize API calls through batching
4. **Concurrent Processing**: Parallel operations where possible
5. **Memory Management**: Efficient memory usage for large datasets

### Scalability
- **Large Projects**: Handle thousands of tickets efficiently
- **Search Performance**: Optimized search algorithms
- **Network Efficiency**: Minimize bandwidth usage
- **Storage Optimization**: Efficient storage formats

## Security Architecture

### Authentication
- **Token Storage**: Secure credential storage
- **Token Refresh**: Automatic token renewal
- **Multi-Factor**: Support for Jira MFA
- **Audit Trail**: Track all operations

### Data Protection
- **Local Encryption**: Encrypt sensitive local data
- **Secure Transmission**: HTTPS for all API calls
- **Permission Validation**: Respect Jira permissions
- **Data Sanitization**: Clean user input

## Testing Strategy

### Test Categories
1. **Unit Tests**: Individual component testing
2. **Integration Tests**: API integration testing
3. **End-to-End Tests**: Full workflow testing
4. **Performance Tests**: Load and stress testing
5. **Security Tests**: Vulnerability scanning

### Test Structure
```
tests/
├── unit/                              # Unit tests
│   ├── commands/
│   ├── storage/
│   └── jira/
├── integration/                       # Integration tests
│   ├── jira_integration_test.go
│   └── ai_integration_test.go
├── e2e/                               # End-to-end tests
│   ├── workflow_test.go
│   └── cli_test.go
└── fixtures/                          # Test data
    ├── tickets/
    └── configs/
```

## Build and Deployment

### Build Process
1. **Go Build**: Native Go compilation
2. **Cross-compilation**: Multi-platform builds
3. **Static Analysis**: Code quality checks
4. **Security Scanning**: Vulnerability detection
5. **Packaging**: Distribution packages

### Installation Methods
1. **Go Install**: Direct from source
2. **Binary Releases**: Pre-built binaries
3. **Package Managers**: Homebrew, apt, yum
4. **Docker**: Containerized deployment
5. **CI/CD Integration**: Automated deployment

## Monitoring and Observability

### Logging
- **Structured Logging**: JSON-formatted logs
- **Log Levels**: Configurable verbosity
- **Log Rotation**: Automatic log management
- **Error Tracking**: Comprehensive error capture

### Metrics
- **Usage Metrics**: Command usage statistics
- **Performance Metrics**: Response times, throughput
- **Error Metrics**: Error rates and patterns
- **User Metrics**: User behavior analytics

### Health Checks
- **Configuration Validation**: Startup health checks
- **API Connectivity**: Jira connection monitoring
- **Storage Health**: File system monitoring
- **Performance Monitoring**: Resource usage tracking

## Future Architecture Considerations

### Planned Enhancements
1. **Plugin System**: Extensible architecture
2. **Real-time Sync**: WebSocket-based updates
3. **Collaboration**: Multi-user support
4. **Mobile API**: REST API for mobile apps
5. **Machine Learning**: Predictive features

### Scalability Roadmap
1. **Database Backend**: Optional database storage
2. **Distributed Caching**: Redis integration
3. **Microservices**: Service decomposition
4. **Cloud Deployment**: SaaS offering
5. **Enterprise Features**: Advanced security and compliance 