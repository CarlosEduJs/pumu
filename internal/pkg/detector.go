// Package pkg provides utilities for detecting, installing, checking,
// and analyzing dependencies across multiple package managers.
package pkg

import (
	"os"
	"path/filepath"
	"strings"
)

// PackageManager represents the type of package manager detected in a project.
type PackageManager string

// Supported package managers.
const (
	Npm      PackageManager = "npm"
	Pnpm     PackageManager = "pnpm"
	Yarn     PackageManager = "yarn"
	Bun      PackageManager = "bun"
	Deno     PackageManager = "deno"
	Cargo    PackageManager = "cargo"
	Go       PackageManager = "go"
	Poetry   PackageManager = "poetry"
	Pip      PackageManager = "pip"
	Composer PackageManager = "composer"
	Mix      PackageManager = "mix"
	Unknown  PackageManager = "unknown"
)

// DetectManager identifies the package manager used in dir by checking for lock files.
func DetectManager(dir string) PackageManager {
	orderedManagers := []struct {
		manager PackageManager
		files   []string
		match   func(string) bool
	}{
		{manager: Bun, files: []string{"bun.lockb", "bun.lock"}},
		{manager: Pnpm, files: []string{"pnpm-lock.yaml"}},
		{manager: Yarn, files: []string{"yarn.lock"}},
		{manager: Npm, files: []string{"package-lock.json"}},
		{manager: Deno, files: []string{"deno.json", "deno.jsonc"}},
		{manager: Cargo, files: []string{"Cargo.toml"}},
		{manager: Go, files: []string{"go.mod"}},
		{manager: Poetry, files: []string{"pyproject.toml"}, match: isPoetryProject},
		{manager: Pip, files: []string{"requirements.txt", "pyproject.toml"}},
		{manager: Composer, files: []string{"composer.json", "composer.lock"}},
		{manager: Mix, files: []string{"mix.exs"}},
	}

	for _, entry := range orderedManagers {
		for _, f := range entry.files {
			path := filepath.Join(dir, f)
			if !fileExists(path) {
				continue
			}
			if entry.match != nil && !entry.match(path) {
				continue
			}
			return entry.manager
		}
	}

	return Unknown
}

func isPoetryProject(pyprojectPath string) bool {
	content, err := os.ReadFile(pyprojectPath)
	if err != nil {
		return false
	}
	return strings.Contains(string(content), "[tool.poetry]")
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
