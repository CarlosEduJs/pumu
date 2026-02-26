package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"pumu/internal/pkg"

	"github.com/fatih/color"
)

// RepairDir scans for projects with broken dependencies and repairs them.
func RepairDir(root string, verbose bool) error {
	color.Cyan("🔧 Scanning for projects with broken dependencies in '%s'...\n", root)

	projects, err := findProjects(root)
	if err != nil {
		return fmt.Errorf("failed to scan projects: %w", err)
	}

	if len(projects) == 0 {
		color.Green("✨ No projects found!\n")
		return nil
	}

	color.Yellow("⏱️  Found %d projects. Checking health...\n", len(projects))

	var repaired, total int

	for _, proj := range projects {
		total++
		result := pkg.CheckHealth(proj.Dir, proj.PM)

		if result.Healthy {
			if verbose {
				fmt.Printf("\n📁 %s (%s)\n", proj.Dir, proj.PM)
				color.Green("   ✅ Healthy, skipping.")
			}
			continue
		}

		// Unhealthy project — show issues and repair
		fmt.Printf("\n📁 %s (%s)\n", proj.Dir, proj.PM)
		for _, issue := range result.Issues {
			color.Red("   ❌ %s", issue)
		}

		// Remove dependency folder
		targetFolder := getTargetFolder(proj.PM)
		targetPath := filepath.Join(proj.Dir, targetFolder)

		if pkg.DirExists(targetPath) {
			fmt.Printf("   🗑️  Removing %s...\n", targetFolder)
			_, err := pkg.RemoveDirectory(targetPath)
			if err != nil {
				color.Red("   ❌ Failed to remove %s: %v", targetFolder, err)
				continue
			}
		}

		// Reinstall
		fmt.Printf("   📦 Reinstalling...\n")
		err := pkg.InstallDependencies(proj.Dir, proj.PM, true)
		if err != nil {
			color.Red("   ❌ Failed to reinstall: %v", err)
			continue
		}

		color.Green("   ✅ Repaired!")
		repaired++
	}

	fmt.Println()
	fmt.Println(strings.Repeat("-", 40))
	color.Green("🔧 Repair complete! Fixed %d/%d projects.", repaired, total)

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


