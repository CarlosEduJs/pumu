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
	defer func() { _ = os.RemoveAll(tempDir) }()

	tests := []struct {
		name         string
		lockfileName string
		content      string
		expected     PackageManager
	}{
		{"Node project with npm", "package-lock.json", "", Npm},
		{"Node project with yarn", "yarn.lock", "", Yarn},
		{"Node project with pnpm", "pnpm-lock.yaml", "", Pnpm},
		{"Bun project", "bun.lockb", "", Bun},
		{"Deno project", "deno.json", "", Deno},
		{"Rust Cargo project", "Cargo.toml", "", Cargo},
		{"Go project", "go.mod", "", Go},
		{"Python project pip", "requirements.txt", "", Pip},
		{"Python pyproject pip", "pyproject.toml", "[build-system]\nrequires = [\"setuptools\"]\n", Pip},
		{"Poetry project", "pyproject.toml", "[tool.poetry]\nname = \"demo\"\n", Poetry},
		{"Composer project", "composer.json", "{\"name\":\"demo/app\"}", Composer},
		{"Mix project", "mix.exs", "defmodule Demo.MixProject do\nend\n", Mix},
		{"Unknown project", "random.txt", "", Unknown},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			caseDir := filepath.Join(tempDir, tt.name)
			err := os.MkdirAll(caseDir, 0o750) //nolint:gosec // test directory
			if err != nil {
				t.Fatalf("failed to create dir %s: %v", caseDir, err)
			}

			lockfilePath := filepath.Join(caseDir, tt.lockfileName)
			file, err := os.Create(lockfilePath) //nolint:gosec // controlled test path
			if err != nil {
				t.Fatalf("failed to create fake lock file %s: %v", lockfilePath, err)
			}
			if tt.content != "" {
				if _, err := file.WriteString(tt.content); err != nil {
					_ = file.Close()
					t.Fatalf("failed to write file: %v", err)
				}
			}
			if err := file.Close(); err != nil {
				t.Fatalf("failed to close file: %v", err)
			}

			pm := DetectManager(caseDir)
			if pm != tt.expected {
				t.Errorf("DetectManager() = %v, want %v", pm, tt.expected)
			}
		})
	}
}
