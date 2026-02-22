package pkg

import (
	"fmt"
	"os"
	"os/exec"
)

// InstallDependencies runs the appropriate install command based on the package manager.
func InstallDependencies(dir string, pm PackageManager, silent bool) error {
	var cmd *exec.Cmd

	switch pm {
	case Bun:
		cmd = exec.Command("bun", "install")
	case Pnpm:
		cmd = exec.Command("pnpm", "install")
	case Yarn:
		cmd = exec.Command("yarn", "install")
	case Npm:
		cmd = exec.Command("npm", "install")
	case Deno:
		cmd = exec.Command("deno", "install") // Deno 2.x supports this
	case Cargo:
		cmd = exec.Command("cargo", "build")
	case Go:
		cmd = exec.Command("go", "mod", "tidy")
	case Pip:
		cmd = exec.Command("pip", "install", "-r", "requirements.txt") // Usually requires careful venv handling but good enough for MVP
	default:
		return fmt.Errorf("unknown package manager, cannot run install")
	}

	cmd.Dir = dir

	if !silent {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		fmt.Printf("📦 Running %s install...\n", pm)
		return cmd.Run()
	}

	// Capture outputs to suppress default stdout mess
	_, err := cmd.CombinedOutput()
	return err
}
