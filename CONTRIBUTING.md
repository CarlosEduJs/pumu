# Contributing to Pumu 🧹

First off, thank you for considering contributing to Pumu! It's people like you that make Pumu such a great tool.

## Code of Conduct

By participating in this project, you agree to abide by our terms. Please be respectful and professional in all interactions.

## How Can I Help?

### Reporting Bugs
If you find a bug, please create an issue on GitHub. Include:
- A clear, descriptive title.
- Steps to reproduce the bug.
- Your OS and Go version.
- Any relevant logs or screenshots.

### Suggesting Enhancements
I love new ideas! Please open an issue and describe the feature you'd like to see, why it's useful, and how it should work.

### Pull Requests
1.  **Fork the repository** and create your branch from `main`.
2.  **Install dependencies**: `go mod download`.
3.  **Make your changes**. Ensure your code is formatted with `go fmt ./...`.
4.  **Write tests** for any new functionality.
5.  **Run existing tests**: `go test ./...` to ensure no regressions.
6.  **Check for lint errors**: We use `golangci-lint`. Make sure it passes.
7.  **Submit a Pull Request** with a clear description of what you did.

## Development Setup

### Requirements
- **Go 1.24.0+** (required for Bubble Tea TUI components)
- A terminal with ANSI color support

### Building locally
```bash
go build -o pumu .
```

### Running tests
```bash
go test -v ./...
```

## Project Structure

- `main.go`: Entry point and CLI command definitions.
- `internal/scanner/`: Core logic for scanning, deleting, and UI components.
- `internal/pkg/`: Shared utilities, package manager detection, and health checkers.
- `internal/ui/`: TUI components and styles.

## Coding Style

- Follow standard Go idioms and conventions.
- Use meaningful variable and function names.
- Keep functions small and focused (respect the cyclomatic complexity limits).
- Add comments for complex logic.

Thank you for your contribution!
