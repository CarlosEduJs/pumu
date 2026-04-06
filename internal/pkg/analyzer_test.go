package pkg

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestAnalyzeFolderPruneDefaultScore(t *testing.T) {
	root := t.TempDir()
	projectDir := filepath.Join(root, "project")
	if err := os.MkdirAll(projectDir, 0o750); err != nil {
		t.Fatalf("failed to create project dir: %v", err)
	}

	// Create lockfile with modification time 15 days ago (not too recent, not too old)
	lockfilePath := filepath.Join(projectDir, "package-lock.json")
	if err := os.WriteFile(lockfilePath, []byte("{}"), 0o640); err != nil {
		t.Fatalf("failed to write lockfile: %v", err)
	}
	fifteenDaysAgo := time.Now().Add(-15 * 24 * time.Hour)
	if err := os.Chtimes(lockfilePath, fifteenDaysAgo, fifteenDaysAgo); err != nil {
		t.Fatalf("failed to set lockfile time: %v", err)
	}

	folderPath := filepath.Join(projectDir, "node_modules")
	if err := os.MkdirAll(folderPath, 0o750); err != nil {
		t.Fatalf("failed to create node_modules: %v", err)
	}
	// Set access time to 15 days ago to avoid "recently accessed" heuristic
	if err := os.Chtimes(folderPath, fifteenDaysAgo, fifteenDaysAgo); err != nil {
		t.Fatalf("failed to set folder time: %v", err)
	}

	result := AnalyzeFolder(folderPath, 0)
	if result.Score != 45 {
		t.Fatalf("expected score 45, got %d (reason: %s)", result.Score, result.Reason)
	}
}
