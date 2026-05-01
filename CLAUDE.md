# gelete Development Guidelines

## Overview

`gelete` is an interactive CLI tool for git branch deletion, built with Go and the Bubble Tea TUI framework.

## Tech Stack

- Go (version specified in `go.mod`)
- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Cobra](https://github.com/spf13/cobra) — CLI framework
- [golangci-lint](https://golangci-lint.run/) v2 — linter

## Project Structure

```
cmd/              CLI entry point (Cobra root command)
internal/
  git/            Git operations (branch, worktree, repository)
  ui/             Bubble Tea TUI (model, update, view, styles)
tests/
  unit/           Unit tests
  integration/    Integration tests
  contract/       Contract tests
specs/            Feature specifications and plans
.github/
  workflows/      CI (ci.yml) and release (release.yml)
```

## Commands

```bash
make ci       # Run all CI checks locally (fmt → vet → lint → test → build)
make test     # Run tests with coverage
make lint     # Run golangci-lint (auto-installs if missing)
make fmt      # Check gofmt formatting
make vet      # Run go vet
make build    # Build binary
make clean    # Remove binary and coverage.txt
```

## Code Style

- Follow standard Go conventions
- No comments unless the WHY is non-obvious

## Commit Messages

- Written in **English**
- Format: `<prefix>: <message>` — a single concise line
- Do not include author information or co-author lines
- Common prefixes: `feat`, `fix`, `chore`, `refactor`, `test`, `docs`

Example:
```
fix: resolve golangci-lint v2 config incompatibilities
```

## Pull Requests

- Title and body must be written in **English**
- Keep both title and body concise
- Follow the PR template (`.github/pull_request_template.md`):

```
## Summary

## Changes

## Test plan

- [ ]
```
