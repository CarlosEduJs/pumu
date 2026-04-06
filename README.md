# Pumu 🧹

**Pumu** (Package Manager Utility) is a fast, concurrent disk space management tool for developers. It helps you find and remove heavy dependency folders across multiple projects and package managers, then optionally reinstall them.

Stop manually hunting down `node_modules` folders eating up gigabytes of disk space. Pumu does it for you, intelligently.

## Features

- 🔍 **Multi-language support** - Works with npm, pnpm, yarn, bun, deno, cargo, go, and pip
- ⚡ **Blazingly fast** - Concurrent scanning and deletion using goroutines with semaphore-based throttling
- 📊 **Visual feedback** - Color-coded output based on folder size with human-readable formatting
- 🎯 **Smart detection** - Automatically identifies package managers via lockfiles and manifests
- 🔄 **Reinstall option** - Optionally reinstall dependencies after cleanup
- 🛡️ **Safe defaults** - Ignores system folders and version control directories
- 📋 **Dry-run mode** - Preview what will be deleted before taking action
- ✅ **Interactive selection** - Choose which folders to delete/reinstall via TUI multi-select
- 🔧 **Repair mode** - Detect and fix corrupted dependencies automatically
- 🌿 **Smart prune** - Score-based intelligent cleanup that only removes what's safe

## Supported Folders

Pumu can detect and clean these dependency/build folders:

| Folder          | Package Manager(s)              | Typical Size |
|-----------------|---------------------------------|--------------|
| `node_modules`  | npm, yarn, pnpm, bun, deno     | 50-500+ MB   |
| `target`        | cargo (Rust)                    | 100-2000+ MB |
| `.venv`         | pip (Python virtual env)        | 50-300+ MB   |
| `.next`         | Next.js                         | 100-500+ MB  |
| `.svelte-kit`   | SvelteKit                       | 50-200+ MB   |
| `dist`          | Various build tools             | 10-100+ MB   |
| `build`         | Various build tools             | 10-100+ MB   |

## Installation

### Homebrew (Recommended)

```bash
brew install carlosedujs/pumu/pumu
```

### From Source

```bash
git clone https://github.com/carlosedujs/pumu.git
cd pumu
go build -o pumu
sudo mv pumu /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/carlosedujs/pumu@latest
```

## Usage

```bash
pumu [command] --help
```

```
pumu scans your filesystem for heavy dependency folders
(node_modules, target, .venv, etc.) and lets you sweep, list,
repair or prune them with ease.

Usage:
  pumu [flags]
  pumu [command]

Available Commands:
  list        List heavy dependency folders (dry-run)
  sweep       Sweep (delete) heavy dependency folders
  repair      Repair dependency folders
  prune       Prune dependency folders by staleness score
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command

Flags:
  -h, --help          help for pumu
  -p, --path string   Root path to scan (default ".")
  -v, --version       version for pumu
```

### 1. Default Mode (Refresh Current Directory)

Detects the package manager in the current directory, removes its dependency folder, and reinstalls:

```bash
pumu
```

**Example Output:**

```
Running refresh in current directory...
🔍 Detected package manager: npm
🗑️  Removing node_modules...
✅ Removed in 1.23s
📦 Running npm install...
[npm install output...]
🎉 Refresh complete!
```

### 2. Version Command

Display the current version of Pumu:

```bash
pumu --version
# or
pumu -v
```

**Example Output:**

```
pumu version v1.2.1
```

### 3. List Mode (Dry Run)

Recursively scans for heavy folders without deleting them. You can specify a path using the global `-p` or `--path` flag:

```bash
pumu list                      # scan current directory
pumu list --path ~/projects    # scan a specific directory
```

**Example Output:**

<pre>
🔎 Listing heavy dependency folders in '.'...
⏱️  Found 3 folders. Calculating sizes concurrently...

<span style="text-decoration: underline;">Folder Path                                                                       | Size</span>
/home/user/projects/webapp/node_modules                                           | <span style="color: #ff0000;">1.23 GB 🚨</span>
/home/user/projects/rust-app/target                                               | <span style="color: #ffff00;">487.50 MB ⚠️</span>
/home/user/projects/api/.venv                                                     | <span style="color: #00ff00;">89.32 MB</span>
----------------------------------------------------------------------------------------------------
<span style="color: #00ff00;">📋 List complete! Found 3 heavy folders.</span>
<span style="color: #00ffff;">💾 Total space that can be freed: 1.79 GB</span>
</pre>

