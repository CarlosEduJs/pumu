package pkg

import (
	"os"
	"path/filepath"
)

type PackageManager string

const (
	Npm     PackageManager = "npm"
	Pnpm    PackageManager = "pnpm"
	Yarn    PackageManager = "yarn"
	Bun     PackageManager = "bun"
	Deno    PackageManager = "deno"
	Cargo   PackageManager = "cargo"
	Go      PackageManager = "go"
	Pip     PackageManager = "pip"
	Unknown PackageManager = "unknown"
)

func DetectManager(dir string) PackageManager {
	managers := map[PackageManager][]string{
		Bun:   {"bun.lockb", "bun.lock"},
		Pnpm:  {"pnpm-lock.yaml"},
		Yarn:  {"yarn.lock"},
		Npm:   {"package-lock.json"},
		Deno:  {"deno.json", "deno.jsonc"},
		Cargo: {"Cargo.toml"},
		Go:    {"go.mod"},
		Pip:   {"requirements.txt", "pyproject.toml"},
	}

	for mgr, files := range managers {
		for _, f := range files {
			if fileExists(filepath.Join(dir, f)) {
				return mgr
			}
		}
	}

	return Unknown
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
