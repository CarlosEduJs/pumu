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

Pumu has three modes of operation:

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

### 2. List Mode (Dry Run)

Recursively scans for heavy folders without deleting them:

```bash
pumu list
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

### 3. Sweep Mode

Recursively scans and **deletes** all heavy folders:

```bash
pumu sweep
```

**Example Output:**

<pre>
🔎 Scanning for heavy dependency folders in '.'...
⏱️  Found 3 folders. Calculating sizes concurrently...

<span style="text-decoration: underline;">Folder Path                                                                       | Size</span>
/home/user/projects/webapp/node_modules                                           | <span style="color: #ff0000;">1.23 GB 🚨</span>
/home/user/projects/rust-app/target                                               | <span style="color: #ffff00;">487.50 MB ⚠️</span>
/home/user/projects/api/.venv                                                     | <span style="color: #00ff00;">89.32 MB</span>

<span style="color: #ffff00;">🗑️  Deleting folders concurrently...</span>
----------------------------------------------------------------------------------------------------
<span style="color: #00ff00;">🧹 Sweep complete! Processed 3 heavy folders.</span>
<span style="color: #00ffff;">💾 Total space actually freed: 1.79 GB</span>
</pre>

### 4. Sweep with Reinstall

Delete folders and automatically reinstall dependencies for each project:

```bash
pumu sweep --reinstall
```

**Example Output:**

<pre>
🔎 Scanning for heavy dependency folders in '.'...
⏱️  Found 2 folders. Calculating sizes concurrently...

<span style="text-decoration: underline;">Folder Path                                                                       | Size</span>
/home/user/projects/webapp/node_modules                                           | <span style="color: #ff0000;">1.23 GB 🚨</span>
/home/user/projects/api/node_modules                                              | <span style="color: #ffff00;">156.78 MB ⚠️</span>

<span style="color: #ffff00;">🗑️  Deleting folders concurrently...</span>
----------------------------------------------------------------------------------------------------
<span style="color: #00ff00;">🧹 Sweep complete! Processed 2 heavy folders.</span>
<span style="color: #00ffff;">💾 Total space actually freed: 1.38 GB</span>

<span style="color: #ffff00;">⚙️  Reinstalling dependencies sequentially...</span>
📦 Reinstalling for /home/user/projects/webapp (npm)...
<span style="color: #00ff00;">✅ Reinstalled /home/user/projects/webapp</span>
📦 Reinstalling for /home/user/projects/api (pnpm)...
<span style="color: #00ff00;">✅ Reinstalled /home/user/projects/api</span>
<span style="color: #00ff00;">🎉 All target reinstallations complete!</span>
</pre>

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
├── main.go                      # CLI entry point and command routing
├── internal/
│   ├── scanner/
│   │   ├── scanner.go           # Core scanning and deletion logic
│   │   └── scanner_test.go      # Scanner tests
│   └── pkg/
│       ├── detector.go          # Package manager detection
│       ├── detector_test.go     # Detector tests
│       ├── installer.go         # Dependency installation
│       └── cleaner.go           # Directory removal utilities
├── go.mod
├── go.sum
├── LICENSE
└── README.md
```

## Requirements

- **Go 1.22.2+** for building from source
- **Package managers** must be installed if using `--reinstall` flag:
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

### Fresh Start on a Project

Remove and reinstall dependencies for a clean build:

```bash
cd my-project
pumu
```

### Bulk Cleanup with Reinstall

Clean multiple projects and reinstall everything:

```bash
cd ~/dev
pumu sweep --reinstall
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Author

**Carlos**

---

Built with ❤️ and Go
