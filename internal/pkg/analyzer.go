package pkg

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// PruneResult holds the analysis result for a folder, deciding if it's safe to prune.
type PruneResult struct {
	Path         string
	Size         int64
	Score        int    // 0-100, higher = safer to delete
	Reason       string // Human readable explanation
	SafeToDelete bool   // Whether score meets the threshold
}

// AnalyzeFolder evaluates whether a dependency/build folder is safe to prune
// based on multiple heuristics: orphan status, build cache, lockfile staleness,
// and uncommitted changes.
func AnalyzeFolder(folderPath string, size int64) PruneResult {
	result := PruneResult{
		Path: folderPath,
		Size: size,
	}

	folderName := filepath.Base(folderPath)
	projectDir := filepath.Dir(folderPath)

	// Build output folders are always re-generable - heuristic 1
	if isBuildCache(folderName) {
		result.Score = 90
		result.Reason = "🟢 Build cache (re-generable)"
		return result
	}

	// Check if project has a lockfile at all - heuristic 2
	pm := DetectManager(projectDir)
	if pm == Unknown {
		result.Score = 95
		result.Reason = "🔴 No lockfile (orphan folder)"
		return result
	}

	// Check lockfile age (staleness) - heuristic 3
	lockfileAge := getLockfileAge(projectDir, pm)
	if lockfileAge > 0 {
		days := int(lockfileAge.Hours() / 24)

		if days > 90 {
			result.Score = 80
			result.Reason = "🟡 Lockfile very stale (" + formatDays(days) + ")"
			return result
		}
		if days > 30 {
			result.Score = 60
			result.Reason = "🟡 Lockfile stale (" + formatDays(days) + ")"
			return result
		}
	}

	// Check git status for uncommitted lockfile changes - heuristic 4
	if hasUncommittedLockfileChanges(projectDir) {
		result.Score = 15
		result.Reason = "⚪ Uncommitted lockfile changes (active work)"
		return result
	}

	// Recent lockfile = active project - heuristic 5
	if lockfileAge > 0 && lockfileAge.Hours()/24 < 7 {
		result.Score = 20
		result.Reason = "⚪ Active project (recently modified)"
		return result
	}

	// default: moderate score for dependency folders with lockfiles
	result.Score = 45
	result.Reason = "🟡 Dependency folder with lockfile"
	return result
}

// isBuildCache returns true for folders that are purely build output.
func isBuildCache(name string) bool {
	caches := []string{".next", ".svelte-kit", "dist", "build"}
	for _, c := range caches {
		if name == c {
			return true
		}
	}
	return false
}

// getLockfileAge returns the age of the lockfile for a project.
// Returns 0 if no lockfile found.
func getLockfileAge(dir string, pm PackageManager) time.Duration {
	lockfiles := getLockfiles(pm)

	for _, lf := range lockfiles {
		path := filepath.Join(dir, lf)
		info, err := os.Stat(path)
		if err == nil {
			return time.Since(info.ModTime())
		}
	}

	return 0
}

// getLockfiles returns the lockfile names for a given package manager.
func getLockfiles(pm PackageManager) []string {
	switch pm {
	case Npm:
		return []string{"package-lock.json"}
	case Pnpm:
		return []string{"pnpm-lock.yaml"}
	case Yarn:
		return []string{"yarn.lock"}
	case Bun:
		return []string{"bun.lockb", "bun.lock"}
	case Deno:
		return []string{"deno.lock"}
	case Cargo:
		return []string{"Cargo.lock"}
	case Go:
		return []string{"go.sum"}
	case Pip:
		return []string{"requirements.txt", "pyproject.toml"}
	default:
		return nil
	}
}

// hasUncommittedLockfileChanges checks if git reports uncommitted changes
// in the project directory.
func hasUncommittedLockfileChanges(dir string) bool {
	// Quick check: is this even a git repo?
	gitDir := filepath.Join(dir, ".git")
	if _, err := os.Stat(gitDir); err != nil {
		// Walk up to find .git
		parent := dir
		found := false
		for i := 0; i < 5; i++ {
			parent = filepath.Dir(parent)
			if _, err := os.Stat(filepath.Join(parent, ".git")); err == nil {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}

	// Use git status to check for uncommitted changes in lockfiles
	cmd := execCommand("git", "status", "--porcelain", dir)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		return false
	}

	return len(output) > 0
}

// formatDays returns a human-readable string for a number of days.
func formatDays(days int) string {
	if days == 1 {
		return "1 day"
	}
	if days < 30 {
		return fmt.Sprintf("%d days", days)
	}
	months := days / 30
	if months == 1 {
		return "~1 month"
	}
	return fmt.Sprintf("~%d months", months)
}

// execCommand is a wrapper around exec.Command for testability.
var execCommand = exec.Command
