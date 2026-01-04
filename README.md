# gelete

[![CI](https://github.com/Kdaito/gelete/actions/workflows/ci.yml/badge.svg)](https://github.com/Kdaito/gelete/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/Kdaito/gelete)](https://goreportcard.com/report/github.com/Kdaito/gelete)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**gelete** is an interactive command-line tool for deleting git branches. It provides a terminal UI (TUI) that allows you to visually select multiple local branches for deletion, with built-in safety features for unmerged branches and git worktree awareness.

## Features

- **Interactive TUI**: Select branches using keyboard navigation with a clean, intuitive interface
- **Multi-select**: Delete multiple branches in a single session
- **Safety First**:
  - Automatic detection of unmerged branches with force delete option
  - Git worktree awareness with automatic worktree removal
  - Confirmation prompts before any destructive operations
- **Smart Filtering**: Current branch is automatically excluded from the deletion list
- **Cross-Platform**: Works on Linux, macOS, and Windows (amd64 and arm64)

## Installation

### Go Install

```bash
go install github.com/Kdaito/gelete@latest
```

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/Kdaito/gelete/releases).

**Available platforms:**
- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

### Homebrew (macOS/Linux)

```bash
brew tap Kdaito/tap
brew install gelete
```

#### macOS Additional Setup

After installation on macOS, you may need to remove the quarantine attribute:

```bash
xattr -d com.apple.quarantine $(which gelete)
```

> **Note**: The path can vary depending on your environment. Use `which gelete` to find the exact location.

## Usage

Simply run `gelete` in any git repository:

```bash
cd /path/to/your/git/repo
gelete
```

### Keyboard Controls

**Branch Selection:**
- `↑/k` - Move cursor up
- `↓/j` - Move cursor down
- `Space/Enter` - Toggle branch selection
- `d` - Delete selected branches
- `q/Ctrl+C` - Quit without deleting

**Confirmation:**
- `y` - Confirm deletion
- `n` - Cancel

**Force Delete (for unmerged branches):**
- `y` - Force delete unmerged branches
- `n` - Skip unmerged branches

## Examples

### Basic Usage

```bash
$ gelete
gelete - Interactive Branch Deletion

  > [✓] feature/old-feature
    [ ] feature/experimental
    [✓] bugfix/issue-123

↑/k: up • ↓/j: down • space/enter: toggle • d: delete selected • q: quit
```

### Handling Unmerged Branches

When you attempt to delete a branch with unmerged changes, gelete will:
1. Detect the unmerged branch
2. Show a clear warning message
3. Offer the option to force delete with `-D` flag

```bash
⚠ Warning: Unmerged Branches Detected

The following branches have unmerged changes:

  • feature/experimental
    error: The branch 'feature/experimental' is not fully merged.

Force delete will permanently remove 1 unmerged branch(es).
This action cannot be undone!

y: force delete • n: cancel and skip these branches
```

### Git Worktree Awareness

gelete automatically detects and handles git worktrees:

```bash
gelete - Interactive Branch Deletion

  > [✓] feature/old-feature [worktree]
    [ ] feature/new-feature
    [✓] bugfix/issue-123

↑/k: up • ↓/j: down • space/enter: toggle • d: delete selected • q: quit
```

When deleting a branch with an active worktree, gelete will:
1. Automatically remove the worktree directory
2. Handle locked worktrees with force removal if needed
3. Then delete the branch

## Requirements

- Git 2.0 or higher
- Go 1.21 or higher (for building from source)

## Development

### Building from Source

```bash
git clone https://github.com/Kdaito/gelete.git
cd gelete
go build -o gelete .
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run with coverage
go test ./... -cover

# Run specific test suites
go test ./tests/unit/...
go test ./tests/integration/...
go test ./tests/contract/...
```

### Project Structure

```
gelete/
├── cmd/              # CLI commands (Cobra)
├── internal/
│   ├── git/          # Git operations (branch, worktree, repository)
│   └── ui/           # TUI components (Bubbletea)
├── tests/
│   ├── unit/         # Unit tests
│   ├── integration/  # Integration tests
│   └── contract/     # Contract tests
└── main.go           # Entry point
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [Bubbletea](https://github.com/charmbracelet/bubbletea) - The TUI framework
- CLI framework by [Cobra](https://github.com/spf13/cobra)
- Terminal styling by [Lipgloss](https://github.com/charmbracelet/lipgloss)

## Author

[Kdaito](https://github.com/Kdaito)

---

**Note**: gelete only deletes **local** branches. Remote branches are not affected.