### 4. Sweep Mode

Recursively scans and **deletes** heavy folders. Shows an interactive multi-select by default so you can choose which folders to delete:

```bash
pumu sweep
pumu sweep -p ~/dev
```

**Example Output:**

<pre>
🔎 Scanning for heavy dependency folders in '.'...
⏱️  Found 3 folders. Calculating sizes concurrently...

🗑️  Select folders to delete:
▸ [✓] /home/user/projects/webapp/node_modules       1.23 GB
  [✓] /home/user/projects/rust-app/target            487.50 MB
  [ ] /home/user/projects/api/.venv                  89.32 MB

  2/3 selected
  press ? for help
</pre>

#### Interactive Selection Shortcuts

| Key | Action |
|-----|--------|
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `g` / `G` | Go to first / last |
| `space` | Toggle item |
| `a` | Select all |
| `n` | Deselect all |
| `i` | Invert selection |
| `enter` | Confirm |
| `q` / `esc` | Cancel |
| `?` | Toggle help |

#### Sweep with Reinstall

Delete folders and choose which projects to reinstall:

```bash
pumu sweep --reinstall
```

#### Sweep without Interactive Selection

Skip the multi-select and delete all found folders directly (old behavior):

```bash
pumu sweep --no-select
pumu sweep --reinstall --no-select
```

### 5. Repair Mode

Scans for projects with corrupted or broken dependencies and automatically fixes them by removing and reinstalling:

```bash
pumu repair
```

**Example Output:**

<pre>
🔧 Scanning for projects with broken dependencies in '.'...
⏱️  Found 3 projects. Checking health...

📁 ./webapp (npm)
   <span style="color: #ff0000;">❌ Missing: react, react-dom</span>
   🗑️  Removing node_modules...
   📦 Reinstalling...
   <span style="color: #00ff00;">✅ Repaired!</span>

📁 ./api (pnpm)
   <span style="color: #00ff00;">✅ Healthy, skipping.</span>

📁 ./rust-cli (cargo)
   <span style="color: #ff0000;">❌ Compilation errors detected</span>
   🗑️  Removing target...
   📦 Rebuilding...
   <span style="color: #00ff00;">✅ Repaired!</span>

-----
<span style="color: #00ff00;">🔧 Repair complete! Fixed 2/3 projects.</span>
</pre>

#### Verbose Mode

Show details for all projects, including healthy ones:

```bash
pumu repair --verbose
```

#### Health Checks per Ecosystem

| Ecosystem | Check Method |
|-----------|-------------|
| npm / pnpm | `npm ls --json` / `pnpm ls --json` |
| yarn | `yarn check --verify-tree` |
| cargo | `cargo check` |
| go | `go mod verify` |
| pip | `pip check` |

### 6. Prune Mode

Smart cleanup — analyzes folders with a safety score (0-100) and only deletes what's truly safe. Less destructive than sweep:

```bash
pumu prune
```

**Example Output:**

<pre>
🌿 Pruning safely deletable folders in '.'...
⏱️  Found 5 folders. Analyzing...

<span style="text-decoration: underline;">Folder Path                                             | Size       | Score | Reason</span>
./old-project/node_modules                              | 456.78 MB  |  <span style="color: #ff0000;">  95</span> | 🔴 No lockfile (orphan)
./webapp/.next                                          | 234.56 MB  |  <span style="color: #ff0000;">  90</span> | 🟢 Build cache (re-generable)
./api/node_modules                                      | 189.00 MB  |  <span style="color: #ffff00;">  60</span> | 🟡 Lockfile stale (45 days)
<span style="color: #808080;">./active-project/node_modules                            | 567.89 MB  |    20 | ⚪ Active project (skipped)</span>
<span style="color: #808080;">./wip/target                                             | 890.12 MB  |    15 | ⚪ Uncommitted changes (skipped)</span>

----------------------------------------------
<span style="color: #00ff00;">🌿 Prune complete! Removed 3 folders (score ≥ 50).</span>
<span style="color: #00ffff;">💾 Space freed: 880.34 MB (of 2.34 GB total found)</span>
</pre>

#### Prune Scoring Heuristics

| Score | Reason |
|-------|--------|
| 90-95 | Orphan folder (no lockfile) or build cache |
| 60-80 | Stale lockfile (30-90+ days without changes) |
| 45 | Dependency folder with lockfile (moderate) |
| 15-20 | Active project or uncommitted lockfile changes |

#### Prune Options

