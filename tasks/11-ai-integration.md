# Task 11: AI Integration Layer

**Status:** [ready]

## Overview
Implement the AI provider interface and OpenAI implementation to enrich ticket descriptions and comments.

## Objectives
- Define `ai.Provider` interface with `Enrich(content, context) (string, error)`
- Implement OpenAI provider using completion API
- Load prompt templates from config templates dir
- Support `{{expression}}` evaluation within templates
- Configurable model, max_tokens, temperature

## Deliverables
- [ready] `internal/ai/provider.go` (interface + factory)
- [ready] `internal/ai/openai.go` implementation
- [ready] `internal/ai/templates.go` for template processing
- [ready] Unit tests using fake OpenAI server

## Dependencies
- Task 02 Config (AI settings)
- net/http

## Implementation Notes
- Respect OpenAI token from env/config
- Provide dry-run mode for tests (`OPENAI_MOCK=1`)
- Cache templates in memory

## Acceptance Criteria
- `go test ./internal/ai/...` passes
- `jit epic --no-create` runs enrichment locally without hitting real API when mock enabled

## Next Tasks
- (Back to Task 00 for implementation cycle) 