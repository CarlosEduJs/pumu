# AGENTS.md

## Project Overview

Pumu is a fast, concurrent disk space management CLI tool written in Go. It helps developers find and remove heavy dependency folders (node_modules, target, .venv, etc.) across multiple projects and package managers, with optional reinstallation.

Key characteristics:
- Go 1.24.2+ project using Cobra for CLI and Bubble Tea for TUI
- Concurrent operations using goroutines with semaphore-based throttling
- Multi-package manager support (npm, pnpm, yarn, bun, deno, cargo, go, pip)
- Interactive TUI components for folder selection

## Setup Commands

```bash
# Install dependencies
go mod download

# Build the binary
go build -o pumu

# Run locally
go run main.go [command]

# Install globally
go install

# Run tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/scanner/...
go test ./internal/pkg/...
```

## Development Workflow

### Running Commands During Development

```bash
# Test list command
go run main.go list

# Test sweep with dry-run
go run main.go sweep --no-select

# Test in examples directory
go run main.go list -p ./examples

# Test repair mode
go run main.go repair -p ./examples

# Test prune mode
go run main.go prune --dry-run
```

### Testing with Example Projects

The `examples/` directory contains sample projects for testing:
- `examples/node-npm/` - npm project
- `examples/node-pnpm/` - pnpm project
- `examples/node-bun/` - bun project
- `examples/nextjs-project/` - Next.js project
- `examples/go-project/` - Go project
- `examples/rust-project/` - Rust project
- `examples/python-project/` - Python project
- `examples/deno-project/` - Deno project

Use `examples/setup.sh` to install dependencies and `examples/reinstall.sh` to test reinstallation.

## Code Style & Conventions

### Go Standards
- Follow standard Go formatting (gofmt, goimports)
- Use golangci-lint for linting (config in `.golangci.yml`)
- Minimum Go version: 1.24.2
- Use meaningful variable names, avoid single-letter names except for common idioms (i, j, err, wg, mu)

### Project-Specific Patterns
- Use Cobra for CLI commands (see `cmd/` directory)
- Use Bubble Tea/Huh for interactive TUI components
- Concurrent operations should use semaphores to limit parallelism (typically 20 concurrent ops)
- Use atomic operations for thread-safe counters
- Error handling: return errors, don't panic; continue processing on individual failures

### File Organization
```
cmd/           - CLI command definitions (Cobra)
internal/
  scanner/     - Core scanning, deletion, repair, prune logic
  pkg/         - Package manager detection, installation, health checks
  ui/          - TUI components (multi-select, progress bars, styling)
main.go        - Entry point
```

### Naming Conventions
- Package managers: use `PackageManager` type from `internal/pkg/detector.go`
- Target folders: use `TargetFolder` struct with `Path` and `Size` fields
- Commands: verb-based (sweep, list, repair, prune)
- Functions: camelCase, exported functions start with uppercase

## Testing Guidelines

### Running Tests
```bash
# Run all tests
go test ./...

# Run with verbose output
go test -v ./...

# Run specific test
go test -run TestDetectManager ./internal/pkg/...

# Run with coverage
go test -cover ./...
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Files
- `internal/pkg/detector_test.go` - Package manager detection tests
- `internal/scanner/scanner_test.go` - Scanner logic tests

### Writing Tests
- Use table-driven tests for multiple scenarios
- Test files should end with `_test.go`
- Use `t.Run()` for subtests
- Mock filesystem operations when possible
- Test edge cases: empty directories, permission errors, unknown package managers

## Linting & Code Quality

```bash
# Run golangci-lint
golangci-lint run

# Run with auto-fix
golangci-lint run --fix

# Run specific linters
golangci-lint run --enable-only=errcheck,gosec
```

Enabled linters (see `.golangci.yml`):
- errcheck - check for unchecked errors
- revive - general Go linting
- gosec - security checks
- unconvert - unnecessary type conversions
- goconst - repeated strings that could be constants
- gocyclo - cyclomatic complexity (max 15)
- dupl - code duplication (threshold 100)
- misspell - spelling errors

## Building & Releasing

### Local Build
```bash
# Build for current platform
go build -o pumu

# Build with version info
go build -ldflags "-X pumu/cmd.version=v1.2.1" -o pumu
```

### Cross-Platform Build
```bash
# Build for Linux
GOOS=linux GOARCH=amd64 go build -o pumu-linux-amd64

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o pumu-darwin-amd64
GOOS=darwin GOARCH=arm64 go build -o pumu-darwin-arm64