```bash
pumu prune --dry-run          # Only analyze, don't delete
pumu prune --threshold 80     # Only prune folders with score ≥ 80
```

## How It Works

### Package Manager Detection

Pumu automatically detects the package manager by checking for specific files in priority order:

1. **Bun** - `bun.lockb` or `bun.lock`
2. **pnpm** - `pnpm-lock.yaml`
3. **Yarn** - `yarn.lock`
4. **npm** - `package-lock.json`
5. **Deno** - `deno.json` or `deno.jsonc`
6. **Cargo** - `Cargo.toml`
7. **Go** - `go.mod`
8. **Pip** - `requirements.txt` or `pyproject.toml`

### Performance Optimizations

- **Concurrent size calculation** - Uses goroutines with a bounded semaphore (max 20 concurrent operations) to calculate folder sizes in parallel
- **Concurrent deletion** - Deletes multiple folders simultaneously while respecting system limits
- **Smart path skipping** - Automatically skips `.git`, `.cache`, IDE folders, and other non-project directories
- **Atomic operations** - Thread-safe accumulation of deleted space using atomic operations

### Ignored Paths

To avoid scanning irrelevant directories, Pumu skips:

- `.Trash`
- `.cache`, `.npm`, `.yarn`, `.cargo`, `.rustup`
- `Library`, `AppData`, `Local`, `Roaming`
- `.vscode`, `.idea`
- `.git` (version control)

## Project Structure

```
pumu/
├── main.go                      # CLI entry point
├── cmd/                         # CLI commands (Cobra)
│   ├── root.go                  # Root command and global flags
│   ├── sweep.go                 # Sweep command definition
│   ├── list.go                  # List command definition
│   ├── repair.go                # Repair command definition
│   └── prune.go                 # Prune command definition
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # Core scanning and deletion logic
│   │   ├── scanner_test.go      # Scanner tests
│   │   ├── repair.go            # Repair command logic
│   │   └── prune.go             # Prune command logic
│   ├── pkg/
│   │   ├── detector.go          # Package manager detection
│   │   ├── detector_test.go     # Detector tests
│   │   ├── installer.go         # Dependency installation
│   │   ├── cleaner.go           # Directory removal utilities
│   │   ├── checker.go           # Health checks per package manager
│   │   └── analyzer.go          # Prune scoring heuristics
│   └── ui/
│       └── multiselect.go       # Interactive TUI multi-select component
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Requirements

- **Go 1.24.0+** for building from source
- **Package managers** must be installed if using `--reinstall` or `repair`:
  - Node.js: npm, yarn, pnpm, or bun
  - Rust: cargo
  - Python: pip
  - Go: go
  - Deno: deno

## Color Coding

Pumu uses visual indicators to help you prioritize cleanup:

- <span style="color: #ff0000;">**Red (🚨)**</span> - Folders larger than 1 GB (critical space users)
- <span style="color: #ffff00;">**Yellow (⚠️)**</span> - Folders between 100 MB and 1 GB (moderate space users)
- <span style="color: #00ff00;">**Green**</span> - Folders smaller than 100 MB (minor space users)

## Safety Features

- **Dry-run by default** with `list` command - preview before deletion
- **Interactive selection** - choose exactly what to delete via TUI multi-select
- **Smart prune** - score-based cleanup that skips active projects
- **Repair before delete** - fix corrupted deps instead of blindly removing
- **Explicit sweep** - requires `sweep` command to actually delete
- **Smart folder detection** - only removes known dependency folders
- **Concurrent safe** - uses mutexes and atomic operations to prevent race conditions
- **Error handling** - continues processing even if individual operations fail

## Use Cases

### Clean Up Development Machine

Free up disk space across all your projects:

```bash
cd ~/projects
pumu sweep
```

### Preview Disk Usage

See what's taking up space without deleting:

```bash
cd ~/workspace
pumu list
```

### Bulk Cleanup with Reinstall

Clean multiple projects and reinstall everything:

```bash
cd ~/dev
pumu sweep --reinstall
```

### Fix Broken Projects

Repair corrupted dependencies across all projects:

```bash
cd ~/projects
pumu repair
```

### Safe Cleanup (Only Remove What's Safe)

Intelligently prune only stale and orphan folders:

```bash
cd ~/projects
pumu prune --dry-run        # preview first
pumu prune                  # prune score ≥ 50
pumu prune --threshold 80   # conservative mode
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Carlos**

---

Built with ❤️ and Go
