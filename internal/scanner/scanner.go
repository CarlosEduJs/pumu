package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"sync/atomic"

	"pumu/internal/pkg"

	"github.com/fatih/color"
)

type TargetFolder struct {
	Path string
	Size int64
}

func isIgnoredPath(name string) bool {
	ignored := []string{
		".Trash", ".cache", ".npm", ".yarn", ".cargo", ".rustup",
		"Library", "AppData", "Local", "Roaming", ".vscode", ".idea",
	}
	for _, ig := range ignored {
		if name == ig {
			return true
		}
	}
	return false
}

func isDeletableTarget(name string) bool {
	targets := []string{
		"node_modules",
		"target",
		".next",
		".svelte-kit",
		".venv",
		"dist",
		"build",
	}
	for _, t := range targets {
		if name == t {
			return true
		}
	}
	return false
}

func getTargetFolder(pm pkg.PackageManager) string {
	switch pm {
	case pkg.Bun, pkg.Pnpm, pkg.Yarn, pkg.Npm:
		return "node_modules"
	case pkg.Cargo:
		return "target"
	case pkg.Pip:
		return ".venv"
	}
	return "node_modules"
}

func RefreshCurrentDir() error {
	dir := "."
	pm := pkg.DetectManager(dir)
	if pm == pkg.Unknown {
		return fmt.Errorf("could not detect package manager in current directory")
	}

	fmt.Printf("🔍 Detected package manager: %s\n", pm)

	targetFolder := getTargetFolder(pm)
	targetPath := filepath.Join(dir, targetFolder)

	if fileExists(targetPath) {
		fmt.Printf("🗑️  Removing %s...\n", targetFolder)
		duration, err := pkg.RemoveDirectory(targetPath)
		if err != nil {
			return fmt.Errorf("failed to remove %s: %v", targetFolder, err)
		}
		fmt.Printf("✅ Removed in %v\n", duration)
	} else {
		fmt.Printf("ℹ️  No %s found, skipping deletion.\n", targetFolder)
	}

	err := pkg.InstallDependencies(dir, pm, false)
	if err != nil {
		return fmt.Errorf("failed to install dependencies: %v", err)
	}

	fmt.Println("🎉 Refresh complete!")
	return nil
}

func SweepDir(root string, dryRun bool, reinstall bool) error {
	printScanMessage(dryRun, root)

	targets, err := findTargetFolders(root)
	if err != nil {
		return err
	}

	if len(targets) == 0 {
		color.Green("✨ No heavy folders found!\n")
		return nil
	}

	folders := calculateFolderSizes(targets)
	totalFreed, totalDeleted := processFolders(folders, dryRun)
	printSummary(dryRun, folders, totalFreed, totalDeleted)

	if !dryRun && reinstall {
		reinstallDependencies(folders)
	}

	return nil
}

func printScanMessage(dryRun bool, root string) {
	if dryRun {
		color.Cyan("🔎 Listing heavy dependency folders in '%s'...\n", root)
	} else {
		color.Cyan("🔎 Scanning for heavy dependency folders in '%s'...\n", root)
	}
}

func findTargetFolders(root string) ([]string, error) {
	var targets []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if d.Name() == ".git" || isIgnoredPath(d.Name()) {
				return filepath.SkipDir
			}

			if isDeletableTarget(d.Name()) {
				targets = append(targets, path)
				return filepath.SkipDir
			}
		}

		return nil
	})

	return targets, err
}

func calculateFolderSizes(targets []string) []TargetFolder {
	color.Yellow("⏱️  Found %d folders. Calculating sizes concurrently...", len(targets))

	var wg sync.WaitGroup
	var mu sync.Mutex
	var folders []TargetFolder

	sem := make(chan struct{}, 20)

	for _, tPath := range targets {
		wg.Add(1)
		go func(p string) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			size, err := dirSize(p)
			if err != nil {
				size = 0
			}

			mu.Lock()
			folders = append(folders, TargetFolder{Path: p, Size: size})
			mu.Unlock()
		}(tPath)
	}

	wg.Wait()

	sort.Slice(folders, func(i, j int) bool {
		return folders[i].Size > folders[j].Size
	})

	return folders
}

