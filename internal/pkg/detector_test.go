package pkg

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDetectManager(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "pumu-test-detect")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	tests := []struct {
		name         string
		lockfileName string
		expected     PackageManager
	}{
		{"Node project with npm", "package-lock.json", Npm},
		{"Node project with yarn", "yarn.lock", Yarn},
		{"Node project with pnpm", "pnpm-lock.yaml", Pnpm},
		{"Bun project", "bun.lockb", Bun},
		{"Deno project", "deno.json", Deno},
		{"Rust Cargo project", "Cargo.toml", Cargo},
		{"Go project", "go.mod", Go},
		{"Python project pip", "requirements.txt", Pip},
		{"Python pyproject", "pyproject.toml", Pip},
		{"Unknown project", "random.txt", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caseDir := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(caseDir, 0755)
			if err != nil {
				t.Fatalf("failed to create dir %s: %v", caseDir, err)
			}

			lockfilePath := filepath.Join(caseDir, tt.lockfileName)
			file, err := os.Create(lockfilePath)
			if err != nil {
				t.Fatalf("failed to create fake lock file %s: %v", lockfilePath, err)
			}
			file.Close()

			pm := DetectManager(caseDir)
			if pm != tt.expected {
				t.Errorf("DetectManager() = %v, want %v", pm, tt.expected)
			}
		})
	}
}
