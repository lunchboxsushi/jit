# Task 12: Test Harness & Continuous Integration

**Status:** [completed]

## Overview
Establish a robust automated testing workflow and continuous integration (CI) pipeline to run linting, unit tests, and build checks on every commit and pull request.

## Objectives
- Configure GitHub Actions workflow (`.github/workflows/ci.yml`)
  - Run `go vet`, `go test ./...`, and `go build ./cmd/jit`
  - Cache Go modules for faster builds
- Integrate static analysis tools:
  - `golangci-lint` with default linters (govet, staticcheck, gocyclo)
- Generate coverage reports and upload as artifact
- Fail build on linter or test errors

## Deliverables
- [completed] `.github/workflows/ci.yml` GitHub Actions config
- [completed] `Makefile` or `go run` scripts for local lint/test
- [completed] `golangci.yml` configuration

## Dependencies
- All prior code tasks (compilable project)

## Implementation Notes
- Use matrix build for latest two Go versions (e.g., 1.21, 1.22)
- Enable PR status checks
- Optional: Badge in README for build status

## Acceptance Criteria
- Pushing to any branch triggers CI and passes
- PR shows green check on successful build

## Next Tasks
- 13-cli-completion-docs.md 