func processFolders(folders []TargetFolder, dryRun bool) (int64, int64) {
	var totalFreed int64
	var deletedWg sync.WaitGroup
	var totalDeleted int64

	fmt.Println()
	color.Set(color.FgWhite, color.Underline)
	fmt.Printf("%-80s | %s\n", "Folder Path", "Size")
	color.Unset()

	sem := make(chan struct{}, 20)

	for _, folder := range folders {
		printFolderInfo(folder)
		totalFreed += folder.Size

		if !dryRun {
			deletedWg.Add(1)
			go func(p string, s int64) {
				defer deletedWg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				_, err := pkg.RemoveDirectory(p)
				if err == nil {
					atomic.AddInt64(&totalDeleted, s)
				}
			}(folder.Path, folder.Size)
		}
	}

	if !dryRun {
		color.Yellow("\n🗑️  Deleting folders concurrently...")
		deletedWg.Wait()
	}

	return totalFreed, totalDeleted
}

func printFolderInfo(folder TargetFolder) {
	sizeMB := float64(folder.Size) / 1024 / 1024
	formattedSize := formatSize(folder.Size)

	var sizeStr string
	if sizeMB > 1000 {
		sizeStr = color.RedString(fmt.Sprintf("%10s 🚨", formattedSize))
	} else if sizeMB > 100 {
		sizeStr = color.YellowString(fmt.Sprintf("%10s ⚠️", formattedSize))
	} else {
		sizeStr = color.GreenString(fmt.Sprintf("%10s", formattedSize))
	}

	displayPath := folder.Path
	if len(displayPath) > 80 {
		displayPath = "..." + displayPath[len(displayPath)-77:]
	}

	fmt.Printf("%-80s | %s\n", displayPath, sizeStr)
}

func printSummary(dryRun bool, folders []TargetFolder, totalFreed, totalDeleted int64) {
	fmt.Println(strings.Repeat("-", 100))
	if dryRun {
		color.Green("📋 List complete! Found %d heavy folders.", len(folders))
		color.Cyan("💾 Total space that can be freed: %s\n", formatSize(totalFreed))
	} else {
		color.Green("🧹 Sweep complete! Processed %d heavy folders.", len(folders))
		color.Cyan("💾 Total space actually freed: %s\n", formatSize(totalDeleted))
	}
}

func reinstallDependencies(folders []TargetFolder) {
	color.Yellow("\n⚙️  Reinstalling dependencies sequentially...")
	seen := make(map[string]bool)

	for _, folder := range folders {
		baseDir := filepath.Dir(folder.Path)
		if seen[baseDir] {
			continue
		}
		seen[baseDir] = true

		pm := pkg.DetectManager(baseDir)
		if pm != pkg.Unknown {
			fmt.Printf("📦 Reinstalling for %s (%s)...\n", baseDir, pm)
			err := pkg.InstallDependencies(baseDir, pm, true)
			if err != nil {
				color.Red("❌ Failed to reinstall %s: %v", baseDir, err)
			} else {
				color.Green("✅ Reinstalled %s", baseDir)
			}
		}
	}
	color.Green("🎉 All target reinstallations complete!")
}

func dirSize(path string) (int64, error) {
	var size int64
	err := filepath.WalkDir(path, func(_ string, d os.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if !d.IsDir() {
			info, err := d.Info()
			if err == nil {
				size += info.Size()
			}
		}
		return nil
	})
	return size, err
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// formatSize converts a byte count into a human-readable string (KB, MB, GB, etc.)
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	// "KMGTPE" represents Kilo, Mega, Giga, Tera, Peta, Exa
	return fmt.Sprintf("%.2f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