# Build for Windows
GOOS=windows GOARCH=amd64 go build -o pumu-windows-amd64.exe
```

### Release Process
- Uses GoReleaser (config in `.goreleaser.yaml`)
- Version is defined in `cmd/root.go` as `const version`
- Homebrew tap: `carlosedujs/pumu/pumu`

## Key Implementation Details

### Package Manager Detection
Priority order (see `internal/pkg/detector.go`):
1. Bun - `bun.lockb` or `bun.lock`
2. pnpm - `pnpm-lock.yaml`
3. yarn - `yarn.lock`
4. npm - `package-lock.json`
5. Deno - `deno.json` or `deno.jsonc`
6. Cargo - `Cargo.toml`
7. Go - `go.mod`
8. Pip - `requirements.txt` or `pyproject.toml`

### Ignored Paths
Scanner skips these directories (see `scanner.go`):
- `.Trash`, `.cache`, `.npm`, `.yarn`, `.cargo`, `.rustup`
- `Library`, `AppData`, `Local`, `Roaming`
- `.vscode`, `.idea`
- `.git`

### Deletable Targets
- `node_modules` (npm, yarn, pnpm, bun, deno)
- `target` (cargo)
- `.venv` (pip)
- `.next` (Next.js)
- `.svelte-kit` (SvelteKit)
- `dist`, `build` (various build tools)

### Concurrency Patterns
- Use semaphores to limit concurrent operations: `sem := make(chan struct{}, 20)`
- Use `sync.WaitGroup` for goroutine coordination
- Use `sync.Mutex` for protecting shared data
- Use `atomic` operations for counters

## Common Tasks

### Adding a New Package Manager
1. Add constant to `PackageManager` type in `internal/pkg/detector.go`
2. Add detection logic to `DetectManager()` function
3. Add installation command to `InstallDependencies()` in `internal/pkg/installer.go`
4. Add health check to `CheckHealth()` in `internal/pkg/checker.go`
5. Update `getTargetFolder()` in `internal/scanner/scanner.go`
6. Add tests to `internal/pkg/detector_test.go`
7. Update README.md documentation

### Adding a New Command
1. Create new file in `cmd/` directory (e.g., `cmd/newcommand.go`)
2. Define Cobra command with flags
3. Add command to root in `init()` function
4. Implement logic in `internal/scanner/` or `internal/pkg/`
5. Add tests
6. Update README.md with usage examples

### Modifying TUI Components
- Multi-select component: `internal/ui/multiselect.go`
- Progress bars: `internal/ui/progress.go`
- Styling/colors: `internal/ui/style.go`
- Uses Bubble Tea framework and Huh library

## Dependencies

Main dependencies (see `go.mod`):
- `github.com/spf13/cobra` - CLI framework
- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/huh` - Form/input components
- `github.com/charmbracelet/lipgloss` - Terminal styling
- `github.com/charmbracelet/bubbles` - TUI components

## Security Considerations

- Never delete files outside of known dependency folders
- Always check if path exists before deletion
- Use `filepath.SkipDir` to avoid descending into system directories
- Validate user input for path flags
- Handle permission errors gracefully
- Don't follow symlinks during scanning

## Performance Tips

- Limit concurrent operations to avoid overwhelming the system (default: 20)
- Use `filepath.WalkDir` instead of `filepath.Walk` (more efficient)
- Calculate sizes concurrently with goroutines
- Skip `.git` directories early to avoid scanning large repos
- Use atomic operations instead of mutexes for simple counters

## Debugging

```bash
# Run with verbose output
go run main.go list -p ./examples

# Use Go's built-in profiler
go run main.go list -cpuprofile=cpu.prof
go tool pprof cpu.prof

# Check for race conditions
go run -race main.go sweep --no-select

# Print debug info
# Add fmt.Printf() statements in code for debugging
```

## Common Issues

### "Permission denied" errors
- Some folders may require elevated permissions
- Scanner continues processing other folders on error

### "Package manager not detected"
- Ensure lockfile or manifest exists in project directory
- Check `DetectManager()` logic in `internal/pkg/detector.go`

### Slow scanning on large directories
- Adjust semaphore limit in scanner.go
- Ensure ignored paths are properly configured

### TUI rendering issues
- Ensure terminal supports ANSI colors
- Check terminal width for proper formatting
