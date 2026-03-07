package scanner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pumu/internal/pkg"

	"pumu/internal/ui"

	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

// RepairDir scans for projects with broken dependencies and repairs them.
func RepairDir(root string, verbose bool) error {
	fmt.Println(ui.InfoStyle.Render(fmt.Sprintf("🔧 Scanning for projects with broken dependencies in '%s'...", root)))

	var projects []project
	var scanErr error
	err := spinner.New().
		Title(ui.InfoStyle.Render("Scanning for projects...")).
		Action(func() {
			projects, scanErr = findProjects(root)
		}).
		Run()
	if err != nil {
		return fmt.Errorf("failed to run scanner: %w", err)
	}

	if scanErr != nil {
		return fmt.Errorf("failed to scan projects: %w", scanErr)
	}

	if len(projects) == 0 {
		fmt.Println(ui.SuccessStyle.Render("✨ No projects found!"))
		return nil
	}

	var repaired, total int
	fmt.Println()

	for _, proj := range projects {
		total++

		// Project header
		fmt.Printf("📁 %s (%s)\n", ui.BoldStyle.Render(proj.Dir), ui.InfoStyle.Render(string(proj.PM)))

		var result pkg.HealthResult
		if err := spinner.New().
			Title(ui.SubtextStyle.Render("   Checking health...")).
			Action(func() {
				result = pkg.CheckHealth(proj.Dir, proj.PM)
			}).
			Run(); err != nil {
			fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("   ❌ Failed to check health: %v", err)))
			continue
		}

		if result.Healthy {
			if verbose {
				fmt.Println(ui.SuccessStyle.Render("   ✅ Healthy, skipping."))
			}
			continue
		}

		// Unhealthy project — show issues and repair
		for _, issue := range result.Issues {
			fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("   ❌ %s", issue)))
		}

		// Remove dependency folder
		targetFolder := getTargetFolder(proj.PM)
		targetPath := filepath.Join(proj.Dir, targetFolder)

		if pkg.DirExists(targetPath) {
			var rmErr error
			if err := spinner.New().
				Title(ui.WarningStyle.Render(fmt.Sprintf("   🗑️  Removing %s...", targetFolder))).
				Action(func() {
					_, rmErr = pkg.RemoveDirectory(targetPath)
				}).
				Run(); err != nil || rmErr != nil {
				if err != nil {
					fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("   ❌ Failed to remove folder: %v", err)))
				} else {
					fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("   ❌ Failed to remove %s: %v", targetFolder, rmErr)))
				}
				continue
			}
		}

		// Reinstall
		installErr := errors.New("not started")
		if err := spinner.New().
			Title(ui.InfoStyle.Render("   📦 Reinstalling dependencies...")).
			Action(func() {
				installErr = pkg.InstallDependencies(proj.Dir, proj.PM, true)
			}).
			Run(); err != nil {
			fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("   ❌ Failed to reinstall: %v", err)))
			continue
		}

		if installErr != nil {
			fmt.Println(ui.AlertStyle.Render(fmt.Sprintf("   ❌ Failed to reinstall: %v", installErr)))
			continue
		}

		fmt.Println(ui.SuccessStyle.Render("   ✅ Repaired!"))
		repaired++
	}

	fmt.Println()
	fmt.Println(lipgloss.NewStyle().Foreground(ui.ColorSubtext).Render(strings.Repeat("-", 40)))
	fmt.Println(ui.SuccessStyle.Render(fmt.Sprintf("🔧 Repair complete! Fixed %d/%d projects.", repaired, total)))

	return nil
}

// project represents a detected project directory with its package manager.
type project struct {
	Dir string
	PM  pkg.PackageManager
}

// findProjects recursively scans for directories containing lockfiles/manifests.
// WalkDir is sequential, so no mutex is needed.
func findProjects(root string) ([]project, error) {
	var projects []project

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if d.Name() == ".git" || isIgnoredPath(d.Name()) || isDeletableTarget(d.Name()) {
				return filepath.SkipDir
			}

			pm := pkg.DetectManager(path)
			if pm != pkg.Unknown {
				projects = append(projects, project{Dir: path, PM: pm})
			}
		}

		return nil
	})

	return projects, err
}
