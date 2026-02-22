package scanner

import (
	"testing"
)

func TestIsIgnoredPath(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{".Trash", true},
		{".cache", true},
		{".vscode", true},
		{".git", false},
		{"my-project", false},
		{"node_modules", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isIgnoredPath(tt.path)
			if result != tt.expected {
				t.Errorf("isIgnoredPath(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}

func TestIsDeletableTarget(t *testing.T) {
	tests := []struct {
		path     string
		expected bool
	}{
		{"node_modules", true},
		{"target", true},
		{".next", true},
		{".svelte-kit", true},
		{".venv", true},
		{"dist", true},
		{"build", true},
		{"src", false},
		{"bin", false},
		{".config", false},
		{"Cargo.toml", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isDeletableTarget(tt.path)
			if result != tt.expected {
				t.Errorf("isDeletableTarget(%s) = %v, want %v", tt.path, result, tt.expected)
			}
		})
	}
}
