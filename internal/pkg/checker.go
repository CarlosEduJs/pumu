package pkg

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

// HealthResult holds the result of a project health check.
type HealthResult struct {
	Dir     string
	PM      PackageManager
	Issues  []string
	Healthy bool
}

// CheckHealth verifies the integrity of a project's dependencies.
func CheckHealth(dir string, pm PackageManager) HealthResult {
	result := HealthResult{Dir: dir, PM: pm, Healthy: true}

	switch pm {
	case Npm:
		result = checkNodeHealth(dir, pm, "npm")
	case Pnpm:
		result = checkNodeHealth(dir, pm, "pnpm")
	case Yarn:
		result = checkNodeHealth(dir, pm, "yarn")
	case Bun:
		result = checkNodeHealth(dir, pm, "bun")
	case Cargo:
		result = checkCargoHealth(dir)
	case Go:
		result = checkGoHealth(dir)
	case Pip:
		result = checkPipHealth(dir)
	case Deno:
		result = checkNodeHealth(dir, pm, "deno")
	default:
		result.Issues = append(result.Issues, "Unknown package manager, cannot check health")
		result.Healthy = false
	}

	return result
}

// checkNodeHealth checks Node.js project health via `<pm> ls` or install --dry-run.
func checkNodeHealth(dir string, pm PackageManager, binary string) HealthResult {
	result := HealthResult{Dir: dir, PM: pm, Healthy: true}

	targetPath := dir + "/node_modules"
	if !FileExists(targetPath) && !DirExists(targetPath) {
		result.Healthy = false
		result.Issues = append(result.Issues, "node_modules not found")
		return result
	}

	// Use npm/pnpm/yarn ls to detect issues
	var cmd *exec.Cmd
	switch binary {
	case "npm":
		cmd = exec.Command("npm", "ls", "--json", "--depth=0")
	case "pnpm":
		cmd = exec.Command("pnpm", "ls", "--json", "--depth=0")
	case "yarn":
		cmd = exec.Command("yarn", "check", "--verify-tree")
	case "bun":
		// Bun doesn't have a native ls health check; try a dry install
		cmd = exec.Command("bun", "install", "--dry-run")
	case "deno":
		cmd = exec.Command("deno", "check", ".")
	default:
		cmd = exec.Command("npm", "ls", "--json", "--depth=0")
	}

	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Healthy = false
		// Try to parse npm/pnpm JSON output for specific issues
		if binary == "npm" || binary == "pnpm" {
			issues := parseNpmLsOutput(output)
			if len(issues) > 0 {
				result.Issues = issues
			} else {
				result.Issues = append(result.Issues, fmt.Sprintf("%s reports dependency issues", binary))
			}
		} else {
			result.Issues = append(result.Issues, fmt.Sprintf("%s health check failed", binary))
		}
	}

	return result
}

// parseNpmLsOutput extracts problem descriptions from npm/pnpm ls --json output.
func parseNpmLsOutput(output []byte) []string {
	var data struct {
		Problems []string `json:"problems"`
	}

	if err := json.Unmarshal(output, &data); err != nil {
		return nil
	}

	// Limit to first 5 problems to avoid noise
	if len(data.Problems) > 5 {
		count := len(data.Problems)
		data.Problems = append(data.Problems[:5], fmt.Sprintf("... and %d more issues", count-5))
	}

	return data.Problems
}

// checkCargoHealth checks Rust project health via `cargo check`.
func checkCargoHealth(dir string) HealthResult {
	result := HealthResult{Dir: dir, PM: Cargo, Healthy: true}

	targetPath := dir + "/target"
	if !DirExists(targetPath) {
		result.Healthy = false
		result.Issues = append(result.Issues, "target/ not found (never built)")
		return result
	}

	cmd := exec.Command("cargo", "check")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Healthy = false
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "error") {
				result.Issues = append(result.Issues, strings.TrimSpace(line))
				if len(result.Issues) >= 5 {
					break
				}
			}
		}
		if len(result.Issues) == 0 {
			result.Issues = append(result.Issues, "cargo check failed")
		}
	}

	return result
}

// checkGoHealth checks Go project health via `go mod verify`.
func checkGoHealth(dir string) HealthResult {
	result := HealthResult{Dir: dir, PM: Go, Healthy: true}

	cmd := exec.Command("go", "mod", "verify")
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Healthy = false
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				result.Issues = append(result.Issues, trimmed)
				if len(result.Issues) >= 5 {
					break
				}
			}
		}
		if len(result.Issues) == 0 {
			result.Issues = append(result.Issues, "go mod verify failed")
		}
	}

	return result
}

// checkPipHealth checks Python project health by verifying installed packages.
func checkPipHealth(dir string) HealthResult {
	result := HealthResult{Dir: dir, PM: Pip, Healthy: true}

	venvPath := dir + "/.venv"
	if !DirExists(venvPath) {
		result.Healthy = false
		result.Issues = append(result.Issues, ".venv not found")
		return result
	}

	// Try pip check inside the venv
	pipBin := filepath.Join(filepath.Clean(venvPath), "bin", "pip")
	cmd := exec.Command(pipBin, "check") //nolint:gosec // path is constructed from known project directory
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()

	if err != nil {
		result.Healthy = false
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			trimmed := strings.TrimSpace(line)
			if trimmed != "" {
				result.Issues = append(result.Issues, trimmed)
				if len(result.Issues) >= 5 {
					break
				}
			}
		}
		if len(result.Issues) == 0 {
			result.Issues = append(result.Issues, "pip check failed")
		}
	}

	return result
}


